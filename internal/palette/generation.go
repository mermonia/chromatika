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

const SEARCH_ITERATIONS = 4
const N_CLUSTERS_STEP = 4

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
		fmt.Printf("found %d suitable colors\n", len(suitableColors))

		if len(suitableColors) >= N_BASE_COLORS {
			break
		}
	}

	for _, color := range suitableColors {
		fmt.Printf("L: %f, C: %f, H: %f\n%s\n\n", color.L, color.C, color.H, color)
	}
	fmt.Println()

	if len(suitableColors) < N_BASE_COLORS {
		suitableColors = fillColorPalette(suitableColors, N_BASE_COLORS)
	}

	fmt.Printf("NEW PALETTE: (%d colors)\n", len(suitableColors))
	for _, color := range suitableColors {
		fmt.Printf("L: %f, C: %f, H: %f\n%s\n\n", color.L, color.C, color.H, color)
	}
	fmt.Println()

	return nil, nil
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

func fillColorPalette(base []*colors.LCHab, targetSize int) []*colors.LCHab {
	nMissingColors := targetSize - len(base)
	if nMissingColors <= 0 {
		return base
	}

	// Newly generated colors
	var newColors []*colors.LCHab

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
	for _, color := range base {
		if color.C < 5 {
			if color.L > 92 {
				lightNeutral = color
			} else if color.L < 15 {
				darkNeutral = color
			}
		}
	}

	// Generate dark neutral color
	if darkNeutral == nil {
		newColors = append(newColors, &colors.LCHab{
			L: 5,
			C: 15,
			H: dominantColor.H,
		})
		nMissingColors--
	}

	// Generate light neutral color
	if lightNeutral == nil {
		newColors = append(newColors, &colors.LCHab{
			L: 92,
			C: 10,
			H: dominantColor.H,
		})
		nMissingColors--
	}

	// Mean deltaE00 for all base colors
	meanDeltaE00 := averageDeltaE00(deltaE00Mat)

	if meanDeltaE00 < 25.0 {
		// contrast-neutral strategy
		newColors = append(newColors, contrastNeutral(dominantColor, nMissingColors)...)
	} else {
		// harmonic expand strategy
		newColors = append(newColors, harmonicExpand(dominantColor, nMissingColors)...)
	}

	return append(base, newColors...)
}

func contrastNeutral(dominant *colors.LCHab, n int) []*colors.LCHab {
	fmt.Println("performing contrast neutral...")
	newColors := make([]*colors.LCHab, n)
	if n <= 0 {
		return newColors
	}

	accentL := dominant.L
	if accentL > 50 {
		accentL += 10
	} else {
		accentL -= 10
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
	fmt.Println("performing harmonic expand...")
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
