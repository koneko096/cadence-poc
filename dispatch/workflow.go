package dispatch

import (
	"cadence-poc/geo"
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
		tripStart geo.GeoPoint
	}

	DispatchState struct {
		driverID  int32
		tripStart geo.GeoPoint
		startTime time.Time
	}
)

func DispatchWorkflow(ctx workflow.Context, start geo.GeoPoint) error {
	var activities *Activities

	future := workflow.ExecuteActivity(ctx, activities.FindNearestDriver, start)
	driver := Driver{}
	if err := future.Get(ctx, &driver); err != nil {
		log.Printf("Cannot find nearest driver\n%v", err)
		return err
	}

	future = workflow.ExecuteActivity(ctx, activities.DispatchDriver, DispatchRequest{driverID: driver.id, tripStart: start})
	state := DispatchState{}
	if err := future.Get(ctx, &state); err != nil {
		log.Printf("Driver dispatch failed\n%v", err)
		return err
	}

	if err := workflow.Sleep(ctx, 10*time.Minute); err != nil {
		log.Printf("Error on getting trip\n%v", err)
		return err
	}

	return nil
}

func (*Activities) FindNearestDriver(ctx context.Context, start geo.GeoPoint) (Driver, error) {
	return Driver{id: 13}, nil
}

func (*Activities) DispatchDriver(ctx context.Context, req DispatchRequest) (DispatchState, error) {
	return DispatchState{driverID: req.driverID, tripStart: req.tripStart, startTime: time.Now()}, nil
}
