package booking

import (
	"context"
	_ "github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
	_ "time"
)

type Activities struct{}

type (
	GeoPoint struct {
		lat  float32
		long float32
	}

	BookingState struct {
		userID string
		start  GeoPoint
		end    GeoPoint
		fare   int64
	}
)

func BookingWorkflow(ctx workflow.Context, booking BookingState) error {
}

func (*Activities) CreateBooking(ctx context.Context, booking BookingState) error {

}

func (*Activities) FinishBooking(ctx context.Context, booking BookingState) error {

}
