package dispatch

import (
	rpc "cadence-poc/grpc"
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/workflow"
)

type Activities struct{}

type (
	Driver struct {
		id int32
	}

	DispatchRequest struct {
		driverID  int32
		tripStart *rpc.GeoPoint
	}

	DispatchState struct {
		tripID    int32
		driverID  int32
		tripStart *rpc.GeoPoint
		startTime time.Time
	}
)

func DispatchDriverWorkflow(ctx workflow.Context, req *rpc.TripRequest) error {
	var activities *Activities
	activityoptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityoptions)

	future := workflow.ExecuteActivity(ctx, activities.FindNearestDriver, &req.Start)
	driver := Driver{}
	if err := future.Get(ctx, &driver); err != nil {
		log.Printf("Cannot find nearest driver\n%v", err)
		return err
	}

	future = workflow.ExecuteActivity(ctx, activities.DispatchDriver, &DispatchRequest{
		driverID:  driver.id,
		tripStart: req.Start,
	})
	state := DispatchState{}
	if err := future.Get(ctx, &state); err != nil {
		log.Printf("Driver dispatch failed\n%v", err)
		return err
	}

	if err := workflow.Sleep(ctx, 5*time.Second); err != nil {
		log.Printf("Error on getting trip\n%v", err)
		return err
	}

	future = workflow.ExecuteActivity(ctx, activities.FinishTrip, state.tripID)
	var s string
	if err := future.Get(ctx, &s); err != nil {
		log.Printf("Finish trip failed\n%v", err)
		return err
	}

	return nil
}

func (*Activities) FindNearestDriver(ctx context.Context, start *rpc.GeoPoint) (Driver, error) {
	return Driver{id: 13}, nil
}

func (*Activities) DispatchDriver(ctx context.Context, req *DispatchRequest) (*DispatchState, error) {
	return &DispatchState{tripID: 44, driverID: req.driverID, tripStart: req.tripStart, startTime: time.Now()}, nil
}

func (*Activities) FinishTrip(ctx context.Context, tripID string) error {
	return nil
}
