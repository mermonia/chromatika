package evaluation

import "github.com/mermonia/chromatika/internal/colors"

func ChoosePrimaryColor(cols []*colors.Lab, partMatrix [][]float64) *colors.Lab {
	var chosenColor *colors.Lab
	var highestEval float64

	for i, color := range cols {
		presence := getColorPresence(i, partMatrix)
		chroma := color.GetChroma()

		eval := presence * chroma
		if eval > highestEval {
			highestEval = eval
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
