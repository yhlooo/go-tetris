package web

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris"
)

// NewGameUI 创建 GameUI
func NewGameUI() *GameUI {
	return &GameUI{
		touchController: &TouchController{},
	}
}

// GameUI 基于浏览器运行的 WebAssembly 的游戏用户交互界面
type GameUI struct {
	app.Compo

	handleKeyDown   app.Func
	touchController *TouchController

	field      *TetrisGrid
	hold       *TetrisGrid
	next       [3]*TetrisGrid
	score      int
	level      int
	clearLines int

	page      string
	showHelp  bool
	showAbout bool

	tetris tetris.Tetris
}

var _ app.Composer = (*GameUI)(nil)

// Render 渲染画面
func (ui *GameUI) Render() app.UI {
	width := app.Window().Get("innerWidth").Int()
	app.Logf("width: %d", width)
	classes := []string{"tetris-container"}
	if width < 560 {
		classes = append(classes, "tetris-xs")
	}
	return app.Div().Class(classes...).Body(
		ui.renderMain(),
	).
		On("touchstart", ui.touchController.HandleTouchStart).
		On("touchmove", ui.touchController.HandleTouchMove).
		On("touchend", ui.touchController.HandleTouchEnd)
}

// OnMount 挂载元素时
func (ui *GameUI) OnMount(ctx app.Context) {
	app.Log("tetris component mount")
	holdAndNextColors := DefaultTetrominoColors
	holdAndNextColors.Background = holdAndNextColors.Border
	ui.hold = NewTetrisGrid(2, 3, holdAndNextColors)
	ui.field = NewTetrisGrid(20, 10, DefaultTetrominoColors)
	ui.next[0] = NewTetrisGrid(2, 3, holdAndNextColors)
	ui.next[1] = NewTetrisGrid(2, 3, holdAndNextColors)
	ui.next[2] = NewTetrisGrid(2, 3, holdAndNextColors)

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
