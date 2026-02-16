package colors

import (
	"fmt"
	"math"

	"github.com/mermonia/chromatika/internal/utils"
)

func RGBtoNRGB(in *Rgb) *NRgb {
	out := &NRgb{
		R: normalizeRGB(in.R),
		G: normalizeRGB(in.G),
		B: normalizeRGB(in.B),
	}

	return out
}

func NRGBtoRGB(in *NRgb) *Rgb {
	out := &Rgb{
		R: scaleNRGB(in.R),
		G: scaleNRGB(in.G),
		B: scaleNRGB(in.B),
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

func scaleNRGB(value float32) uint8 {
	return uint8(value * 255.0)
}

func gammaCorrect(value float32) float32 {
	if value <= 0.04045 {
		return value / 12.92
	}

	return float32(math.Pow(float64((value+0.055)/1.055), 2.4))
}

