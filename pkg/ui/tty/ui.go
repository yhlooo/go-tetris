package tty

import (
	"context"
	"fmt"

	"github.com/bombsimon/logrusr/v4"
	"github.com/gdamore/tcell/v2"
	"github.com/go-logr/logr"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"

	"github.com/yhlooo/go-tetris/pkg/tetris"
	"github.com/yhlooo/go-tetris/pkg/tetris/common"
)

// NewGameUI 创建 GameUI
func NewGameUI() *GameUI {
	return &GameUI{}
}

// GameUI 基于终端的游戏用户交互界面
type GameUI struct {
	app                                             *tview.Application
	pages                                           *tview.Pages
	holdBox, scoreBox, levelBox, linesBox, stateBox *tview.TextView
	fieldBox                                        *tview.TextView
	nextBox                                         *tview.TextView
	logBox                                          *tview.TextView
	gameOverBox                                     *tview.TextView

	tetris       tetris.Tetris
	logrusLogger *logrus.Logger
	logger       logr.Logger
}

// Run 初始化并开始运行
func (ui *GameUI) Run() error {
	root := ui.newRoot()
	ui.app = tview.NewApplication().SetRoot(root, true).SetFocus(root)

	return ui.app.Run()
}

// newRoot 创建根元素
func (ui *GameUI) newRoot() tview.Primitive {
	ui.pages = tview.NewPages().
		AddPage("help", ui.newHelpPage(), true, false).
		AddPage("about", ui.newAboutPage(), true, false).
		AddPage("main", ui.newMainPage(), true, true).
		AddPage("pause", ui.newPauseMenuPage(), true, false).
		AddPage("menu", ui.newMainMenuPage(), true, true).
		AddPage("over", ui.newGameOverPage(), true, false)

	return tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(ui.pages, 46, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)
}

// newMainPage 创建主页
func (ui *GameUI) newMainPage() tview.Primitive {
	ui.holdBox = tview.NewTextView()
	ui.holdBox.SetDynamicColors(true).SetBorder(true).SetTitle("Hold")

	ui.scoreBox = tview.NewTextView()
	ui.scoreBox.SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("Score")

	ui.levelBox = tview.NewTextView()
	ui.levelBox.SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("Level")

	ui.linesBox = tview.NewTextView()
	ui.linesBox.SetTextAlign(tview.AlignCenter).SetBorder(true).SetTitle("Lines")

	ui.stateBox = tview.NewTextView()
	ui.stateBox.SetTextAlign(tview.AlignCenter).SetDynamicColors(true)

	ui.fieldBox = tview.NewTextView()
	ui.fieldBox.SetDynamicColors(true).SetBorder(true).SetInputCapture(ui.handleGameInput)

	ui.nextBox = tview.NewTextView()
	ui.nextBox.SetDynamicColors(true).SetBorder(true).SetTitle("Next")

	ui.logBox = tview.NewTextView().SetScrollable(true).SetDynamicColors(true)
	ui.logrusLogger = logrus.New()
	ui.logrusLogger.Out = ui.logBox
	ui.logrusLogger.Formatter = &logFormatter{}
	ui.logger = logrusr.New(ui.logrusLogger)
	ui.logBox.SetChangedFunc(func() {
		ui.logBox.ScrollToEnd()
	})

	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.holdBox, 6, 1, false).
		AddItem(tview.NewBox(), 3, 1, false).
		AddItem(ui.scoreBox, 3, 1, false).
		AddItem(ui.levelBox, 3, 1, false).
		AddItem(ui.linesBox, 3, 1, false).
		AddItem(ui.stateBox, 0, 1, false)
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.nextBox, 12, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)

	mainPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				AddItem(leftFlex, 12, 1, false).
				AddItem(ui.fieldBox, 22, 1, true).
				AddItem(rightFlex, 12, 1, false),
			22, 1, true,
		).
		AddItem(ui.logBox, 0, 1, false)

	return mainPage
}

// newMainMenuPage 创建主菜单页
func (ui *GameUI) newMainMenuPage() tview.Primitive {
	mainMenu := tview.NewTable().SetSelectable(true, true).
		SetCell(0, 0, tview.NewTableCell("   Play   ").SetAlign(tview.AlignCenter)).
		SetCell(1, 0, tview.NewTableCell("   Help   ").SetAlign(tview.AlignCenter)).
		SetCell(2, 0, tview.NewTableCell("  !About  ").SetAlign(tview.AlignCenter))
	mainMenu.SetBorder(true)
	mainMenu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
		default:
			return event
		}
		// 切换页面
		row, _ := mainMenu.GetSelection()
		switch row {
		case 0:
			// 开始游戏
			ui.startGame()
		case 1:
			ui.pages.SwitchToPage("help")
		case 2:
			ui.pages.SwitchToPage("about")
		}
		return event
	})
	mainMenuPage := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(mainMenu, 5, 1, true)
	mainMenuPage.SetBorderPadding(8, 0, 17, 17)

	return mainMenuPage
}

