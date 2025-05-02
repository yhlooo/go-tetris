package tetris

import (
	"context"
)

// Tetris 游戏实例
type Tetris interface {
	// Start 开始游戏
	//
	// 对于每个 Tetris 对象只能被调用一次
	Start(ctx context.Context) error
	// Stop 停止游戏
	//
	// 对于每个 Tetris 对象只能被调用一次
	Stop(ctx context.Context) error
	// Pause 暂停游戏
	Pause(ctx context.Context) error
	// Resume 继续游戏
	Resume(ctx context.Context) error

	// Input 输入操作指令
	Input(ctx context.Context, op Op) error

	// Frames 获取帧通道
	//
	// 游戏运行时，若 Tetris 对象中用于构成画面的信息发生变化，将通过该通道发送新的帧，用于更新画面。
	// 游戏结束后该通道会被关闭。
	//
	// 通道满时产生的帧会被丢弃。
	//
	// 每个 Tetris 对象只有一个通道，多次调用该方法返回的是同一通道
	Frames() <-chan Frame
}

// Op 操作指令
type Op byte

const (
	// OpMoveRight 向右移动
	OpMoveRight Op = iota
	// OpMoveLeft 向左移动
	OpMoveLeft
	// OpRotateRight 顺时针旋转
	OpRotateRight
	// OpRotateLeft 逆时针旋转
	OpRotateLeft
	// OpSoftDrop 软下落（下落一格）
	OpSoftDrop
	// OpHardDrop 硬下落（下落到底）
	OpHardDrop
	// OpHold 暂存当前方块
	OpHold
)

// Frame 帧
//
// 包含某时刻游戏画面应显示的信息，如方块位置、得分等
type Frame struct {
	// 行列数
	Rows, Columns int
	// 场上方块填充情况，不含活跃方块
	Field FieldStatus
	// 活跃的方块
	ActiveBlock Block
	// 暂存的方块
	HoldingBlock *BlockType
	// 下一个方块
	NextBlock BlockType
	// 级别
	Level int
	// 分数
	Score int
	// 已消除的行数
	ClearLines int
}

// FieldStatus 场上方块填充情况
type FieldStatus [][]BlockType

// Block 获取指定位置填充的方块类型
func (f FieldStatus) Block(row, col int) (BlockType, bool) {
	if row < 0 || len(f) <= row {
		return 0, false
	}
	if col < 0 || len(f[row]) <= col {
		return 0, false
	}
	return f[row][col], true
}

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

// BlockDir 方块方向
type BlockDir byte

// BlockDir 的枚举
const (
	Dir1 BlockDir = iota
	Dir2
	Dir3
	Dir4
)
