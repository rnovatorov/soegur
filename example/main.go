package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"

	"github.com/rnovatorov/go-eventsource/pkg/eventstore/eventstorepostgres"
	"github.com/rnovatorov/soegur"
	"github.com/rnovatorov/soegur/example/postgresadapter"
	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("new database pool: %w", err)
	}
	defer pool.Close()

	eventStore := eventstorepostgres.Start(pool,
		eventstorepostgres.WithLogger(logger))
	defer eventStore.Stop()

	taskQueue, err := postgresadapter.StartRiverTaskQueue(pool)
	if err != nil {
		return fmt.Errorf("start river task queue: %w", err)
	}
	defer taskQueue.Stop()

	orchestrator := soegur.NewOrchestratorService(soegur.OrchestratorServiceParams{
		EventStore: eventStore,
		TaskQueue:  taskQueue,
	})

	if err := eventStore.Subscribe(
		ctx, "orchestrator", orchestrator.HandleSagaEvent,
	); err != nil {
		return fmt.Errorf("subscribe orchestrator event handler: %w", err)
	}

	taskWorker, err := postgresadapter.StartRiverTaskWorker(ctx,
		postgresadapter.RiverTaskWorkerParams{
			Pool:                   pool,
			Executor:               executeTask,
			ExecutionResultHandler: orchestrator,
		})
	if err != nil {
		return fmt.Errorf("start task worker: %w", err)
	}
	defer taskWorker.Stop()

	id := uuid.NewString()

	spec, err := loadSagaSpec()
	if err != nil {
		return fmt.Errorf("load saga spec: %w", err)
	}

	config, err := structpb.NewStruct(map[string]any{
		"max_delay": 1000,
		"fail":      0.25,
	})
	if err != nil {
		return fmt.Errorf("new config struct: %w", err)
	}

	if err := orchestrator.BeginSaga(ctx, id, spec, config); err != nil {
		return fmt.Errorf("begin saga: %w", err)
	}

	<-ctx.Done()
	return nil
}

func executeTask(
	ctx context.Context, taskType string, input *structpb.Struct,
) (*structpb.Value, error) {
	delay := time.Millisecond * time.Duration(
		rand.Float64()*input.Fields["max_delay"].GetNumberValue())
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	fail := input.Fields["fail"].GetNumberValue()
	if fail > rand.Float64() {
		return nil, errors.New("oops")
	}

	switch taskType {
	case "noop":
		return structpb.NewNullValue(), nil
	case "add":
		a := input.Fields["a"].GetNumberValue()
		b := input.Fields["b"].GetNumberValue()
		return structpb.NewNumberValue(a + b), nil
	case "multiply":
		a := input.Fields["a"].GetNumberValue()
		b := input.Fields["b"].GetNumberValue()
		return structpb.NewNumberValue(a * b), nil
	default:
		return nil, fmt.Errorf("not implemented: %s", taskType)
	}
}

func loadSagaSpec() (*sagaspecpb.Saga, error) {
	var payload map[string]any
	if err := yaml.Unmarshal(sagaBytes, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal YAML: %w", err)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal JSON: %w", err)
	}

	var spec sagaspecpb.Saga
	if err := protojson.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("unmarshal spec: %w", err)
	}

	return &spec, nil
}

//go:embed saga.yml
var sagaBytes []byte
