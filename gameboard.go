// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

import (
	"math/rand"
	"sync"
)

var _ GameBoard = (*gameBoard)(nil)

// GameBoard is the game-level controller that manages the falling-piece queue
// and applies operations to the underlying Board grid.
type GameBoard interface {
	Board
	// Apply attempts one or more Controls on the provided piece; returns true only if all succeed.
	// Controls are applied sequentially and may mutate the piece in-place.
	Apply(piece *Piece, fn ...Control) (ok bool)
	// Next consumes and returns the next piece from the queue (refilling as needed).
	Next() *Piece
	// Preview returns the upcoming piece shapes without consuming the queue.
	Preview() []Grid
	// IsTopOut returns true if the provided piece overlaps occupied cells when spawned.
	IsTopOut(piece Piece) bool
}

// Piece represents a falling tetromino instance with a shape and a top-left
// location in the board coordinate system.
type Piece struct {
	Grid     Grid
	Location Location
}

// Control is an operation applied to a piece in the context of the underlying Board.
// It returns true if the operation is valid and was applied.
// Typical controls are Move(Left/Right/Down) and Rotate(RotateLeft/RotateRight).
type Control func(board Board, p *Piece) (ok bool)

// HardDrop moves the piece down until it collides by repeatedly applying Move(Down).
func HardDrop() Control {
	return func(board Board, p *Piece) (ok bool) {
		for {
			if !Move(Down)(board, p) {
				return true
			}
		}
	}
}

// Move attempts to move the piece using the provided Direction and validates
// the new location against the underlying Board using Test.
func Move(fn Direction) Control {
	return func(board Board, p *Piece) (ok bool) {
		var l = p.Location
		fn(&l)
		if ok := board.Test(p.Grid, l); ok {
			p.Location = l
			return true
		}
		return false
	}
}

// Rotate attempts to rotate the piece using the provided Rotation and validates
// the rotated shape against the Board using Test.
func Rotate(fn Rotation) Control {
	return func(board Board, p *Piece) (ok bool) {
		var rotated = fn(p.Grid)
		if ok := board.Test(rotated, p.Location); ok {
			p.Grid = rotated
			return true
		}
		return false
	}
}

// NewGameBoard constructs a new GameBoard (controller) with the given dimensions.
func NewGameBoard(rows, cols int) GameBoard {
	var g = &gameBoard{
		Board: NewBoard(rows, cols),
		rows:  rows,
		cols:  cols,
	}
	return g
}

// gameBoard is the default implementation of GameBoard.
type gameBoard struct {
	mu sync.RWMutex
	Board
	rows, cols int
	upcoming   []Tetromino
}

func (b *gameBoard) Apply(piece *Piece, fn ...Control) (ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, f := range fn {
		if !f(b.Board, piece) {
			return false
		}
	}
	return true
}

func (b *gameBoard) Next() *Piece {
	const minQueueSize = 3
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.upcoming) <= minQueueSize {
		var seq = []Tetromino{NewI(), NewJ(), NewL(), NewO(), NewS(), NewT(), NewZ()}
		rand.Shuffle(len(seq), func(i, j int) { seq[i], seq[j] = seq[j], seq[i] })
		b.upcoming = append(b.upcoming, seq...)
	}
	if len(b.upcoming) == 0 {
		return nil
	}
	head := b.upcoming[0]
	b.upcoming = b.upcoming[1:]
	shape := head.Grid()
	return &Piece{
		Grid:     shape,
		Location: Location{X: (b.cols - len(shape[0])) / 2, Y: 0},
	}
}

// Preview returns a snapshot of upcoming piece shapes without consuming the queue.
func (b *gameBoard) Preview() []Grid {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var shapes []Grid
	for _, tetromino := range b.upcoming {
		shapes = append(shapes, tetromino.Grid())
	}
	return shapes
}

func (b *gameBoard) IsTopOut(piece Piece) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	shape := piece.Grid.Clone()
	var topOut = false
	shape.Walk(func(row, col int, state State) {
		l := Location{X: piece.Location.X + col, Y: piece.Location.Y + row}
		if state != Empty && b.Board.Contains(l) && !b.Board.Empty(l) {
			topOut = true
		}
	})
	return topOut
}
