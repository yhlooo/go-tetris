package tetris

import (
	"fmt"
	"strconv"

	"github.com/go-logr/logr"

	"github.com/yhlooo/go-tetris/pkg/tetris/randomizer"
	"github.com/yhlooo/go-tetris/pkg/tetris/rotationsystems"
)

// Options 游戏选项
type Options struct {
	// 行列数
	Rows, Columns int

	// 是否开启暂存方块功能
	HoldEnabled bool
	// 提示的下个方块数量
	ShowNextTetriminos int

	// 初始级别
	InitialLevel int
	// 每级别需要消除多少行
	LinesPerLevel int
	// 下落速度控制器
	SpeedController SpeedController
	// 处理频率（单位： ticket/s ）
	Frequency int

	// 随机生成器
	Randomizer randomizer.Randomizer
	// 评分器
	Scorer Scorer
	// 旋转系统
	RotationSystem rotationsystems.RotationSystem

	Logger logr.Logger
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
	if opts.SpeedController == nil {
		opts.SpeedController = DefaultSpeedController
	}
	if opts.Frequency == 0 {
		opts.Frequency = 1000
	}

	if opts.Randomizer == nil {
		opts.Randomizer = &randomizer.Bag7{}
	}
	if opts.Scorer == nil {
		opts.Scorer = DefaultScorer()
	}
	if opts.RotationSystem == nil {
		opts.RotationSystem = rotationsystems.SuperRotationSystem{}
	}
}

// SpeedController 返回指定级别下落速度（单位：格/s ）
type SpeedController func(level int) float64

// Scorer 评分器
type Scorer func(level int, event ScoreEvent) (score int, reason []string)

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

	HoldEnabled:        true,
	ShowNextTetriminos: 3,

	InitialLevel:    1,
	LinesPerLevel:   10,
	SpeedController: DefaultSpeedController,
	Frequency:       60,

	Randomizer:     &randomizer.Bag7{},
	Scorer:         DefaultScorer(),
	RotationSystem: rotationsystems.SuperRotationSystem{},

	Logger: logr.Discard(),
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
	return func(level int, event ScoreEvent) (int, []string) {
		score := 0
		var reason []string

		if event.SoftDrop > 0 {
			score += event.SoftDrop
			d := ""
			if event.SoftDrop > 1 {
				d = " " + strconv.Itoa(event.SoftDrop)
			}
			reason = append(reason, "Soft Drop"+d)
		}
		if event.HardDrop > 0 {
			score += event.HardDrop * 2
			reason = append(reason, fmt.Sprintf("Hard Drop %d", event.HardDrop))
		}

		// 清行分
		clearScore := 0
		difficult := false
		if event.TSpin {
			switch event.ClearLines {
			case 1:
				// T-Spin Single
				clearScore = 800
				difficult = true
				reason = append(reason, "T-Spin Single")
			case 2:
				// T-Spin Double
				clearScore = 1200
				difficult = true
				reason = append(reason, "T-Spin Double")
			case 3:
				// T-Spin Triple
				clearScore = 1600
				difficult = true
				reason = append(reason, "T-Spin Triple")
			}
		} else {
			switch event.ClearLines {
			case 1:
				// Single Line
				clearScore = 100
				reason = append(reason, "Single Line Clear")
			case 2:
				// Double Line
				clearScore = 300
				reason = append(reason, "Double Line Clear")
			case 3:
				// Triple Line
				clearScore = 500
				reason = append(reason, "Triple Line Clear")
			case 4:
				// Tetris
				clearScore = 800
				difficult = true
				reason = append(reason, "Tetris")
			}
		}

		if b2b && difficult {
			// Back-to-Back
			clearScore += clearScore / 2
			reason = append(reason, "Back-to-Back")
		}
		if event.ClearLines > 0 {
			b2b = difficult
		}

		score += clearScore

		if event.ClearLines == 0 && event.TSpin {
			// T-Spin
			score += 400
			reason = append(reason, "T-Spin")
		}

		return score * level, reason
	}
}
