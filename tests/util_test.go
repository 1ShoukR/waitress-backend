package tests

import (
	"fmt"
	"math"
	"testing"
	"waitress-backend/internal/utilities"

	"github.com/stretchr/testify/assert"
)

// Note: Calculating to int or float64 is still up for debate. expectedDistance is in meters
func RoundToNearestMeter(d float64) float64 {
	return math.Round(d)
}

func assertAlmostEqual(t *testing.T, actual, expected, tolerance float64, msgAndArgs ...interface{}) bool {
	if math.Abs(actual-expected) > tolerance {
		return assert.Fail(t, fmt.Sprintf("Values are not within %v tolerance: %v != %v", tolerance, actual, expected), msgAndArgs...)
	}
	return true
}

func TestHaversine__distance_yields_0m(t *testing.T) {
	// Given
	lat1 := 40.73060805009797
	long1 := -73.93520326445689
	expectedDistance := 0.0

	// When
	sut := utilities.Haversine(lat1, long1, lat1, long1)
	res := RoundToNearestMeter(sut)

	// Then
	assert.Equal(t, res, expectedDistance, "dev config")
}

func TestHaversine__distance_yields_1000m(t *testing.T) {
	// Given
	lat1 := 40.73060805009797   // latitude of location 1
	long1 := -73.93520326445689 // longitude of location 1

	lat2 := 40.73060805009797   // latitude of location 2
	long2 := -73.92310326445689 // longitude of location 2
	expectedDistance := 1020.0

	// When
	sut := utilities.Haversine(lat1, long1, lat2, long2)
	res := RoundToNearestMeter(sut)

	// Then
	assert.Equal(t, res, expectedDistance)
}

func TestHaversine__distance_yields_5000m(t *testing.T) {
	// Given
	tolerance := 10.0
	lat1 := 40.73060805009797
	long1 := -73.93520326445689

	lat2 := 40.77560805009797
	long2 := -73.93520326445689
	expectedDistance := 5000.0

	// When
	sut := utilities.Haversine(lat1, long1, lat2, long2)
	res := RoundToNearestMeter(sut)

	// Then
	assertAlmostEqual(t, res, expectedDistance, tolerance)
}
