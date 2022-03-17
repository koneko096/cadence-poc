package main

import (
	"cadence-poc/booking"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Worker and Activity functions
	w := worker.New(c, "BOOKING_QUEUE", worker.Options{})

	a := &booking.Activities{}

	w.RegisterActivity(a.CreateBooking)
	w.RegisterActivity(a.FinishBooking)

	w.RegisterWorkflow(booking.BookingWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
