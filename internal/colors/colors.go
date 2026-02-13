package colors

import (
	"fmt"
	"image/color"
	"math"

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

func PixelToLab(in color.Color) (*Lab, error) {
	rgb := PixelToRGB(in)
	lab, err := RGBtoLab(rgb)

	if err != nil {
		return nil, fmt.Errorf("could not convert rgb to lab: %w", err)
	}

	return lab, nil
}

func PixelToRGB(in color.Color) *Rgb {
	r, g, b, _ := in.RGBA()

	out := &Rgb{
		R: uint8(r / 257),
		G: uint8(g / 257),
		B: uint8(b / 257),
	}

	return out
}

func (col *NRgb) ToMatrix() *utils.Matrix {
	data := [][]float32{
		{col.R},
		{col.G},
		{col.B},
	}

	return utils.NewMatrix(data)
}

func RGBtoNRGB(in *Rgb) *NRgb {
	out := &NRgb{
		R: normalizeRGB(in.R),
		G: normalizeRGB(in.G),
		B: normalizeRGB(in.B),
	}

	return out
}

func RGBtoXYZ(in *Rgb) (*Xyz, error) {
	workingMatData := [][]float32{
		{0.4124564, 0.3575761, 0.1804375},
		{0.2126729, 0.7151522, 0.0721750},
		{0.0193339, 0.1191920, 0.9503041},
	}
	workingMat := utils.NewMatrix(workingMatData)

	nrgb := RGBtoNRGB(in)
	corrected := gammaCorrectNRGB(nrgb)

	colorMat := corrected.ToMatrix()

	var convertedMat utils.Matrix
	err := convertedMat.Mul(workingMat, colorMat)
	if err != nil {
		return nil, fmt.Errorf("could not multiply working matrix and nrgb color: %w", err)
	}

	xyz := &Xyz{
		X: convertedMat.At(0, 0),
		Y: convertedMat.At(1, 0),
		Z: convertedMat.At(2, 0),
	}

	return xyz, nil
}

func RGBtoLab(in *Rgb) (*Lab, error) {
	xyz, err := RGBtoXYZ(in)
	if err != nil {
		return nil, fmt.Errorf("could not convert RGB to Xyz")
	}

	// reference white D65
	var wX, wY, wZ float32
	wX, wY, wZ = 0.95047003, 1.0000001, 1.08883

	xTrans := labTransform(xyz.X / wX)
	yTrans := labTransform(xyz.Y / wY)
	zTrans := labTransform(xyz.Z / wZ)

	lab := &Lab{
		L: 116*yTrans - 16,
		A: 500 * (xTrans - yTrans),
		B: 200 * (yTrans - zTrans),
	}

	return lab, nil
}

func labTransform(t float32) float32 {
	if t > 0.008856 {
		return float32(math.Pow(float64(t), 1.0/3.0))
	}

	return t*7.787 + 16.0/116.0
}

func gammaCorrectNRGB(in *NRgb) *NRgb {
	out := &NRgb{
		R: gammaCorrect(in.R),
		G: gammaCorrect(in.G),
		B: gammaCorrect(in.B),
	}

	return out
}

func normalizeRGB(value uint8) float32 {
	return float32(value) / 255.0
}

func gammaCorrect(value float32) float32 {
	if value <= 0.04045 {
		return value / 12.92
	}

	return float32(math.Pow(float64((value+0.055)/1.055), 2.4))
}
