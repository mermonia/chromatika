package format

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/mermonia/chromatika/internal/palette"
)

type PaletteFormatter interface {
	Format(*palette.Palette) string
}

type TOMLPaletteFormatter struct{}
type ASCIIPaletteFormatter struct{}

func (*TOMLPaletteFormatter) Format(pal *palette.Palette) string {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(pal); err != nil {
		return ""
	}
	return buf.String()
}

func (*ASCIIPaletteFormatter) Format(pal *palette.Palette) string {
	res := ""

	res += fmt.Sprintf("Background: %s\n", pal.Background)
	res += fmt.Sprintf("Foreground: %s\n", pal.Foreground)

	res += fmt.Sprintf("Primary:    %s\n", pal.Primary)
	res += fmt.Sprintf("Secondary:  %s\n", pal.Secondary)
	res += fmt.Sprintf("Accent:     %s\n", pal.Accent)

	for i, color := range pal.ANSIBase {
		res += fmt.Sprintf("Color[%d]: %s\n", i, color)
	}

	for i, color := range pal.ANSILighter {
		res += fmt.Sprintf("Derived Color[%d]: %s\n", i, color)
	}

	return res
}
