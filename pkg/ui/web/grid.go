package web

import (
	"sync"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris"
)

// NewTetrisGrid 创建 TetrisGrid
func NewTetrisGrid(rows, cols int, colors BlockColors) *TetrisGrid {
	data := make([][]tetris.BlockType, rows)
	for i := range data {
		data[i] = make([]tetris.BlockType, cols)
	}
	grid := &TetrisGrid{
		cellWidth:   20,
		borderWidth: 2,
		colors:      colors,
		data:        data,
		canvas:      app.Canvas(),
	}
	grid.resize(rows, cols)
	return grid
}

// TetrisGrid Tetris 网格
type TetrisGrid struct {
	app.Compo
	lock sync.Mutex

	cellWidth   int
	borderWidth int
	colors      BlockColors

	data [][]tetris.BlockType

	rows, cols    int
	width, height int
	canvas        app.HTMLCanvas
}

var _ app.Composer = (*TetrisGrid)(nil)

// Render 渲染画面
func (grid *TetrisGrid) Render() app.UI {
	return grid.canvas.Width(grid.width).Height(grid.height)
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

// Size 获取当前大小
func (grid *TetrisGrid) Size() (width, height int) {
	return grid.width, grid.height
}

// resize 调整大小
func (grid *TetrisGrid) resize(rows, cols int) {
	grid.rows = rows
	grid.cols = cols
	grid.width = cols*grid.cellWidth + (cols-1)*grid.borderWidth
	grid.height = rows*grid.cellWidth + (rows-1)*grid.borderWidth
	if canvas := grid.canvas.JSValue(); canvas != nil {
		canvas.Set("width", grid.width)
		canvas.Set("height", grid.height)
		grid.paintBorder()
	}
}

// paintBlocks 绘制方块
func (grid *TetrisGrid) paintBlocks() {
	if grid.canvas.JSValue() == nil {
		return
	}
	canvasCTX := grid.canvas.JSValue().Call("getContext", "2d")
	currentColor := ""
	for i, row := range grid.data {
		for j, block := range row {
			color := grid.colors.Block(block)
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
	if grid.canvas.JSValue() == nil {
		return
	}
	canvasCTX := grid.canvas.JSValue().Call("getContext", "2d")
	canvasCTX.Set("strokeStyle", grid.colors.Border)
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

// BlockColors 方块颜色
type BlockColors struct {
	Border     string
	Background string
	BlockI     string
	BlockJ     string
	BlockL     string
	BlockO     string
	BlockS     string
	BlockT     string
	BlockZ     string
}

// Block 获取指定方块颜色
func (colors BlockColors) Block(blockType tetris.BlockType) string {
	switch blockType {
	case tetris.BlockNone:
		return colors.Background
	case tetris.BlockI:
		return colors.BlockI
	case tetris.BlockJ:
		return colors.BlockJ
	case tetris.BlockL:
		return colors.BlockL
	case tetris.BlockO:
		return colors.BlockO
	case tetris.BlockS:
		return colors.BlockS
	case tetris.BlockT:
		return colors.BlockT
	case tetris.BlockZ:
		return colors.BlockZ
	}
	return colors.Background
}

// DefaultBlockColors 默认颜色
var DefaultBlockColors = BlockColors{
	Border:     "#1b1b1b",
	Background: "#000000",
	BlockI:     "#67c4ec",
	BlockJ:     "#5f64a9",
	BlockL:     "#df8136",
	BlockO:     "#f0d543",
	BlockS:     "#62b451",
	BlockT:     "#a25399",
	BlockZ:     "#db3e32",
}
