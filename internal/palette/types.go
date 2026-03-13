package palette

import (
	"github.com/mermonia/chromatika/internal/colors"
)

type Palette struct {
	// Special colors
	Background,
	Foreground *colors.LCHab

	// Base colors
	BaseColors [8]*colors.LCHab

	// Derived colors, variants of the 0-7 base colors
	DerivedColors [8]*colors.LCHab
}

type RawColors struct {
	// Neutral Colors
	DarkNeutral,
	LightNeutral *colors.LCHab

	Colors [8]*colors.LCHab
}

func (rc *RawColors) String() string {
	res := rc.DarkNeutral.String() + rc.LightNeutral.String() + " "
	for i := range rc.Colors {
		res += rc.Colors[i].String()
	}
	return res
}
