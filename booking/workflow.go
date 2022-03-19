package booking

import (
	"cadence-poc/grpc"
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	PricingClient grpc.PricingClient
	GeoClient     grpc.GeoClient
	PaymentClient grpc.PaymentClient
}

type (
	BookingRequest struct {
		UserID int32
		Trip   grpc.TripRequest
	}

	TripFare struct {
		Distance float32
		Fare     int64
	}

	BookingState struct {
		Request  BookingRequest
		TripFare TripFare
	}

	BookingStatus struct {
		bookingID string
	}
)

func BookingWorkflow(ctx workflow.Context, req *BookingRequest) error {
	var activities *Activities
	activityoptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityoptions)

	fare := TripFare{}
	future := workflow.ExecuteActivity(ctx, activities.CalculateFare, &req.Trip)
	if err := future.Get(ctx, &fare); err != nil {
		log.Printf("Fare calculation failed\n%v", err)
		return err
	}

	bookingState := BookingState{}
	future = workflow.ExecuteActivity(ctx, activities.CreateBooking, &bookingState)
	bookingStatus := BookingStatus{}
	if err := future.Get(ctx, &bookingStatus); err != nil {
		log.Printf("Booking creation failed\n%v", err)
		return err
	}

	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue: "DISPATCH_QUEUE",
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	chFuture := workflow.ExecuteChildWorkflow(ctx, "DispatchDriverWorkflow", req)
	if err := chFuture.Get(ctx, &bookingStatus); err != nil {
		log.Printf("Driver dispatch failed\n%v", err)
		return err
	}

	future = workflow.ExecuteActivity(ctx, activities.DeductFare, fare)
	if err := future.Get(ctx, &bookingStatus); err != nil {
		log.Printf("Fare deduction failed\n%v", err)
		return err
	}

	future = workflow.ExecuteActivity(ctx, activities.FinishBooking, bookingStatus)
	if err := future.Get(ctx, &fare); err != nil {
		log.Printf("Booking update failed\n%v", err)
		return err
	}

	return nil
}

func (s *Activities) CalculateFare(ctx context.Context, req *grpc.TripRequest) (TripFare, error) {
	tripDetail, err := s.GeoClient.ComputeRoute(ctx, req)
	if err != nil {
		log.Printf("Cannot find trip route\n%v", err)
		return TripFare{}, err
	}

	rate, err := s.PricingClient.BidRate(ctx, req.Start)
	if err != nil {
		log.Printf("Cannot compute fare rate in the region\n%v", err)
		return TripFare{}, err
	}

	distance := tripDetail.Length
	return TripFare{
		Distance: distance,
		Fare:     int64(rate.Value) * int64(distance+0.5),
	}, nil
}

func (*Activities) CreateBooking(ctx context.Context, booking *BookingState) (BookingStatus, error) {
	return BookingStatus{"GRABJEK"}, nil
}

func (s *Activities) DeductFare(ctx context.Context, bill *grpc.Billing) error {
	_, err := s.PaymentClient.DeductFare(ctx, bill)
	if err != nil {
		log.Printf("Deduct fare failed\n%v", err)
		return err
	}

	return nil
}

func (*Activities) FinishBooking(ctx context.Context, booking *BookingStatus) error {
	return nil
}
