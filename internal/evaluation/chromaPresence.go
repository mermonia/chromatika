package evaluation

import "github.com/mermonia/chromatika/internal/colors"

func ChoosePrimaryColor(cols []*colors.Lab, partMatrix [][]float64) *colors.Lab {
	chosenColor := chooseFromColorPresence(
		cols,
		partMatrix,
		func(a, b float64) bool {
			return a > b
		})

	return chosenColor
}

func ChooseBackgroundColor(cols []*colors.Lab, partMatrix [][]float64) *colors.Lab {
	chosenColor := chooseFromColorPresence(
		cols,
		partMatrix,
		func(a, b float64) bool {
			return a < b
		})

	return chosenColor
}

func chooseFromColorPresence(cols []*colors.Lab, partMatrix [][]float64, comp func(a, b float64) bool) *colors.Lab {
	var chosenColor *colors.Lab
	var bestEval float64

	for i, color := range cols {
		presence := getColorPresence(i, partMatrix)
		chroma := color.GetChroma()
		eval := presence * chroma

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

func getColorPresence(idx int, partMatrix [][]float64) float64 {
	sum := 0.0
	for _, color := range partMatrix {
		sum += color[idx]
	}
	return sum
}
