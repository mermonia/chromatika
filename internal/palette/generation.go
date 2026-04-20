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
const SEARCH_ITERATIONS = 3

const N_CANDIDATES = 8
const MIN_DELTAE00 = 20

func GeneratePalette(cfg *GenerationConfig) (*Palette, error) {
	/*
		Extract dominant colors multiple times, increasing the number of colors extracted
		each iteration.
	*/
	var suitableColors []*colors.LCHab
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

	if cfg.DarkMode {
		return getDarkModePalette(rawColors)
	}
	return getLightModePalette(rawColors)

}

func getDarkModePalette(rawColors *RawColors) (*Palette, error) {
	primary, secondary, accent, err := accentColors(rawColors)
	if err != nil {
		return nil, fmt.Errorf("could not obtain primary colors: %w", err)
	}

	ansiColors := ansiColors(rawColors)
	ansiLighterColors := ansiLighterColors(ansiColors, true)

	return &Palette{
		Background: rawColors.DarkNeutral,
		Foreground: rawColors.LightNeutral,

		Primary: primary,
		Secondary: secondary,
		Accent: accent,

		ANSIBase: [8]*colors.LCHab(ansiColors),
		ANSILighter: [8]*colors.LCHab(ansiLighterColors),
	}, nil
}

func getLightModePalette(rawColors *RawColors) (*Palette, error) {
	primary, secondary, accent, err := accentColors(rawColors)
	if err != nil {
		return nil, fmt.Errorf("could not obtain primary colors: %w", err)
	}

	ansiColors := ansiColors(rawColors)
	ansiLighterColors := ansiLighterColors(ansiColors, false)

	return &Palette{
		Background: rawColors.LightNeutral,
		Foreground: rawColors.DarkNeutral,

		Primary: primary,
		Secondary: secondary,
		Accent: accent,

		ANSIBase: ansiColors,
		ANSILighter: ansiLighterColors,
	}, nil
}

func sigmoid(x, center, k float64) float32 {
	return float32(1.0 / (1.0 + math.Exp(-k*(x-center))))
}

func bell(x, center, sigma float64) float32 {
	d := (x - center) / sigma
	return float32(math.Exp(-0.5 * d * d))
}

func scorePrimary(c *colors.LCHab, darkContrast, lightContrast float32) float32 {
	// Best contrast against either neutral — reward versatility
	bestContrast := math.Max(float64(darkContrast), float64(lightContrast))
	contrastScore := sigmoid(bestContrast, 4.5, 1.2)

	// Chroma: vivid but not garish
	chromaScore := bell(float64(c.C), 55, 25)

	return contrastScore*0.6 + chromaScore*0.4
}

func scoreSecondary(candidate, primary *colors.LCHab) float32 {
	deltaC := math.Abs(float64(candidate.C - primary.C))
	deltaH := float64(utils.HueDifference(candidate.H, primary.H))

	// Similar but slightly less chromatic
	chromaScore := bell(deltaC, 8, 10)

	// Analogous hue — close but not identical
	hueScore := bell(deltaH, 30, 15)

	return chromaScore*0.5 + hueScore*0.5
}

func scoreAccent(candidate, primary, secondary *colors.LCHab, targetHue float32) float32 {
	deltaC := math.Abs(float64(candidate.C - primary.C))
	deltaH := float64(utils.HueDifference(candidate.H, targetHue))

	// Similar chroma to primary
	chromaScore := bell(deltaC, 5, 15)

	// Close to target hue
	hueScore := bell(deltaH, 0, 12)

	// Penalize if too similar to secondary
	secondaryDeltaH := float64(utils.HueDifference(candidate.H, secondary.H))
	separationScore := sigmoid(secondaryDeltaH, 20, 0.3)

	return chromaScore*0.45 + hueScore*0.45 + separationScore*0.1
}

