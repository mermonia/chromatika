package clustering

import (
	"fmt"
	"math"

	"github.com/mermonia/chromatika/internal/colors"
)

type FCMParameters struct {
	M, E float64
	B, K int
}

func FCM(cols []*colors.Lab, np []int, params FCMParameters) ([]*colors.Lab, [][]float64, error) {
	// Pre-calculation of color distances
	distanceMat := colors.DistanceMatrix(cols)

	// Inital center selection
	initialCenters, err := SelectInitialCenters(params.K, distanceMat)
	if err != nil {
		return nil, nil, fmt.Errorf("could not calculate initial cluster centers: %w", err)
	}

	centers := make([]*colors.Lab, params.K)
	for i, cidx := range initialCenters {
		centers[i] = cols[cidx]
	}

	// Iteration counter and difference metric init
	var iters int = 0
	var difference float64 = math.MaxFloat64

	// Partition matrix
	var partitionMat [][]float64

	// Main algorithm
	for iters <= params.B && difference >= params.E {
		partitionMat = calculatePartitionMatrix(cols, centers, params.M)
		newCenters := calculateCenters(cols, partitionMat, np, params.M)

		difference = calculateDifference(newCenters, centers)
		centers = newCenters
		iters++
	}

	return centers, partitionMat, nil
}

func calculateDifference(a, b []*colors.Lab) float64 {
	var maxDiff float64 = 0.0

	for i := range len(a) {
		dist := colors.DistanceLab(a[i], b[i])
		if dist > maxDiff {
			maxDiff = dist
		}
	}

	return maxDiff
}

func calculatePartitionMatrix(cols []*colors.Lab, centers []*colors.Lab, m float64) [][]float64 {
	ncols := len(cols)
	k := len(centers)

	U := make([][]float64, ncols)
	for i := range U {
		U[i] = make([]float64, k)
	}

	exponent := -2.0 / (m - 1)

	for i, color := range cols {
		// Check if the current sample is the center of any clusters
		ownClusters := []int{}
		for j, center := range centers {
			if center.Equals(color) {
				ownClusters = append(ownClusters, j)
			}
		}

		// If the current sample is the center of any clusters, it's a member of those clusters only
		if len(ownClusters) > 0 {
			for _, c := range ownClusters {
				U[i][c] = 1.0 / float64(len(ownClusters))
			}
			continue
		}

		// If the current sample is not a center, calculate its membership to all clusters
		denom := 0.0
		weights := make([]float64, k)

		for j, center := range centers {
			dist := colors.DistanceLab(color, center)
			weig := math.Pow(dist, exponent)
			weights[j] = weig
			denom += weig
		}

		for j := range centers {
			U[i][j] = weights[j] / denom
		}
	}

	return U
}

func calculateCenters(cols []*colors.Lab, U [][]float64, np []int, m float64) []*colors.Lab {
	k := len(U[0])
	centers := make([]*colors.Lab, k)

	// for each cluster, find its center
	for c := range k {
		var sumL, sumA, sumB float32
		var sumWeights float64

		for i, color := range cols {
			weight := math.Pow(U[i][c], m) * float64(np[i])
			sumWeights += weight
			sumL += color.L * float32(weight)
			sumA += color.A * float32(weight)
			sumB += color.B * float32(weight)
		}

		centers[c] = &colors.Lab{
			L: sumL / float32(sumWeights),
			A: sumA / float32(sumWeights),
			B: sumB / float32(sumWeights),
		}
	}

	return centers
}
