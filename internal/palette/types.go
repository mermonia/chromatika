package palette

import (
	"github.com/mermonia/chromatika/internal/colors"
)

type Palette struct {
	// Special colors
	Background,
	Foreground,

	Primary,
	Secondary,
	Accent *colors.LCHab

	ANSIBase,
	ANSILighter [8]*colors.LCHab
}

type RawColors struct {
	// Neutral Colors
	DarkNeutral,
	LightNeutral *colors.LCHab

	Colors []*colors.LCHab
}
