package model

import (
	"fmt"
	"sort"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/soegur/internal/api/sagaeventspb"
)

type SagaAggregate = eventsource.Aggregate[Saga, *Saga]

type Saga struct {
	steps           map[string]*Step
	dependencyGraph *dependencyGraph
	config          *structpb.Struct
	begun           bool
	aborted         bool
	ended           bool
}

func (s *Saga) Config() *structpb.Struct {
	return s.config
}

func (s *Saga) Copy() *Saga {
	steps := make(map[string]*Step, len(s.steps))
	for id, step := range s.steps {
		steps[id] = step.copy()
	}
	return &Saga{
		steps:           s.steps,
		dependencyGraph: s.dependencyGraph,
		config:          s.config,
		begun:           s.begun,
		aborted:         s.aborted,
		ended:           s.ended,
	}
}

func (s *Saga) Begun() bool {
	return s.begun
}

func (s *Saga) Ended() bool {
	return s.ended
}

func (s *Saga) Aborted() bool {
	return s.aborted
}

func (s *Saga) Step(id string) *Step {
	return s.steps[id]
}

func (s *Saga) InProgress() bool {
	return s.begun && !s.ended && !s.aborted
}

func (s *Saga) StepTransitiveDependencies(stepID string) []string {
	return s.dependencyGraph.transitiveDependencies(stepID)
}

func (s *Saga) StepInProgress(stepID string) bool {
	step, ok := s.steps[stepID]
	if !ok {
		return false
	}
	return step.inProgress()
}

func (s *Saga) StepEnded(stepID string) bool {
	step, ok := s.steps[stepID]
	if !ok {
		return false
	}
	return step.ended
}

func (s *Saga) TaskType(stepID string) string {
	step, ok := s.steps[stepID]
	if !ok {
		return ""
	}
	return step.task.spec.Type
}

func (s *Saga) CompensatingTaskType(stepID string) string {
	step, ok := s.steps[stepID]
	if !ok {
		return ""
	}
	return step.compensatingTask.spec.Type
}

func (s *Saga) TaskInput(stepID string) (*structpb.Struct, error) {
	step, ok := s.steps[stepID]
	if !ok {
		return nil, nil
	}

	return s.evalInput(step.task)
}

func (s *Saga) CompensatingTaskInput(stepID string) (*structpb.Struct, error) {
	step, ok := s.steps[stepID]
	if !ok {
		return nil, nil
	}

	return s.evalInput(step.compensatingTask)
}

func (s *Saga) evalInput(task *Task) (*structpb.Struct, error) {
	outputs := make(map[string]any)

	for _, dep := range task.dependencies {
		outputs[dep] = s.steps[dep].output.AsInterface()
	}

	input, err := task.EvalInput(EvalContext{
		config:  s.config.AsMap(),
		outputs: outputs,
		data:    nil,
	})
	if err != nil {
		return nil, err
	}

	return structpb.NewStruct(input)
}

func (s *Saga) ProcessCommand(
	command eventsource.Command,
) (eventsource.StateChanges, error) {
	switch cmd := command.(type) {
	case BeginSaga:
		return s.processBegin(cmd)
	case EndStep:
		return s.processEndStep(cmd)
	case AbortStep:
		return s.processAbortStep(cmd)
	case EndStepCompensation:
		return s.processEndStepCompensation(cmd)
	case AbortStepCompensation:
		return s.processAbortStepCompensation(cmd)
	default:
		return nil, fmt.Errorf("%w: %T", eventsource.ErrCommandUnknown, cmd)
	}
}

