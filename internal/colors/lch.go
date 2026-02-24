package colors

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
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
