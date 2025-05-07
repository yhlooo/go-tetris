package web

import (
	"context"
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris"
)

// NewGameUI 创建 GameUI
func NewGameUI() *GameUI {
	return &GameUI{}
}

// GameUI 基于浏览器运行的 WebAssembly 的游戏用户交互界面
type GameUI struct {
	app.Compo

	field      *TetrisGrid
	hold       *TetrisGrid
	next       [3]*TetrisGrid
	score      int
	level      int
	clearLines int

	tetris tetris.Tetris
}

var _ app.Composer = (*GameUI)(nil)

// Render 渲染画面
func (ui *GameUI) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.Div().Body(
				app.Div().Body(
					app.P().Text("Hold:"),
					ui.hold,
				),
				app.Div().Body(
					app.Div().Text(fmt.Sprintf("Score: %d", ui.score)),
					app.Div().Text(fmt.Sprintf("Level: %d", ui.level)),
					app.Div().Text(fmt.Sprintf("Lines: %d", ui.clearLines)),
				),
			).Style("padding", "10px"),
			app.Div().Body(ui.field),
			app.Div().Body(
				app.P().Text("Next:"),
				app.Range(ui.next).Slice(func(i int) app.UI {
					return app.Div().Body(ui.next[i])
				}),
			).Style("padding", "10px"),
		).Style("display", "flex"),
		app.Div().Body(
			app.P().Body(app.B().Text("Help:")),
			app.P().Body(
				app.Text("Up / w / i : Rotate Right"), app.Br(),
				app.Text("Left / a / j : Move Left"), app.Br(),
				app.Text("Right / d / l : Move Right"), app.Br(),
				app.Text("Down / s / k : Soft Drop"), app.Br(),
				app.Text("z : Rotate Left"), app.Br(),
				app.Text("c : Hold"), app.Br(),
				app.Text("Space : Hard Drop"), app.Br(),
				app.Text("ESC : Pause"),
			),
		),
	).Style("outline", "0").
		TabIndex(0).
		OnKeyDown(ui.handleInput)
}

// OnMount 挂载元素时
func (ui *GameUI) OnMount(ctx app.Context) {
	app.Log("tetris component mount")
	ui.hold = NewTetrisGrid(2, 4)
	ui.field = NewTetrisGrid(20, 10)
	ui.next[0] = NewTetrisGrid(2, 4)
	ui.next[1] = NewTetrisGrid(2, 4)
	ui.next[2] = NewTetrisGrid(2, 4)

	time.Sleep(time.Second)
	opts := tetris.DefaultOptions
	opts.Frequency = 60
	ui.tetris = tetris.NewTetris(opts)
	go ui.paintFrameLoop(ctx, ui.tetris.Frames())
	_ = ui.tetris.Start(context.Background())
}

// handleInput 处理用户输入事件
func (ui *GameUI) handleInput(ctx app.Context, e app.Event) {
	keyCode := e.Get("key").String()
	app.Logf("key down: %q\n", keyCode)
	switch keyCode {
	case "ArrowUp", "w", "i":
		ui.tetris.Input(tetris.OpRotateRight)
	case "ArrowDown", "s", "k":
		ui.tetris.Input(tetris.OpSoftDrop)
	case "ArrowLeft", "a", "j":
		ui.tetris.Input(tetris.OpMoveLeft)
	case "ArrowRight", "d", "l":
		ui.tetris.Input(tetris.OpMoveRight)
	case " ":
		ui.tetris.Input(tetris.OpHardDrop)
	case "z":
		ui.tetris.Input(tetris.OpRotateLeft)
	case "c":
		ui.tetris.Input(tetris.OpHold)
	case "Enter":
		_ = ui.tetris.Resume()
	case "Escape":
		if ui.tetris.State() == tetris.StateRunning {
			_ = ui.tetris.Pause()
		} else {
			_ = ui.tetris.Resume()
		}
	}
}

// paintFrameLoop 绘制游戏帧循环
func (ui *GameUI) paintFrameLoop(ctx app.Context, ch <-chan tetris.Frame) {
	for frame := range ch {
		ui.paintFrame(ctx, frame)
	}
}

// paintFrame 绘制帧
func (ui *GameUI) paintFrame(ctx app.Context, frame tetris.Frame) {
	ui.field.UpdateBlocks(frame.Field.BlocksWithActiveBlock())
	ui.next[0].UpdateBlocks(newBlockGridData(frame.NextBlocks[0]))
	ui.next[1].UpdateBlocks(newBlockGridData(frame.NextBlocks[1]))
	ui.next[2].UpdateBlocks(newBlockGridData(frame.NextBlocks[2]))
	if frame.HoldingBlock != nil {
		ui.hold.UpdateBlocks(newBlockGridData(*frame.HoldingBlock))
	} else {
		ui.hold.UpdateBlocks(newBlockGridData(tetris.BlockNone))
	}

	ui.score = frame.Score
	ui.level = frame.Level
	ui.clearLines = frame.ClearLines

	ctx.Update()
}
