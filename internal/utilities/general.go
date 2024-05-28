// This file contains utilities that are more general in nature
//
// The utilities here are as follows:
// - Haversine


package utilities

import (
	"math"
)

// This function calculates the distance between two points on the Earth's surface
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	latRad1 := lat1 * math.Pi / 180
	latRad2 := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latRad1)*math.Cos(latRad2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}
