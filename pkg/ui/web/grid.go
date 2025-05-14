package web

import (
	"sync"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris/common"
)

// NewTetrisGrid 创建 TetrisGrid
func NewTetrisGrid(rows, cols int, colors TetriminoColors) *TetrisGrid {
	data := make([][]common.TetriminoType, rows)
	for i := range data {
		data[i] = make([]common.TetriminoType, cols)
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
	colors      TetriminoColors

	data [][]common.TetriminoType

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
	grid.paintTetriminos()
}

// UpdateTetriminos 更新方块
func (grid *TetrisGrid) UpdateTetriminos(data [][]common.TetriminoType) {
	grid.data = data
	rows := len(data)
	cols := grid.cols
	if rows > 0 {
		cols = len(data[0])
	}
	if rows != grid.rows || cols != grid.cols {
		grid.resize(rows, cols)
	}
	grid.paintTetriminos()
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

// paintTetriminos 绘制方块
func (grid *TetrisGrid) paintTetriminos() {
	if grid.canvas.JSValue() == nil {
		return
	}
	canvasCTX := grid.canvas.JSValue().Call("getContext", "2d")
	currentColor := ""
	for i, row := range grid.data {
		for j, tetrimino := range row {
			color := grid.colors.Tetrimino(tetrimino)
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

// newTetriminoGridData 创建方块网格数据
func newTetriminoGridData(tetriminoType common.TetriminoType) [][]common.TetriminoType {
	switch tetriminoType {
	case common.I:
		return [][]common.TetriminoType{
			{common.I, common.I, common.I, common.I},
			{common.TetriminoNone, common.TetriminoNone, common.TetriminoNone, common.TetriminoNone},
		}
	case common.J:
		return [][]common.TetriminoType{
			{common.J, common.J, common.J},
			{common.J, common.TetriminoNone, common.TetriminoNone},
		}
	case common.L:
		return [][]common.TetriminoType{
			{common.L, common.L, common.L},
			{common.TetriminoNone, common.TetriminoNone, common.L},
		}
	case common.O:
		return [][]common.TetriminoType{
			{common.O, common.O},
			{common.O, common.O},
		}
	case common.S:
		return [][]common.TetriminoType{
			{common.S, common.S, common.TetriminoNone},
			{common.TetriminoNone, common.S, common.S},
		}
	case common.T:
		return [][]common.TetriminoType{
			{common.T, common.T, common.T},
			{common.TetriminoNone, common.T, common.TetriminoNone},
		}
	case common.Z:
		return [][]common.TetriminoType{
			{common.TetriminoNone, common.Z, common.Z},
			{common.Z, common.Z, common.TetriminoNone},
		}
	default:
		return [][]common.TetriminoType{
			{common.TetriminoNone, common.TetriminoNone, common.TetriminoNone},
			{common.TetriminoNone, common.TetriminoNone, common.TetriminoNone},
		}
	}
}

// TetriminoColors 方块颜色
type TetriminoColors struct {
	Border     string
	Background string
	TetriminoI string
	TetriminoJ string
	TetriminoL string
	TetriminoO string
	TetriminoS string
	TetriminoT string
	TetriminoZ string
}

// Tetrimino 获取指定方块颜色
func (colors TetriminoColors) Tetrimino(tetriminoType common.TetriminoType) string {
	switch tetriminoType {
	case common.TetriminoNone:
		return colors.Background
	case common.I:
		return colors.TetriminoI
	case common.J:
		return colors.TetriminoJ
	case common.L:
		return colors.TetriminoL
	case common.O:
		return colors.TetriminoO
	case common.S:
		return colors.TetriminoS
	case common.T:
		return colors.TetriminoT
	case common.Z:
		return colors.TetriminoZ
	}
	return colors.Background
}

// DefaultTetriminoColors 默认颜色
var DefaultTetriminoColors = TetriminoColors{
	Border:     "#1b1b1b",
	Background: "#000000",
	TetriminoI: "#67c4ec",
	TetriminoJ: "#5f64a9",
	TetriminoL: "#df8136",
	TetriminoO: "#f0d543",
	TetriminoS: "#62b451",
	TetriminoT: "#a25399",
	TetriminoZ: "#db3e32",
}
