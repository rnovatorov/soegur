package model

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
)

type Step struct {
	spec              *sagaspecpb.Step
	task              *Task
	compensatingTask  *Task
	begun             bool
	ended             bool
	output            *structpb.Value
	aborted           bool
	abortReason       string
	compensationBegun bool
	compensationEnded bool
}

func newStep(spec *sagaspecpb.Step, graph *dependencyGraph) (*Step, error) {
	deps := graph.transitiveDependencies(spec.GetId())

	task, err := newTask(spec.GetTask(), deps)
	if err != nil {
		return nil, fmt.Errorf("new task: %w", err)
	}

	compDeps := append(deps, spec.GetId())

	compensatingTask, err := newTask(spec.GetCompensatingTask(), compDeps)
	if err != nil {
		return nil, fmt.Errorf("new compensating task: %w", err)
	}

	return &Step{
		spec:             spec,
		task:             task,
		compensatingTask: compensatingTask,
	}, nil
}

func (s *Step) inProgress() bool {
	return s.begun && !s.aborted && !s.ended && !s.compensationBegun
}

func (s *Step) compensationInProgress() bool {
	return s.compensationBegun && !s.compensationEnded
}
