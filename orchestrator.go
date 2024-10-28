package soegur

import (
	"context"
	"fmt"

	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/rnovatorov/go-eventsource/pkg/eventstore"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/soegur/internal/api/sagaeventspb"
	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
	"github.com/rnovatorov/soegur/internal/api/sagataskqueuepb"
	"github.com/rnovatorov/soegur/internal/model"
)

type TaskQueue interface {
	EnqueueTask(context.Context, *sagataskqueuepb.EnqueueRequest) error
}

type OrchestratorService struct {
	sagaRepository *eventsource.AggregateRepository[model.Saga, *model.Saga]
	taskQueue      TaskQueue
}

type OrchestratorServiceParams struct {
	EventStore eventstore.Interface
	TaskQueue  TaskQueue
}

func NewOrchestratorService(p OrchestratorServiceParams) *OrchestratorService {
	return &OrchestratorService{
		sagaRepository: eventsource.NewAggregateRepository[model.Saga](p.EventStore),
		taskQueue:      p.TaskQueue,
	}
}

func (s *OrchestratorService) BeginSaga(
	ctx context.Context, sagaID string, spec *sagaspecpb.Saga,
	config *structpb.Struct,
) error {
	_, err := s.sagaRepository.Create(ctx, sagaID, model.BeginSaga{
		Spec:   spec,
		Config: config,
	})
	return err
}

func (s *OrchestratorService) HandleSagaEvent(
	ctx context.Context, event *eventstore.Event,
) error {
	data, err := event.Data.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	switch d := data.(type) {
	case *sagaeventspb.StepBegun:
		return s.handleStepBegun(ctx, event, d)
	case *sagaeventspb.StepCompensationBegun:
		return s.handleStepCompensationBegun(ctx, event, d)
	default:
		return nil
	}
}

func (s *OrchestratorService) handleStepBegun(
	ctx context.Context, e *eventstore.Event, d *sagaeventspb.StepBegun,
) error {
	saga, err := s.sagaRepository.Get(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("get saga: %w", err)
	}

	input, err := saga.Root().TaskInput(d.Id)
	if err != nil {
		if err := saga.ProcessCommand(ctx, model.AbortStep{
			ID:     d.Id,
			Reason: fmt.Sprintf("task input: %v", err),
		}); err != nil {
			return fmt.Errorf("abort step: %w", err)
		}
		if err := s.sagaRepository.Save(ctx, saga); err != nil {
			return fmt.Errorf("save saga: %w", err)
		}
		return nil
	}

	return s.taskQueue.EnqueueTask(ctx, &sagataskqueuepb.EnqueueRequest{
		SagaId:       saga.ID(),
		StepId:       d.Id,
		Compensation: false,
		TaskType:     saga.Root().TaskType(d.Id),
		TaskInput:    input,
	})
}

func (s *OrchestratorService) handleStepCompensationBegun(
	ctx context.Context, e *eventstore.Event, d *sagaeventspb.StepCompensationBegun,
) error {
	saga, err := s.sagaRepository.Get(ctx, e.AggregateID)
	if err != nil {
		return fmt.Errorf("get saga: %w", err)
	}

	input, err := saga.Root().CompensatingTaskInput(d.Id)
	if err != nil {
		if err := saga.ProcessCommand(ctx, model.AbortStepCompensation{
			ID:     d.Id,
			Reason: fmt.Sprintf("compensating task input: %v", err),
		}); err != nil {
			return fmt.Errorf("abort step compensation: %w", err)
		}
		if err := s.sagaRepository.Save(ctx, saga); err != nil {
			return fmt.Errorf("save saga: %w", err)
		}
		return nil
	}

	return s.taskQueue.EnqueueTask(ctx, &sagataskqueuepb.EnqueueRequest{
		SagaId:       saga.ID(),
		StepId:       d.Id,
		Compensation: true,
		TaskType:     saga.Root().CompensatingTaskType(d.Id),
		TaskInput:    input,
	})
}

func (s *OrchestratorService) HandleTaskExecutionResult(
	ctx context.Context, res *sagataskqueuepb.ExecutionResult,
) error {
	if res.Compensation {
		return s.handleCompensatingTaskExecutionResult(ctx, res)
	}
	return s.handleTaskExecutionResult(ctx, res)
}

func (s *OrchestratorService) handleTaskExecutionResult(
	ctx context.Context, res *sagataskqueuepb.ExecutionResult,
) error {
	if res.Error != "" {
		_, err := s.sagaRepository.Update(ctx, res.SagaId, model.AbortStep{
			ID:     res.StepId,
			Reason: res.Error,
		})
		return err
	}

	_, err := s.sagaRepository.Update(ctx, res.SagaId, model.EndStep{
		ID:     res.StepId,
		Output: res.Output,
	})
	return err
}

func (s *OrchestratorService) handleCompensatingTaskExecutionResult(
	ctx context.Context, res *sagataskqueuepb.ExecutionResult,
) error {
	if res.Error != "" {
		_, err := s.sagaRepository.Update(ctx, res.SagaId,
			model.AbortStepCompensation{
				ID:     res.StepId,
				Reason: res.Error,
			})
		return err
	}

	_, err := s.sagaRepository.Update(ctx, res.SagaId, model.EndStepCompensation{
		ID: res.StepId,
	})
	return err
}
