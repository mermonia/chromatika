package colors

import "math"

func DistanceLab(a, b *Lab) float64 {
	deltaL := float64(b.L - a.L)
	deltaA := float64(b.A - a.A)
	deltaB := float64(b.B - a.B)

	return math.Sqrt(deltaL*deltaL + deltaA*deltaA + deltaB*deltaB)
}
