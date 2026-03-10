package palette

import (
	"github.com/mermonia/chromatika/internal/colors"
)

type Palette struct {
	// Special colors
	Background,
	Foreground,
	Cursor,

	// Base colors
	Color0,
	Color1,
	Color2,
	Color3,
	Color4,
	Color5,
	Color6,
	Color7,

	// Derived colors, variants of the 0-7 base colors
	Color8,
	Color9,
	Color10,
	Color11,
	Color12,
	Color13,
	Color14,
	Color15 colors.LCHab
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
