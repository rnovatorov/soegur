package model

import (
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
)

type BeginSaga struct {
	Spec   *sagaspecpb.Saga
	Config *structpb.Struct
}

type TriggerNextSteps struct{}

type EndStep struct {
	ID     string
	Output *structpb.Value
}

type AbortStep struct {
	ID     string
	Reason string
}

type EndStepCompensation struct {
	ID string
}

type AbortStepCompensation struct {
	ID     string
	Reason string
}
