// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

// I returns the shape matrix for the I tetromino.
func I(x State) Grid {
	return [][]State{{x, x, x, x}}
}

// J returns the shape matrix for the J tetromino.
func J(x State) Grid {
	return [][]State{
		{x, 0, 0},
		{x, x, x},
	}
}

// L returns the shape matrix for the L tetromino.
func L(x State) Grid {
	return [][]State{
		{0, 0, x},
		{x, x, x},
	}
}

// O returns the shape matrix for the O tetromino.
func O(x State) Grid {
	return [][]State{
		{x, x},
		{x, x},
	}
}

// S returns the shape matrix for the S tetromino.
func S(x State) Grid {
	return [][]State{
		{0, x, x},
		{x, x, 0},
	}
}

// T returns the shape matrix for the T tetromino.
func T(x State) Grid {
	return [][]State{
		{0, x, 0},
		{x, x, x},
	}
}

// Z returns the shape matrix for the Z tetromino.
func Z(x State) Grid {
	return [][]State{
		{x, x, 0},
		{0, x, x},
	}
}

// NewI returns a new I tetromino with its canonical spawn offset.
// Shapes are padded into a square matrix via newTetromino so rotations keep dimensions.
func NewI() Tetromino { return newTetromino(I(Cyan), Location{X: 1, Y: 2}) }

// NewJ returns a new J tetromino with its canonical spawn offset.
func NewJ() Tetromino { return newTetromino(J(Blue), Location{X: 0, Y: 0}) }

// NewL returns a new L tetromino with its canonical spawn offset.
func NewL() Tetromino { return newTetromino(L(Orange), Location{X: 0, Y: 0}) }

// NewO returns a new O tetromino with its canonical spawn offset.
func NewO() Tetromino { return newTetromino(O(Yellow), Location{X: 0, Y: 0}) }

// NewS returns a new S tetromino with its canonical spawn offset.
func NewS() Tetromino { return newTetromino(S(Green), Location{X: 0, Y: 0}) }

// NewT returns a new T tetromino with its canonical spawn offset.
func NewT() Tetromino { return newTetromino(T(Purple), Location{X: 0, Y: 0}) }

// NewZ returns a new Z tetromino with its canonical spawn offset.
func NewZ() Tetromino { return newTetromino(Z(Red), Location{X: 0, Y: 0}) }
