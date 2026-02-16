package colors

import (
	"fmt"
	"image/color"
)

func PixelToLab(in color.Color) (*Lab, error) {
	rgb := PixelToRGB(in)
	lab, err := RGBtoLab(rgb)

	if err != nil {
		return nil, fmt.Errorf("could not convert rgb to lab: %w", err)
	}

	return lab, nil
}

func PixelToRGB(in color.Color) *Rgb {
	r, g, b, _ := in.RGBA()

	out := &Rgb{
		R: uint8(r / 257),
		G: uint8(g / 257),
		B: uint8(b / 257),
	}

	return out
}

