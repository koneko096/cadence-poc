package main

import (
	"context"
	"flag"
	"log"
	"net"

	rpc "cadence-poc/grpc"

	"google.golang.org/grpc"
)

type pricingServer struct {
	rpc.UnimplementedPricingServer
}

func newServer() rpc.PricingServer {
	return &pricingServer{}
}

func (*pricingServer) BidRate(ctx context.Context, req *rpc.GeoPoint) (*rpc.Fare, error) {
	return &rpc.Fare{Value: 666}, nil
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
	log.Printf("Listening on localhost:10886")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("error occured: %v", err)
	}
}
