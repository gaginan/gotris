# gotris


[English](./README.md) | 简体中文

一个用 Go 实现的、简洁的类俄罗斯方块代码。通过添加渲染和响应用户输入的代码，可以实现一个可运行的游戏。初学者可将其作为一个原型来了解和学习编程和软件设计。

## 快速开始

实现一个自己的 Renderer（实现 `Renderer` 接口）的最小示例：

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

    // 发送控制指令，然后在稍后退出。
    ctrls <- gotris.Move(gotris.Left)
    ctrls <- gotris.Rotate(gotris.RotateRight)
    time.AfterFunc(100*time.Millisecond, cancel)

    <-ctx.Done()
}
```

## 开发命令

- 构建：`go build`
- 测试：`go test ./...`

## 设计理念

- 通过接口实现清晰的职责分离：
  - Board：网格存储、碰撞检测、堆叠与行消除（压缩）
  - GameBoard：在 Board 之上加入“下一个方块”队列与控制指令的应用
  - Game：主循环（重力 + 输入）
  - Renderer：只读的渲染器，消费 GameState 快照
- `Grid` 提供旋转、覆盖（overlay）、遍历与克隆等工具方法。
- 方块形状由函数返回 `Grid`，并在创建时填充为正方，确保旋转过程维度稳定。
