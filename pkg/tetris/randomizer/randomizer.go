package randomizer

import "github.com/yhlooo/go-tetris/pkg/tetris/common"

// Randomizer 方块随机生成器
type Randomizer interface {
	// Next 获取下一个方块类型
	Next() common.TetrominoType
}
