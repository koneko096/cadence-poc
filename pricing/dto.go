package pricing

import "time"

type (
	TripFare struct {
		estimateFare  int64
		estimatedTrip time.Time
	}
)
