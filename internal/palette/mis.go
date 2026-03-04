package palette

func MIS_Complete(g *Graph) []int {
	if g.Empty() {
		return []int{}
	}

	// if there is a degree 0/1 vertex, add it to the solution
	for v, edges := range g.adjList {
		if len(edges) <= 1 {
			v_inc := g.Clone()
			v_inc.RemoveNeighbors(v)
			return append(MIS_Complete(v_inc), v)
		}
	}

	// if there is a degree >= 3 vertex, choose between the solution
	// with the vertex in it and the solution without the vertex in it
	for v, edges := range g.adjList {
		if len(edges) >= 3 {
			v_inc := g.Clone()
			v_inc.RemoveNeighbors(v)
			v_inc_result := append(MIS_Complete(v_inc), v)

			v_exc := g.Clone()
			v_exc.RemoveVertex(v)
			v_exc_result := MIS_Complete(v_exc)

			if len(v_inc_result) > len(v_exc_result) {
				return v_inc_result
			}
			return v_exc_result
		}
	}

	// if there are only vertexes with degree 2, we can directly compute
	// the remaining additions to MIS
	return MIS_Cycles(g)
}

func MIS_Cycles(g *Graph) []int {
	visited := make(map[int]bool)
	result := []int{}

	for start := range g.adjList {
		if visited[start] {
			continue
		}
		cycle := []int{}
		prev, curr := -1, start
		for !visited[curr] {
			visited[curr] = true
			cycle = append(cycle, curr)
			for _, neighbor := range g.adjList[curr] {
				if neighbor != prev {
					prev, curr = curr, neighbor
					break
				}
			}
		}
		for i := 0; i < len(cycle); i += 2 {
			result = append(result, cycle[i])
		}
	}
	return result
}
