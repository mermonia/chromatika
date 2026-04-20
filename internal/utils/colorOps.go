package utils

import "math"

func HueDifference(a, b float64) float64 {
	rawDifference := math.Abs(a - b)
	return min(rawDifference, 360-rawDifference)
}

func SignedHueDifference(a, b float64) float64 {
	diff := a - b
	for diff < -180 {
		diff += 360
	}
	for diff > 180 {
		diff -= 360
	}
	return diff
}
