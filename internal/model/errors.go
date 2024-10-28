package model

import (
	"errors"
)

var (
	ErrSagaAlreadyBegun = errors.New("saga already begun")
	ErrSagaEnded        = errors.New("saga ended")
	ErrSagaStepsMissing = errors.New("saga steps missing")

	ErrDependencyGraphCyclic             = errors.New("dependency graph cyclic")
	ErrDependencyGraphVertexDefinedTwice = errors.New("dependency graph vertex defined twice")
	ErrDependencyGraphVertexNotFound     = errors.New("dependency graph vertex not found")

	ErrTaskInputReturnedTooManyValues = errors.New("task input returned too many values")

	ErrStepAbortReasonMissing       = errors.New("step abort reason missing")
	ErrStepAborted                  = errors.New("step aborted")
	ErrStepAlreadyAborted           = errors.New("step already aborted")
	ErrStepEnded                    = errors.New("step ended")
	ErrStepAlreadyEnded             = errors.New("step already ended")
	ErrStepNotDefined               = errors.New("step not defined")
	ErrStepNotBegun                 = errors.New("step not begun")
	ErrStepOutputMissing            = errors.New("step output missing")
	ErrStepCompensationAlreadyEnded = errors.New("step compensation already ended")
	ErrStepCompensationNotBegun     = errors.New("step compensation not begun")
)
