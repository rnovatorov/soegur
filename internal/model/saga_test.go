package model_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rnovatorov/go-eventsource/pkg/eventsource"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
	"github.com/rnovatorov/soegur/internal/model"
)

type SagaSuite struct {
	suite.Suite
	saga *model.SagaAggregate
}

func (s *SagaSuite) SetupTest() {
	s.saga = eventsource.NewAggregate[model.Saga](uuid.NewString())
}

func (s *SagaSuite) TestSequence() {
	s.shouldProcessCommand(model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("a"),
			s.step("b", "a"),
			s.step("c", "b"),
		}},
	})
	s.Require().Equal([]string(nil), s.saga.Root().StepTransitiveDependencies("a"))
	s.Require().Equal([]string{"a"}, s.saga.Root().StepTransitiveDependencies("b"))
	s.Require().Equal([]string{"a", "b"}, s.saga.Root().StepTransitiveDependencies("c"))

	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepInProgress("a"))
	s.Require().False(s.saga.Root().StepInProgress("b"))
	s.Require().False(s.saga.Root().StepInProgress("c"))

	s.shouldProcessCommand(model.EndStep{ID: "a", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepInProgress("b"))
	s.Require().False(s.saga.Root().StepInProgress("c"))

	s.shouldProcessCommand(model.EndStep{ID: "b", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("b"))
	s.Require().True(s.saga.Root().StepInProgress("c"))

	s.shouldProcessCommand(model.EndStep{ID: "c", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("b"))
	s.Require().True(s.saga.Root().StepEnded("c"))
}

func (s *SagaSuite) TestDiamond() {
	s.shouldProcessCommand(model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("d", "b", "c"),
			s.step("b", "a"),
			s.step("c", "a"),
			s.step("a"),
		}},
	})
	s.Require().Equal([]string(nil), s.saga.Root().StepTransitiveDependencies("a"))
	s.Require().Equal([]string{"a"}, s.saga.Root().StepTransitiveDependencies("b"))
	s.Require().Equal([]string{"a"}, s.saga.Root().StepTransitiveDependencies("c"))
	s.Require().Equal([]string{"a", "b", "c"}, s.saga.Root().StepTransitiveDependencies("d"))

	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepInProgress("a"))
	s.Require().False(s.saga.Root().StepInProgress("b"))
	s.Require().False(s.saga.Root().StepInProgress("c"))
	s.Require().False(s.saga.Root().StepInProgress("d"))

	s.shouldProcessCommand(model.EndStep{ID: "a", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepInProgress("b"))
	s.Require().True(s.saga.Root().StepInProgress("c"))
	s.Require().False(s.saga.Root().StepInProgress("d"))

	s.shouldProcessCommand(model.EndStep{ID: "b", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("b"))
	s.Require().True(s.saga.Root().StepInProgress("c"))
	s.Require().False(s.saga.Root().StepInProgress("d"))

	s.shouldProcessCommand(model.EndStep{ID: "c", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("b"))
	s.Require().True(s.saga.Root().StepEnded("c"))
	s.Require().True(s.saga.Root().StepInProgress("d"))

	s.shouldProcessCommand(model.EndStep{ID: "d", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("b"))
	s.Require().True(s.saga.Root().StepEnded("c"))
	s.Require().True(s.saga.Root().StepEnded("d"))
}

func (s *SagaSuite) TestTwoTasksWithSameID() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("a"),
			s.step("b"),
			s.step("c"),
			s.step("a"),
		}},
	})
	s.Require().ErrorIs(err, model.ErrDependencyGraphVertexDefinedTwice)
}

func (s *SagaSuite) TestDependencyNotFound() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("a"),
			s.step("b", "a"),
			s.step("c", "d"),
		}},
	})
	s.Require().ErrorIs(err, model.ErrDependencyGraphVertexNotFound)
}

func (s *SagaSuite) TestShortCycle() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("a", "c"),
			s.step("b", "a"),
			s.step("c", "b"),
		}},
	})
	s.Require().ErrorIs(err, model.ErrDependencyGraphCyclic)
}

func (s *SagaSuite) TestLongCycle() {
	err := s.saga.ProcessCommand(context.Background(), model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("a", "b"),
			s.step("b", "c", "d"),
			s.step("c", "e"),
			s.step("d", "e"),
			s.step("e", "b"),
		}},
	})
	s.Require().ErrorIs(err, model.ErrDependencyGraphCyclic)
}

func (s *SagaSuite) TestDisconnectedSequences() {
	s.shouldProcessCommand(model.BeginSaga{
		Spec: &sagaspecpb.Saga{Steps: []*sagaspecpb.Step{
			s.step("a"),
			s.step("b", "a"),
			s.step("x"),
			s.step("y", "x"),
		}},
	})
	s.Require().Equal([]string(nil), s.saga.Root().StepTransitiveDependencies("a"))
	s.Require().Equal([]string{"a"}, s.saga.Root().StepTransitiveDependencies("b"))
	s.Require().Equal([]string(nil), s.saga.Root().StepTransitiveDependencies("x"))
	s.Require().Equal([]string{"x"}, s.saga.Root().StepTransitiveDependencies("y"))

	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepInProgress("a"))
	s.Require().True(s.saga.Root().StepInProgress("x"))
	s.Require().False(s.saga.Root().StepInProgress("b"))
	s.Require().False(s.saga.Root().StepInProgress("y"))

	s.shouldProcessCommand(model.EndStep{ID: "a", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepInProgress("x"))
	s.Require().True(s.saga.Root().StepInProgress("b"))
	s.Require().False(s.saga.Root().StepInProgress("y"))

	s.shouldProcessCommand(model.EndStep{ID: "x", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("x"))
	s.Require().True(s.saga.Root().StepInProgress("b"))
	s.Require().True(s.saga.Root().StepInProgress("y"))

	s.shouldProcessCommand(model.EndStep{ID: "y", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().InProgress())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("x"))
	s.Require().True(s.saga.Root().StepInProgress("b"))
	s.Require().True(s.saga.Root().StepEnded("y"))

	s.shouldProcessCommand(model.EndStep{ID: "b", Output: &structpb.Value{}})
	s.shouldProcessCommand(model.TriggerNextSteps{})
	s.Require().True(s.saga.Root().Ended())
	s.Require().True(s.saga.Root().StepEnded("a"))
	s.Require().True(s.saga.Root().StepEnded("x"))
	s.Require().True(s.saga.Root().StepEnded("b"))
	s.Require().True(s.saga.Root().StepEnded("y"))
}

func (s *SagaSuite) shouldProcessCommand(cmd eventsource.Command) {
	err := s.saga.ProcessCommand(context.Background(), cmd)
	s.Require().NoError(err)
}

func (s *SagaSuite) step(id string, deps ...string) *sagaspecpb.Step {
	return &sagaspecpb.Step{
		Id:           id,
		Dependencies: deps,
	}
}

func TestSaga(t *testing.T) {
	suite.Run(t, new(SagaSuite))
}
