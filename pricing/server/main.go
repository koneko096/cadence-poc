package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	rpc "cadence-poc/grpc"
)

type PricingBox struct {
	StartHour int32
	EndHour   int32
	BaseRate  int32
}

var pricingBoxes = []PricingBox{
	{
		StartHour: 0,
		EndHour:   6,
		BaseRate:  1,
	},
	{
		StartHour: 6,
		EndHour:   10,
		BaseRate:  4,
	},
	{
		StartHour: 10,
		EndHour:   14,
		BaseRate:  3,
	},
	{
		StartHour: 14,
		EndHour:   16,
		BaseRate:  2,
	},
	{
		StartHour: 16,
		EndHour:   20,
	},
	{
		StartHour: 21,
		EndHour:   24,
		BaseRate:  1,
	},
}

type pricingServer struct {
	rpc.UnimplementedPricingServer
}

func newServer() rpc.PricingServer {
	return &pricingServer{}
}

func (*pricingServer) BidRate(ctx context.Context, req *rpc.TripRequest) (*rpc.Fare, error) {
	now := time.Now()
	hourNow := int32(now.Hour())
	var baseRate int32
	for _, timeBox := range pricingBoxes {
		if hourNow >= timeBox.EndHour {
			continue
		}
		baseRate = timeBox.BaseRate
		break
	}
	return &rpc.Fare{Value: baseRate, Request: req}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "localhost:10886")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterPricingServer(grpcServer, newServer())
	reflection.Register(grpcServer)
	log.Printf("Listening on localhost:10886")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("error occured: %v", err)
	}
}
