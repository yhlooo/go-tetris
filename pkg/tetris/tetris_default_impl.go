package tetris

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"

	"github.com/yhlooo/go-tetris/pkg/tetris/common"
	"github.com/yhlooo/go-tetris/pkg/tetris/randomizer"
	"github.com/yhlooo/go-tetris/pkg/tetris/rotationsystems"
)

const framesChLen = 16

// NewTetris 创建 Tetris 游戏实例
func NewTetris(opts Options) Tetris {
	opts.Complete()
	t := &defaultTetris{
		rows:  opts.Rows,
		cols:  opts.Columns,
		level: opts.InitialLevel,

		holdEnabled: opts.HoldEnabled,

		linesPerLevel: opts.LinesPerLevel,
		ticker:        time.NewTicker(time.Second / time.Duration(opts.Frequency)),
		speed:         opts.SpeedController,
		freq:          opts.Frequency,

		randomizer:     opts.Randomizer,
		scorer:         opts.Scorer,
		rotationSystem: opts.RotationSystem,

		state:    StatePending,
		framesCh: make(chan Frame, framesChLen),

		logger: opts.Logger,
	}
	t.field = common.NewField(opts.Rows, opts.Columns, t.newTetromino(common.TetrominoNone))
	for i := 0; i < opts.ShowNextTetrominoes+1; i++ {
		t.nextTetrominoes = append(t.nextTetrominoes, t.randomizer.Next())
	}
	return t
}

// defaultTetris 是 Tetris 的默认实现
type defaultTetris struct {
	lock      sync.Mutex
	startOnce sync.Once
	cancel    context.CancelFunc

	rows, cols       int
	field            *common.Field
	nextTetrominoes  []common.TetrominoType
	holdingTetromino *common.TetrominoType
	level            int
	score            int
	clearLines       int

	holed   bool
	tickets int64
	notMove bool

	holdEnabled bool

	linesPerLevel int
	ticker        *time.Ticker
	speed         SpeedController
	freq          int

	randomizer     randomizer.Randomizer
	scorer         Scorer
	rotationSystem rotationsystems.RotationSystem

	debug    bool
	state    GameState
	framesCh chan Frame
	logger   logr.Logger
}

var _ Tetris = (*defaultTetris)(nil)

// State 返回当前游戏状态
func (t *defaultTetris) State() GameState {
	return t.state
}

// Start 开始游戏
func (t *defaultTetris) Start(ctx context.Context) error {
	var err error
	t.startOnce.Do(func() {
		t.lock.Lock()
		defer t.lock.Unlock()

		ctx, t.cancel = context.WithCancel(ctx)
		t.state = StateRunning
		go t.run(ctx)
		t.sendFrame()
		t.logger.Info("started")
	})
	return err
}

// Stop 停止游戏
func (t *defaultTetris) Stop() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.cancel != nil {
		t.cancel()
	}
	t.state = StateFinished
	t.logger.Info("stoped")
	return nil
}

// Pause 暂停游戏
func (t *defaultTetris) Pause() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.state != StateRunning {
		return fmt.Errorf("not in running state: %s", t.state)
	}
	t.state = StatePaused
	t.logger.Info("paused")
	return nil
}

// Resume 继续游戏
func (t *defaultTetris) Resume() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.state != StatePaused {
		return fmt.Errorf("not in paused state: %s", t.state)
	}
	t.state = StateRunning
	t.logger.Info("resumed")
	return nil
}

// SetDebug 设置调试模式
func (t *defaultTetris) SetDebug(enabled bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.debug = enabled
	t.logger.Info(fmt.Sprintf("set debug mode: %t", enabled))
}

// Debug 返回是否调试模式
func (t *defaultTetris) Debug() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.debug
}

// ChangeActiveTetrominoType 更换活跃方块类型
func (t *defaultTetris) ChangeActiveTetrominoType(tetrominoType common.TetrominoType) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.debug {
		return fmt.Errorf("not in debug mode")
	}
	if t.state != StateRunning && t.state != StatePaused {
		return fmt.Errorf("not in running or paused state: %s", t.state)
	}

	b := t.field.ActiveTetromino()
	oldType := b.Type
	b.Type = tetrominoType

	// 不合法，还原
	if !t.field.IsValid() {
		b.Type = oldType
		return fmt.Errorf("insufficient space")
	}

	t.logger.V(1).Info(fmt.Sprintf("change tetromino type: %s -> %s", oldType, tetrominoType))
	t.sendFrame()

	return nil
}