// newPauseMenuPage 创建暂停菜单页
func (ui *GameUI) newPauseMenuPage() tview.Primitive {
	menu := tview.NewTable().SetSelectable(true, true).
		SetCell(0, 0, tview.NewTableCell("  Resume  ").SetAlign(tview.AlignCenter)).
		SetCell(1, 0, tview.NewTableCell("   Help   ").SetAlign(tview.AlignCenter)).
		SetCell(2, 0, tview.NewTableCell("   Quit   ").SetAlign(tview.AlignCenter))
	menu.SetBorder(true)
	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
		case tcell.KeyEsc:
			// 继续游戏
			ui.resumeGame()
		default:
			return event
		}
		// 切换页面
		row, _ := menu.GetSelection()
		switch row {
		case 0:
			// 继续游戏
			ui.resumeGame()
		case 1:
			ui.pages.SwitchToPage("help")
		case 2:
			// 结束游戏
			ui.stopGame()
		}
		return event
	})
	menuPage := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(menu, 5, 1, true)
	menuPage.SetBorderPadding(8, 0, 17, 17)

	return menuPage
}

// newGameOverPage 创建游戏结束页
func (ui *GameUI) newGameOverPage() tview.Primitive {
	ui.gameOverBox = tview.NewTextView()
	ui.gameOverBox.
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBackgroundColor(tcell.ColorBlue).
		SetBorder(true).
		SetTitle("Game Over").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEnter, tcell.KeyEsc:
				ui.stopGame()
			default:
			}
			return event
		})

	gameOverPage := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(ui.gameOverBox, 6, 1, true)
	gameOverPage.SetBorderPadding(8, 0, 0, 0)
	return gameOverPage
}

