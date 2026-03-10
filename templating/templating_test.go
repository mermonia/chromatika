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

color0 = "{{.Color0.ToHex}}"
color1 = "{{.Color1.ToHex}}"
color2 = "{{.Color2.ToHex}}"
color3 = "{{.Color3.ToHex}}"
color4 = "{{.Color4.ToHex}}"
color5 = "{{.Color5.ToHex}}"
color6 = "{{.Color6.ToHex}}"
color7 = "{{.Color7.ToHex}}"

color8 = "{{.Color8.ToHex}}"
color9 = "{{.Color9.ToHex}}"
color10 = "{{.Color10.ToHex}}"
color11 = "{{.Color11.ToHex}}"
color12 = "{{.Color12.ToHex}}"
color13 = "{{.Color13.ToHex}}"
color14 = "{{.Color14.ToHex}}"
color15 = "{{.Color15.ToHex}}"
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
		Background: cols[0],
		Foreground: cols[1],
		Color0: cols[2],
		Color1: cols[3],
		Color2: cols[4],
		Color3: cols[5],
		Color4: cols[6],
		Color5: cols[7],
		Color6: cols[8],
		Color7: cols[9],
		Color8: cols[10],
		Color9: cols[11],
		Color10: cols[12],
		Color11: cols[13],
		Color12: cols[14],
		Color13: cols[15],
		Color14: cols[16],
		Color15: cols[17],
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
