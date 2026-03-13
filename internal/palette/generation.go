package palette

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/extraction"
	"github.com/mermonia/chromatika/internal/utils"
)

var defaultScaleWidth int = 512
var defaultQuantInterval int = 4
var defaultFuzziness float64 = 2.0
var defaultThreshold float64 = 0.001
var defaultMaxIterations int = 100

const MIN_DELTAE00 = 20

const SEARCH_ITERATIONS = 1
const N_CLUSTERS_STEP = 0

const N_BASE_COLORS = 8

func GeneratePalette(imagePath string, darkmode bool) (*Palette, error) {
	/*
		Extract dominant colors multiple times, increasing the number of colors extracted
		each iteration.
	*/
	var suitableColors []*colors.LCHab
	for i := range SEARCH_ITERATIONS {
		labCols, _, err := extraction.GetDominantColors(
			imagePath,
			defaultScaleWidth,
			defaultQuantInterval,
			clustering.FCMParameters{
				M: defaultFuzziness,
				E: defaultThreshold,
				B: defaultMaxIterations,
				K: 16 + i*N_CLUSTERS_STEP,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("could not generate a %d color palette: %w", i+16, err)
		}

		lchCols := utils.Map(labCols, colors.LabToLCH)

		suitableColors = filterSuitableColors(lchCols, false)

		if len(suitableColors) >= N_BASE_COLORS {
			break
		}
	}

	rawColors := generateRawPalette(suitableColors)

	if darkmode {
		return getDarkModePalette(rawColors)
	}
	return getLightModePalette(rawColors)
}

func getDarkModePalette(rawColors *RawColors) (*Palette, error) {
	derivedColors := [8]*colors.LCHab{}
	for i, baseCol := range rawColors.Colors {
		derivedColors[i] = colors.Darker(baseCol)
	}

	return &Palette{
		Background:    rawColors.DarkNeutral,
		Foreground:    rawColors.LightNeutral,
		BaseColors:    rawColors.Colors,
		DerivedColors: derivedColors,
	}, nil
}

func getLightModePalette(rawColors *RawColors) (*Palette, error) {
	derivedColors := [8]*colors.LCHab{}
	for i, baseCol := range rawColors.Colors {
		derivedColors[i] = colors.Lighter(baseCol)
	}

	return &Palette{
		Background:    rawColors.LightNeutral,
		Foreground:    rawColors.DarkNeutral,
		BaseColors:    rawColors.Colors,
		DerivedColors: derivedColors,
	}, nil
}

/*
Returns a set of colors that are distinct enough to form a palette.
*/
func filterSuitableColors(cols []*colors.LCHab, fast bool) []*colors.LCHab {
	g := getColorSimilarityGraph(cols, MIN_DELTAE00)
	suitableColors := make([]*colors.LCHab, 0, len(cols))

	var mis []int
	if fast {
		mis = MIS_Fast(g)
	} else {
		mis = MIS_Complete(g)
	}

	for _, idx := range mis {
		suitableColors = append(suitableColors, cols[idx])
	}

	return suitableColors
}

func getColorSimilarityGraph(cols []*colors.LCHab, minDelta float64) *Graph {
	g := NewGraph(len(cols))

	for i := range cols {
		for j := range i {
			delta := deltaE00(cols[i], cols[j])
			if delta < minDelta {
				g.AddEdge(i, j)
			}
		}
	}

	return g
}

func generateRawPalette(base []*colors.LCHab) *RawColors {
	// Base colors
	newBaseColors := make([]*colors.LCHab, len(base))
	copy(newBaseColors, base)

	// Precompute deltaE00 matrix
	deltaE00Mat := make([][]float64, len(base))
	for i := range deltaE00Mat {
		deltaE00Mat[i] = make([]float64, len(base))
		for j := range i {
			d := deltaE00(base[i], base[j])
			deltaE00Mat[i][j] = d
			deltaE00Mat[j][i] = d
		}
	}

	// Get the dominanant color out of the provided base colors
	dominantColor := base[dominantColor(deltaE00Mat)]

	// Ensure neutral light and neutral dark colors
	var darkNeutral, lightNeutral *colors.LCHab
	for i := len(newBaseColors) - 1; i >= 0; i-- {
		color := newBaseColors[i]

		if color.C < 8 && color.L > 92 {
			lightNeutral = color
			newBaseColors = append(newBaseColors[:i], newBaseColors[i+1:]...)
		}

		if color.C < 15 && color.L < 15 {
			darkNeutral = color
			newBaseColors = append(newBaseColors[:i], newBaseColors[i+1:]...)
		}
	}

	// Generate dark neutral color
	if darkNeutral == nil {
		darkNeutral = &colors.LCHab{
			L: 8,
			C: 10,
			H: dominantColor.H,
		}
	}

	// Generate light neutral color
	if lightNeutral == nil {
		lightNeutral = &colors.LCHab{
			L: 92,
			C: 6,
			H: dominantColor.H,
		}
	}

	rawColors := &RawColors{
		LightNeutral: lightNeutral,
		DarkNeutral:  darkNeutral,
	}

	// Generate missing colors
	nMissingColors := len(rawColors.Colors) - len(newBaseColors)
	// meanDeltaE00 := averageDeltaE00(deltaE00Mat)

	newBaseColors = append(newBaseColors, harmonicExpand(dominantColor, min(3, nMissingColors))...)
	nMissingColors -= 3
	newBaseColors = append(newBaseColors, contrastNeutral(dominantColor, nMissingColors)...)

	copy((*rawColors).Colors[:], newBaseColors)

	return rawColors
}

func contrastNeutral(dominant *colors.LCHab, n int) []*colors.LCHab {
	if n <= 0 {
		return []*colors.LCHab{}
	}

	newColors := make([]*colors.LCHab, n)

	accentL := dominant.L
	if accentL > 50 {
		accentL += 5
	}

	accentC := utils.Clamp(dominant.C*1.2, 0, 100)
	accentH := float32(math.Mod(float64(dominant.H)+180, 360))

	// add accent color
	newColors[0] = &colors.LCHab{
		L: accentL,
		C: accentC,
		H: accentH,
	}
	n--

	for i := range n {
		if i%2 == 0 {
			// add light neutral
			newColors[i+1] = &colors.LCHab{
				L: 85 - float32(i)*15,
				C: 20,
				H: dominant.H,
			}
		} else {
			// add dark neutral
			newColors[i+1] = &colors.LCHab{
				L: 30 + float32(i)*15,
				C: 30,
				H: dominant.H,
			}
		}
	}

	return newColors
}

func harmonicExpand(dominant *colors.LCHab, n int) []*colors.LCHab {
	newColors := make([]*colors.LCHab, n)
	for i := range n {
		var hueOffsetSign float32 = 1
		if i%2 == 0 {
			hueOffsetSign = -1
		}
		newColors[i] = &colors.LCHab{
			L: dominant.L + 30*(rand.Float32()-0.5),
			C: dominant.C * 0.85,
			H: dominant.H + hueOffsetSign*30*float32(i/2+1),
		}
	}
	return newColors
}

func dominantColor(distanceMat [][]float64) int {
	var bestColor int
	var bestScore float64 = math.MaxFloat64

	for i := range distanceMat {
		var score float64
		for j := range distanceMat {
			if i == j {
				continue
			}
			score += distanceMat[i][j]
		}
		if score < bestScore {
			bestColor = i
			bestScore = score
		}
	}

	return bestColor
}

func averageDeltaE00(distanceMat [][]float64) float64 {
	var meanDeltaE00 float64 = 0
	count := 0
	for i := range distanceMat {
		for j := range distanceMat {
			count++
			meanDeltaE00 += distanceMat[i][j]
		}
	}
	return meanDeltaE00
}
