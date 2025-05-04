package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/yhlooo/go-tetris/pkg/tetris"
)

func main() {
	holdBox := tview.NewTextView()
	holdBox.SetDynamicColors(true).SetBorder(true).SetTitle("Hold")
	scoreBox := tview.NewTextView()
	scoreBox.SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("Score")
	levelBox := tview.NewTextView()
	levelBox.SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("Level")
	linesBox := tview.NewTextView()
	linesBox.SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("Lines")
	fieldBox := tview.NewTextView()
	fieldBox.SetDynamicColors(true).SetBorder(true)
	nextBox := tview.NewTextView()
	nextBox.SetDynamicColors(true).SetBorder(true).SetTitle("Next")
	stateBox := tview.NewTextView()
	stateBox.SetTextAlign(tview.AlignCenter).SetDynamicColors(true)

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(holdBox, 6, 1, false).
				AddItem(tview.NewBox(), 3, 1, false).
				AddItem(scoreBox, 3, 1, false).
				AddItem(levelBox, 3, 1, false).
				AddItem(linesBox, 3, 1, false).
				AddItem(stateBox, 3, 1, false),
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
		case tetris.StateRunning, tetris.StatePaused:
			switch event.Key() {
			case tcell.KeyEnter:
				// 继续游戏
				t.Resume(ctx)
				paintFrame(t.CurrentFrame(), holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox)
			case tcell.KeyEsc:
				if t.State() == tetris.StatePaused {
					// 结束游戏
					t.Stop(ctx)
					clearScreen(holdBox, scoreBox, levelBox, linesBox, stateBox, nextBox, fieldBox)
					t = nil
				} else {
					// 暂停游戏
					t.Pause(ctx)
					if !t.Debug() {
						paintPause(holdBox, nextBox, fieldBox)
					}
				}
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
				case 'w', 'i':
					t.Input(ctx, tetris.OpRotateRight)
				case 'a', 'j':
					t.Input(ctx, tetris.OpMoveLeft)
				case 's', 'k':
					t.Input(ctx, tetris.OpSoftDrop)
				case 'd', 'l':
					t.Input(ctx, tetris.OpMoveRight)
				case 'z':
					t.Input(ctx, tetris.OpRotateLeft)
				case 'c':
					t.Input(ctx, tetris.OpHold)
				case ' ':
					t.Input(ctx, tetris.OpHardDrop)
				case 'X':
					t.SetDebug(!t.Debug())
					stateBox.Clear()
					if t.Debug() {
						_, _ = fmt.Fprint(stateBox, "[red]DEBUG MODE[black]")
					}
				case 'I':
					_ = t.ChangeActiveBlockType(tetris.BlockI)
				case 'J':
					_ = t.ChangeActiveBlockType(tetris.BlockJ)
				case 'L':
					_ = t.ChangeActiveBlockType(tetris.BlockL)
				case 'O':
					_ = t.ChangeActiveBlockType(tetris.BlockO)
				case 'S':
					_ = t.ChangeActiveBlockType(tetris.BlockS)
				case 'T':
					_ = t.ChangeActiveBlockType(tetris.BlockT)
				case 'Z':
					_ = t.ChangeActiveBlockType(tetris.BlockZ)
				}
			default:
			}
		case tetris.StateFinished:
			switch event.Key() {
			case tcell.KeyEnter, tcell.KeyEsc:
				clearScreen(holdBox, scoreBox, levelBox, linesBox, stateBox, nextBox, fieldBox)
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
		paintFrame(frame, holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox)
		app.Draw()
	}
}

// paintFrame 绘制一帧
func paintFrame(
	frame tetris.Frame,
	holdBox, scoreBox, levelBox, linesBox, nextBox, fieldBox *tview.TextView,
) {
	fieldContent := ""
	for i := 19; i >= 0; i-- {
		for j := 0; j < 10; j++ {
			switch frame.Field.BlockWithActiveBlock(i, j) {
			case tetris.BlockNone:
				fieldContent += "  "
			case tetris.BlockI:
				fieldContent += "[:darkcyan]  [:black]"
			case tetris.BlockJ:
				fieldContent += "[:blue]  [:black]"
			case tetris.BlockL:
				fieldContent += "[:darkorange]  [:black]"
			case tetris.BlockO:
				fieldContent += "[:orange]  [:black]"
			case tetris.BlockS:
				fieldContent += "[:lightgreen]  [:black]"
			case tetris.BlockT:
				fieldContent += "[:mediumpurple]  [:black]"
			case tetris.BlockZ:
				fieldContent += "[:red]  [:black]"
			}
		}
	}
	fieldBox.Clear()
	_, _ = fmt.Fprint(fieldBox, fieldContent)

	holdBox.Clear()
	if frame.HoldingBlock != nil {
		_, _ = fmt.Fprint(holdBox, paintBlock(*frame.HoldingBlock))
	}
	scoreBox.Clear()
	_, _ = fmt.Fprintf(scoreBox, "%d", frame.Score)
	levelBox.Clear()
	_, _ = fmt.Fprintf(levelBox, "%d", frame.Level)
	linesBox.Clear()
	_, _ = fmt.Fprintf(linesBox, "%d", frame.ClearLines)
	nextBox.Clear()
	for _, b := range frame.NextBlocks {
		_, _ = fmt.Fprint(nextBox, paintBlock(b))
	}
}

// paintPause 绘制暂停画面
func paintPause(holdBox, nextBox, fieldBox *tview.TextView) {
	holdBox.Clear()
	nextBox.Clear()
	fieldBox.Clear()
	_, _ = fmt.Fprint(fieldBox, strings.Repeat("\n", 10)+"     - PAUSED -")
}

// paintBlock 绘制方块
func paintBlock(blockType tetris.BlockType) string {
	switch blockType {
	case tetris.BlockNone:
	case tetris.BlockI:
		return `

 [:darkcyan]        [:black]
`
	case tetris.BlockJ:
		return `
  [:blue]  [:black]
  [:blue]      [:black]
`
	case tetris.BlockL:
		return `
      [:darkorange]  [:black]
  [:darkorange]      [:black]
`
	case tetris.BlockO:
		return `
   [:orange]    [:black]
   [:orange]    [:black]
`
	case tetris.BlockS:
		return `
    [:lightgreen]    [:black]
  [:lightgreen]    [:black]
`
	case tetris.BlockT:
		return `
    [:mediumpurple]  [:black]
  [:mediumpurple]      [:black]
`
	case tetris.BlockZ:
		return `
  [:red]    [:black]
    [:red]    [:black]
`
	}
	return ""
}

// clearScreen 清空画面
func clearScreen(holdBox, scoreBox, levelBox, linesBox, stateBox, nextBox, fieldBox *tview.TextView) {
	holdBox.Clear()
	scoreBox.Clear()
	levelBox.Clear()
	linesBox.Clear()
	stateBox.Clear()
	nextBox.Clear()
	fieldBox.Clear()
}
