package tetris

import "time"

// Options 游戏选项
type Options struct {
	// 行列数
	Rows, Columns int
	// 是否开启暂存方块功能
	HoldEnabled bool
	// 初始级别
	InitialLevel int
	// 每级别需要消除多少行
	LinesPerLevel int
	// 下落速度控制器
	SpeedController SpeedController
	// 处理频率（单位： ticket/s ）
	Frequency int
	// 评分器
	Scorer Scorer
	// 随机数种子
	RandSeed int64
}

// Complete 补全选项
func (opts *Options) Complete() {
	if opts.Rows == 0 {
		opts.Rows = 20
	}
	if opts.Columns == 0 {
		opts.Columns = 10
	}
	if opts.InitialLevel == 0 {
		opts.InitialLevel = 1
	}
	if opts.LinesPerLevel == 0 {
		opts.LinesPerLevel = 10
	}
	if opts.Frequency == 0 {
		opts.Frequency = 1000
	}
	if opts.RandSeed == 0 {
		opts.RandSeed = time.Now().UnixNano()
	}
}

// SpeedController 返回指定级别下落速度（单位：格/s ）
type SpeedController func(level int) float64

// Scorer 评分器
type Scorer func(level int, event ScoreEvent) int

// ScoreEvent 评分事件
type ScoreEvent struct {
	// 软下落行数
	SoftDrop int
	// 硬下落行数
	HardDrop int
	// 清除行数
	ClearLines int
	// T-Spin
	TSpin bool
}
