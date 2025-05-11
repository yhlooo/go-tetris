package web

import (
	"math"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris"
)

// TouchController 触摸控制器
type TouchController struct {
	lock sync.Mutex

	tetris tetris.Tetris

	startTime            time.Time
	lastMoveTime         time.Time
	lastMoveX, lastMoveY int
	lastX, lastY         int
	offsetX, offsetY     int
	opCnt                int
}

// SetTetris 设置控制的游戏对象
func (c *TouchController) SetTetris(tetris tetris.Tetris) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.tetris = tetris
}

// HandleTouchStart 处理触摸开始事件
func (c *TouchController) HandleTouchStart(_ app.Context, e app.Event) {
	touch := c.getTouch(e)
	if touch == nil {
		return
	}
	if c.tetris == nil || c.tetris.State() != tetris.StateRunning {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.startTime = time.Now()
	c.lastMoveX, c.lastMoveY = 0, 0
	c.offsetX, c.offsetY = 0, 0
	c.lastX = touch.Get("screenX").Int()
	c.lastY = touch.Get("screenY").Int()
	c.opCnt = 0
}

// HandleTouchMove 处理触摸移动事件
func (c *TouchController) HandleTouchMove(ctx app.Context, e app.Event) {
	touch := c.getTouch(e)
	if touch == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// 当前坐标
	x := touch.Get("screenX").Int()
	y := touch.Get("screenY").Int()
	// 更新上次移动情况
	if time.Since(c.lastMoveTime) < 200*time.Millisecond {
		c.lastMoveX += x - c.lastX
		c.lastMoveY += y - c.lastY
	} else {
		c.lastMoveX = x - c.lastX
		c.lastMoveY = y - c.lastY
		c.lastMoveTime = time.Now()
	}
	// 累计偏移
	c.offsetX += x - c.lastX
	c.offsetY += y - c.lastY
	// 更新上一坐标
	c.lastX = x
	c.lastY = y

	// 计算移动位移
	moveY := c.offsetY / 20
	c.offsetY %= 20
	moveX := c.offsetX / 20
	c.offsetX %= 20
	if math.Abs(float64(moveX)) > 2*math.Abs(float64(moveY)) {
		moveY = 0
	} else if moveY != 0 {
		moveX = 0
	}

	if c.tetris == nil {
		return
	}

	// 操作位移
	if moveX != 0 {
		op := tetris.OpMoveRight
		if moveX < 0 {
			op = tetris.OpMoveLeft
			moveX = -moveX
		}
		for i := 0; i < moveX; i++ {
			app.Logf("touch: %s", op)
			c.tetris.Input(op)
			c.opCnt++
		}
	}
	if moveY > 0 {
		for i := 0; i < moveY; i++ {
			app.Logf("touch: %s", tetris.OpSoftDrop)
			c.tetris.Input(tetris.OpSoftDrop)
			c.opCnt++
		}
	}

	e.PreventDefault()
}

// HandleTouchEnd 处理触摸结束事件
func (c *TouchController) HandleTouchEnd(_ app.Context, e app.Event) {
	touch := c.getTouch(e)
	if touch == nil {
		return
	}
	if c.tetris == nil || c.tetris.State() != tetris.StateRunning {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	switch {
	case int64(c.lastMoveY) > time.Since(c.lastMoveTime).Milliseconds()/5:
		app.Logf("touch: %s", tetris.OpHardDrop)
		c.tetris.Input(tetris.OpHardDrop)
	case int64(-c.lastMoveY) > time.Since(c.lastMoveTime).Milliseconds()/5:
		app.Logf("touch: %s", tetris.OpHold)
		c.tetris.Input(tetris.OpHold)
	case time.Since(c.startTime) < 3*time.Second && c.opCnt == 0:
		app.Logf("touch: %s", tetris.OpRotateRight)
		c.tetris.Input(tetris.OpRotateRight)
	}
}

// getTouch 获取触控信息
func (c *TouchController) getTouch(e app.Event) app.Value {
	changedTouches := e.Get("changedTouches")
	if changedTouches.Length() > 0 {
		return changedTouches.Index(0)
	}
	return nil
}
