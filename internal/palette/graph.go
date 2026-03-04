package palette

type Graph struct {
	adjList map[int][]int
}

func NewGraph(n int) *Graph {
	g := &Graph{
		adjList: make(map[int][]int, n),
	}

	for i := range n {
		g.adjList[i] = make([]int, 0)
	}

	return g
}

func (g *Graph) AddVertex(v int) {
	if _, exists := g.adjList[v]; !exists {
		g.adjList[v] = make([]int, 0)
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adjList[u] = append(g.adjList[u], v)
	g.adjList[v] = append(g.adjList[v], u)
}

func removeFromSlice(s []int, n int) []int {
	for i, v := range s {
		if v == n {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (g *Graph) RemoveEdge(u, v int) {
	g.adjList[u] = removeFromSlice(g.adjList[u], v)
	g.adjList[v] = removeFromSlice(g.adjList[v], u)
}

func (g *Graph) RemoveVertex(v int) {
	delete(g.adjList, v)

	for u := range g.adjList {
		g.adjList[u] = removeFromSlice(g.adjList[u], v)
	}
}

func (g *Graph) Empty() bool {
	return len(g.adjList) <= 0
}

func (g *Graph) Clone() *Graph {
	clone := &Graph{
		adjList: make(map[int][]int, len(g.adjList)),
	}

	for v, edges := range g.adjList {
		newSlice := make([]int, len(edges))
		copy(newSlice, edges)
		clone.adjList[v] = newSlice
	}

	return clone
}

func (g *Graph) RemoveNeighbors(v int) {
	neighbors := g.adjList[v]
	for _, n := range neighbors {
		delete(g.adjList, n)
	}
	delete(g.adjList, v)
}

func (g *Graph) Size() (v, e int) {
	e2 := 0
	for i := range g.adjList {
		e2 += len(g.adjList[i])
	}
	return len(g.adjList), e2 / 2
}
