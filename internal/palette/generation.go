package palette

import (
	"fmt"
	"math"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/extraction"
	"github.com/mermonia/chromatika/internal/utils"
)

const N_CLUSTERS = 8
const N_CANDIDATES = 8
const MIN_DELTAE00 = 15

const (
	MIN_BG_SCORE        = 0.6
	MIN_FG_SCORE        = 0.6
	MIN_SECONDARY_SCORE = 0.3
	MIN_ACCENT_SCORE    = 0.5
)

type Candidate struct {
	Color      colors.LCHab
	BgContrast float64
	Presence   float64
}

func GeneratePalette(cfg *GenerationConfig) (*Palette, error) {
	labCols, weights, err := extraction.GetDominantColors(
		cfg.ImagePath,
		cfg.ScaleWidth,
		cfg.QuantInterval, clustering.FCMParameters{
			M: cfg.Fuzziness,
			E: cfg.Threshold,
			B: cfg.MaxIter,
			K: N_CLUSTERS,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not generate a color palette: %w", err)
	}

	lchCols := utils.Map(labCols, colors.LabToLCH)

	candidates := make([]Candidate, len(lchCols))
	for i := range candidates {
		candidates[i].Color = lchCols[i]
		candidates[i].Presence = weights[i]
	}

	candidates = filterDistinctCandidates(candidates, false)

	bg := generateBackgroundColor(candidates, cfg.DarkMode)
	candidates = removeNearby(candidates, bg, MIN_DELTAE00)

	for i, c := range candidates {
		contrast, err := wcag(c.Color, bg)
		if err != nil {
			return nil, fmt.Errorf("could not get contrast with bg: %w", err)
		}
		candidates[i].BgContrast = contrast
	}

	primary := generatePrimaryColor(candidates)
	candidates = removeNearby(candidates, primary, MIN_DELTAE00)

	fg := generateForegroundColor(candidates, primary, cfg.DarkMode)
	candidates = removeNearby(candidates, fg, MIN_DELTAE00)

	secondary := generateSecondaryColor(candidates, primary)
	candidates = removeNearby(candidates, secondary, MIN_DELTAE00)

	accent := generateAccentColor(candidates, primary)
	candidates = removeNearby(candidates, accent, MIN_DELTAE00)

	ansiColors := generateAnsiColors(primary, bg, cfg.DarkMode)
	ansiVariants := generateAnsiVariants(ansiColors, cfg.DarkMode)

	return &Palette{
		Darkmode:    cfg.DarkMode,
		Background:  bg,
		Foreground:  fg,
		Primary:     primary,
		Secondary:   secondary,
		Accent:      accent,
		ANSIBase:    ansiColors,
		ANSILighter: ansiVariants,
	}, nil
}

/*
Primary color selection
*/
func generatePrimaryColor(candidates []Candidate) colors.LCHab {
	var bestScore float64 = -1
	var bestCandidate Candidate = candidates[0]

	for _, c := range candidates {
		score := scorePrimaryColor(c)
		if score > bestScore {
			bestScore = score
			bestCandidate = c
		}
	}

	return bestCandidate.Color
}

func scorePrimaryColor(candidate Candidate) float64 {
	const (
		weightPresence  = .42
		weightChroma    = .33
		weightContrast  = .2
		weightLightness = .05
	)

	presence := highPresenceScore(candidate.Presence)
	chroma := highChromaScore(candidate.Color.C)
	contrast := highContrastScore(candidate.BgContrast)

	midLightness := 1 - math.Abs(candidate.Color.L-50)/50
	midLightness = utils.Clamp01(midLightness)

	return weightPresence*presence +
		weightChroma*chroma +
		weightContrast*contrast +
		weightLightness*midLightness
}

/*
Background color generation
*/
func generateBackgroundColor(candidates []Candidate, darkmode bool) colors.LCHab {
	var bestScore float64 = -1
	var bestCandidate Candidate = candidates[0]

	for _, c := range candidates {
		score := scoreBackground(c, darkmode)
		if score > bestScore {
			bestScore = score
			bestCandidate = c
		}
	}

	if bestScore < MIN_BG_SCORE {
		return fallbackBackgroundColor(darkmode)
	}

	return bestCandidate.Color
}

func scoreBackground(candidate Candidate, darkmode bool) float64 {
	const (
		weightLightness = .60
		weightChroma    = .25
		weightPresence  = .15
	)

	targetL := 94.0
	targetSpread := 10.0
	if darkmode {
		targetL = 15.0
		targetSpread = 12.0
	}

	lightness := lerpScore(candidate.Color.L, targetL, targetSpread)
	chroma := lowChromaScore(candidate.Color.C)
	presence := lowPresenceScore(candidate.Presence)

	return weightLightness*lightness +
		weightChroma*chroma +
		weightPresence*presence
}

func fallbackBackgroundColor(darkmode bool) colors.LCHab {
	if darkmode {
		return colors.LCHab{
			L: 15,
			C: 6,
			H: 0,
		}
	}

	return colors.LCHab{
		L: 94,
		C: 6,
		H: 0,
	}
}

/*
Foreground color generation
*/
func generateForegroundColor(candidates []Candidate, primary colors.LCHab, darkmode bool) colors.LCHab {
	if len(candidates) <= 0 {
		return fallbackForegroundColor(primary, darkmode)
	}

	var bestScore float64 = -1
	var bestCandidate Candidate = candidates[0]

	for _, c := range candidates {
		score := scoreForegroundColor(c, darkmode)
		if score > bestScore {
			bestScore = score
			bestCandidate = c
		}
	}

	if bestScore < MIN_FG_SCORE {
		return fallbackForegroundColor(primary, darkmode)
	}

	return bestCandidate.Color
}

func scoreForegroundColor(candidate Candidate, darkmode bool) float64 {
	const (
		lightnessSpread = 14.0

		weightContrast  = .45
		weightLightness = .30
		weightChroma    = .15
		weightPresence  = .10
	)

	targetL := 90.0
	if !darkmode {
		targetL = 10.0
	}

	lightness := lerpScore(candidate.Color.L, targetL, lightnessSpread)
	chroma := lowChromaScore(candidate.Color.C)
	contrast := veryHighContrastScore(candidate.BgContrast)
	presence := lowPresenceScore(candidate.Presence)

	return weightContrast*contrast +
		weightLightness*lightness +
		weightChroma*chroma +
		weightPresence*presence
}

func fallbackForegroundColor(primary colors.LCHab, darkmode bool) colors.LCHab {
	if darkmode {
		return colors.LCHab{
			L: 90,
			C: 6,
			H: primary.H,
		}
	}
	return colors.LCHab{
		L: 10,
		C: 6,
		H: primary.H,
	}
}

/*
Secondary color generation
*/
func generateSecondaryColor(candidates []Candidate, primary colors.LCHab) colors.LCHab {
	if len(candidates) <= 0 {
		return fallbackSecondaryColor(primary)
	}

	var bestScore float64 = -1
	var bestCandidate Candidate = candidates[0]

	for _, c := range candidates {
		score := scoreSecondaryColor(c, primary)
		if score > bestScore {
			bestScore = score
			bestCandidate = c
		}
	}

	if bestScore < MIN_SECONDARY_SCORE {
		return fallbackSecondaryColor(primary)
	}

	return bestCandidate.Color
}

func scoreSecondaryColor(candidate Candidate, primary colors.LCHab) float64 {
	const (
		hueSpread       = 75.0
		lightnessSpread = 20.0
		chromaSpread    = 18.0

		weightHue       = .35
		weightLightness = .25
		weightChroma    = .20
		weightContrast  = .10
		weightPresence  = .10
	)

	targetHue := preferredSecondaryHue(primary.H)
	hue := circularScore(candidate.Color.H, targetHue, hueSpread)

	lightness := lerpScore(candidate.Color.L, primary.L, lightnessSpread)
	targetChroma := colors.RegularizeChroma(primary.C - 15)
	chroma := lerpScore(candidate.Color.C, targetChroma, chromaSpread)

	contrast := highContrastScore(candidate.BgContrast)
	presence := highPresenceScore(candidate.Presence)

	return weightHue*hue +
		weightLightness*lightness +
		weightChroma*chroma +
		weightContrast*contrast +
		weightPresence*presence
}

func fallbackSecondaryColor(primary colors.LCHab) colors.LCHab {
	hue := primary.H
	if hue >= 90 && hue <= 300 || hue > 330 {
		hue = colors.RegularizeHue(hue - 60)
	} else {
		hue = colors.RegularizeHue(hue + 60)
	}

	return colors.LCHab{
		L: primary.L,
		C: colors.RegularizeChroma(primary.C - 15),
		H: hue,
	}
}

/*
Accent color generation
*/
func generateAccentColor(candidates []Candidate, primary colors.LCHab) colors.LCHab {
	if len(candidates) <= 0 {
		return fallbackAccentColor(primary)
	}

	var bestScore float64 = -1
	var bestCandidate Candidate = candidates[0]

	for _, c := range candidates {
		score := scoreAccentColor(c, primary)
		if score > bestScore {
			bestScore = score
			bestCandidate = c
		}
	}

	if bestScore < MIN_ACCENT_SCORE {
		return fallbackAccentColor(primary)
	}

	return bestCandidate.Color
}

func scoreAccentColor(candidate Candidate, primary colors.LCHab) float64 {
	const (
		hueSpread       = 60.0
		lightnessSpread = 18.0

		weightHue       = 0.38
		weightChroma    = 0.20
		weightLightness = 0.18
		weightContrast  = 0.14
		weightPresence  = 0.10
	)

	targetHue := colors.RegularizeHue(primary.H + 180)

	hue := circularScore(candidate.Color.H, targetHue, hueSpread)
	lightness := lerpScore(candidate.Color.L, primary.L, lightnessSpread)
	contrast := highContrastScore(candidate.BgContrast)

	chroma := highChromaScore(candidate.Color.C)
	saturationBonus := utils.Clamp01((candidate.Color.C - primary.C + 5) / 20.0)

	presence := highPresenceScore(candidate.Presence)

	return weightHue*hue +
		weightChroma*chroma +
		weightLightness*lightness +
		weightContrast*contrast +
		weightPresence*(0.5*presence+0.5*saturationBonus)
}

func fallbackAccentColor(primary colors.LCHab) colors.LCHab {
	return colors.LCHab{
		L: primary.L,
		C: math.Max(colors.RegularizeChroma(primary.C+10), 50),
		H: colors.RegularizeHue(primary.H + 180),
	}
}

/*
ANSI color generation
*/
func generateAnsiColors(primary, bg colors.LCHab, darkmode bool) [8]colors.LCHab {
	var ansiColors [8]colors.LCHab

	lightness := ansiLightness(bg, darkmode)
	chroma := ansiChroma(primary)
	for i, h := range ansiHues {
		if h < 0 {
			continue
		}
		warpedHue := warpAnsiHue(h, primary.H)
		ansiColors[i] = colors.LCHab{
			L: lightness,
			C: chroma,
			H: warpedHue,
		}
	}

	ansiColors[0] = ansiBlack(bg, darkmode)
	ansiColors[7] = ansiWhite(bg, darkmode)

	return ansiColors
}

func generateAnsiVariants(originals [8]colors.LCHab, darkmode bool) [8]colors.LCHab {
	var ansiVariants [8]colors.LCHab
	for i, c := range originals {
		ansiVariants[i] = ansiVariant(c, darkmode)
	}
	return ansiVariants
}

func ansiLightness(bg colors.LCHab, darkmode bool) float64 {
	if darkmode {
		return utils.Clamp(bg.L+45, 50, 70)
	}
	return utils.Clamp(bg.L-45, 30, 45)
}

func ansiChroma(primary colors.LCHab) float64 {
	return math.Max(32, math.Min(primary.C+5, 60))
}

func warpAnsiHue(targetHue, primaryHue float64) float64 {
	distance := colors.SignedHueDistance(targetHue, primaryHue)
	return colors.RegularizeHue(targetHue + 0.15*distance)
}

func ansiBlack(bg colors.LCHab, darkmode bool) colors.LCHab {
	if darkmode {
		return colors.LCHab{
			L: utils.Clamp(bg.L-5, 2, 20),
			C: bg.C * .5,
			H: bg.H,
		}
	}
	return colors.LCHab{
		L: 20.0,
		C: bg.C * .35,
		H: bg.H,
	}
}

func ansiWhite(bg colors.LCHab, darkmode bool) colors.LCHab {
	if darkmode {
		return colors.LCHab{
			L: 90,
			C: bg.C * 0.4,
			H: bg.H,
		}
	}

	return colors.LCHab{
		L: 92,
		C: bg.C * 0.35,
		H: bg.H,
	}
}

func ansiVariant(original colors.LCHab, darkmode bool) colors.LCHab {
	if darkmode {
		return colors.LCHab{
			L: utils.Clamp(original.L+15, 0, 100),
			C: utils.Clamp(original.C+5, 0, 130),
			H: original.H,
		}
	}

	return colors.LCHab{
		L: utils.Clamp(original.L-15, 0, 100),
		C: utils.Clamp(original.C+5, 0, 130),
		H: original.H,
	}
}

/*
Utility functions
*/
func filterDistinctCandidates(candidates []Candidate, fast bool) []Candidate {
	g := getColorSimilarityGraph(candidates, MIN_DELTAE00)
	suitableColors := make([]Candidate, 0, len(candidates))

	var mis []int
	if fast {
		mis = MIS_Fast(g)
	} else {
		mis = MIS_Complete(g)
	}

	for _, idx := range mis {
		suitableColors = append(suitableColors, candidates[idx])
	}

	return suitableColors
}

func getColorSimilarityGraph(candidates []Candidate, minDelta float64) *Graph {
	g := NewGraph(len(candidates))
	for i := range candidates {
		for j := 0; j < i; j++ {
			delta := deltaE00(candidates[i].Color, candidates[j].Color)
			if delta < minDelta {
				g.AddEdge(i, j)
			}
		}
	}
	return g
}

func removeNearby(candidates []Candidate, centerColor colors.LCHab, minDelta float64) []Candidate {
	var newCandidates []Candidate
	for _, c := range candidates {
		if deltaE00(centerColor, c.Color) > minDelta {
			newCandidates = append(newCandidates, c)
		}
	}
	return newCandidates
}

func (c Candidate) String() string {
	return c.Color.String()
}

/*
Scoring helpers
*/
func lerpScore(value, target, spread float64) float64 {
	if spread <= 0 {
		return 0
	}
	d := math.Abs(value - target)
	return utils.Clamp01(1 - d/spread)
}

func circularScore(actual, target, spread float64) float64 {
	if spread <= 0 {
		return 0
	}
	return utils.Clamp01(1 - colors.HueDistance(actual, target)/spread)
}

func highPresenceScore(p float64) float64 {
	return p
}

func lowPresenceScore(p float64) float64 {
	return 1 - p
}

func highChromaScore(c float64) float64 {
	return utils.Clamp01((c - 8.0) / 24.0)
}

func lowChromaScore(c float64) float64 {
	return 1 - highChromaScore(c)
}

func highContrastScore(c float64) float64 {
	return utils.Clamp01((c - 1.0) / 8.0)
}

func veryHighContrastScore(c float64) float64 {
	return utils.Clamp01((c - 2.0) / 6.0)
}

func preferredSecondaryHue(primaryHue float64) float64 {
	hue := colors.RegularizeHue(primaryHue)
	if (hue >= 90 && hue <= 300) || hue > 330 {
		return colors.RegularizeHue(hue - 60)
	}
	return colors.RegularizeHue(hue + 60)
}
