package palette

import (
	"fmt"

	"github.com/mermonia/chromatika/internal/colors"
)

// The following are opinionated values, and are subject to change
const DARK_MODE_THRESHOLD float32 = 0.5
const MAXIMUM_NEUTRAL_CHROMA float32 = 0.12

var LUMINOSITY_DELTAS []float32 = []float32{
	0,  // base
	-2, // mantle
	-4, // crust
	3,  // surface-0
	6,  // surface-1
	9,  // surface-2
	12, // overlay-0
	15, // overlay-1
	55, // text
	40, // subtext-0
	30, // subtext-1
}

var CHROMA_LIGHT_MODE []float32 = []float32{
	3,   // base
	2,   // mantle
	1.5, // crust
	3,   // surface-0
	3.5, // surface-1
	4,   // surface-2
	3,   // overlay-0
	2,   // overlay-1
	1,   // text
	1,   // subtext-0
	1,   // subtext-1
}

var CHROMA_DARK_MODE []float32 = []float32{
	4,    // base
	3,    // mantle
	2,    // crust
	4.5,  // surface-0
	5,    // surface-1
	5.5,  // surface-2
	3,    // overlay-0
	2,    // overlay-1
	1,    // text
	1.5,  // subtext-0
	1.75, // subtext-1
}

func GenerateNeutrals(baseline *colors.LCHab) *NeutralColors {
	if baseline.L > DARK_MODE_THRESHOLD {
		return generateLightModeNeutrals(baseline)
	}
	return generateDarkModeNeutrals(baseline)
}

func generateLightModeNeutrals(baseline *colors.LCHab) *NeutralColors {
	cols := make([]*colors.LCHab, len(LUMINOSITY_DELTAS))
	for i := range LUMINOSITY_DELTAS {
		cols[i] = &colors.LCHab{
			L: clampLuminosity(baseline.L - LUMINOSITY_DELTAS[i]),
			C: CHROMA_LIGHT_MODE[i],
			H: baseline.H,
		}
	}

	return &NeutralColors{
		Base:     cols[0],
		Mantle:   cols[1],
		Crust:    cols[2],
		Surface0: cols[3],
		Surface1: cols[4],
		Surface2: cols[5],
		Overlay0: cols[6],
		Overlay1: cols[7],
		Text:     cols[8],
		Subtext0: cols[9],
		Subtext1: cols[10],
	}
}

func generateDarkModeNeutrals(baseline *colors.LCHab) *NeutralColors {
	cols := make([]*colors.LCHab, len(LUMINOSITY_DELTAS))
	for i := range LUMINOSITY_DELTAS {
		cols[i] = &colors.LCHab{
			L: clampLuminosity(baseline.L + LUMINOSITY_DELTAS[i]),
			C: CHROMA_DARK_MODE[i],
			H: baseline.H,
		}
	}

	return &NeutralColors{
		Base:     cols[0],
		Mantle:   cols[1],
		Crust:    cols[2],
		Surface0: cols[3],
		Surface1: cols[4],
		Surface2: cols[5],
		Overlay0: cols[6],
		Overlay1: cols[7],
		Text:     cols[8],
		Subtext0: cols[9],
		Subtext1: cols[10],
	}
}

func clampLuminosity(l float32) float32 {
	newL := min(l, 100)
	newL = max(l, 0)

	return newL
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
		"\nText Colors:\n"+
		"Text: %v\n"+
		"Subtext0: %v\n"+
		"Subtext1: %v\n",
		n.Base, n.Mantle, n.Crust,
		n.Surface0, n.Surface1, n.Surface2,
		n.Overlay0, n.Overlay1,
		n.Text, n.Subtext0, n.Subtext1,
	)
}