func (s *Saga) processBegin(cmd BeginSaga) (eventsource.StateChanges, error) {
	if s.begun {
		return nil, ErrSagaAlreadyBegun
	}

	if len(cmd.Spec.GetSteps()) == 0 {
		return nil, ErrSagaStepsMissing
	}

	graph, err := buildDependencyGraph(cmd.Spec.Steps)
	if err != nil {
		return nil, fmt.Errorf("build dependency graph: %w", err)
	}

	for _, spec := range cmd.Spec.Steps {
		if _, err := newStep(spec, graph); err != nil {
			return nil, fmt.Errorf("new step: %w", err)
		}
	}

	stateChanges := eventsource.StateChanges{&sagaeventspb.SagaBegun{
		Spec:   cmd.Spec,
		Config: cmd.Config,
	}}

	eventsource.Given(s, stateChanges, func() {
		for _, id := range s.beginnableSteps() {
			stateChanges = append(stateChanges, &sagaeventspb.StepBegun{
				Id: id,
			})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processEndStep(cmd EndStep) (eventsource.StateChanges, error) {
	if s.ended {
		return nil, ErrSagaEnded
	}

	step, ok := s.steps[cmd.ID]
	if !ok {
		return nil, ErrStepNotDefined
	}
	if !step.begun {
		return nil, ErrStepNotBegun
	}
	if step.ended {
		return nil, ErrStepAlreadyEnded
	}
	if step.aborted {
		return nil, ErrStepAborted
	}

	if cmd.Output == nil {
		return nil, ErrStepOutputMissing
	}

	stateChanges := eventsource.StateChanges{&sagaeventspb.StepEnded{
		Id:     cmd.ID,
		Output: cmd.Output,
	}}

	eventsource.Given(s, stateChanges, func() {
		if s.aborted {
			for _, id := range s.beginnableCompensations() {
				stateChanges = append(stateChanges,
					&sagaeventspb.StepCompensationBegun{
						Id: id,
					})
			}
			return
		}
		taskBegun := false
		for _, id := range s.beginnableSteps() {
			stateChanges = append(stateChanges, &sagaeventspb.StepBegun{
				Id: id,
			})
			taskBegun = true
		}
		if !taskBegun && !s.hasStepsInProgress() {
			stateChanges = append(stateChanges, &sagaeventspb.SagaEnded{})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processAbortStep(cmd AbortStep) (eventsource.StateChanges, error) {
	if s.ended {
		return nil, ErrSagaEnded
	}

	t, ok := s.steps[cmd.ID]
	if !ok {
		return nil, ErrStepNotDefined
	}
	if !t.begun {
		return nil, ErrStepNotBegun
	}
	if t.ended {
		return nil, ErrStepEnded
	}
	if t.aborted {
		return nil, ErrStepAlreadyAborted
	}

	if cmd.Reason == "" {
		return nil, ErrStepAbortReasonMissing
	}

	stateChanges := eventsource.StateChanges{&sagaeventspb.StepAborted{
		Id:     cmd.ID,
		Reason: cmd.Reason,
	}}

	eventsource.Given(s, stateChanges, func() {
		compensationBegun := false
		for _, id := range s.beginnableCompensations() {
			stateChanges = append(stateChanges,
				&sagaeventspb.StepCompensationBegun{
					Id: id,
				})
			compensationBegun = true
		}
		if !compensationBegun && !s.hasStepsInProgress() && !s.hasCompensationsInProgress() {
			stateChanges = append(stateChanges, &sagaeventspb.SagaEnded{})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processEndStepCompensation(
	cmd EndStepCompensation,
) (eventsource.StateChanges, error) {
	if s.ended {
		return nil, ErrSagaEnded
	}

	step, ok := s.steps[cmd.ID]
	if !ok {
		return nil, ErrStepNotDefined
	}
	if !step.compensationBegun {
		return nil, ErrStepCompensationNotBegun
	}
	if step.compensationEnded {
		return nil, ErrStepCompensationAlreadyEnded
	}

	stateChanges := eventsource.StateChanges{&sagaeventspb.StepCompensationEnded{
		Id: cmd.ID,
	}}

	eventsource.Given(s, stateChanges, func() {
		compensationBegun := false
		for _, id := range s.beginnableCompensations() {
			stateChanges = append(stateChanges,
				&sagaeventspb.StepCompensationBegun{
					Id: id,
				})
			compensationBegun = true
		}
		if !compensationBegun && !s.hasStepsInProgress() && !s.hasCompensationsInProgress() {
			stateChanges = append(stateChanges, &sagaeventspb.SagaEnded{})
		}
	})

	return stateChanges, nil
}

func (s *Saga) processAbortStepCompensation(
	cmd AbortStepCompensation,
) (eventsource.StateChanges, error) {
	// TODO
	panic(cmd)
}

func (s *Saga) hasStepsInProgress() bool {
	for _, step := range s.steps {
		if step.inProgress() {
			return true
		}
	}
	return false
}

func (s *Saga) hasCompensationsInProgress() bool {
	for _, step := range s.steps {
		if step.compensationInProgress() {
			return true
		}
	}
	return false
}

func (s *Saga) beginnableSteps() (ids []string) {
	if !s.begun {
		return nil
	}
	if s.aborted {
		return nil
	}
	for id, step := range s.steps {
		if step.begun {
			continue
		}
		depsSatisfied := true
		for _, dep := range step.spec.Dependencies {
			if !s.steps[dep].ended {
				depsSatisfied = false
				break
			}
		}
		if depsSatisfied {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (s *Saga) beginnableCompensations() (ids []string) {
	if !s.aborted {
		return nil
	}
	for id, step := range s.steps {
		if !step.begun {
			continue
		}
		if step.aborted {
			continue
		}
		if !step.ended {
			continue
		}
		if step.compensationBegun {
			continue
		}
		depsSatisfied := true
		for dep := range s.steps {
			if !s.dependencyGraph.dependsOn(dep, id) {
				continue
			}
			if !s.steps[dep].begun {
				continue
			}
			if s.steps[dep].aborted {
				continue
			}
			if s.steps[dep].compensationEnded {
				continue
			}
			depsSatisfied = false
			break
		}
		if depsSatisfied {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (s *Saga) ApplyStateChange(stateChange eventsource.StateChange) {
	switch sc := stateChange.(type) {
	case *sagaeventspb.SagaBegun:
		s.applySagaBegun(sc)
	case *sagaeventspb.SagaEnded:
		s.applySagaEnded(sc)
	case *sagaeventspb.StepBegun:
		s.applyStepBegun(sc)
	case *sagaeventspb.StepEnded:
		s.applyStepEnded(sc)
	case *sagaeventspb.StepAborted:
		s.applyStepAborted(sc)
	case *sagaeventspb.StepCompensationBegun:
		s.applyStepCompensationBegun(sc)
	case *sagaeventspb.StepCompensationEnded:
		s.applyStepCompensationEnded(sc)
	default:
		panic(fmt.Sprintf("unexpected state change: %T", sc))
	}
}

func (s *Saga) applySagaBegun(sc *sagaeventspb.SagaBegun) {
	s.config = sc.Config
	s.begun = true
	s.dependencyGraph, _ = buildDependencyGraph(sc.Spec.Steps)
	s.steps = make(map[string]*Step, len(sc.Spec.Steps))
	for _, spec := range sc.Spec.Steps {
		step, err := newStep(spec, s.dependencyGraph)
		if err != nil {
			panic(err)
		}
		s.steps[spec.Id] = step
	}
}

func (s *Saga) applySagaEnded(*sagaeventspb.SagaEnded) {
	s.ended = true
}

func (s *Saga) applyStepBegun(sc *sagaeventspb.StepBegun) {
	s.steps[sc.Id].begun = true
}

func (s *Saga) applyStepEnded(sc *sagaeventspb.StepEnded) {
	s.steps[sc.Id].output = sc.Output
	s.steps[sc.Id].ended = true
}

func (s *Saga) applyStepAborted(sc *sagaeventspb.StepAborted) {
	s.steps[sc.Id].aborted = true
	s.steps[sc.Id].abortReason = sc.Reason
	if !s.aborted {
		s.aborted = true
	}
}

func (s *Saga) applyStepCompensationBegun(sc *sagaeventspb.StepCompensationBegun) {
	s.steps[sc.Id].compensationBegun = true
}

func (s *Saga) applyStepCompensationEnded(sc *sagaeventspb.StepCompensationEnded) {
	s.steps[sc.Id].compensationEnded = true
}
