package clustering

import (
	"fmt"
	"math/rand"
)

func SelectCenters(k int, distances [][]float64) ([]int, error) {
	rows := len(distances)
	if rows < k {
		return []int{}, fmt.Errorf("at least %d points are needed to calculate %d cluster centers", k, k)
	}

	for _, row := range distances {
		if len(row) != rows {
			return []int{}, fmt.Errorf("wrong distance matrix format: the matrix should be square")
		}
	}

	centroids := []int{}
	squareDistances := make([]float64, rows)

	// Select first center
	centroids = append(centroids, rand.Intn(rows))

	// Select the rest of the centers
	for len(centroids) < k {
		// Compute square distance to closest selected center
		for i := range rows {
			closest := centroids[0]
			for j := 1; j < len(centroids); j++ {
				centerPoint := centroids[j]
				if distances[i][centerPoint] < distances[i][closest] {
					closest = centerPoint
				}
			}
			squareDistances[i] = distances[i][closest] * distances[i][closest]
		}

		// Choose new center with squareDistances as a weight
		newCenter := weightedSelection(squareDistances)

		centroids = append(centroids, newCenter)
	}

	return centroids, nil
}

func weightedSelection(weights []float64) int {
	chosenIndex := 0
	var total, threshold, cumsum float64

	for _, distance := range weights {
		total += distance
	}
	threshold = rand.Float64() * total
	cumsum = 0

	for i, d := range weights {
		cumsum += d
		if cumsum >= threshold {
			chosenIndex = i
			break
		}
	}

	return chosenIndex
}
