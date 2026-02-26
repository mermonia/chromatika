package utils

import "golang.org/x/exp/constraints"

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
