package common

import (
	"fmt"
)

// TetrominoType 方块类型
type TetrominoType byte

// TetrominoType 的枚举
const (
	TetrominoNone TetrominoType = iota
	I
	J
	L
	O
	S
	T
	Z
)

// String 返回字符串表示
func (t TetrominoType) String() string {
	switch t {
	case TetrominoNone:
		return "None"
	case I:
		return "I"
	case J:
		return "J"
	case L:
		return "L"
	case O:
		return "O"
	case S:
		return "S"
	case T:
		return "T"
	case Z:
		return "Z"
	}
	return fmt.Sprintf("Invalid(%d)", t)
}

// TetrominoDir 方块方向
type TetrominoDir byte

// TetrominoDir 的枚举
const (
	// Dir0 初始状态
	Dir0 TetrominoDir = iota
	// DirR 顺时针旋转 90 度
	DirR
	// Dir2 旋转 180 度
	Dir2
	// DirL 逆时针旋转 90 度
	DirL
)

// String 返回字符串表示
func (d TetrominoDir) String() string {
	switch d {
	case Dir0:
		return "0"
	case DirR:
		return "R"
	case Dir2:
		return "2"
	case DirL:
		return "L"
	}
	return fmt.Sprintf("Invalid(%d)", d)
}

var (
	// tetrominoShapes 方块形状
	tetrominoShapes = [7][4][4]Location{
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

// Tetromino 方块
type Tetromino struct {
	// 方块类型
	Type TetrominoType
	// 方块位置
	Row, Column int
	// 方块方向
	Dir TetrominoDir
}

// Cells 获取方块各格坐标
//
// 每个元素是一个方格的坐标
func (t Tetromino) Cells() [4]Location {
	// 获取相对方块定位点的偏移
	if t.Type < 1 || t.Type > 7 || t.Dir < 0 || t.Dir > 3 {
		return [4]Location{}
	}
	ret := tetrominoShapes[t.Type-1][t.Dir]

	// 加上方块本身位置
	for i := range ret {
		ret[i][0] += t.Row
		ret[i][1] += t.Column
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
