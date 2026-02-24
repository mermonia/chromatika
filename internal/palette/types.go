package palette

import "github.com/mermonia/chromatika/internal/colors"

type Palette struct {
	Neutrals NeutralColors
	Accents  AccentColors
}

type NeutralColors struct {
	Base,
	Mantle,
	Crust,

	Surface0,
	Surface1,
	Surface2,

	Overlay0,
	Overlay1,

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
