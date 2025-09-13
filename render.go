// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

// GameState contains the information required by a Renderer to draw the game.
type GameState struct {
	Board   Grid
	Current Piece
	Next    []Grid
	Lines   int
	Level   int
}

// Renderer renders snapshots of the block-stacking game state.
type Renderer interface {
	// Update receives a snapshot of the current game state.
	// Implementations should be fast and non-blocking; avoid heavy work.
	Update(state GameState)
	// Clear clears the rendering surface.
	Clear()
}
