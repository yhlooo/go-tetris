package tetris

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-logr/logr"
)

const framesChLen = 16

// NewTetris 创建 Tetris 游戏实例
func NewTetris(opts Options) Tetris {
	opts.Complete()
	t := &defaultTetris{
		rand:     rand.New(rand.NewSource(opts.RandSeed)),
		framesCh: make(chan Frame, framesChLen),
		level:    opts.InitialLevel,
		ticker:   time.NewTicker(time.Second / time.Duration(opts.Frequency)),
		freq:     opts.Frequency,
		speed:    opts.SpeedController,
	}
	t.field = NewField(opts.Rows, opts.Columns, nil)
	t.field.ChangeActiveBlock(t.newBlock(BlockNone))
	t.nextBlock = t.newBlockType()
	return t
}

// defaultTetris 是 Tetris 的默认实现
type defaultTetris struct {
	lock      sync.Mutex
	startOnce sync.Once
	cancel    context.CancelFunc

	field        *Field
	nextBlock    BlockType
	holdingBlock *BlockType
	holed        bool
	level        int
	score        int
	clearLines   int

	rand    *rand.Rand
	ticker  *time.Ticker
	freq    int
	speed   SpeedController
	tickets int64

	pause    bool
	framesCh chan Frame
}

var _ Tetris = (*defaultTetris)(nil)

// Start 开始游戏
func (t *defaultTetris) Start(ctx context.Context) error {
	var err error
	t.startOnce.Do(func() {
		t.lock.Lock()
		defer t.lock.Unlock()

		ctx, t.cancel = context.WithCancel(ctx)
		go t.run(ctx)
		t.sendFrame()
		logr.FromContextOrDiscard(ctx).Info("started")
	})
	return err
}

// Stop 停止游戏
func (t *defaultTetris) Stop(ctx context.Context) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.cancel != nil {
		t.cancel()
	}
	logr.FromContextOrDiscard(ctx).Info("stoped")
}

// Pause 暂停游戏
func (t *defaultTetris) Pause(ctx context.Context) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.pause = true
	logr.FromContextOrDiscard(ctx).Info("paused")
}

// Resume 继续游戏
func (t *defaultTetris) Resume(ctx context.Context) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.pause = false
	logr.FromContextOrDiscard(ctx).Info("resumed")
}

// Input 输入操作指令
func (t *defaultTetris) Input(ctx context.Context, op Op) {
	logger := logr.FromContextOrDiscard(ctx)

	t.lock.Lock()
	defer t.lock.Unlock()

	if t.pause {
		return
	}

	changed := false
	switch op {
	case OpMoveRight:
		changed = t.field.MoveActiveBlock(0, 1)
		logger.V(1).Info(fmt.Sprintf("move right, ret: %t", changed))
	case OpMoveLeft:
		changed = t.field.MoveActiveBlock(0, -1)
		logger.V(1).Info(fmt.Sprintf("move left, ret: %t", changed))
	case OpRotateRight:
		changed = t.field.RotateActiveBlock(1)
		logger.V(1).Info(fmt.Sprintf("rotate right, ret: %t", changed))
	case OpRotateLeft:
		changed = t.field.RotateActiveBlock(-1)
		logger.V(1).Info(fmt.Sprintf("rotate left, ret: %t", changed))
	case OpSoftDrop:
		logger.V(1).Info("soft drop")
		if ok := t.field.MoveActiveBlock(-1, 0); !ok {
			t.pinBlock()
			logger.V(1).Info("pin block")
		}
		t.tickets = 0
		changed = true
	case OpHardDrop:
		logger.V(1).Info("hard drop")
		for t.field.MoveActiveBlock(-1, 0) {
		}
		t.pinBlock()
		t.tickets = 0
		logger.V(1).Info("pin block")
		changed = true
	case OpHold:
		if !t.holed {
			logger.V(1).Info("hold block")
			oldActive := t.field.ActiveBlock().Type
			if t.holdingBlock != nil {
				t.field.ChangeActiveBlock(t.newBlock(*t.holdingBlock))
			} else {
				t.field.ChangeActiveBlock(t.newBlock(t.nextBlock))
				t.nextBlock = t.newBlockType()
			}
			t.holdingBlock = &oldActive
			changed = true
		} else {
			logger.V(1).Info("can not hold block: block already holed")
		}
	}
	if changed {
		if ok := t.sendFrame(); !ok {
			logger.Info("WARN: frames channel busy, frame dropped")
		}
	}
}

// Frames 获取帧通道
func (t *defaultTetris) Frames() <-chan Frame {
	return t.framesCh
}

// run 运行
func (t *defaultTetris) run(ctx context.Context) {
	logger := logr.FromContextOrDiscard(ctx)

	defer close(t.framesCh)
	defer t.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.ticker.C:
		}

		if t.pause {
			continue
		}

		t.tickets++
		t.lock.Lock()
		speed := t.speed(t.level)
		if t.tickets > int64(float64(t.freq)/speed) {
			// 到时间下落一格了
			logger.V(1).Info("auto drop")
			if ok := t.field.MoveActiveBlock(-1, 0); !ok {
				t.pinBlock()
				logger.V(1).Info("pin block")
			}
			t.sendFrame()
			t.tickets = 0
		}
		t.lock.Unlock()

	}
}

// pinBlock 钉住当前活跃方块
func (t *defaultTetris) pinBlock() {
	t.field.PinActiveBlock(t.newBlock(t.nextBlock))
	t.nextBlock = t.newBlockType()
	t.holed = false
}

// sendFrame 发送帧
func (t *defaultTetris) sendFrame() bool {
	select {
	case t.framesCh <- Frame{
		Field:        t.field,
		HoldingBlock: t.holdingBlock,
		NextBlock:    t.nextBlock,
		Level:        t.level,
		Score:        t.score,
		ClearLines:   t.clearLines,
	}:
	default:
		return false
	}
	return true
}

// newBlockType 创建新方块类型
func (t *defaultTetris) newBlockType() BlockType {
	return BlockType(t.rand.Int())%7 + 1
}

// newBlock 创建新方块
func (t *defaultTetris) newBlock(blockType BlockType) *Block {
	if blockType == BlockNone {
		blockType = t.newBlockType()
	}
	rows, cols := t.field.Size()
	col := cols/2 - 2
	row := rows - 3
	switch blockType {
	case BlockO:
		col = cols/2 - 1
		row = rows - 2
	default:
	}
	return &Block{
		Type:   blockType,
		Row:    row,
		Column: col,
		Dir:    Dir1,
	}
}
