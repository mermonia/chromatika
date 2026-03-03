package palette

import (
	"fmt"

	"github.com/mermonia/chromatika/internal/clustering"
	"github.com/mermonia/chromatika/internal/colors"
	"github.com/mermonia/chromatika/internal/extraction"
	"github.com/mermonia/chromatika/internal/utils"
)

var defaultScaleWidth int = 512
var defaultQuantInterval int = 4
var defaultFuzziness float64 = 2.0
var defaultThreshold float64 = 0.001
var defaultMaxIterations int = 100

func GeneratePalette(imagePath string, darkmode bool) (*Palette, error) {
	/*
	Extract up to 32 colors from an image, check if there are enough
	suitable colors to form a palette in each iteration.
	*/
	for i := range 16 {
		labCols, _, err := extraction.GetDominantColors(
			imagePath,
			defaultScaleWidth,
			defaultQuantInterval,
			clustering.FCMParameters{
				M: defaultFuzziness,
				E: defaultThreshold,
				B: defaultMaxIterations,
				K: i + 16,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("could not generate a %d color palette: %w", i+16, err)
		}

		lchCols := utils.Map(labCols, colors.LabToLCH)
		suitableColors, err := filterSuitableColors(lchCols)

		if len(suitableColors) >= 16 {
			break
		}
	}

	return nil, nil
}

/*
Returns a set of colors that are distinct enough to form a palette.
*/
func filterSuitableColors(colors []*colors.LCHab) ([]*colors.LCHab, error) {
	return nil, nil
}
