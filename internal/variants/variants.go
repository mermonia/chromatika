package variants

import (
	"github.com/mermonia/chromatika/internal/colors"
)

const N_VARIANTS int = 10
const PREFERRED_ANCHOR int = 5
const DELTA_L float64 = 7.5

func GenerateVariants(base colors.LCHab) *Variants {
	anchor := PREFERRED_ANCHOR

	for base.L <= float64(N_VARIANTS-anchor-1)*DELTA_L && anchor < N_VARIANTS-1 {
		anchor++
	}

	for base.L >= 100-float64(anchor)*DELTA_L && anchor > 0 {
		anchor--
	}

	cols := generateVariantsFromAnchor(base, anchor)

	return &Variants{
		BaseIndex:   anchor,
		Variant_50:  cols[0],
		Variant_100: cols[1],
		Variant_200: cols[2],
		Variant_300: cols[3],
		Variant_400: cols[4],
		Variant_500: cols[5],
		Variant_600: cols[6],
		Variant_700: cols[7],
		Variant_800: cols[8],
		Variant_900: cols[9],
	}
}

func generateVariantsFromAnchor(base colors.LCHab, anchor int) []colors.LCHab {
	result := make([]colors.LCHab, N_VARIANTS)
	result[anchor] = base

	for i := anchor + 1; i < N_VARIANTS; i++ {
		result[i] = darkerVariant(result[i-1])
	}

	for i := anchor - 1; i >= 0; i-- {
		result[i] = lighterVariant(result[i+1])
	}

	return result
}

func darkerVariant(base colors.LCHab) colors.LCHab {
	return colors.LCHab{
		L: colors.RegularizeLuminosity(base.L - DELTA_L),
		C: base.C,
		H: base.H,
	}
}

func lighterVariant(base colors.LCHab) colors.LCHab {
	return colors.LCHab{
		L: colors.RegularizeLuminosity(base.L + DELTA_L),
		C: base.C,
		H: base.H,
	}
}
