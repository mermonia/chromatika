package palette

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"time"
)

func TestMIS_Complete(t *testing.T) {
	g := randomGraph(16, 5, time.Now().UnixNano())
	fmt.Println(MIS_Complete(g))
}

func randomGraph(v, maxDeg int, seed int64) *Graph {
	rng := rand.New(rand.NewSource(seed))
	g := NewGraph(v)
	for u := range v {
		for _, candidate := range rng.Perm(v) {
			if candidate == u {
				continue
			}
			if len(g.adjList[u]) >= maxDeg {
				break
			}
			if len(g.adjList[candidate]) >= maxDeg {
				continue
			}
			alreadyConnected := slices.Contains(g.adjList[u], candidate)
			if !alreadyConnected {
				g.AddEdge(u, candidate)
			}
		}
	}
	return g
}
