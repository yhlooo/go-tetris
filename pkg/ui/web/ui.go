package web

import (
	"context"
	"fmt"
	"strconv"
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

	Field *TetrisGrid
	Hold  *TetrisGrid
	Next  [3]*TetrisGrid

	tetris tetris.Tetris
}

var _ app.Composer = (*GameUI)(nil)

// Render 渲染画面
func (ui *GameUI) Render() app.UI {
	return app.Div().Body(
		ui.Hold,
		ui.Field,
		app.Range(ui.Next).Slice(func(i int) app.UI {
			return ui.Next[i]
		}),
	).Style("outline", "0").
		Attr("tabindex", "0").
		OnKeyDown(ui.handleInput)
}

// OnMount 挂载元素时
func (ui *GameUI) OnMount(_ app.Context) {
	fmt.Println("component mount")
	ui.Hold = NewTetrisGrid(2, 4)
	ui.Field = NewTetrisGrid(20, 10)
	ui.Next[0] = NewTetrisGrid(2, 4)
	ui.Next[1] = NewTetrisGrid(2, 4)
	ui.Next[2] = NewTetrisGrid(2, 4)

	time.Sleep(time.Second)
	opts := tetris.DefaultOptions
	opts.Frequency = 60
	ui.tetris = tetris.NewTetris(opts)
	go ui.paintFrameLoop(ui.tetris.Frames())
	_ = ui.tetris.Start(context.Background())
}

// handleInput 处理用户输入事件
func (ui *GameUI) handleInput(ctx app.Context, e app.Event) {
	keyCode := e.Get("key").String()
	fmt.Printf("key down: %q\n", keyCode)
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

	fmt.Printf("key donw: %q %s %d\n", e.Get("key"), e.Get("code"), e.Get("keyCode").Int())
}

// paintFrameLoop 绘制游戏帧循环
func (ui *GameUI) paintFrameLoop(ch <-chan tetris.Frame) {
	for frame := range ch {
		ui.paintFrame(frame)
	}
}

// paintFrame 绘制帧
func (ui *GameUI) paintFrame(frame tetris.Frame) {
	ui.Field.UpdateBlocks(frame.Field.BlocksWithActiveBlock())
	ui.Next[0].UpdateBlocks(newBlockGridData(frame.NextBlocks[0]))
	ui.Next[1].UpdateBlocks(newBlockGridData(frame.NextBlocks[1]))
	ui.Next[2].UpdateBlocks(newBlockGridData(frame.NextBlocks[2]))
	if frame.HoldingBlock != nil {
		ui.Hold.UpdateBlocks(newBlockGridData(*frame.HoldingBlock))
	} else {
		ui.Hold.UpdateBlocks(newBlockGridData(tetris.BlockNone))
	}
}

// newBlockGridData 创建方块网格数据
func newBlockGridData(blockType tetris.BlockType) [][]tetris.BlockType {
	switch blockType {
	case tetris.BlockI:
		return [][]tetris.BlockType{
			{tetris.BlockI, tetris.BlockI, tetris.BlockI, tetris.BlockI},
			{tetris.BlockNone, tetris.BlockNone, tetris.BlockNone, tetris.BlockNone},
		}
	case tetris.BlockJ:
		return [][]tetris.BlockType{
			{tetris.BlockJ, tetris.BlockJ, tetris.BlockJ},
			{tetris.BlockJ, tetris.BlockNone, tetris.BlockNone},
		}
	case tetris.BlockL:
		return [][]tetris.BlockType{
			{tetris.BlockL, tetris.BlockL, tetris.BlockL},
			{tetris.BlockNone, tetris.BlockNone, tetris.BlockL},
		}
	case tetris.BlockO:
		return [][]tetris.BlockType{
			{tetris.BlockO, tetris.BlockO},
			{tetris.BlockO, tetris.BlockO},
		}
	case tetris.BlockS:
		return [][]tetris.BlockType{
			{tetris.BlockS, tetris.BlockS, tetris.BlockNone},
			{tetris.BlockNone, tetris.BlockS, tetris.BlockS},
		}
	case tetris.BlockT:
		return [][]tetris.BlockType{
			{tetris.BlockT, tetris.BlockT, tetris.BlockT},
			{tetris.BlockNone, tetris.BlockT, tetris.BlockNone},
		}
	case tetris.BlockZ:
		return [][]tetris.BlockType{
			{tetris.BlockNone, tetris.BlockZ, tetris.BlockZ},
			{tetris.BlockZ, tetris.BlockZ, tetris.BlockNone},
		}
	default:
		return [][]tetris.BlockType{
			{tetris.BlockNone, tetris.BlockNone, tetris.BlockNone},
			{tetris.BlockNone, tetris.BlockNone, tetris.BlockNone},
		}
	}
}

