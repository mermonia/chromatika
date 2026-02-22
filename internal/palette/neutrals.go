package palette

import (
	"github.com/mermonia/chromatika/internal/colors"
)

// The following are opinionated values, and are subject to change
const DARK_MODE_THRESHOLD float32 = 0.5
const MAXIMUM_NEUTRAL_CHROMA float32 = 0.12

func GenerateNeutrals(baseline *colors.LCHab) *NeutralColors {
	if baseline.L > DARK_MODE_THRESHOLD {
		return generateLightModeNeutrals(baseline)
	}
	return generateDarkModeNeutrals(baseline)
}

func generateLightModeNeutrals(baseline *colors.LCHab) *NeutralColors {
	return nil
}

func generateDarkModeNeutrals(baseline *colors.LCHab) *NeutralColors {
	return nil
}

func clampNeutralBase(original *colors.LCHab) *colors.LCHab {
	out := &colors.LCHab{
		L: original.L,
		C: min(original.C, 0.12),
		H: original.H,
	}

	return out
}
