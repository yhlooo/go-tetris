package tetris

import (
	"context"
	"fmt"
)

// Tetris 游戏实例
type Tetris interface {
	// State 返回当前游戏状态
	State() GameState
	// Start 开始游戏
	//
	// 对于每个 Tetris 对象只能被调用一次
	Start(ctx context.Context) error
	// Stop 停止游戏
	//
	// 对于每个 Tetris 对象只能被调用一次
	Stop() error
	// Pause 暂停游戏
	Pause() error
	// Resume 继续游戏
	Resume() error
	// SetDebug 设置调试模式
	SetDebug(enabled bool)
	// Debug 返回是否调试模式
	Debug() bool
	// ChangeActiveBlockType 更换活跃方块类型
	//
	// 仅在调试模式下生效
	ChangeActiveBlockType(blockType BlockType) error

	// Input 输入操作指令
	Input(op Op)

	// Frames 获取帧通道
	//
	// 游戏运行时，若 Tetris 对象中用于构成画面的信息发生变化，将通过该通道发送新的帧，用于更新画面。
	// 游戏结束后该通道会被关闭。
	//
	// 通道满时产生的帧会被丢弃。
	//
	// 每个 Tetris 对象只有一个通道，多次调用该方法返回的是同一通道
	Frames() <-chan Frame
	// CurrentFrame 获取当前帧
	CurrentFrame() Frame
}

// GameState 游戏状态
type GameState byte

// GameState 的枚举值
const (
	StatePending GameState = iota
	StateRunning
	StatePaused
	StateFinished
)

// String 返回字符串表示
func (s GameState) String() string {
	switch s {
	case StatePending:
		return "Pending"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateFinished:
		return "Finished"
	}
	return "Invalid"
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

// String 返回字符串表示
func (op Op) String() string {
	switch op {
	case OpMoveRight:
		return "MoveRight"
	case OpMoveLeft:
		return "MoveLeft"
	case OpRotateRight:
		return "RotateRight"
	case OpRotateLeft:
		return "RotateLeft"
	case OpSoftDrop:
		return "SoftDrop"
	case OpHardDrop:
		return "HardDrop"
	case OpHold:
		return "Hold"
	}
	return fmt.Sprintf("Invalid(%d)", op)
}

// Frame 帧
//
// 包含某时刻游戏画面应显示的信息，如方块位置、得分等
type Frame struct {
	// 场上方块填充情况
	Field FieldReader
	// 暂存的方块
	HoldingBlock *BlockType
	// 下几个方块
	NextBlocks []BlockType
	// 级别
	Level int
	// 分数
	Score int
	// 已消除的行数
	ClearLines int
	// 游戏结束
	GameOver bool
}
