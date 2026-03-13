package templating

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/palette"
)

func TestCreateFileFromFile(t *testing.T) {
	textToRender := `
background = "{{.Background.ToHex}}"
foreground = "{{.Foreground.ToHex}}"

color0 = "{{(index .BaseColors 0).ToHex}}"
color1 = "{{(index .BaseColors 1).ToHex}}"
color2 = "{{(index .BaseColors 2).ToHex}}"
color3 = "{{(index .BaseColors 3).ToHex}}"
color4 = "{{(index .BaseColors 4).ToHex}}"
color5 = "{{(index .BaseColors 5).ToHex}}"
color6 = "{{(index .BaseColors 6).ToHex}}"
color7 = "{{(index .BaseColors 7).ToHex}}"

color8 = "{{(index .DerivedColors 0).ToHex}}"
color9 = "{{(index .DerivedColors 1).ToHex}}"
color10 = "{{(index .DerivedColors 2).ToHex}}"
color11 = "{{(index .DerivedColors 3).ToHex}}"
color12 = "{{(index .DerivedColors 4).ToHex}}"
color13 = "{{(index .DerivedColors 5).ToHex}}"
color14 = "{{(index .DerivedColors 6).ToHex}}"
color15 = "{{(index .DerivedColors 7).ToHex}}"
	`

	expectedResult := `
background = "c14e78"
foreground = "c34e71"

color0 = "c44e6a"
color1 = "c44f63"
color2 = "c4505c"
color3 = "c35255"
color4 = "c1544e"
color5 = "bf5648"
color6 = "bc5942"
color7 = "b95c3c"

color8 = "b65e36"
color9 = "b26131"
color10 = "ae642c"
color11 = "a96727"
color12 = "a46a23"
color13 = "9e6d1f"
color14 = "996f1c"
color15 = "937219"
	`

	baseColor := &colors.LCHab{
		L: 50,
		C: 50,
		H: 0,
	}

	cols := make([]*colors.LCHab, 18)
	for i := range cols {
		cols[i] = &colors.LCHab{
			L: baseColor.L,
			C: baseColor.C,
			H: baseColor.H + 5*float32(i),
		}
	}

	pal := &palette.Palette{
		Background:    cols[0],
		Foreground:    cols[1],
		BaseColors:    [8]*colors.LCHab(cols[2:10]),
		DerivedColors: [8]*colors.LCHab(cols[10:]),
	}

	tempDir := t.TempDir()

	src, err := os.CreateTemp(tempDir, "srcFile")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	dst, err := os.CreateTemp(tempDir, "dstFile")
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	// write test text to file
	if _, err := src.WriteString(textToRender); err != nil {
		log.Fatal(err)
	}

	// render to destination file
	if err := RenderFileToFile(src.Name(), dst.Name(), pal); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(dst)
	if err != nil {
		log.Fatal(err)
	}

	result := string(data)

	if expectedResult != result {
		diff := cmp.Diff(expectedResult, result)
		fmt.Print(diff)
		t.Fatal("template result is different than expected")
	}
}
