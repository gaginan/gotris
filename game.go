// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

import (
	"context"
	"sync"
	"time"
)

// Game represents the block-stacking game.
type Game interface {
	// Run starts the game loop (gravity + input processing) until ctx is done.
	Run(ctx context.Context)
}

type game struct {
	mu        sync.RWMutex
	gameBoard GameBoard
	level     int
	render    Renderer
	controls  <-chan Control
	lines     int
	steps     time.Duration
}

// New creates a new block-stacking game instance.
// Provide a Renderer implementation and a channel of Control values.
// The game starts with a 20x10 board and level 1.
func New(ctx context.Context, r Renderer, ctrls <-chan Control) Game {
	return &game{
		gameBoard: NewGameBoard(20, 10),
		level:     1,
		render:    r,
		controls:  ctrls,
		steps:     time.Second / 20,
	}
}

func (g *game) Run(ctx context.Context) {
	var t = time.NewTicker(time.Second)
	defer t.Stop()
	var piece = g.gameBoard.Next()
	var gameOver = false
	g.render.Clear()
	for {
		select {
		case <-ctx.Done():
			return
		case ctrl := <-g.controls:
			g.gameBoard.Apply(piece, ctrl)
			g.update(*piece)
		case <-t.C:
			if !g.gameBoard.Apply(piece, Move(Down)) {
				g.gameBoard.Stack(piece.Grid, piece.Location)
				g.lines += g.gameBoard.Compact()
				piece = g.gameBoard.Next()
				gameOver = g.gameBoard.IsTopOut(*piece)
				if g.levelUp() > 0 {
					t.Reset(g.levelUp())
				}
			}
			g.update(*piece)
		}
		if gameOver {
			g.render.Clear()
			break
		}
	}
}

func (g *game) levelUp() time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.lines/10 >= g.level {
		g.level++
	}
	n := time.Second - time.Millisecond*time.Duration(g.level*20)
	if n > 0 {
		return n
	}
	return 0
}

func (g *game) update(piece Piece) {
	g.mu.Lock()
	defer g.mu.Unlock()
	board := g.gameBoard.State()
	g.render.Update(GameState{
		Board:   board,
		Current: piece,
		Next:    g.gameBoard.Preview(),
		Lines:   g.lines,
		Level:   g.level,
	})
}
