package input

import (
	"strconv"
	"strings"

	"github.com/mermonia/chromatika/internal/colors"
)

func RGBfromString(str string) *colors.Rgb {
	trimmed := strings.TrimSpace(str)
	noOctothorpe := strings.TrimPrefix(trimmed, "#")

	r, _ := strconv.ParseUint(noOctothorpe[0:2], 16, 8)
	g, _ := strconv.ParseUint(noOctothorpe[2:4], 16, 8)
	b, _ := strconv.ParseUint(noOctothorpe[4:6], 16, 8)

	return &colors.Rgb{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}
}
