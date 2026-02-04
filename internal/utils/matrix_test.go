package utils

import (
	"testing"
)

func TestMul(t *testing.T) {
	a := [][]float32{
		{1, 2},
		{3, 4},
	}

	b := [][]float32{
		{5, 6},
		{7, 8},
	}

	c := [][]float32{
		{19, 22},
		{43, 50},
	}

	FailIfNotProduct(a, b, c, t)
}

func BenchmarkMatrixMul(b *testing.B) {
	size := 4

	A := randomMatrix(size, size)
	B := randomMatrix(size, size)

	b.ResetTimer()

	for b.Loop() {
		var C Matrix
		_ = C.Mul(A, B)
	}
}

func randomMatrix(r, c int) *Matrix {
	data := make([][]float32, r)
	for i := range r {
		row := make([]float32, c)
		for j := range c {
			row[j] = float32(i*j + j)
		}
		data[i] = row
	}
	return &Matrix{Data: data}
}

func FailIfNotProduct(a, b, res [][]float32, t *testing.T) {
	t.Helper()

	A := &Matrix{Data: a}
	B := &Matrix{Data: b}

	var C Matrix
	C.Mul(A, B)

	rC, cC := C.Dims()

	if len(res) != rC {
		t.Fatalf("The number of rows of the product matrix is different than expected: got %d, expected %d",
			rC, len(res))
	}

	if len(res[0]) != cC {
		t.Fatalf("The number of columns of the product matrix is different than expected: got %d, expected %d",
			rC, len(res))
	}

	for i, row := range C.Data {
		for j, value := range row {
			if value != res[i][j] {
				t.Fatalf("The value at position %d, %d is not expected: %f, expected %f",
					i, j, value, res[i][j])
			}
		}
	}
}
