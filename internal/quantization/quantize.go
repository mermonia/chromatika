package quantization

import (
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/utils"
)

func Quantize(q int, cols []*colors.Lab) ([]*colors.Lab, []int) {
	colorArray := generateQuantizedColorArray(q)
	naiveNP := make([]int, len(colorArray))

	for _, color := range cols {
		naiveNP[getColorIndex(q, color)]++
	}

	filteredColors := make([]*colors.Lab, 0, len(colorArray))
	filteredNP := make([]int, 0, len(naiveNP))

	for i, color := range colorArray {
		if naiveNP[i] > 0 {
			filteredColors = append(filteredColors, color)
			filteredNP = append(filteredNP, naiveNP[i])
		}
	}

	return filteredColors, filteredNP
}

func generateQuantizedColorArray(q int) []*colors.Lab {
	levelsL := 100 / q
	levelsA := 240 / q
	levelsB := 240 / q

	cols := make([]*colors.Lab, levelsL*levelsA*levelsB)

	for l := range levelsL {
		for a := range levelsA {
			for b := range levelsB {
				// The center of the quantization bin is the representative color
				cols[l*levelsA*levelsB+a*levelsB+b] = &colors.Lab{
					L: float32(l*q) + float32(q)/2.0,
					A: float32(a*q-120.0) + float32(q)/2.0,
					B: float32(b*q-120.0) + float32(q)/2.0,
				}
			}
		}
	}

	return cols
}

func getColorIndex(q int, color *colors.Lab) int {
	levelsA := 240 / q
	levelsB := 240 / q

	l := int(color.L / float32(q))
	a := int((color.A + 120.0) / float32(q))
	b := int((color.B + 120.0) / float32(q))

	l = utils.Clamp(l, 0, 100/q-1)
	a = utils.Clamp(a, 0, 240/q-1)
	b = utils.Clamp(b, 0, 240/q-1)

	return l*levelsA*levelsB + a*levelsB + b
}
