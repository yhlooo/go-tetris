package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/yhlooo/go-tetris/pkg/tetris"
)

func main() {
	holdBox := tview.NewTextView()
	holdBox.SetBorder(true).SetTitle("Hold")
	scoreBox := tview.NewTextView()
	scoreBox.SetBorder(true).SetTitle("Score")
	levelBox := tview.NewTextView()
	levelBox.SetBorder(true).SetTitle("Level")
	linesBox := tview.NewTextView()
	linesBox.SetBorder(true).SetTitle("Lines")
	fieldBox := tview.NewTextView()
	fieldBox.SetDynamicColors(true).SetBorder(true)
	nextBox := tview.NewTextView()
	nextBox.SetBorder(true).SetTitle("Next")

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(holdBox, 6, 1, false).
				AddItem(tview.NewBox(), 3, 1, false).
				AddItem(scoreBox, 3, 1, false).
				AddItem(levelBox, 3, 1, false).
				AddItem(linesBox, 3, 1, false),
			12, 1, false,
		).
		AddItem(fieldBox, 22, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).AddItem(nextBox, 12, 1, false),
			12, 1, false,
		)
	flex.SetRect(0, 0, 46, 22)

	ctx := context.Background()
	var t tetris.Tetris
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if t == nil {
			switch event.Key() {
			case tcell.KeyEnter:
				// 开始游戏
				t = tetris.NewTetris(tetris.DefaultOptions)
				go paintLoop(t.Frames(), app, holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox)
				if err := t.Start(ctx); err != nil {
					log.Fatalf("start tetris error: %v", err)
				}
			}
			return event
		}
		switch t.State() {
		case tetris.StatePending:
			if err := t.Start(ctx); err != nil {
				log.Fatalf("start tetris error: %v", err)
			}
		case tetris.StateRunning:
			switch event.Key() {
			case tcell.KeyEsc:
				// 暂停游戏
				t.Pause(ctx)
			case tcell.KeyUp:
				t.Input(ctx, tetris.OpRotateRight)
			case tcell.KeyDown:
				t.Input(ctx, tetris.OpSoftDrop)
			case tcell.KeyLeft:
				t.Input(ctx, tetris.OpMoveLeft)
			case tcell.KeyRight:
				t.Input(ctx, tetris.OpMoveRight)
			case tcell.KeyRune:
				switch event.Rune() {
				case 'w':
					t.Input(ctx, tetris.OpRotateRight)
				case 'a':
					t.Input(ctx, tetris.OpMoveLeft)
				case 's':
					t.Input(ctx, tetris.OpSoftDrop)
				case 'd':
					t.Input(ctx, tetris.OpMoveRight)
				case 'z':
					t.Input(ctx, tetris.OpRotateLeft)
				case 'c':
					t.Input(ctx, tetris.OpHold)
				case ' ':
					t.Input(ctx, tetris.OpHardDrop)
				}
			default:
			}
		case tetris.StatePaused:
			switch event.Key() {
			case tcell.KeyEnter:
				// 继续游戏
				t.Resume(ctx)
			case tcell.KeyEsc:
				// 结束游戏
				t.Stop(ctx)
				clearScreen(holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox)
				t = nil
			}
		case tetris.StateFinished:
			switch event.Key() {
			case tcell.KeyEnter, tcell.KeyEsc:
				clearScreen(holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox)
				t = nil
			}
		}

		return event
	})

	if err := app.SetRoot(flex, false).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

// paintLoop 绘制画面的循环
func paintLoop(
	ch <-chan tetris.Frame,
	app *tview.Application,
	holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox *tview.TextView,
) {
	for frame := range ch {
		fieldContent := ""
		for i := 19; i >= 0; i-- {
			for j := 0; j < 10; j++ {
				switch frame.Field.BlockWithActiveBlock(i, j) {
				case tetris.BlockNone:
					fieldContent += "  "
				case tetris.BlockI:
					fieldContent += "[:lightcyan]  [:black]"
				case tetris.BlockJ:
					fieldContent += "[:blue]  [:black]"
				case tetris.BlockL:
					fieldContent += "[:orange]  [:black]"
				case tetris.BlockO:
					fieldContent += "[:yellow]  [:black]"
				case tetris.BlockS:
					fieldContent += "[:green]  [:black]"
				case tetris.BlockT:
					fieldContent += "[:purple]  [:black]"
				case tetris.BlockZ:
					fieldContent += "[:red]  [:black]"
				}
			}
		}
		fieldBox.Clear()
		_, _ = fmt.Fprint(fieldBox, fieldContent)

		holdBox.Clear()
		if frame.HoldingBlock != nil {
			_, _ = fmt.Fprintf(holdBox, "%s", frame.HoldingBlock)
		}
		scoreBox.Clear()
		_, _ = fmt.Fprintf(scoreBox, "%d", frame.Score)
		levelBox.Clear()
		_, _ = fmt.Fprintf(levelBox, "%d", frame.Level)
		linesBox.Clear()
		_, _ = fmt.Fprintf(linesBox, "%d", frame.ClearLines)
		nextBox.Clear()
		_, _ = fmt.Fprintf(nextBox, "%s", frame.NextBlock)

		app.Draw()
	}
}

// clearScreen 清空画面
func clearScreen(holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox *tview.TextView) {
	holdBox.Clear()
	scoreBox.Clear()
	levelBox.Clear()
	linesBox.Clear()
	nextBox.Clear()
	fieldBox.Clear()
}
