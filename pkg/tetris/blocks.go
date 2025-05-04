package tetris

import "fmt"

// Block 方块
type Block struct {
	// 方块类型
	Type BlockType
	// 方块位置
	Row, Column int
	// 方块方向
	Dir BlockDir
}

// BlockType 方块类型
type BlockType byte

// BlockType 的枚举
const (
	BlockNone BlockType = iota
	BlockI
	BlockJ
	BlockL
	BlockO
	BlockS
	BlockT
	BlockZ
)

// String 返回字符串表示
func (t BlockType) String() string {
	switch t {
	case BlockNone:
		return "None"
	case BlockI:
		return "I"
	case BlockJ:
		return "J"
	case BlockL:
		return "L"
	case BlockO:
		return "O"
	case BlockS:
		return "S"
	case BlockT:
		return "T"
	case BlockZ:
		return "Z"
	}
	return fmt.Sprintf("Invalid(%d)", t)
}

// BlockDir 方块方向
type BlockDir byte

// BlockDir 的枚举
const (
	Dir1 BlockDir = iota
	Dir2
	Dir3
	Dir4
)

var (
	// blockEdges 方块边缘
	blockEdges = [7][4][4]int{
		// I
		{{2, 2, 0, 3}, {3, 0, 2, 2}, {1, 1, 0, 3}, {3, 0, 1, 1}},
		// J
		{{2, 1, 0, 2}, {2, 0, 1, 2}, {1, 0, 0, 2}, {2, 0, 0, 1}},
		// L
		{{2, 1, 0, 2}, {2, 0, 1, 2}, {1, 0, 0, 2}, {2, 0, 0, 1}},
		// O
		{{1, 0, 0, 1}, {1, 0, 0, 1}, {1, 0, 0, 1}, {1, 0, 0, 1}},
		// S
		{{2, 1, 0, 2}, {2, 0, 1, 2}, {1, 0, 0, 2}, {2, 0, 0, 1}},
		// T
		{{2, 1, 0, 2}, {2, 0, 1, 2}, {1, 0, 0, 2}, {2, 0, 0, 1}},
		// Z
		{{2, 1, 0, 2}, {2, 0, 1, 2}, {1, 0, 0, 2}, {2, 0, 0, 1}},
	}
	// blockShapes 方块形状
	blockShapes = [7][4][4]Location{
		// I
		{
			{{2, 0}, {2, 1}, {2, 2}, {2, 3}},
			{{0, 2}, {1, 2}, {2, 2}, {3, 2}},
			{{1, 0}, {1, 1}, {1, 2}, {1, 3}},
			{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
		},
		// J
		{
			{{1, 0}, {2, 0}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 1}, {2, 1}, {2, 2}},
			{{0, 2}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
		},
		// L
		{
			{{1, 0}, {1, 1}, {1, 2}, {2, 2}},
			{{0, 1}, {0, 2}, {1, 1}, {2, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 1}, {2, 0}, {2, 1}},
		},
		// O
		{
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		},
		// S
		{
			{{1, 0}, {1, 1}, {2, 1}, {2, 2}},
			{{0, 2}, {1, 1}, {1, 2}, {2, 1}},
			{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
		},
		// T
		{
			{{1, 0}, {1, 1}, {1, 2}, {2, 1}},
			{{0, 1}, {1, 1}, {1, 2}, {2, 1}},
			{{0, 1}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 1}},
		},
		// Z
		{
			{{1, 1}, {1, 2}, {2, 0}, {2, 1}},
			{{0, 1}, {1, 1}, {1, 2}, {2, 2}},
			{{0, 1}, {0, 2}, {1, 0}, {1, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
		},
	}
)

// Edge 获取方块边缘行列号
func (b Block) Edge() (topRow, bottomRow, leftCol, rightCol int) {
	// 获取相对方块定位点的偏移
	if b.Type < 1 || b.Type > 7 || b.Dir < 0 || b.Dir > 3 {
		return
	}
	edge := blockEdges[b.Type-1][b.Dir]

	// 加上方块本身位置
	return edge[0] + b.Row, edge[1] + b.Row, edge[2] + b.Column, edge[3] + b.Column
}

// Cells 获取方块各格坐标
//
// 每个元素是一个方格的坐标
func (b Block) Cells() [4]Location {
	// 获取相对方块定位点的偏移
	if b.Type < 1 || b.Type > 7 || b.Dir < 0 || b.Dir > 3 {
		return [4]Location{}
	}
	ret := blockShapes[b.Type-1][b.Dir]

	// 加上方块本身位置
	for i := range ret {
		ret[i][0] += b.Row
		ret[i][1] += b.Column
	}

	return ret
}

// Location 格子位置 {row, col}
//
// 基于坐标原点，从下往上第一行为 0 ，从左往右第一列为 0
type Location [2]int

// Row 返回行序号
func (loc Location) Row() int {
	return loc[0]
}

// Column 返回列序号
func (loc Location) Column() int {
	return loc[1]
}
