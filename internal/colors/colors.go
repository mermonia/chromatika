package colors

import (
	"github.com/mermonia/chromatika/internal/utils"
)

type Lab struct {
	L, A, B float32
}

type Rgb struct {
	R, G, B uint8
}

type NRgb struct {
	R, G, B float32
}

type Xyz struct {
	X, Y, Z float32
}

type LCHab struct {
	L, C, H float32
}

func (col *NRgb) ToMatrix() *utils.Matrix {
	data := [][]float32{
		{col.R},
		{col.G},
		{col.B},
	}

	return utils.NewMatrix(data)
}

func (col *Xyz) ToMatrix() *utils.Matrix {
	data := [][]float32{
		{col.X},
		{col.Y},
		{col.Z},
	}

	return utils.NewMatrix(data)
}
