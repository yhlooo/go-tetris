package tetris

import (
	"context"
)

// NewTetris 创建 Tetris 游戏实例
func NewTetris(opts Options) Tetris {
	return &defaultTetris{}
}

// defaultTetris 是 Tetris 的默认实现
type defaultTetris struct{}

var _ Tetris = (*defaultTetris)(nil)

func (t *defaultTetris) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *defaultTetris) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *defaultTetris) Pause(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *defaultTetris) Resume(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *defaultTetris) Input(ctx context.Context, op Op) error {
	//TODO implement me
	panic("implement me")
}

func (t *defaultTetris) Frames() <-chan Frame {
	//TODO implement me
	panic("implement me")
}
