package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	rpc "cadence-poc/grpc"

	"google.golang.org/grpc"
)

type geoServer struct {
	rpc.UnimplementedGeoServer
}

func newServer() rpc.GeoServer {
	return &geoServer{}
}

func (*geoServer) ComputeRoute(ctx context.Context, req *rpc.TripRequest) (*rpc.TripDetail, error) {
	return &rpc.TripDetail{
		Request:  req,
		Length:   255.78,
		Duration: durationpb.New(25 * time.Minute),
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "localhost:10887")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterGeoServer(grpcServer, newServer())
	log.Printf("Listening on localhost:10887")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("error occured: %v", err)
	}
}
