package geo

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
