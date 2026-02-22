package evaluation

import (
	"github.com/mermonia/chromatika/internal/colors"
)

func ChoosePrimaryColor(cols []*colors.Lab, partMatrix [][]float64) *colors.Lab {
	chosenColor := chooseFromChroma(
		cols,
		partMatrix,
		func(a, b float64) bool {
			return a > b
		})

	return chosenColor
}

func ChooseBackgroundColor(cols []*colors.Lab, partMatrix [][]float64) *colors.Lab {
	chosenColor := chooseFromChroma(
		cols,
		partMatrix,
		func(a, b float64) bool {
			return a < b
		})

	return chosenColor
}

func chooseFromChroma(cols []*colors.Lab, partMatrix [][]float64, comp func(a, b float64) bool) *colors.Lab {
	var chosenColor *colors.Lab
	var bestEval float64

	for i, color := range cols {
		chroma := color.GetChroma()
		eval := chroma

		if i == 0 {
			bestEval = eval
			chosenColor = color
		}

		if comp(eval, bestEval) {
			bestEval = eval
			chosenColor = color
		}
	}

	return chosenColor
}
