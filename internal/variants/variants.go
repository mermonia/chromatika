package variants

import (
	"github.com/mermonia/chromatika/internal/colors"
)

const N_VARIANTS int = 10

func GenerateVariants(base colors.LCHab) *Variants {
	startL := 50.0 / float64(N_VARIANTS)
	deltaL := 100.0 / float64(N_VARIANTS)
	cols := make([]colors.LCHab, N_VARIANTS)
	for i := range N_VARIANTS {
		cols[i] = colors.LCHab{
			L: startL + deltaL*float64(i),
			C: base.C,
			H: base.H,
		}
	}

	// Assume 10 variants
	return &Variants{
		Variant50:  cols[9],
		Variant100: cols[8],
		Variant200: cols[7],
		Variant300: cols[6],
		Variant400: cols[5],
		Variant500: cols[4],
		Variant600: cols[3],
		Variant700: cols[2],
		Variant800: cols[1],
		Variant900: cols[0],
	}
}
