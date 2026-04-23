package palette

import (
	"fmt"
	"math"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/extraction"
	"github.com/mermonia/chromatika/internal/utils"
)

const N_INITIAL_CLUSTERS = 16
const N_CLUSTERS_STEP = 4
const SEARCH_ITERATIONS = 2

const N_CANDIDATES = 8
const MIN_DELTAE00 = 15

func GeneratePalette(cfg *GenerationConfig) (*Palette, error) {
	/*
		Extract dominant colors multiple times, increasing the number of colors extracted
		each iteration.
	*/
	var suitableColors []colors.LCHab
	for i := range SEARCH_ITERATIONS {
		labCols, _, err := extraction.GetDominantColors(
			cfg.ImagePath,
			cfg.ScaleWidth,
			cfg.QuantInterval, clustering.FCMParameters{
				M: cfg.Fuzziness,
				E: cfg.Threshold,
				B: cfg.MaxIter,
				K: N_INITIAL_CLUSTERS + i*N_CLUSTERS_STEP,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("could not generate a %d color palette: %w", i+16, err)
		}

		lchCols := utils.Map(labCols, colors.LabToLCH)

		suitableColors = filterSuitableColors(lchCols, false)

		if len(suitableColors) >= N_CANDIDATES {
			break
		}
	}

	rawColors := generateCompleteCandidates(suitableColors)

	return getPaletteFromCandidates(rawColors, cfg.DarkMode)
}

func getPaletteFromCandidates(rawColors *RawColors, darkmode bool) (*Palette, error) {
	result := &Palette{Darkmode: darkmode}

	if darkmode {
		result.Background = rawColors.DarkNeutral
		result.Foreground = rawColors.LightNeutral
	} else {
		result.Background = rawColors.LightNeutral
		result.Foreground = rawColors.DarkNeutral
	}

	contrasts := make([]float64, len(rawColors.Colors))
	for i, c := range rawColors.Colors {
		contrast, err := wcag(result.Background, c)
		if err != nil {
			return nil, fmt.Errorf("could not calculate wcag: %w", err)
		}
		contrasts[i] = contrast
	}

	primary, err := getPrimaryColor(rawColors, contrasts)
	if err != nil {
		return nil, fmt.Errorf("could not find a suitable primary color: %w", err)
	}

	result.Primary = primary
	result.Secondary = getSecondaryColor(rawColors, contrasts, primary)
	result.Accent = getAccentColor(rawColors, contrasts, primary)
	result.ANSIBase = ansiColors(rawColors)
	result.ANSILighter = ansiLighterColors(result.ANSIBase, darkmode)

	return result, nil
}

func sigmoid(x, center, k float64) float64 {
	return 1.0 / (1.0 + math.Exp(-k*(x-center)))
}

func bell(x, center, sigma float64) float64 {
	d := (x - center) / sigma
	return math.Exp(-0.5 * d * d)
}

func scorePrimary(c colors.LCHab, bgContrast float64) float64 {
	contrastScore := sigmoid(bgContrast, 4.5, 1.2)
	chromaScore := bell(c.C, 55, 25)

	return contrastScore*0.6 + chromaScore*0.4
}

func scoreSecondary(candidate, primary colors.LCHab, bgContrast float64) float64 {
	deltaC := math.Abs(candidate.C - primary.C)
	deltaH := utils.HueDifference(candidate.H, primary.H)

	contrastScore := sigmoid(bgContrast, 4.5, 1.2)
	chromaScore := bell(deltaC, 8, 10)
	hueScore := bell(deltaH, 30, 15)

	return contrastScore*0.4 + chromaScore*0.3 + hueScore*0.3
}

func scoreAccent(candidate, primary colors.LCHab, targetHue float64, bgContrast float64) float64 {
	deltaC := math.Abs(candidate.C - primary.C)
	deltaH := utils.HueDifference(candidate.H, targetHue)

	contrastScore := sigmoid(bgContrast, 4.5, 1.2)
	chromaScore := bell(deltaC, 5, 15)
	hueScore := bell(deltaH, 0, 12)

	return contrastScore*0.2 + chromaScore*0.3 + hueScore*0.5
}

func getPrimaryColor(rawColors *RawColors, contrasts []float64) (colors.LCHab, error) {
	if len(rawColors.Colors) <= 0 {
		return colors.LCHab{}, fmt.Errorf("no colors available")
	}

	bestPrimaryScore := math.Inf(-1)
	var bestPrimary colors.LCHab
	for i, c := range rawColors.Colors {
		score := scorePrimary(c, contrasts[i])
		if score > bestPrimaryScore {
			bestPrimaryScore = score
			bestPrimary = c
		}
	}

	return bestPrimary, nil
}

func getSecondaryColor(rawColors *RawColors, contrasts []float64, primary colors.LCHab) colors.LCHab {
	bestSecondaryScore := math.Inf(-1)
	var bestSecondary colors.LCHab
	for i, c := range rawColors.Colors {
		if c == primary {
			continue
		}
		score := scoreSecondary(c, primary, contrasts[i])
		if score > bestSecondaryScore {
			bestSecondaryScore = score
			bestSecondary = c
		}
	}

	if bestSecondaryScore > 0.5 {
		return bestSecondary
	}

	return colors.LCHab{
		L: colors.RegularizeLuminosity(primary.L - 15),
		C: colors.RegularizeChroma(primary.C - 8),
		H: colors.RegularizeHue(primary.H + 30),
	}
}

func getAccentColor(rawColors *RawColors, contrasts []float64, primary colors.LCHab) colors.LCHab {
	targetHue := colors.RegularizeHue(primary.H + 180)

	bestAccentScore := math.Inf(-1)
	var bestAccent colors.LCHab
	for i, c := range rawColors.Colors {
		if c == primary {
			continue
		}

		score := scoreAccent(c, primary, targetHue, contrasts[i])
		if score > bestAccentScore {
			bestAccentScore = score
			bestAccent = c
		}
	}

	if bestAccentScore > 0.5 {
		return bestAccent
	}

	return colors.LCHab{
		L: colors.RegularizeLuminosity(primary.L - 5),
		C: colors.RegularizeChroma(primary.C - 5),
		H: targetHue,
	}
}

func ansiColors(rawColors *RawColors) [8]colors.LCHab {
	var result [8]colors.LCHab

	for i, hue := range ansiHues {
		if hue < 0 {
			continue
		}
		result[i] = getSimilarColorWithHue(rawColors, hue)
	}

	result[0] = rawColors.DarkNeutral
	result[7] = rawColors.LightNeutral

	return result
}

func getSimilarColorWithHue(rawColors *RawColors, hue float64) colors.LCHab {
	var avgL float64 = 0
	var avgC float64 = 0
	for _, color := range rawColors.Colors {
		avgL += color.L
		avgC += color.C
	}
	avgL /= float64(len(rawColors.Colors))
	avgC /= float64(len(rawColors.Colors))

	return colors.LCHab{
		L: min(80, max(45, avgL)),
		C: min(60, max(40, avgC)),
		H: hue,
	}
}

func ansiLighterColors(baseAnsi [8]colors.LCHab, darkmode bool) [8]colors.LCHab {
	var result [8]colors.LCHab

	for i, color := range baseAnsi {
		newColor := derivedColor(color, darkmode)
		result[i] = newColor
	}

	return result
}

func derivedColor(base colors.LCHab, darkmode bool) colors.LCHab {
	if darkmode {
		return colors.LCHab{
			L: base.L * 0.94,
			C: colors.RegularizeChroma(base.C + 8),
			H: colors.RegularizeHue(base.H + 2),
		}
	}

	return colors.LCHab{
		L: colors.RegularizeLuminosity(base.L * 1.09),
		C: base.C,
		H: colors.RegularizeHue(base.H + 2),
	}
}

/*
Returns a set of colors that are distinct enough to form a palette.
*/
func filterSuitableColors(cols []colors.LCHab, fast bool) []colors.LCHab {
	g := getColorSimilarityGraph(cols, MIN_DELTAE00)
	suitableColors := make([]colors.LCHab, 0, len(cols))

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

func getColorSimilarityGraph(cols []colors.LCHab, minDelta float64) *Graph {
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

func generateCompleteCandidates(base []colors.LCHab) *RawColors {
	// Base colors
	newBaseColors := make([]colors.LCHab, len(base))
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
	var darkNeutral, lightNeutral colors.LCHab
	foundDarkNeutral := false
	foundLightNeutral := false
	for i := len(newBaseColors) - 1; i >= 0; i-- {
		color := newBaseColors[i]

		if color.C < 8 && color.L > 92 {
			foundLightNeutral = true
			lightNeutral = color
			newBaseColors = append(newBaseColors[:i], newBaseColors[i+1:]...)
		}

		if color.C < 15 && color.L < 15 {
			foundDarkNeutral = true
			darkNeutral = color
			newBaseColors = append(newBaseColors[:i], newBaseColors[i+1:]...)
		}
	}

	// Generate dark neutral color
	if !foundDarkNeutral {
		darkNeutral = colors.LCHab{
			L: 8,
			C: 10,
			H: dominantColor.H,
		}
	}

	// Generate light neutral color
	if !foundLightNeutral {
		lightNeutral = colors.LCHab{
			L: 92,
			C: 6,
			H: dominantColor.H,
		}
	}

	rawColors := &RawColors{
		LightNeutral: lightNeutral,
		DarkNeutral:  darkNeutral,
	}

	rawColors.Colors = newBaseColors
	return rawColors
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
