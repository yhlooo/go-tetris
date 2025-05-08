package web

import (
	"fmt"
	"strconv"

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

	handleKeyDown app.Func

	field      *TetrisGrid
	hold       *TetrisGrid
	next       [3]*TetrisGrid
	score      int
	level      int
	clearLines int

	page string

	tetris tetris.Tetris
}

var _ app.Composer = (*GameUI)(nil)

// Render 渲染画面
func (ui *GameUI) Render() app.UI {
	fieldWidth := 100
	fieldHeight := 100
	if ui.field != nil {
		fieldWidth, fieldHeight = ui.field.Size()
	}
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
			).Styles(map[string]string{
				"padding": "10px",
				"width":   "120px",
			}),
			app.Div().Body(
				app.If(ui.page == "", func() app.UI {
					return app.Div().Body(
						app.Button().Text("Start").OnClick(func(ctx app.Context, _ app.Event) { ui.toGame(ctx) }),
					)
				}).ElseIf(ui.page == "paused", func() app.UI {
					return app.Div().Body(
						app.Button().Text("Resume").OnClick(func(ctx app.Context, _ app.Event) { ui.toGame(ctx) }),
						app.Button().Text("Quit").OnClick(func(ctx app.Context, _ app.Event) { ui.toStartMenu(ctx) }),
					)
				}).ElseIf(ui.page == "over", func() app.UI {
					return app.Div().Body(
						app.Div().Text("Game Over"),
						app.P().Text(fmt.Sprintf("Score: %d", ui.score)),
						app.Button().Text("Ok").OnClick(func(ctx app.Context, _ app.Event) { ui.toStartMenu(ctx) }),
					)
				}).Else(func() app.UI {
					return ui.field
				}),
			).Styles(map[string]string{
				"width":  strconv.Itoa(fieldWidth) + "px",
				"height": strconv.Itoa(fieldHeight) + "px",
			}),
			app.Div().Body(
				app.P().Text("Next:"),
				app.Range(ui.next).Slice(func(i int) app.UI {
					return app.Div().Body(ui.next[i])
				}),
			).Style("padding", "10px"),
		).Styles(map[string]string{
			"display": "flex",
		}),
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
	)
}

// OnMount 挂载元素时
func (ui *GameUI) OnMount(ctx app.Context) {
	app.Log("tetris component mount")
	ui.hold = NewTetrisGrid(2, 3)
	ui.field = NewTetrisGrid(20, 10)
	ui.next[0] = NewTetrisGrid(2, 3)
	ui.next[1] = NewTetrisGrid(2, 3)
	ui.next[2] = NewTetrisGrid(2, 3)

	ui.handleKeyDown = app.FuncOf(func(this app.Value, args []app.Value) any {
		ui.handleInput(ctx, args[0])
		return nil
	})
	app.Window().Call("addEventListener", "keydown", ui.handleKeyDown)
}

// OnDismount 卸载元素时
func (ui *GameUI) OnDismount() {
	app.Log("tetris component dismount")
	app.Window().Call("removeEventListener", "keydown", ui.handleKeyDown)
}

// handleInput 处理用户输入事件
func (ui *GameUI) handleInput(ctx app.Context, e app.Value) {
	if ui.tetris == nil {
		return
	}

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
			ui.toPaused(ctx)
		} else {
			ui.toGame(ctx)
		}
	}

	ctx.Update()
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

	if frame.GameOver {
		ui.toGameOver(ctx)
	}

	ctx.Update()
}

// toStartMenu 回到开始菜单
func (ui *GameUI) toStartMenu(_ app.Context) {
	ui.page = ""
	if ui.tetris != nil {
		if err := ui.tetris.Stop(); err != nil {
			app.Logf("stop tetris error: %v", err)
		}
		ui.tetris = nil
	}
}

// toGame 开始或回到游戏
func (ui *GameUI) toGame(ctx app.Context) {
	if ui.tetris == nil {
		ui.tetris = tetris.NewTetris(tetris.DefaultOptions)
		go ui.paintFrameLoop(ctx, ui.tetris.Frames())
		if err := ui.tetris.Start(ctx); err != nil {
			app.Logf("start tetris error: %v", err)
			return
		}
	}
	if ui.tetris.State() == tetris.StatePaused {
		if err := ui.tetris.Resume(); err != nil {
			app.Logf("resume tetris error: %v", err)
			return
		}
	}

	ui.page = "game"
}

// toGameOver 游戏结束
func (ui *GameUI) toGameOver(_ app.Context) {
	ui.page = "over"
}

// toPaused 暂停
func (ui *GameUI) toPaused(_ app.Context) {
	if ui.tetris != nil {
		if err := ui.tetris.Pause(); err != nil {
			app.Logf("pause tetris error: %v", err)
			return
		}
	}
	ui.page = "paused"
}
