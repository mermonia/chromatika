package palette

import (
	"fmt"
	"math"

	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/utils"
)

const PREFERRED_LUMINOSITY_DARK_MODE float32 = 70
const PREFERRED_LUMINOSITY_LIGHT_MODE float32 = 60
const LUMINOSITY_BIAS float32 = 0.75

const PREFERRED_CHROMA_DARK_MODE float32 = 50
const PREFERRED_CHROMA_LIGHT_MODE float32 = 70
const CHROMA_BIAS float32 = 0.75

const ANALOGOUS_DEGREE_SHIFT = 30
const TRIADIC_DEGREE_SHIFT = 120

const PRIMARY_STEP_SIZE = 30

func GenerateAccents(baseline *colors.LCHab, darkmode bool) *AccentColors {
	if darkmode {
		return generateDarkModeAccents(baseline)
	}
	return generateLightModeAccents(baseline)
}

func generateDarkModeAccents(baseline *colors.LCHab) *AccentColors {
	l := biasedLerp(PREFERRED_LUMINOSITY_DARK_MODE, baseline.L, LUMINOSITY_BIAS)
	c := biasedLerp(PREFERRED_CHROMA_DARK_MODE, baseline.C, CHROMA_BIAS)
	h := flattenHue(baseline.H, PRIMARY_STEP_SIZE)

	return &AccentColors{
		Primary: &colors.LCHab{
			L: l,
			C: c,
			H: h,
		},
		Secondary: &colors.LCHab{
			L: l,
			C: c * 0.9,
			H: colors.RegularizeHue(h + ANALOGOUS_DEGREE_SHIFT),
		},
		Tertiary: &colors.LCHab{
			L: utils.Clamp(l*1.2, 0, 100),
			C: c * 0.8,
			H: colors.RegularizeHue(h - ANALOGOUS_DEGREE_SHIFT),
		},
		Error: &colors.LCHab{
			L: PREFERRED_LUMINOSITY_DARK_MODE,
			C: utils.Clamp(c*1.2, 0, 100),
			H: 0,
		},
		Warning: &colors.LCHab{
			L: PREFERRED_LUMINOSITY_DARK_MODE,
			C: utils.Clamp(c*1.3, 0, 100),
			H: 60,
		},
		Success: &colors.LCHab{
			L: PREFERRED_LUMINOSITY_DARK_MODE,
			C: utils.Clamp(c*1.3, 0, 100),
			H: 120,
		},
		ExtraAccent0: &colors.LCHab{
			L: utils.Clamp(l*1.2, 0, 100),
			C: c * 0.8,
			H: colors.RegularizeHue(h + TRIADIC_DEGREE_SHIFT),
		},
		ExtraAccent1: &colors.LCHab{
			L: utils.Clamp(l*1.2, 0, 100),
			C: c * 0.8,
			H: colors.RegularizeHue(h - TRIADIC_DEGREE_SHIFT),
		},
	}
}

func generateLightModeAccents(baseline *colors.LCHab) *AccentColors {
	l := biasedLerp(PREFERRED_LUMINOSITY_LIGHT_MODE, baseline.L, LUMINOSITY_BIAS)
	c := biasedLerp(PREFERRED_CHROMA_LIGHT_MODE, baseline.C, CHROMA_BIAS)
	h := flattenHue(baseline.H, PRIMARY_STEP_SIZE)

	return &AccentColors{
		Primary: &colors.LCHab{
			L: l,
			C: c,
			H: h,
		},
		Secondary: &colors.LCHab{
			L: l,
			C: c * 0.9,
			H: colors.RegularizeHue(h + ANALOGOUS_DEGREE_SHIFT),
		},
		Tertiary: &colors.LCHab{
			L: l * 0.8,
			C: c * 0.8,
			H: colors.RegularizeHue(h - ANALOGOUS_DEGREE_SHIFT),
		},
		Error: &colors.LCHab{
			L: PREFERRED_LUMINOSITY_LIGHT_MODE,
			C: utils.Clamp(c*1.2, 0, 100),
			H: 0,
		},
		Warning: &colors.LCHab{
			L: PREFERRED_LUMINOSITY_LIGHT_MODE,
			C: utils.Clamp(c*1.3, 0, 100),
			H: 60,
		},
		Success: &colors.LCHab{
			L: PREFERRED_LUMINOSITY_LIGHT_MODE,
			C: utils.Clamp(c*1.3, 0, 100),
			H: 120,
		},
		ExtraAccent0: &colors.LCHab{
			L: l,
			C: c * 0.8,
			H: colors.RegularizeHue(h + TRIADIC_DEGREE_SHIFT),
		},
		ExtraAccent1: &colors.LCHab{
			L: l,
			C: c * 0.8,
			H: colors.RegularizeHue(h - TRIADIC_DEGREE_SHIFT),
		},
	}
}

func flattenHue(original, stepSize float32) float32 {
	return float32(math.Round(float64(original)/float64(stepSize))) * stepSize
}

func biasedLerp(a, b, bias float32) float32 {
	return a*(1-bias) + b*bias
}

func (a *AccentColors) String() string {
	return fmt.Sprintf(
		"Primary: %v\n"+
			"Secondary: %v\n"+
			"Tertiary: %v\n"+
			"Error: %v\n"+
			"Warning: %v\n"+
			"Success: %v\n"+
			"Extra Accent 0: %v\n"+
			"Extra Accent 1: %v\n",
		a.Primary, a.Secondary, a.Tertiary,
		a.Error, a.Warning, a.Success,
		a.ExtraAccent0, a.ExtraAccent1,
	)
}
