package tetris

import "time"

// Options 游戏选项
type Options struct {
	// 行列数
	Rows, Columns int

	// 是否开启暂存方块功能
	HoldEnabled bool
	// 提示的下个方块数量
	NextBlock int

	// 随机数种子
	RandSeed int64

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
	// 旋转系统
	RotationSystem RotationSystem
}

// Complete 补全选项
func (opts *Options) Complete() {
	if opts.Rows == 0 {
		opts.Rows = 20
	}
	if opts.Columns == 0 {
		opts.Columns = 10
	}

	if opts.RandSeed == 0 {
		opts.RandSeed = time.Now().UnixNano()
	}

	if opts.InitialLevel == 0 {
		opts.InitialLevel = 1
	}
	if opts.LinesPerLevel == 0 {
		opts.LinesPerLevel = 10
	}
	if opts.SpeedController == nil {
		opts.SpeedController = DefaultSpeedController
	}
	if opts.Frequency == 0 {
		opts.Frequency = 1000
	}

	if opts.Scorer == nil {
		opts.Scorer = DefaultScorer()
	}
	if opts.RotationSystem == nil {
		opts.RotationSystem = SuperRotationSystem{}
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

// DefaultOptions 默认选项
var DefaultOptions = Options{
	Rows:    20,
	Columns: 10,

	HoldEnabled: true,
	NextBlock:   3,

	RandSeed: 0,

	InitialLevel:    1,
	LinesPerLevel:   10,
	SpeedController: DefaultSpeedController,
	Frequency:       1000,

	Scorer:         DefaultScorer(),
	RotationSystem: SuperRotationSystem{},
}

// DefaultSpeedController 默认速度控制器
func DefaultSpeedController(level int) float64 {
	switch level {
	case 1:
		return 1
	case 2:
		return 1.26102
	case 3:
		return 1.61862
	case 4:
		return 2.11536
	case 5:
		return 2.8158
	case 6:
		return 3.8166
	case 7:
		return 5.274
	case 8:
		return 7.416
	case 9:
		return 10.65
	case 10:
		return 15.588
	case 11:
		return 23.28
	case 12:
		return 35.4
	case 13:
		return 55.2
	case 14:
		return 87.6
	default:
		return 141.6
	}
}

// DefaultScorer 默认评分器
func DefaultScorer() Scorer {
	b2b := false
	return func(level int, event ScoreEvent) int {
		score := 0
		score += event.SoftDrop
		score += event.HardDrop * 2

		// 清行分
		clearScore := 0
		difficult := false
		if event.TSpin {
			switch event.ClearLines {
			case 1:
				// T-Spin Single
				clearScore = 800
				difficult = true
			case 2:
				// T-Spin Double
				clearScore = 1200
				difficult = true
			case 3:
				// T-Spin Triple
				clearScore = 1600
				difficult = true
			}
		} else {
			switch event.ClearLines {
			case 1:
				// Single Line
				clearScore = 100
			case 2:
				// Double Line
				clearScore = 300
			case 3:
				// Triple Line
				clearScore = 500
			case 4:
				// Tetris
				clearScore = 800
				difficult = true
			}
		}

		if b2b && difficult {
			// Back-to-Back
			clearScore += clearScore / 2
		}
		b2b = difficult

		score += clearScore

		if event.ClearLines == 0 && event.TSpin {
			// T-Spin
			score += 400
		}

		return score * level
	}
}
