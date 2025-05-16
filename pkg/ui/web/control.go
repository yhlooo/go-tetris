package web

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/tetris"
	"github.com/yhlooo/go-tetris/pkg/tetris/common"
)

// handleInput 处理用户输入事件
func (ui *GameUI) handleInput(ctx app.Context, e app.Value) {
	if ui.tetris == nil {
		return
	}

	keyCode := e.Get("key").String()
	app.Logf("key down: %q\n", keyCode)
	switch keyCode {
	case "ArrowUp", "w", "i":
		ui.tetris.Input(tetris.OpRotateRight)
	case "ArrowDown", "s", "k":
		ui.tetris.Input(tetris.OpSoftDrop)
	case "ArrowLeft", "a", "j":
		ui.tetris.Input(tetris.OpMoveLeft)
	case "ArrowRight", "d", "l":
		ui.tetris.Input(tetris.OpMoveRight)
	case " ":
		ui.tetris.Input(tetris.OpHardDrop)
	case "z":
		ui.tetris.Input(tetris.OpRotateLeft)
	case "c":
		ui.tetris.Input(tetris.OpHold)
	case "Enter":
		_ = ui.tetris.Resume()
	case "Escape":
		if ui.tetris.State() == tetris.StateRunning {
			ui.toPaused(ctx)
		} else {
			ui.toGame(ctx)
		}
	}

	ctx.Update()
}

// paintFrameLoop 绘制游戏帧循环
func (ui *GameUI) paintFrameLoop(ctx app.Context, ch <-chan tetris.Frame) {
	for frame := range ch {
		ui.paintFrame(ctx, frame)
	}
}

// paintFrame 绘制帧
func (ui *GameUI) paintFrame(ctx app.Context, frame tetris.Frame) {
	ui.field.UpdateTetriminos(frame.Field.Cells())
	ui.next[0].UpdateTetriminos(newTetriminoGridData(frame.NextTetriminos[0]))
	ui.next[1].UpdateTetriminos(newTetriminoGridData(frame.NextTetriminos[1]))
	ui.next[2].UpdateTetriminos(newTetriminoGridData(frame.NextTetriminos[2]))
	if frame.HoldingTetrimino != nil {
		ui.hold.UpdateTetriminos(newTetriminoGridData(*frame.HoldingTetrimino))
	} else {
		ui.hold.UpdateTetriminos(newTetriminoGridData(common.TetriminoNone))
	}

	ui.score = frame.Score
	ui.level = frame.Level
	ui.clearLines = frame.ClearLines

	if frame.GameOver {
		ui.toGameOver(ctx)
	}

	ctx.Update()
}

// toStartMenu 回到开始菜单
func (ui *GameUI) toStartMenu(_ app.Context) {
	ui.page = ""
	if ui.tetris != nil {
		if err := ui.tetris.Stop(); err != nil {
			app.Logf("stop tetris error: %v", err)
		}
		ui.touchController.SetTetris(nil)
		ui.tetris = nil
	}
}

// toGame 开始或回到游戏
func (ui *GameUI) toGame(ctx app.Context) {
	if ui.tetris == nil {
		ui.tetris = tetris.NewTetris(tetris.DefaultOptions)
		ui.touchController.SetTetris(ui.tetris)
		go ui.paintFrameLoop(ctx, ui.tetris.Frames())
		if err := ui.tetris.Start(ctx); err != nil {
			app.Logf("start tetris error: %v", err)
			return
		}
	}
	if ui.tetris.State() == tetris.StatePaused {
		if err := ui.tetris.Resume(); err != nil {
			app.Logf("resume tetris error: %v", err)
			return
		}
	}

	ui.page = "game"
}

// toGameOver 游戏结束
func (ui *GameUI) toGameOver(_ app.Context) {
	ui.page = "over"
}

// toPaused 暂停
func (ui *GameUI) toPaused(_ app.Context) {
	if ui.tetris != nil {
		if err := ui.tetris.Pause(); err != nil {
			app.Logf("pause tetris error: %v", err)
			return
		}
	}
	ui.page = "paused"
}
