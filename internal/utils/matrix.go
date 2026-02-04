package utils

import "fmt"

type Matrix struct {
	Data [][]float32
}

func (mat *Matrix) Dims() (int, int) {
	if len(mat.Data) == 0 || len(mat.Data[0]) == 0 {
		return 0, 0
	}

	return len(mat.Data), len(mat.Data[0])
}

func (mat *Matrix) At(x, y int) float32 {
	return mat.Data[x][y]
}

func NewMatrix(data [][]float32) *Matrix {
	mat := &Matrix{
		Data: data,
	}

	return mat
}

func (mat *Matrix) Mul(A, B *Matrix) error {
	rA, cA := A.Dims()
	rB, cB := B.Dims()

	if cA != rB {
		fmt.Println(cA, rB)
		return fmt.Errorf("the number of columns of the first matrix must match the number of rows of the second one")
	}

	result := make([][]float32, rA)
	for i := range result {
		result[i] = make([]float32, cB)
	}

	for i, rowA := range A.Data {
		Ri := result[i]
		for k, aik := range rowA {
			Bk := B.Data[k]
			for j, val := range Bk {
				Ri[j] += aik * val
			}
		}
	}

	mat.Data = result
	return nil
}
