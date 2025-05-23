package web

import (
	"sync"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris/common"
)

// NewTetrisGrid 创建 TetrisGrid
func NewTetrisGrid(rows, cols int, colors TetrominoColors) *TetrisGrid {
	data := make([][]common.Cell, rows)
	for i := range data {
		data[i] = make([]common.Cell, cols)
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
	colors      TetrominoColors

	data [][]common.Cell

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
	grid.paintTetrominoes()
}

// UpdateTetrominoes 更新方块
func (grid *TetrisGrid) UpdateTetrominoes(data [][]common.Cell) {
	grid.data = data
	rows := len(data)
	cols := grid.cols
	if rows > 0 {
		cols = len(data[0])
	}
	if rows != grid.rows || cols != grid.cols {
		grid.resize(rows, cols)
	}
	grid.paintTetrominoes()
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

// paintTetrominoes 绘制方块
func (grid *TetrisGrid) paintTetrominoes() {
	if grid.canvas.JSValue() == nil {
		return
	}
	canvasCTX := grid.canvas.JSValue().Call("getContext", "2d")

	canvasCTX.Set("lineWidth", grid.borderWidth)
	borderOffset := grid.borderWidth / 2

	for i, row := range grid.data {
		for j, cell := range row {
			color := grid.colors.Tetromino(cell.Type)
			x := j * (grid.cellWidth + grid.borderWidth)
			y := (grid.rows - i - 1) * (grid.cellWidth + grid.borderWidth)
			if cell.Shadow {
				canvasCTX.Set("strokeStyle", color)
				canvasCTX.Set("fillStyle", grid.colors.Background)
				canvasCTX.Call("fillRect", x, y, grid.cellWidth, grid.cellWidth)
				canvasCTX.Call(
					"strokeRect",
					x+borderOffset, y+borderOffset,
					grid.cellWidth-grid.borderWidth, grid.cellWidth-grid.borderWidth,
				)
			} else {
				canvasCTX.Set("fillStyle", color)
				canvasCTX.Call("fillRect", x, y, grid.cellWidth, grid.cellWidth)
			}
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

// newTetrominoGridData 创建方块网格数据
func newTetrominoGridData(tetrominoType common.TetrominoType) [][]common.Cell {
	i := common.Cell{Type: common.I}
	j := common.Cell{Type: common.J}
	l := common.Cell{Type: common.L}
	o := common.Cell{Type: common.O}
	s := common.Cell{Type: common.S}
	t := common.Cell{Type: common.T}
	z := common.Cell{Type: common.Z}
	none := common.Cell{}
	switch tetrominoType {
	case common.I:
		return [][]common.Cell{
			{i, i, i, i},
			{none, none, none, none},
		}
	case common.J:
		return [][]common.Cell{
			{j, j, j},
			{j, none, none},
		}
	case common.L:
		return [][]common.Cell{
			{l, l, l},
			{none, none, l},
		}
	case common.O:
		return [][]common.Cell{
			{o, o},
			{o, o},
		}
	case common.S:
		return [][]common.Cell{
			{s, s, none},
			{none, s, s},
		}
	case common.T:
		return [][]common.Cell{
			{t, t, t},
			{none, t, none},
		}
	case common.Z:
		return [][]common.Cell{
			{none, z, z},
			{z, z, none},
		}
	default:
		return [][]common.Cell{
			{none, none, none},
			{none, none, none},
		}
	}
}

// TetrominoColors 方块颜色
type TetrominoColors struct {
	Border     string
	Background string
	TetrominoI string
	TetrominoJ string
	TetrominoL string
	TetrominoO string
	TetrominoS string
	TetrominoT string
	TetrominoZ string
}

// Tetromino 获取指定方块颜色
func (colors TetrominoColors) Tetromino(tetrominoType common.TetrominoType) string {
	switch tetrominoType {
	case common.TetrominoNone:
		return colors.Background
	case common.I:
		return colors.TetrominoI
	case common.J:
		return colors.TetrominoJ
	case common.L:
		return colors.TetrominoL
	case common.O:
		return colors.TetrominoO
	case common.S:
		return colors.TetrominoS
	case common.T:
		return colors.TetrominoT
	case common.Z:
		return colors.TetrominoZ
	}
	return colors.Background
}

// DefaultTetrominoColors 默认颜色
var DefaultTetrominoColors = TetrominoColors{
	Border:     "#1b1b1b",
	Background: "#000000",
	TetrominoI: "#67c4ec",
	TetrominoJ: "#5f64a9",
	TetrominoL: "#df8136",
	TetrominoO: "#f0d543",
	TetrominoS: "#62b451",
	TetrominoT: "#a25399",
	TetrominoZ: "#db3e32",
}
