package model

import (
	"context"
	"fmt"
	"time"

	"github.com/itchyny/gojq"

	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
)

const varConfig = "$_config"

type Task struct {
	spec                 *sagaspecpb.Task
	dependencies         []string
	compiledInputQueries map[string]*gojq.Code
}

func newTask(spec *sagaspecpb.Task, deps []string) (*Task, error) {
	variables := []string{varConfig}
	for _, dep := range deps {
		variables = append(variables, "$"+dep)
	}

	input := make(map[string]*gojq.Code, len(spec.GetInput()))
	for name, src := range spec.GetInput() {
		query, err := gojq.Parse(src)
		if err != nil {
			return nil, fmt.Errorf("parse input: %s: %w", name, err)
		}
		code, err := gojq.Compile(query, gojq.WithVariables(variables))
		if err != nil {
			return nil, fmt.Errorf("compile input: %s: %w", name, err)
		}
		input[name] = code
	}

	return &Task{
		spec:                 spec,
		dependencies:         deps,
		compiledInputQueries: input,
	}, nil
}

type EvalContext struct {
	config  map[string]any
	outputs map[string]any
	data    any
}

func (t *Task) EvalInput(evalCtx EvalContext) (map[string]any, error) {
	input := make(map[string]any)

	for name, code := range t.compiledInputQueries {
		value, err := t.eval(evalCtx, code)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		input[name] = value
	}

	return input, nil
}

func (t *Task) eval(evalCtx EvalContext, code *gojq.Code) (any, error) {
	var value any

	vars := []any{evalCtx.config}
	for _, dep := range t.dependencies {
		vars = append(vars, evalCtx.outputs[dep])
	}

	const timeout = time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	iter := code.RunWithContext(ctx, evalCtx.data, vars...)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			return nil, err
		}
		if value != nil {
			return nil, ErrTaskInputReturnedTooManyValues
		}
		value = v
	}

	return value, nil
}
