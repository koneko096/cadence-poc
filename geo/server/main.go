package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"github.com/jftuga/geodist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	rpc "cadence-poc/grpc"
)

type geoServer struct {
	rpc.UnimplementedGeoServer
}

func newServer() rpc.GeoServer {
	return &geoServer{}
}

func (*geoServer) ComputeRoute(ctx context.Context, req *rpc.TripRequest) (*rpc.TripDetail, error) {
	orig := geodist.Coord{Lat: float64(req.Start.Latitude), Lon: float64(req.Start.Longitude)}
	dest := geodist.Coord{Lat: float64(req.End.Latitude), Lon: float64(req.End.Longitude)}
	_, km, err := geodist.VincentyDistance(orig, dest)
	if err != nil {
		log.Println(err)
		km = 0
	}
	estimatedDuration := time.Duration(math.Ceil(km)*5) * time.Minute
	return &rpc.TripDetail{
		Request:  req,
		Length:   km,
		Duration: fmt.Sprintf("%s", estimatedDuration),
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
	reflection.Register(grpcServer)
	log.Printf("Listening on localhost:10887")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("error occured: %v", err)
	}
}
