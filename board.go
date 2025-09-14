// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

import (
	"slices"
	"sync"
)

// Board represents the grid where pieces fall. Board stores stacked cells and
// provides collision tests and compaction logic.
type Board interface {
	// Contains reports whether the location is within the bounds of the board.
	Contains(l Location) bool
	// Empty reports whether the given location contains an Empty cell.
	// The caller must ensure the location is within bounds (use Contains),
	// otherwise behavior is undefined.
	Empty(l Location) bool
	// Stack merges non-empty cells of the shape into the board at the given location.
	// Cells outside the board bounds are ignored.
	Stack(grid Grid, l Location)
	// Test returns whether the shape can be placed at the location without collisions or bounds errors.
	Test(grid Grid, l Location) (ok bool)
	// State returns a copy of the board's cells.
	State() Grid
	// Compact removes complete rows and moves rows above downward; returns number of lines cleared.
	Compact() int
	// Clear resets all cells in the board to Empty state.
	Clear()
}

type board struct {
	mu         sync.RWMutex
	rows, cols int
	grid       Grid
}

var _ Board = (*board)(nil)

// NewBoard creates a new Board (grid) with the specified number of rows and columns.
func NewBoard(rows, cols int) Board {
	return &board{
		rows: rows,
		cols: cols,
		grid: NewGrid(rows, cols),
	}
}

func (b *board) Empty(l Location) bool {
	return b.grid[l.Y][l.X] == Empty
}

func (b *board) Contains(l Location) bool {
	return l.Y < b.rows &&
		l.X < b.cols &&
		l.X >= 0 &&
		l.Y >= 0
}

func (b *board) State() Grid {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.grid.Clone()
}

func (b *board) Test(grid Grid, l Location) (ok bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ok = true
	grid.Walk(func(row, col int, state State) {
		location := Location{X: col + l.X, Y: row + l.Y}
		ok = ok && (state == Empty || b.Contains(location) && b.Empty(location))
	})
	return
}

func (b *board) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.grid = NewGrid(b.rows, b.cols)
}

func (b *board) Stack(grid Grid, l Location) {
	b.mu.Lock()
	defer b.mu.Unlock()
	grid.Walk(func(row, col int, state State) {
		location := Location{X: col + l.X, Y: row + l.Y}
		if state != Empty && b.Contains(location) && b.Empty(location) {
			b.grid[location.Y][location.X] = state
		}
	})
}

func (b *board) Compact() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	grid := NewGrid(b.rows, b.cols)
	var line = b.rows - 1
	var rows = 0
	for i := b.rows - 1; i >= 0; i-- {
		if slices.Contains(b.grid[i], Empty) {
			copy(grid[line], b.grid[i])
			line--
			continue
		}
		rows++
	}
	b.grid = grid
	return rows
}
