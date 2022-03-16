package dispatch

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
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

	a := &Activities{}

	w.RegisterActivity(a.CreateBooking)
	w.RegisterActivity(a.FinishBooking)

	w.RegisterWorkflow(BookingWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
