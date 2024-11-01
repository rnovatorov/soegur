package postgresadapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/soegur"
)

type taskExecutor func(
	ctx context.Context, taskType string, taskInput *structpb.Struct,
) (*structpb.Value, error)

type taskExecutionResultHandler interface {
	HandleTaskExecutionResult(context.Context, *soegur.TaskExecutionResult) error
}

type RiverTaskWorker struct {
	client *river.Client[pgx.Tx]
}

type RiverTaskWorkerParams struct {
	Pool                   *pgxpool.Pool
	Executor               taskExecutor
	ExecutionResultHandler taskExecutionResultHandler
}

func StartRiverTaskWorker(
	ctx context.Context, p RiverTaskWorkerParams,
) (*RiverTaskWorker, error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &riverTaskWorker{
		executor:               p.Executor,
		executionResultHandler: p.ExecutionResultHandler,
	})

	client, err := river.NewClient(riverpgxv5.New(p.Pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, fmt.Errorf("new river client: %w", err)
	}

	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("start river client: %w", err)
	}

	return &RiverTaskWorker{
		client: client,
	}, nil
}

func (w *RiverTaskWorker) Stop() {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	w.client.Stop(ctx)
}

type riverTaskWorkerArgs struct {
	SagaID       string           `json:"saga_id"`
	StepID       string           `json:"step_id"`
	Compensation bool             `json:"compensation"`
	TaskType     string           `json:"task_type"`
	TaskInput    *structpb.Struct `json:"task_input"`
}

func (riverTaskWorkerArgs) Kind() string { return "saga_task" }

type riverTaskWorker struct {
	river.WorkerDefaults[riverTaskWorkerArgs]
	executor               taskExecutor
	executionResultHandler taskExecutionResultHandler
}

func (w *riverTaskWorker) Work(
	ctx context.Context, job *river.Job[riverTaskWorkerArgs],
) error {
	var result *soegur.TaskExecutionResult

	output, err := w.executor(ctx, job.Args.TaskType, job.Args.TaskInput)
	if err != nil {
		result = &soegur.TaskExecutionResult{
			SagaID:       job.Args.SagaID,
			StepID:       job.Args.StepID,
			Compensation: job.Args.Compensation,
			Output:       nil,
			Error:        err.Error(),
		}
	} else {
		result = &soegur.TaskExecutionResult{
			SagaID:       job.Args.SagaID,
			StepID:       job.Args.StepID,
			Compensation: job.Args.Compensation,
			Output:       output,
			Error:        "",
		}
	}

	return w.executionResultHandler.HandleTaskExecutionResult(ctx, result)
}
