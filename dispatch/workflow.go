package dispatch

import (
	"context"
	"go.temporal.io/sdk/workflow"
)

type Activities struct{}

type (
	GeoPoint struct {
		lat  float32
		long float32
	}
)

func DispatchWorkflow(ctx workflow.Context, start GeoPoint) error {
}

func (*Activities) FindNearestDriver(ctx context.Context) error {}

func (*Activities) DispatchDriver(ctx context.Context) error {

}
