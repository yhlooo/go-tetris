package web

import (
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
	return app.Div().Body(ui.renderMain()).Styles(map[string]string{
		"display":         "flex",
		"justify-content": "center",
	})
}

// OnMount 挂载元素时
func (ui *GameUI) OnMount(ctx app.Context) {
	app.Log("tetris component mount")
	holdAndNextColors := DefaultBlockColors
	holdAndNextColors.Background = holdAndNextColors.Border
	ui.hold = NewTetrisGrid(2, 3, holdAndNextColors)
	ui.field = NewTetrisGrid(20, 10, DefaultBlockColors)
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
