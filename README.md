# gotris

English | [简体中文](./README.zh-CN.md)

A clean and simple Tetris-Like implementation library written in Go. By adding rendering and user input handling code, you can create a playable game. Beginners can use it as a prototype to understand and learn programming and software design.

## Quick start

Minimal program sketch using your own Renderer (implementing `Renderer`):

```go
package main

import (
    "context"
    "time"

    gotris "github.com/gaginan/gotris"
)

type consoleRenderer struct{}

func (consoleRenderer) Update(s gotris.GameState) {}
func (consoleRenderer) Clear()                   {}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    ctrls := make(chan gotris.Control, 8)
    g := gotris.New(ctx, consoleRenderer{}, ctrls)

    go g.Run(ctx)

    // Drive a few inputs then stop shortly after.
    ctrls <- gotris.Move(gotris.Left)
    ctrls <- gotris.Rotate(gotris.RotateRight)
    time.AfterFunc(100*time.Millisecond, cancel)

    <-ctx.Done()
}
```

## Dev commands

- Build: `go build`
- Test: `go test ./...`

## Design

- Clear separation of concerns via interfaces:
  - Board: grid, collisions, stacking, compaction
  - GameBoard: queue + control application on top of Board
  - Game: loop (gravity + inputs)
  - Renderer: pure observer of GameState snapshots
- `Grid` helpers include rotation, overlays, walking, cloning.
- Shapes are functions returning `Grid` and are padded into square matrices for rotation stability。

## AI Contribution

- **Comments**: Code comments.
- **Documentation**: This README content.
- **Reviews**: Code reviews following the official Go style and best practices.
