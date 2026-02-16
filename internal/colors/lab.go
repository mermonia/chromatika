package colors

import (
	"fmt"
	"math"

	"github.com/mermonia/chromatika/internal/utils"
)

func LabToRGB(in *Lab) (*Rgb, error) {
	workingMatData := [][]float32{
		{3.2404542, -1.5371385, -0.4985314},
		{-0.9692660, 1.8760108, 0.0415560},
		{0.0556434, -0.2040259, 1.0572252},
	}
	workingMat := utils.NewMatrix(workingMatData)

	xyz := LabToXyz(in)
	colorMat := xyz.ToMatrix()

	var convertedMat utils.Matrix
	err := convertedMat.Mul(workingMat, colorMat)
	if err != nil {
		return nil, fmt.Errorf("could not multiply working matrix and xyz color: %w", err)
	}

	linearNRgb := &NRgb{
		R: convertedMat.At(0,0),
		G: convertedMat.At(1,0),
		B: convertedMat.At(2,0),
	}

	compandedNRgb := compandedNRGB(linearNRgb)
	rgb := NRGBtoRGB(compandedNRgb)

	return rgb, nil
}

func LabToXyz(in *Lab) *Xyz {
	// Reference white D65
	var wX, wY, wZ float32
	wX, wY, wZ = 0.95047003, 1.0000001, 1.08883

	// Transformations
	fY := (in.L + 16) / 116
	fX := (in.A/500) + fY
	fZ := fY - (in.B/200)

	// Threshold values
	var e float32 = 0.008856
	var k float32 = 903.3

	// Normal values
	x := fX*fX*fX
	y := fY*fY*fY
	z := fZ*fZ*fZ

	// Over-theshold transformations
	if x <= e {
		x = (116*fX - 16)/k
	}

	if in.L <= k*e {
		y = in.L/k
	}

	if z <= e {
		z = (116*fZ - 16)/k
	}

	xyz := &Xyz{
		X: x * wX,
		Y: y * wY,
		Z: z * wZ,
	}

	return xyz
}

func compandedNRGB(in *NRgb) *NRgb {
	out := &NRgb{
		R: compand(in.R),
		G: compand(in.G),
		B: compand(in.B),
	}

	return out
}

func compand(value float32) float32 {
	if value <= 0.0031308 {
		return value * 12.92
	}
	return float32(1.055 * math.Pow(float64(value), 1/2.4) - 0.055)
}

