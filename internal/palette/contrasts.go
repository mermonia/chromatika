package palette

import (
	"math"

	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/utils"
)


func wcag(a, b *colors.LCHab) (float32, error) {
	labA := colors.LCHtoLab(a)
	labB := colors.LCHtoLab(b)

	xyzA := colors.LabToXyz(labA)
	xyzB := colors.LabToXyz(labB)

	l1 := max(xyzA.Y, xyzB.Y)
	l2 := min(xyzA.Y, xyzB.Y)

	return (l1 + 0.05) / (l2 + 0.05), nil
}

// xY: Y value for color x
func deltaE00(a, b *colors.LCHab) float64 {
	labA := colors.LCHtoLab(a)
	labB := colors.LCHtoLab(b)

	// We need a components to create new intermediate colors
	aA := float64(labA.A)
	bA := float64(labB.A)

	// Average chroma and luminosity
	avgL := float64((a.L + b.L) / 2)
	avgC := float64((a.C + b.C) / 2)

	// G parameter for future calculations
	avgCto7 := math.Pow(avgC, 7)
	G := 0.5 * (1 - math.Sqrt(avgCto7/(avgCto7+6103515625)))

	// Intermediate color generation
	aAprime := float32(aA * (1 + G))
	bAprime := float32(bA * (1 + G))

	labAprime := &colors.Lab{
		L: a.L,
		A: aAprime,
		B: labA.B,
	}

	labBprime := &colors.Lab{
		L: b.L,
		A: bAprime,
		B: labB.B,
	}

	aPrime := colors.LabToLCH(labAprime)
	bPrime := colors.LabToLCH(labBprime)

	// T parameter for future calculations
	avgHprime := float64(aPrime.H + bPrime.H)
	if math.Abs(float64(aPrime.H-bPrime.H)) > 180 {
		avgHprime += 360
	}
	avgHprime /= 2

	T := 1 - 0.17*utils.DegCos(avgHprime-30) + 0.24*utils.DegCos(2*avgHprime) + 0.32*utils.DegCos(3*avgHprime+6) - 0.20*utils.DegCos(4*avgHprime-63)

	// Delta calculations
	deltaL := float64(bPrime.L - aPrime.L)
	deltaC := float64(bPrime.C - aPrime.C)

	hPrimeDiff := float64(bPrime.H - aPrime.H)
	var imDeltaH float64
	if math.Abs(hPrimeDiff) <= 180 {
		imDeltaH = hPrimeDiff
	} else if math.Abs(hPrimeDiff) > 180 && bPrime.H <= aPrime.H {
		imDeltaH = hPrimeDiff + 360
	} else {
		imDeltaH = hPrimeDiff - 360
	}

	deltaH := 2 * math.Sqrt(float64(aPrime.C*bPrime.C)) * utils.DegSin(imDeltaH/2)

	// Calculation of S values
	avgCprime := float64(aPrime.C+bPrime.C) / 2

	sL := 1 + (0.015*math.Pow(avgL-50, 2))/math.Sqrt(20+math.Pow(avgL-50, 2))
	sC := 1 + 0.045*avgCprime
	sH := 1 + 0.015*avgCprime*T

	// Extra parameters
	deltaTheta := 30 * math.Exp(-math.Pow((avgHprime-275)/25, 2))

	avgCprimeTo7 := math.Pow(avgCprime, 7)
	Rc := 2 * math.Sqrt(avgCprimeTo7/(avgCprimeTo7+6103515625))
	Rt := -Rc * utils.DegSin(2*deltaTheta)

	// Final component calculation
	compL := deltaL / sL
	compC := deltaC / sC
	compH := deltaH / sH

	// Result
	res := math.Sqrt(
		compL*compL +
			compC*compC +
			compH*compH +
			Rt*compC*compH,
	)

	return res
}
