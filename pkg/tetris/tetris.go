package tetris

import "context"

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
type Frame struct{}
