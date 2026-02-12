package colors

import "math"

func DistanceLab(a, b *Lab) float64 {
	deltaL := float64(b.L - a.L)
	deltaA := float64(b.A - a.A)
	deltaB := float64(b.B - a.B)

	return math.Sqrt(deltaL*deltaL + deltaA*deltaA + deltaB*deltaB)
}

func DistanceMatrix(colors []*Lab) [][]float64 {
	n := len(colors)

	res := make([][]float64, n)

	// Fill distance matrix
	for i := range n {
		res[i] = make([]float64, n)
		res[i][i] = 0
		for j := i + 1; j < n; j++ {
			dist := DistanceLab(colors[i], colors[j])
			res[i][j] = dist
			res[j][i] = dist
		}
	}

	return res
}

func (a *Lab) Equals(b *Lab) bool {
	if a.L != b.L || a.A != b.A || a.B != b.B {
		return false
	}

	return true
}
