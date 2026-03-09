package palette

import "github.com/mermonia/chromatika/internal/colors"

type Palette struct {
	// Special colors
	Background,
	Foreground,

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
