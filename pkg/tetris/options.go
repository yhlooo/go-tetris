package tetris

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
