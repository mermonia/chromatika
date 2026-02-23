package input

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mermonia/chromatika/internal/colors"
)

func RGBfromString(str string) *colors.Rgb {
	trimmed := strings.Trim(str, " #")
	padded := fmt.Sprintf("%-6s", trimmed)

	r, _ := strconv.ParseUint(padded[0:2], 16, 8)
	g, _ := strconv.ParseUint(padded[2:4], 16, 8)
	b, _ := strconv.ParseUint(padded[4:6], 16, 8)

	return &colors.Rgb{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}
}