// newHelpPage 创建帮助页
func (ui *GameUI) newHelpPage() tview.Primitive {
	helpBox := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[black:lightgray]                  Control                   [white:black]
          Up / w / i : Rotate Right
        Left / a / j : Move Left
       Right / d / l : Move Right
        Down / s / k : Soft Drop
                   z : Rotate Left
                   c : Hold
               Space : Hard Drop
                 ESC : Pause
[black:lightgray]                   Debug                    [white:black]
                   X : On/Off Debug Mode
       O/I/J/L/S/T/Z : Change Tetrimino

[black:lightgray]                   Score                    [white:black]
    Soft Dro                1 * Distance
    Hard Drop               2 * Distance
    Single Line Clear                100
    Double Line Clear                300
    Triple Line Clear                500
    Tetris (4 Line Clear)            800
    T-Spin                           400
    T-Spin Single                    800
    T-Spin Double                   1200
    T-Spin Triple                   1600
    Back-to-Back            0.5 * Tetris
                               or T-Spin

   [lightgray](Press ENTER or ESC to back to menu)[white]
`)
	helpBox.SetBorder(true).SetTitle("Help").SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 回到主页
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEsc:
			ui.pages.SwitchToPage("main")
			if ui.tetris != nil {
				ui.pages.ShowPage("pause")
			} else {
				ui.pages.ShowPage("menu")
			}
		default:
		}
		return event
	})
	return tview.NewFlex().SetDirection(tview.FlexRow).AddItem(helpBox, 30, 1, true)
}

// newAboutPage 创建关于页
func (ui *GameUI) newAboutPage() tview.Primitive {
	aboutBox := tview.NewTextView().SetDynamicColors(true).
		SetText(`Tetris is the addictive puzzle game created by Alexey Pajitnov in 1984. In the decades to follow, Tetris became one of the most successful and recognizable video games, appearing on nearly every gaming platform available.

This version is an open source implementation of Tetris, created by yhlooo in 2025,
see https://github.com/yhlooo/go-tetris .








   [lightgray](Press ENTER or ESC to back to menu)[white]
`)
	aboutBox.
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle("About").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// 回到主页
			switch event.Key() {
			case tcell.KeyEnter, tcell.KeyEsc:
				ui.pages.SwitchToPage("main")
				ui.pages.ShowPage("menu")
			default:
			}
			return event
		})
	return tview.NewFlex().SetDirection(tview.FlexRow).AddItem(aboutBox, 22, 1, true)
}

// startGame 开始游戏
func (ui *GameUI) startGame() {
	ui.logrusLogger.SetLevel(logrus.InfoLevel)
	opts := tetris.DefaultOptions
	opts.Logger = ui.logger
	ui.tetris = tetris.NewTetris(opts)
	go ui.paintGameLoop(ui.tetris.Frames())
	if err := ui.tetris.Start(context.Background()); err != nil {
		ui.logger.Error(err, "start tetris error")
		return
	}
	ui.pages.SwitchToPage("main")
}

// pauseGame 暂停游戏
func (ui *GameUI) pauseGame() {
	_ = ui.tetris.Pause()
	if !ui.tetris.Debug() {
		// 清空显示
		ui.holdBox.Clear()
		ui.nextBox.Clear()
		ui.fieldBox.Clear()
		ui.pages.ShowPage("pause")
	}
}

// resumeGame 继续游戏
func (ui *GameUI) resumeGame() {
	_ = ui.tetris.Resume()
	ui.paintGameFrame(ui.tetris.CurrentFrame())
	ui.pages.SwitchToPage("main")
}

// stopGame 结束游戏
func (ui *GameUI) stopGame() {
	_ = ui.tetris.Stop()
	ui.tetris = nil
	ui.clearGameInfo()
	ui.pages.SwitchToPage("main")
	ui.pages.ShowPage("menu")
}

// handleGameInput 处理游戏输入
func (ui *GameUI) handleGameInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		// 继续游戏
		ui.resumeGame()
	case tcell.KeyEsc:
		// 暂停/继续游戏
		if ui.tetris.State() == tetris.StatePaused {
			ui.resumeGame()
		} else {
			ui.pauseGame()
		}
	case tcell.KeyUp:
		ui.tetris.Input(tetris.OpRotateRight)
	case tcell.KeyDown:
		ui.tetris.Input(tetris.OpSoftDrop)
	case tcell.KeyLeft:
		ui.tetris.Input(tetris.OpMoveLeft)
	case tcell.KeyRight:
		ui.tetris.Input(tetris.OpMoveRight)
	case tcell.KeyRune:
		switch event.Rune() {
		case 'w', 'i':
			ui.tetris.Input(tetris.OpRotateRight)
		case 'a', 'j':
			ui.tetris.Input(tetris.OpMoveLeft)
		case 's', 'k':
			ui.tetris.Input(tetris.OpSoftDrop)
		case 'd', 'l':
			ui.tetris.Input(tetris.OpMoveRight)
		case 'z':
			ui.tetris.Input(tetris.OpRotateLeft)
		case 'c':
			ui.tetris.Input(tetris.OpHold)
		case ' ':
			ui.tetris.Input(tetris.OpHardDrop)
		case 'X':
			ui.tetris.SetDebug(!ui.tetris.Debug())
			ui.stateBox.Clear()
			if ui.tetris.Debug() {
				ui.logrusLogger.SetLevel(logrus.DebugLevel)
				_, _ = fmt.Fprint(ui.stateBox, "[red]DEBUG MODE[black]")
			} else {
				ui.logrusLogger.SetLevel(logrus.InfoLevel)
			}
		case 'I':
			_ = ui.tetris.ChangeActiveTetriminoType(common.I)
		case 'J':
			_ = ui.tetris.ChangeActiveTetriminoType(common.J)
		case 'L':
			_ = ui.tetris.ChangeActiveTetriminoType(common.L)
		case 'O':
			_ = ui.tetris.ChangeActiveTetriminoType(common.O)
		case 'S':
			_ = ui.tetris.ChangeActiveTetriminoType(common.S)
		case 'T':
			_ = ui.tetris.ChangeActiveTetriminoType(common.T)
		case 'Z':
			_ = ui.tetris.ChangeActiveTetriminoType(common.Z)
		}
	default:
	}

	return event
}

// paintGameLoop 绘制游戏画面的循环
func (ui *GameUI) paintGameLoop(ch <-chan tetris.Frame) {
	for frame := range ch {
		ui.paintGameFrame(frame)
		ui.app.Draw()
	}
}

// paintGameFrame 绘制游戏一帧
func (ui *GameUI) paintGameFrame(frame tetris.Frame) {
	fieldContent := ""
	cells := frame.Field.Cells()
	for i := 19; i >= 0; i-- {
		for j := 0; j < 10; j++ {
			if cells[i][j].Shadow {
				switch cells[i][j].Type {
				case common.TetriminoNone:
					fieldContent += "  "
				case common.I:
					fieldContent += "[darkcyan]..[black]"
				case common.J:
					fieldContent += "[blue]..[black]"
				case common.L:
					fieldContent += "[darkorange]..[black]"
				case common.O:
					fieldContent += "[orange]..[black]"
				case common.S:
					fieldContent += "[lightgreen]..[black]"
				case common.T:
					fieldContent += "[mediumpurple]..[black]"
				case common.Z:
					fieldContent += "[red]..[black]"
				}
			} else {
				switch cells[i][j].Type {
				case common.TetriminoNone:
					fieldContent += "  "
				case common.I:
					fieldContent += "[:darkcyan]  [:black]"
				case common.J:
					fieldContent += "[:blue]  [:black]"
				case common.L:
					fieldContent += "[:darkorange]  [:black]"
				case common.O:
					fieldContent += "[:orange]  [:black]"
				case common.S:
					fieldContent += "[:lightgreen]  [:black]"
				case common.T:
					fieldContent += "[:mediumpurple]  [:black]"
				case common.Z:
					fieldContent += "[:red]  [:black]"
				}
			}

		}
	}
	ui.fieldBox.Clear()
	_, _ = fmt.Fprint(ui.fieldBox, fieldContent)

	ui.holdBox.Clear()
	if frame.HoldingTetrimino != nil {
		_, _ = fmt.Fprint(ui.holdBox, paintTetrisTetrimino(*frame.HoldingTetrimino))
	}
	ui.scoreBox.Clear()
	_, _ = fmt.Fprintf(ui.scoreBox, "%d", frame.Score)
	ui.levelBox.Clear()
	_, _ = fmt.Fprintf(ui.levelBox, "%d", frame.Level)
	ui.linesBox.Clear()
	_, _ = fmt.Fprintf(ui.linesBox, "%d", frame.ClearLines)
	ui.nextBox.Clear()
	for _, b := range frame.NextTetriminos {
		_, _ = fmt.Fprint(ui.nextBox, paintTetrisTetrimino(b))
	}

	// 游戏结束
	if frame.GameOver {
		ui.gameOverBox.SetText(fmt.Sprintf(
			"\nScore: %d\n\n[lightgray](Press ENTER or ESC to continue)[white]",
			frame.Score,
		))
		ui.pages.ShowPage("over")
	}
}

// clearGameInfo 清除画面中的游戏信息
func (ui *GameUI) clearGameInfo() {
	ui.holdBox.Clear()
	ui.scoreBox.Clear()
	ui.levelBox.Clear()
	ui.linesBox.Clear()
	ui.stateBox.Clear()
	ui.nextBox.Clear()
	ui.fieldBox.Clear()
}

// paintTetrisTetrimino 绘制方块
func paintTetrisTetrimino(tetriminoType common.TetriminoType) string {
	switch tetriminoType {
	case common.TetriminoNone:
	case common.I:
		return `

 [:darkcyan]        [:black]
`
	case common.J:
		return `
  [:blue]  [:black]
  [:blue]      [:black]
`
	case common.L:
		return `
      [:darkorange]  [:black]
  [:darkorange]      [:black]
`
	case common.O:
		return `
   [:orange]    [:black]
   [:orange]    [:black]
`
	case common.S:
		return `
    [:lightgreen]    [:black]
  [:lightgreen]    [:black]
`
	case common.T:
		return `
    [:mediumpurple]  [:black]
  [:mediumpurple]      [:black]
`
	case common.Z:
		return `
  [:red]    [:black]
    [:red]    [:black]
`
	}
	return ""
}

// logFormatter 日志格式化器
type logFormatter struct{}

var _ logrus.Formatter = (*logFormatter)(nil)

// Format 格式化日志
func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(f.level(entry.Level) + entry.Message + "\n"), nil
}

// level 返回日志级别标识
func (f *logFormatter) level(l logrus.Level) string {
	switch l {
	case logrus.TraceLevel:
		return "[gray]T[white] "
	case logrus.DebugLevel:
		return "[green]D[white] "
	case logrus.InfoLevel:
		return "[blue]I[white] "
	case logrus.WarnLevel:
		return "[darkorange]W[white] "
	case logrus.ErrorLevel:
		return "[red]E[white] "
	case logrus.FatalLevel:
		return "[:red]FATAL[:black] "
	case logrus.PanicLevel:
		return "[:red]PANIC[:black] "
	}
	return "? "
}
