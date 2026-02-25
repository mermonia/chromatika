package palette

import "github.com/mermonia/chromatika/internal/colors"

type Palette struct {
	Neutrals NeutralColors
	Accents  AccentColors
}

const NEUTRAL_COLORS int = 12
const NEUTRAL_BG_COLORS int = 3

const ACCENT_COLORS int = 8

type NeutralColors struct {
	Base,
	Mantle,
	Crust,

	Surface0,
	Surface1,
	Surface2,

	Overlay0,
	Overlay1,
	Overlay2,

	Text,
	Subtext0,
	Subtext1 *colors.LCHab
}

type AccentColors struct {
	Primary,
	Secondary,
	Tertiary,
	Quaternary,

	Error,
	Success,
	Warning,

	ExtraAccent *colors.LCHab
}
