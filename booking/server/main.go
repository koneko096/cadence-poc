package main

import (
	"cadence-poc/booking"
	rpc "cadence-poc/grpc"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client\n", err)
	}
	defer c.Close()

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

	w := booking.Worker{
		Activities: a,
		Client:     &c,
	}
	s := booking.Server{
		Port:   "8090",
		Client: &c,
	}

	w.Run()
	s.Run()
}

func generateConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Connection failure", err)
	}
	return conn
}
