package colors

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/mermonia/chromatika/internal/utils"
)

func (c *LCHab) Render(width int) (string, error) {
	rgb, err := LCHtoRGB(c)
	if err != nil {
		return "", fmt.Errorf("could not convert from lab to rgb: %w", err)
	}
	style := lipgloss.NewStyle().
		Background(lipgloss.Color(
			fmt.Sprintf("#%02x%02x%02x", rgb.R, rgb.G, rgb.B),
		)).Width(width)

	return style.Render(" "), nil
}

func (c *LCHab) String() string {
	block, _ := c.Render(3)
	return block
}

func (c *LCHab) GetTemperature() float64 {
	return math.Cos(float64(c.H) - 60)
}

func (c *LCHab) ToHex() string {
	rgb, err := LCHtoRGB(c)
	if err != nil {
		return ""
	}
	return rgb.ToHex()
}

func Lighter(in *LCHab) *LCHab {
	return &LCHab{
		L: RegularizeLuminosity(in.L * 1.09),
		C: in.C,
		H: RegularizeHue(in.H + 2),
	}
}

func Darker(in *LCHab) *LCHab {
	return &LCHab{
		L: RegularizeLuminosity(in.L * 0.94),
		C: RegularizeChroma(in.C + 8),
		H: RegularizeHue(in.H + 2),
	}
}

func LCHtoLab(in *LCHab) *Lab {
	out := &Lab{
		L: in.L,
		A: in.C * float32(math.Cos(float64(in.H)*math.Pi/180)),
		B: in.C * float32(math.Sin(float64(in.H)*math.Pi/180)),
	}

	return out
}

func LCHtoRGB(in *LCHab) (*Rgb, error) {
	lab := LCHtoLab(in)
	return LabToRGB(lab)
}

func RegularizeHue(h float32) float32 {
	if h < 0 {
		return h + 360
	}
	if h > 360 {
		return h - 360
	}
	return h
}

func RegularizeLuminosity(l float32) float32 {
	return utils.Clamp(l, 0, 100)
}

func RegularizeChroma(c float32) float32 {
	return utils.Clamp(c, 0, 100)
}
