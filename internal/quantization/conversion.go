package quantization

import (
	"fmt"
	"image/color"

	"github.com/mermonia/chromatika/internal/colors"
)

func GetRawLabArray(pixels []color.Color) ([]*colors.Lab, error) {
	res := make([]*colors.Lab, len(pixels))

	for i, pixel := range pixels {
		lab, err := colors.PixelToLab(pixel)
		if err != nil {
			return nil, fmt.Errorf("could not convert pixel to lab: %w", err)
		}
		res[i] = lab
	}

	return res, nil
}
