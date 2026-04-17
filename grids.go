package gotris

// Hash computes a string representation of the grid's cell states, useful for identifying unique configurations.
func (g Grid) Hash() string {
	hash := make([]byte, 0, len(g)*len(g[0]))
	g.Walk(func(row, col int, state State) {
		hash = append(hash, byte(state))
	})
	return string(hash)
}

// Rotations returns a slice of unique grids representing all distinct rotations of the original grid.
func (g Grid) Rotations() []Grid {
	var rotations []Grid
	var signatures = make(map[string]struct{})
	var candidate = g
	for i := 0; i < 4; i++ {
		if i > 0 {
			candidate = candidate.RotateRight()
		}
		hash := candidate.Hash()
		if _, seen := signatures[hash]; !seen {
			signatures[hash] = struct{}{}
			rotations = append(rotations, candidate)
		}
	}
	return rotations
}

// Where returns all locations in the grid that match the specified state.
func (g Grid) Where(state State) (locations []Location) {
	g.Walk(func(row, col int, s State) {
		if s == state {
			locations = append(locations, Location{X: col, Y: row})
		}
	})
	return
}

// Has checks if the cell at the given location equals the specified state, returning false if the location is out of bounds.
func (g Grid) Has(state State, at Location) bool {
	if g.Contains(at) {
		return g[at.Y][at.X] == state
	}
	return false
}

// Groups finds connected regions of cells that match the specified state,
// returning a slice of locations for each region.
func (g Grid) Groups(state State) (groups [][]Location) {
	var visited locationMap = make(locationMap)
	for _, location := range g.Where(state) {
		if visited.Contains(location) {
			continue
		}
		visited.Add(location)
		queue := []Location{location}
		var group []Location
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			group = append(group, current)
			for _, move := range []Direction{Up, Right, Down, Left} {
				next := current
				move(&next)
				if visited.Contains(next) {
					continue
				}
				if g.Has(state, next) {
					visited.Add(next)
					queue = append(queue, next)
				}
			}
		}
		groups = append(groups, group)
	}
	return
}

type locationMap map[Location]struct{}

func (m locationMap) Contains(loc Location) bool {
	_, exists := m[loc]
	return exists
}

func (m locationMap) Add(loc Location) {
	m[loc] = struct{}{}
}
