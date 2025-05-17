package randomizer

import (
	"math/rand/v2"
	"sync"
	"time"

	"github.com/yhlooo/go-tetris/pkg/tetris/common"
)

// New7Bag 创建 7-Bag 生成器
func New7Bag(s rand.Source) *Bag7 {
	return &Bag7{
		rand: rand.New(s),
	}
}

// Bag7 7-Bag 生成器
//
// 以包为单位生成，每次生成含 7 种方块的 7 个方块，打乱顺序依次发出
type Bag7 struct {
	lock   sync.Mutex
	rand   *rand.Rand
	buffer [7]common.TetrominoType
	i      int
}

var _ Randomizer = (*Bag7)(nil)

// Next 获取下一个方块类型
func (b *Bag7) Next() common.TetrominoType {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.rand == nil {
		b.rand = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	}

	if b.buffer[0] == common.TetrominoNone {
		b.buffer = [7]common.TetrominoType{
			common.I, common.J, common.L, common.O,
			common.S, common.T, common.Z,
		}
		b.rand.Shuffle(7, func(i, j int) {
			b.buffer[i], b.buffer[j] = b.buffer[j], b.buffer[i]
		})
	}
	if b.i >= 7 {
		b.rand.Shuffle(7, func(i, j int) {
			b.buffer[i], b.buffer[j] = b.buffer[j], b.buffer[i]
		})
		b.i = 0
	}

	ret := b.buffer[b.i]
	b.i++
	return ret
}
