package postgresadapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/rnovatorov/soegur"
)

type RiverTaskQueue struct {
	client *river.Client[pgx.Tx]
}

func StartRiverTaskQueue(pool *pgxpool.Pool) (*RiverTaskQueue, error) {
	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		return nil, err
	}

	return &RiverTaskQueue{
		client: client,
	}, nil
}

func (q *RiverTaskQueue) Stop() {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q.client.Stop(ctx)
}

func (q *RiverTaskQueue) EnqueueTask(
	ctx context.Context, req *soegur.EnqueueTaskRequest,
) error {
	if _, err := q.client.Insert(ctx, &riverTaskWorkerArgs{
		SagaID:       req.SagaID,
		StepID:       req.StepID,
		Compensation: req.Compensation,
		TaskType:     req.TaskType,
		TaskInput:    req.TaskInput,
	}, nil); err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}
