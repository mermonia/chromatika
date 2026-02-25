package palette

import (
	"fmt"

	"github.com/mermonia/chromatika/internal/colors"
)

// The following are opinionated values, and are subject to change

// The base luminosity value threshold for dark mode; any palette generated from
// a color with a higher luminosity will be a light mode palette.
const DARK_MODE_LUMINOSITY_THRESHOLD float32 = 50.0

// The luminosity values for the darkest and brightest colors of the neutral
// palette, respectively
const DARK_MINIMUM_NEUTRAL_LUMINOSITY float32 = 10.0
const DARK_MAXIMUM_NEUTRAL_LUMINOSITY float32 = 90.0

const LIGHT_MINIMUM_NEUTRAL_LUMINOSITY float32 = 5.0
const LIGHT_MAXIMUM_NEUTRAL_LUMINOSITY float32 = 90.0

// For each generated color, the proportion between its luminosity and chroma values.
const DARK_FIRST_LUMINOSITY_CHROMA_PROPORTION float32 = 0.5
const DARK_LAST_LUMINOSITY_CHROMA_PROPORTION float32 = 4.0

const LIGHT_FIRST_LUMINOSITY_CHROMA_PROPORTION float32 = 1.0
const LIGHT_LAST_LUMINOSITY_CHORMA_PROPORTION float32 = 15.0

// The luminosity difference between background colors (base, mantle and crust).
const BG_COLORS_LUMINOSITY_DELTA float32 = 4.0

func GenerateNeutrals(baseline *colors.LCHab) *NeutralColors {
	if baseline.L > DARK_MODE_LUMINOSITY_THRESHOLD {
		return generateLightModeNeutrals(baseline)
	}
	return generateDarkModeNeutrals(baseline)
}

func generateLightModeNeutrals(baseline *colors.LCHab) *NeutralColors {
	luminosities := generateLuminosityGradient(
		LIGHT_MINIMUM_NEUTRAL_LUMINOSITY,
		LIGHT_MAXIMUM_NEUTRAL_LUMINOSITY,
		BG_COLORS_LUMINOSITY_DELTA,
		NEUTRAL_COLORS,
		NEUTRAL_BG_COLORS,
	)

	chromas := generateChromaGradient(
		luminosities,
		LIGHT_FIRST_LUMINOSITY_CHROMA_PROPORTION,
		LIGHT_LAST_LUMINOSITY_CHORMA_PROPORTION,
	)

	cols := make([]*colors.LCHab, len(luminosities))
	for i := range len(luminosities) {
		cols[i] = &colors.LCHab{
			L: luminosities[i],
			C: chromas[i],
			H: baseline.H,
		}
	}

	return &NeutralColors{
		Text:     cols[0],
		Subtext0: cols[1],
		Subtext1: cols[2],
		Overlay2: cols[3],
		Overlay1: cols[4],
		Overlay0: cols[5],
		Surface2: cols[6],
		Surface1: cols[7],
		Surface0: cols[8],
		Crust:    cols[9],
		Mantle:   cols[10],
		Base:     cols[11],
	}
}

func generateDarkModeNeutrals(baseline *colors.LCHab) *NeutralColors {
	luminosities := generateLuminosityGradient(
		DARK_MINIMUM_NEUTRAL_LUMINOSITY,
		DARK_MAXIMUM_NEUTRAL_LUMINOSITY,
		BG_COLORS_LUMINOSITY_DELTA,
		NEUTRAL_COLORS,
		NEUTRAL_BG_COLORS,
	)

	chromas := generateChromaGradient(
		luminosities,
		DARK_FIRST_LUMINOSITY_CHROMA_PROPORTION,
		DARK_LAST_LUMINOSITY_CHROMA_PROPORTION,
	)

	cols := make([]*colors.LCHab, len(luminosities))
	for i := range len(luminosities) {
		cols[i] = &colors.LCHab{
			L: luminosities[i],
			C: chromas[i],
			H: baseline.H,
		}
	}

	return &NeutralColors{
		Crust:    cols[0],
		Mantle:   cols[1],
		Base:     cols[2],
		Surface0: cols[3],
		Surface1: cols[4],
		Surface2: cols[5],
		Overlay0: cols[6],
		Overlay1: cols[7],
		Overlay2: cols[8],
		Subtext0: cols[9],
		Subtext1: cols[10],
		Text:     cols[11],
	}
}

func generateLuminosityGradient(minL, maxL, bgDelta float32, nColors, nBgColors int) []float32 {
	luminosities := make([]float32, nColors)
	currentLuminosity := minL

	// bg colors luminosities
	for i := range nBgColors {
		luminosities[i] = currentLuminosity
		currentLuminosity += bgDelta
	}
	currentLuminosity -= bgDelta

	luminosityDelta := (maxL - currentLuminosity) / (float32(nColors - nBgColors))
	// other colors luminosities
	for i := nBgColors; i < nColors; i++ {
		currentLuminosity += luminosityDelta
		luminosities[i] = currentLuminosity
	}

	return luminosities
}

func generateChromaGradient(luminosities []float32, firstLC, lastLC float32) []float32 {
	length := len(luminosities)
	chromas := make([]float32, length)

	proportionDelta := (lastLC - firstLC) / float32(length)
	currentProportion := firstLC

	for i := range length {
		chromas[i] = luminosities[i] / currentProportion
		currentProportion += proportionDelta
	}

	return chromas
}

func (n *NeutralColors) String() string {
	return fmt.Sprintf("Neutral Colors\n"+
		"\nBackground Colors:\n"+
		"Base: %v\n"+
		"Mantle: %v\n"+
		"Crust: %v\n"+
		"\nSurface Colors:\n"+
		"Surface0: %v\n"+
		"Surface1: %v\n"+
		"Surface2: %v\n"+
		"\nOverlay Colors:\n"+
		"Overlay0: %v\n"+
		"Overlay1: %v\n"+
		"Overlay2: %v\n"+
		"\nText Colors:\n"+
		"Subtext0: %v\n"+
		"Subtext1: %v\n"+
		"Text: %v\n",
		n.Base, n.Mantle, n.Crust,
		n.Surface0, n.Surface1, n.Surface2,
		n.Overlay0, n.Overlay1, n.Overlay2,
		n.Subtext0, n.Subtext1, n.Text,
	)
}