func accentColors(rawColors *RawColors) (*colors.LCHab, *colors.LCHab, *colors.LCHab, error) {
	bestPrimaryScore := float32(math.Inf(-1))
	bestPrimaryIdx := -1

	fmt.Printf("available candidates: %s\n", rawColors.Colors)

	for i, c := range rawColors.Colors {
		darkContrast, err := wcag(rawColors.DarkNeutral, c)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("wcag dark: %w", err)
		}
		lightContrast, err := wcag(rawColors.LightNeutral, c)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("wcag light: %w", err)
		}
		score := scorePrimary(c, darkContrast, lightContrast)
		if score > bestPrimaryScore {
			bestPrimaryScore = score
			bestPrimaryIdx = i
		}
	}

	if bestPrimaryIdx == -1 {
		return nil, nil, nil, fmt.Errorf("no colors available")
	}

	primary := rawColors.Colors[bestPrimaryIdx]

	// --- Secondary: scored, analogous hue, different lightness ---
	bestSecondaryScore := float32(math.Inf(-1))
	var bestSecondary *colors.LCHab

	for i, c := range rawColors.Colors {
		if i == bestPrimaryIdx {
			continue
		}
		score := scoreSecondary(c, primary)
		if score > bestSecondaryScore {
			bestSecondaryScore = score
			bestSecondary = c
		}

		fmt.Printf("color %d sec score: %f\n", i, score)
	}

	var secondary *colors.LCHab
	if bestSecondary != nil && bestSecondaryScore > 0.3 {
		secondary = bestSecondary
	} else {
		secondary = &colors.LCHab{
			L: colors.RegularizeLuminosity(primary.L - 15),
			C: colors.RegularizeChroma(primary.C - 8),
			H: colors.RegularizeHue(primary.H + 30),
		}
	}

	// --- Accent: triadic offset from primary, away from secondary ---
	hueDelta := float32(30)
	diff := utils.SignedHueDifference(primary.H, secondary.H)
	if diff < 0 {
		// Secondary is on the +H side, push accent the other way
		hueDelta = -30
	}
	targetHue := colors.RegularizeHue(primary.H + hueDelta)
	fmt.Printf("pH: %f, sH: %f, diff: %f, target: %f\n", primary.H, secondary.H, diff, targetHue)

	bestAccentScore := float32(math.Inf(-1))
	var bestAccent *colors.LCHab

	for i, c := range rawColors.Colors {
		if i == bestPrimaryIdx {
			continue
		}
		score := scoreAccent(c, primary, secondary, targetHue)
		if score > bestAccentScore {
			bestAccentScore = score
			bestAccent = c
		}

		fmt.Printf("color %d accent score: %f\n", i, score)
	}

	var accent *colors.LCHab
	if bestAccent != nil && bestAccentScore > 0.3 {
		accent = bestAccent
	} else {
		accent = &colors.LCHab{
			L: colors.RegularizeLuminosity(primary.L - 10),
			C: colors.RegularizeChroma(primary.C),
			H: targetHue,
		}
	}

	return primary, secondary, accent, nil
}

func ansiColors(rawColors *RawColors) [8]*colors.LCHab {
	var result [8]*colors.LCHab

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

func getSimilarColorWithHue(rawColors *RawColors, hue float32) *colors.LCHab {
	var avgL float32 = 0
	var avgC float32 = 0
	for _, color := range rawColors.Colors {
		avgL += color.L
		avgC += color.C
	}
	avgL /= float32(len(rawColors.Colors))
	avgC /= float32(len(rawColors.Colors))

	return &colors.LCHab{
		L: min(80, max(45,avgL)),
		C: min(60, max(40,avgC)),
		H: hue,
	}
}

func ansiLighterColors(baseAnsi [8]*colors.LCHab, darkmode bool) [8]*colors.LCHab {
	var result [8]*colors.LCHab

	for i, color := range baseAnsi {
		newColor := derivedColor(color, darkmode)
		result[i] = newColor
	}

	return result
}

func derivedColor(base *colors.LCHab, darkmode bool) *colors.LCHab  {
	if darkmode {
		return &colors.LCHab{
			L: base.L * 0.94,
			C: colors.RegularizeChroma(base.C + 8),
			H: colors.RegularizeHue(base.H + 2),
		}
	}

	return &colors.LCHab{
		L: colors.RegularizeLuminosity(base.L * 1.09),
		C: base.C,
		H: colors.RegularizeHue(base.H + 2),
	}
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

func generateCompleteCandidates(base []*colors.LCHab) *RawColors {
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
