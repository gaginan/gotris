// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

import (
	"sync"
)

var _ Tetromino = (*tetromino)(nil)

// Tetromino represents a falling block piece and exposes rotation and shape retrieval.
type Tetromino interface {
	// Rotate applies one or more Rotation transforms to the tetromino.
	Rotate(fn ...Rotation)
	// Grid returns a cloned grid of the tetromino, optionally transformed.
	Grid(fn ...Rotation) Grid
}

// Rotation transforms a Grid (typically a tetromino's matrix) and returns a new Grid.
type Rotation func(g Grid) Grid

func newTetromino(grid Grid, off Location) *tetromino {
	rows, cols := grid.Size()
	height := rows + off.Y
	width := cols + off.X
	side := max(height, width)
	result := newGrid(side, side)

	for i := range grid {
		for j := range grid[i] {
			if result.Contains(Location{X: off.X + j, Y: off.Y + i}) {
				result[off.Y+i][off.X+j] = grid[i][j]
			}
		}
	}
	return &tetromino{grid: result}
}

type tetromino struct {
	mu   sync.RWMutex
	grid Grid
}

func (t *tetromino) Grid(fn ...Rotation) Grid {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := t.grid.Clone()
	for _, f := range fn {
		result = f(result)
	}
	return result
}

func (t *tetromino) Rotate(fn ...Rotation) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, f := range fn {
		t.grid = f(t.grid)
	}
}
