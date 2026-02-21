package extraction

import (
	"fmt"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/input"
	"github.com/mermonia/chromatika/internal/quantization"
)

func GetDominantColors(path string, scaleW, quantInterval int, paramsFCM clustering.FCMParameters) ([]*colors.Lab, [][]float64, error) {
	pixels, err := input.ReadImage(path, scaleW)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read image: %w", err)
	}

	labColors, err := quantization.GetRawLabArray(pixels)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get lab color array: %w", err)
	}

	quantizedColors, np := quantization.Quantize(quantInterval, labColors)

	extractedColors, partMatrix, err := clustering.FCM(quantizedColors, np, paramsFCM)
	if err != nil {
		return nil, nil, fmt.Errorf("could not execute fcm: %w", err)
	}

	return extractedColors, partMatrix, nil
}
