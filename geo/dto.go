package geo

import "cadence-poc/grpc"

type Message interface {
	ToMsg() *grpc.GeoPoint
}

type (
	GeoPoint struct {
		lat  float32
		long float32
	}

	Trip struct {
		start GeoPoint
		end   GeoPoint
	}
)

func (s *GeoPoint) ToMsg() *grpc.GeoPoint {
	return &grpc.GeoPoint{
		Latitude:  s.lat,
		Longitude: s.long,
	}
}