// Input 输入操作指令
func (t *defaultTetris) Input(op Op) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.state != StateRunning && (!t.debug || t.state != StatePaused) {
		t.logger.V(1).Info(fmt.Sprintf("ignore input %q: not running: %s", op, t.state))
		return
	}

	changed := false
	switch op {
	case OpMoveRight:
		changed = t.field.MoveActiveTetromino(0, 1)
		if changed {
			t.notMove = false
		}
		t.logger.V(1).Info(fmt.Sprintf("move right, ret: %t", changed))
	case OpMoveLeft:
		changed = t.field.MoveActiveTetromino(0, -1)
		if changed {
			t.notMove = false
		}
		t.logger.V(1).Info(fmt.Sprintf("move left, ret: %t", changed))
	case OpRotateRight:
		changed = t.rotationSystem.RotateRight(t.field)
		if changed {
			t.notMove = true
		}
		t.logger.V(1).Info(fmt.Sprintf("rotate right, ret: %t", changed))
	case OpRotateLeft:
		changed = t.rotationSystem.RotateLeft(t.field)
		if changed {
			t.notMove = true
		}
		t.logger.V(1).Info(fmt.Sprintf("rotate left, ret: %t", changed))
	case OpSoftDrop:
		t.logger.V(1).Info("soft drop")
		if ok := t.field.MoveActiveTetromino(-1, 0); ok {
			t.notMove = false
			t.calcScore(ScoreEvent{SoftDrop: 1})
		}
		t.tickets = 0
		changed = true
	case OpHardDrop:
		t.logger.V(1).Info("hard drop")
		dropLines := 0
		for t.field.MoveActiveTetromino(-1, 0) {
			t.notMove = false
			dropLines++
		}
		t.calcScore(ScoreEvent{HardDrop: dropLines})
		t.lockDown()
		t.tickets = 0
		t.logger.V(1).Info("lock down tetromino")
		changed = true
	case OpHold:
		if !t.holed && t.holdEnabled {
			oldActive := t.field.ActiveTetromino().Type
			var ok bool
			if t.holdingTetromino != nil {
				ok = t.field.ChangeActiveTetromino(t.newTetromino(*t.holdingTetromino))
			} else {
				ok = t.field.ChangeActiveTetromino(t.newTetromino(t.nextTetrominoes[0]))
				if ok {
					t.nextTetrominoes = append(t.nextTetrominoes[1:], t.randomizer.Next())
				}
			}
			if ok {
				t.holdingTetromino = &oldActive
				t.holed = true
				t.notMove = false
			}
			changed = ok
			t.logger.V(1).Info(fmt.Sprintf("hold tetromino, ret: %t", ok))
		} else {
			t.logger.V(1).Info("can not hold tetromino: tetromino already holed")
		}
	}
	if changed {
		if ok := t.sendFrame(); !ok {
			t.logger.Info("WARN: frames channel busy, frame dropped")
		}
	}
}

// Frames 获取帧通道
func (t *defaultTetris) Frames() <-chan Frame {
	return t.framesCh
}

// CurrentFrame 获取当前帧
func (t *defaultTetris) CurrentFrame() Frame {
	return Frame{
		Field:            t.field,
		HoldingTetromino: t.holdingTetromino,
		NextTetrominoes:  t.nextTetrominoes[:len(t.nextTetrominoes)-1],
		Level:            t.level,
		Score:            t.score,
		ClearLines:       t.clearLines,
		GameOver:         t.state == StateFinished,
	}
}

// run 运行
func (t *defaultTetris) run(ctx context.Context) {
	defer close(t.framesCh)
	defer t.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.ticker.C:
		}

		if t.state != StateRunning {
			continue
		}

		t.tickets++
		t.lock.Lock()
		speed := t.speed(t.level)
		if t.tickets > int64(float64(t.freq)/speed) {
			// 到时间下落一格了
			t.logger.V(1).Info("auto drop")
			if ok := t.field.MoveActiveTetromino(-1, 0); !ok {
				t.lockDown()
				t.logger.V(1).Info("lock down tetromino")
			} else {
				t.notMove = false
			}
			t.sendFrame()
			t.tickets = 0
		}
		t.lock.Unlock()
	}
}

// lockDown 锁定当前活跃方块
func (t *defaultTetris) lockDown() {
	tSpin, clearLines, ok := t.field.LockDown(t.newTetromino(t.nextTetrominoes[0]))
	t.calcScore(ScoreEvent{TSpin: tSpin && t.notMove, ClearLines: clearLines})
	t.clearLines += clearLines
	t.level = t.clearLines/t.linesPerLevel + 1
	if !ok {
		t.state = StateFinished
	}
	t.nextTetrominoes = append(t.nextTetrominoes[1:], t.randomizer.Next())
	t.holed = false
}

// sendFrame 发送帧
func (t *defaultTetris) sendFrame() bool {
	select {
	case t.framesCh <- t.CurrentFrame():
	default:
		return false
	}
	return true
}

// calcScore 计算分数
func (t *defaultTetris) calcScore(event ScoreEvent) {
	score, reason := t.scorer(t.level, event)
	if score > 0 {
		t.score += score
		t.logger.Info(fmt.Sprintf("SCORE %s: +%d", strings.Join(reason, ", "), score))
	}
}

// newTetromino 创建新方块
func (t *defaultTetris) newTetromino(tetrominoType common.TetrominoType) *common.Tetromino {
	if tetrominoType == common.TetrominoNone {
		tetrominoType = t.randomizer.Next()
	}

	// 确定位置，放在居中上方刚好露出完整方块的位置
	col := t.cols/2 - 2
	row := t.rows - 3
	switch tetrominoType {
	case common.O:
		col = t.cols/2 - 1
		row = t.rows - 2
	default:
	}

	return &common.Tetromino{
		Type:   tetrominoType,
		Row:    row,
		Column: col,
		Dir:    common.Dir0,
	}
}
