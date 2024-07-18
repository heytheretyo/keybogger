package main

import "math"

func calculate_distance(x1, y1, x2, y2 int16) float64 {
	distanceInPixels := math.Sqrt(math.Abs(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))))
	distanceInInches := distanceInPixels / dpi
	distanceInMeters := distanceInInches * inchesToMeters
	return distanceInMeters
}
