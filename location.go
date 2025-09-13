// SPDX-License-Identifier: MIT
// Copyright (c) 2025 gaginan
package gotris

// Direction represents a function that moves a Location.
type Direction func(l *Location)

// Location represents a position in the board grid.
type Location struct {
	X int
	Y int
}

// Left moves the location left.
func Left(l *Location) { l.X-- }

// Right moves the location right.
func Right(l *Location) { l.X++ }

// Down moves the location down.
func Down(l *Location) { l.Y++ }

// Up moves the location up.
func Up(l *Location) { l.Y-- }
