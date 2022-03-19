package main

import (
	"cadence-poc/dispatch"
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
	w := worker.New(c, "DISPATCH_QUEUE", worker.Options{})

	a := &dispatch.Activities{}

	w.RegisterActivity(a.DispatchDriver)
	w.RegisterActivity(a.FindNearestDriver)

	w.RegisterWorkflow(dispatch.DispatchDriverWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
