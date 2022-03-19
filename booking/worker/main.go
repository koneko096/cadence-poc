package main

import (
	"cadence-poc/booking"
	rpc "cadence-poc/grpc"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client\n", err)
	}
	defer c.Close()

	w := worker.New(c, "BOOKING_QUEUE", worker.Options{})
	pc := generateConn("localhost:10886")
	defer pc.Close()
	gc := generateConn("localhost:10887")
	defer gc.Close()
	yc := generateConn("localhost:10885")
	defer yc.Close()
	log.Printf("Connection(s) made")

	a := &booking.Activities{
		PricingClient: rpc.NewPricingClient(pc),
		GeoClient:     rpc.NewGeoClient(gc),
		PaymentClient: rpc.NewPaymentClient(yc),
	}

	w.RegisterActivity(a.CalculateFare)
	w.RegisterActivity(a.CreateBooking)
	w.RegisterActivity(a.DeductFare)
	w.RegisterActivity(a.FinishBooking)

	w.RegisterWorkflow(booking.BookingWorkflow)

	// Start listening to the Task Queue
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func generateConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Connection failure", err)
	}
	return conn
}