// NewTetrisGrid 创建 TetrisGrid
func NewTetrisGrid(rows, cols int) *TetrisGrid {
	data := make([][]tetris.BlockType, rows)
	for i := range data {
		data[i] = make([]tetris.BlockType, cols)
	}
	cellWidth := 20
	borderWidth := 2
	grid := &TetrisGrid{
		cellWidth:   cellWidth,
		borderWidth: borderWidth,
		borderColor: "#1B1B1B",
		blockColors: map[tetris.BlockType]string{
			tetris.BlockNone: "#000000",
			tetris.BlockI:    "#67C4EC",
			tetris.BlockJ:    "#5F64A9",
			tetris.BlockL:    "#DF8136",
			tetris.BlockO:    "#F0D543",
			tetris.BlockS:    "#62B451",
			tetris.BlockT:    "#A25399",
			tetris.BlockZ:    "#DB3E32",
		},
		data:   data,
		width:  cols*cellWidth + (cols-1)*borderWidth,
		height: rows*cellWidth + (rows-1)*borderWidth,
	}
	grid.canvas = app.Canvas().
		Attr("width", strconv.Itoa(grid.width)+"px").
		Attr("height", strconv.Itoa(grid.height)+"px")
	return grid
}

// TetrisGrid Tetris 网格
type TetrisGrid struct {
	app.Compo

	cellWidth   int
	borderWidth int
	borderColor string
	blockColors map[tetris.BlockType]string

	data [][]tetris.BlockType

	rows, cols    int
	width, height int
	canvas        app.HTMLCanvas
}

var _ app.Composer = (*TetrisGrid)(nil)

// Render 渲染画面
func (grid *TetrisGrid) Render() app.UI {
	return grid.canvas
}

// OnMount 挂载元素时
func (grid *TetrisGrid) OnMount(_ app.Context) {
	grid.paintBorder()
	grid.paintBlocks()
}

// UpdateBlocks 更新方块
func (grid *TetrisGrid) UpdateBlocks(data [][]tetris.BlockType) {
	grid.data = data
	rows := len(data)
	cols := grid.cols
	if rows > 0 {
		cols = len(data[0])
	}
	if rows != grid.rows || cols != grid.cols {
		grid.resize(rows, cols)
	}
	grid.paintBlocks()
}

// resize 调整大小
func (grid *TetrisGrid) resize(rows, cols int) {
	grid.rows = rows
	grid.cols = cols
	grid.width = cols*grid.cellWidth + (cols-1)*grid.borderWidth
	grid.height = rows*grid.cellWidth + (rows-1)*grid.borderWidth
	grid.canvas.
		Attr("width", strconv.Itoa(grid.width)+"px").
		Attr("height", strconv.Itoa(grid.height)+"px")
	grid.paintBorder()
}

// paintBlocks 绘制方块
func (grid *TetrisGrid) paintBlocks() {
	canvasCTX := grid.canvas.JSValue().Call("getContext", "2d")
	currentColor := ""
	for i, row := range grid.data {
		for j, block := range row {
			color, ok := grid.blockColors[block]
			if !ok {
				color = "#000000"
			}
			if color != currentColor {
				canvasCTX.Set("fillStyle", color)
				currentColor = color
			}

			x := j * (grid.cellWidth + grid.borderWidth)
			y := (grid.rows - i - 1) * (grid.cellWidth + grid.borderWidth)
			canvasCTX.Call("fillRect", x, y, grid.cellWidth, grid.cellWidth)
		}
	}
}

// paintBorder 绘制网格边框
func (grid *TetrisGrid) paintBorder() {
	canvasCTX := grid.canvas.JSValue().Call("getContext", "2d")
	canvasCTX.Set("strokeStyle", grid.borderColor)
	canvasCTX.Set("lineWidth", grid.borderWidth)
	for i := 0; i < grid.rows-1; i++ {
		y := (grid.cellWidth+grid.borderWidth)*(i+1) - grid.borderWidth/2
		canvasCTX.Call("moveTo", 0, y)
		canvasCTX.Call("lineTo", grid.width, y)
	}

	for i := 0; i < grid.cols-1; i++ {
		x := (grid.cellWidth+grid.borderWidth)*(i+1) - grid.borderWidth/2
		canvasCTX.Call("moveTo", x, 0)
		canvasCTX.Call("lineTo", x, grid.height)
	}
	canvasCTX.Call("stroke")
}
