package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

func Clamp[T Number](n, min, max T) T {
	if n > max {
		return max
	}
	if n < min {
		return min
	}
	return n
}

func DegSin(x float64) float64 {
	return math.Sin(x * math.Pi / 180)
}

func DegCos(x float64) float64 {
	return math.Cos(x * math.Pi / 180)
}
