package utils

import "math"

func HueDifference(a, b float32) float32 {
	rawDifference := float32(math.Abs(float64(a) - float64(b)))
	return min(rawDifference, 360 - rawDifference)
}

func SignedHueDifference(a, b float32) float32 {
    diff := a - b
    for diff < -180 { diff += 360 }
    for diff > 180 { diff -= 360 }
    return diff
}
