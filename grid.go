// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

// State represents the content of a cell in the board grid or piece shape.
// Values other than Empty are considered filled and may carry color semantics.
type State int

const (
	// Empty indicates an unoccupied cell.
	Empty State = iota
	// Red is used by the Z tetromino.
	Red
	// Orange is used by the L tetromino.
	Orange
	// Yellow is used by the O tetromino.
	Yellow
	// Green is used by the S tetromino.
	Green
	// Cyan is used by the I tetromino.
	Cyan
	// Blue is used by the J tetromino.
	Blue
	// Purple is used by the T tetromino.
	Purple
	// Gray may be used for UI or garbage lines.
	Gray
	// Solid is a sentinel value for walls/borders when needed.
	Solid State = 0xFFFFFFFF
)

// Grid is a 2D matrix of cell states representing a game board or tetromino shape.
type Grid [][]State

// Size returns the number of rows and columns in the grid.
func (g Grid) Size() (rows, cols int) {
	if len(g) == 0 {
		return 0, 0
	}
	return len(g), len(g[0])
}

// Contains reports whether the given location is within the bounds of the grid.
func (g Grid) Contains(l Location) bool {
	return len(g) > 0 &&
		l.Y >= 0 && l.Y < len(g) &&
		l.X >= 0 && l.X < len(g[0])
}

// Walk iterates over the grid, calling fn(row, col, state) for each cell.
func (g Grid) Walk(fn func(row, col int, state State)) {
	for row := range g {
		for col := range g[row] {
			fn(row, col, g[row][col])
		}
	}
}

// Zero returns a new grid with the same dimensions but all cells set to Empty.
func (g Grid) Zero() Grid {
	if len(g) == 0 {
		return Grid{}
	}
	rows := len(g)
	cols := len(g[0])
	return NewGrid(rows, cols)
}

// Clone creates a deep copy of the grid.
func (g Grid) Clone() Grid {
	result := g.Zero()
	g.Walk(func(row, col int, state State) {
		result[row][col] = state
	})
	return result
}

// Fill sets all cells in the grid to the specified state.
func (g Grid) Fill(state State) {
	g.Walk(func(row, col int, s State) {
		g[row][col] = state
	})
}

// Combine overlays another grid onto a clone of the current grid at the specified offset.
// Only non-empty cells from the overlay are applied; out-of-bounds cells are ignored.
func (g Grid) Combine(overlay Grid, off Location) Grid {
	cells := g.Clone()
	overlay.Walk(func(row, col int, state State) {
		location := Location{X: col + off.X, Y: row + off.Y}
		if state != Empty && cells.Contains(location) {
			cells[location.Y][location.X] = state
		}
	})
	return cells
}

// Union merges two grids into a new grid sized to contain both entirely.
// The other grid is placed at the specified offset location.
func (g Grid) Union(other Grid, off Location) Grid {
	if len(g) == 0 {
		return other.Clone()
	}
	if len(other) == 0 {
		return g.Clone()
	}
	rows := max(len(g), len(other)+off.Y)
	cols := max(len(g[0]), len(other[0])+off.X)
	result := NewGrid(rows, cols)
	g.Walk(func(row, col int, state State) {
		if state != Empty {
			result[row][col] = state
		}
	})
	other.Walk(func(row, col int, state State) {
		if state != Empty {
			result[row+off.Y][col+off.X] = state
		}
	})
	return result
}

// RotateRight returns a new grid rotated 90 degrees clockwise.
func (g Grid) RotateRight() Grid { return RotateRight(g) }

// RotateLeft returns a new grid rotated 90 degrees counter-clockwise.
func (g Grid) RotateLeft() Grid { return RotateLeft(g) }

// RotateRight rotates the grid 90 degrees clockwise and returns a new grid.
func RotateRight(g Grid) Grid {
	if len(g) == 0 {
		return Grid{}
	}
	rows, cols := g.Size()
	result := NewGrid(cols, rows) // Swap dimensions for 90-degree rotation
	g.Walk(func(row, col int, state State) {
		result[col][rows-1-row] = state
	})
	return result
}

// RotateLeft rotates the grid 90 degrees counter-clockwise and returns a new grid.
func RotateLeft(g Grid) Grid {
	if len(g) == 0 {
		return Grid{}
	}
	rows, cols := g.Size()
	result := NewGrid(cols, rows) // Swap dimensions for 90-degree rotation
	g.Walk(func(row, col int, state State) {
		result[cols-1-col][row] = state
	})
	return result
}

// NewGrid creates a new grid with the specified dimensions, initialized with Empty cells.
func NewGrid(rows, cols int) Grid {
	result := make(Grid, rows)
	for i := range result {
		result[i] = make([]State, cols)
	}
	return result
}
