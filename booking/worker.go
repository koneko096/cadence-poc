package booking

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	Client     *client.Client
	Activities *Activities
}

func (s *Worker) Run() {
	w := worker.New(*s.Client, "BOOKING_QUEUE", worker.Options{})

	w.RegisterActivity(s.Activities.CalculateFare)
	w.RegisterActivity(s.Activities.CreateBooking)
	w.RegisterActivity(s.Activities.DeductFare)
	w.RegisterActivity(s.Activities.FinishBooking)

	w.RegisterWorkflow(BookingWorkflow)

	// Start listening to the Task Queue
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
