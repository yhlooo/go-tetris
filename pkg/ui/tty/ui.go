package tty

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/go-logr/logr"
	"github.com/rivo/tview"

	"github.com/yhlooo/go-tetris/pkg/tetris"
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

	gameCtx context.Context
	tetris  tetris.Tetris
	logger  logr.Logger
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
		AddPage("options", ui.newOptionsPage(), true, false).
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

	ui.logBox = tview.NewTextView()

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
		SetCell(0, 0, tview.NewTableCell("  Play    ").SetAlign(tview.AlignCenter)).
		SetCell(1, 0, tview.NewTableCell(" Options  ").SetAlign(tview.AlignCenter)).
		SetCell(2, 0, tview.NewTableCell("  Help    ").SetAlign(tview.AlignCenter)).
		SetCell(3, 0, tview.NewTableCell("  About   ").SetAlign(tview.AlignCenter))
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
			ui.pages.SwitchToPage("options")
		case 2:
			ui.pages.SwitchToPage("help")
		case 3:
			ui.pages.SwitchToPage("about")
		}
		return event
	})
	mainMenuPage := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(mainMenu, 6, 1, true)
	mainMenuPage.SetBorderPadding(8, 0, 17, 17)

	return mainMenuPage
}

// newPauseMenuPage 创建暂停菜单页
func (ui *GameUI) newPauseMenuPage() tview.Primitive {
	menu := tview.NewTable().SetSelectable(true, true).
		SetCell(0, 0, tview.NewTableCell(" Resume   ").SetAlign(tview.AlignCenter)).
		SetCell(1, 0, tview.NewTableCell(" Options  ").SetAlign(tview.AlignCenter)).
		SetCell(2, 0, tview.NewTableCell("  Help    ").SetAlign(tview.AlignCenter)).
		SetCell(3, 0, tview.NewTableCell("  Quit    ").SetAlign(tview.AlignCenter))
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
			ui.pages.SwitchToPage("options")
		case 2:
			ui.pages.SwitchToPage("help")
		case 3:
			// 结束游戏
			ui.stopGame()
		}
		return event
	})
	menuPage := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(menu, 6, 1, true)
	menuPage.SetBorderPadding(8, 0, 17, 17)

	return menuPage
}

// newGameOverPage 创建游戏结束页
func (ui *GameUI) newGameOverPage() tview.Primitive {
	ui.gameOverBox = tview.NewTextView()
	ui.gameOverBox.
		SetTextAlign(tview.AlignCenter).
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

// newOptionsPage 创建选项页
func (ui *GameUI) newOptionsPage() tview.Primitive {
	optionsBox := tview.NewTextView()
	optionsBox.SetBorder(true).SetTitle("Options").SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 回到主页
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEsc:
			ui.pages.ShowPage("main")
			if ui.tetris != nil {
				ui.pages.ShowPage("pause")
			} else {
				ui.pages.ShowPage("menu")
			}
		default:
		}
		return event
	})
	return tview.NewFlex().SetDirection(tview.FlexRow).AddItem(optionsBox, 22, 1, true)
}

// newHelpPage 创建帮助页
func (ui *GameUI) newHelpPage() tview.Primitive {
	helpBox := tview.NewTextView()
	helpBox.SetBorder(true).SetTitle("Help").SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 回到主页
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEsc:
			ui.pages.ShowPage("main")
			if ui.tetris != nil {
				ui.pages.ShowPage("pause")
			} else {
				ui.pages.ShowPage("menu")
			}
		default:
		}
		return event
	})
	return tview.NewFlex().SetDirection(tview.FlexRow).AddItem(helpBox, 22, 1, true)
}

// newAboutPage 创建关于页
func (ui *GameUI) newAboutPage() tview.Primitive {
	aboutBox := tview.NewTextView()
	aboutBox.SetBorder(true).SetTitle("About").SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 回到主页
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEsc:
			ui.pages.ShowPage("main")
			ui.pages.ShowPage("menu")
		default:
		}
		return event
	})
	return tview.NewFlex().SetDirection(tview.FlexRow).AddItem(aboutBox, 22, 1, true)
}

// startGame 开始游戏
func (ui *GameUI) startGame() {
	ui.gameCtx = logr.NewContext(context.Background(), ui.logger)
	ui.tetris = tetris.NewTetris(tetris.DefaultOptions)
	go ui.paintGameLoop(ui.tetris.Frames())
	if err := ui.tetris.Start(ui.gameCtx); err != nil {
		ui.logger.Error(err, "start tetris error")
		return
	}
	ui.pages.SwitchToPage("main")
}

// pauseGame 暂停游戏
func (ui *GameUI) pauseGame() {
	ui.tetris.Pause(ui.gameCtx)
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
	ui.tetris.Resume(ui.gameCtx)
	ui.paintGameFrame(ui.tetris.CurrentFrame())
	ui.pages.SwitchToPage("main")
}

// stopGame 结束游戏
func (ui *GameUI) stopGame() {
	ui.tetris.Stop(ui.gameCtx)
	ui.gameCtx = nil
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
		ui.tetris.Input(ui.gameCtx, tetris.OpRotateRight)
	case tcell.KeyDown:
		ui.tetris.Input(ui.gameCtx, tetris.OpSoftDrop)
	case tcell.KeyLeft:
		ui.tetris.Input(ui.gameCtx, tetris.OpMoveLeft)
	case tcell.KeyRight:
		ui.tetris.Input(ui.gameCtx, tetris.OpMoveRight)
	case tcell.KeyRune:
		switch event.Rune() {
		case 'w', 'i':
			ui.tetris.Input(ui.gameCtx, tetris.OpRotateRight)
		case 'a', 'j':
			ui.tetris.Input(ui.gameCtx, tetris.OpMoveLeft)
		case 's', 'k':
			ui.tetris.Input(ui.gameCtx, tetris.OpSoftDrop)
		case 'd', 'l':
			ui.tetris.Input(ui.gameCtx, tetris.OpMoveRight)
		case 'z':
			ui.tetris.Input(ui.gameCtx, tetris.OpRotateLeft)
		case 'c':
			ui.tetris.Input(ui.gameCtx, tetris.OpHold)
		case ' ':
			ui.tetris.Input(ui.gameCtx, tetris.OpHardDrop)
		case 'X':
			ui.tetris.SetDebug(!ui.tetris.Debug())
			ui.stateBox.Clear()
			if ui.tetris.Debug() {
				_, _ = fmt.Fprint(ui.stateBox, "[red]DEBUG MODE[black]")
			}
		case 'I':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockI)
		case 'J':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockJ)
		case 'L':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockL)
		case 'O':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockO)
		case 'S':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockS)
		case 'T':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockT)
		case 'Z':
			_ = ui.tetris.ChangeActiveBlockType(tetris.BlockZ)
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
	ui.fieldBox.Clear()
	_, _ = fmt.Fprint(ui.fieldBox, fieldContent)

	ui.holdBox.Clear()
	if frame.HoldingBlock != nil {
		_, _ = fmt.Fprint(ui.holdBox, paintTetrisBlock(*frame.HoldingBlock))
	}
	ui.scoreBox.Clear()
	_, _ = fmt.Fprintf(ui.scoreBox, "%d", frame.Score)
	ui.levelBox.Clear()
	_, _ = fmt.Fprintf(ui.levelBox, "%d", frame.Level)
	ui.linesBox.Clear()
	_, _ = fmt.Fprintf(ui.linesBox, "%d", frame.ClearLines)
	ui.nextBox.Clear()
	for _, b := range frame.NextBlocks {
		_, _ = fmt.Fprint(ui.nextBox, paintTetrisBlock(b))
	}

	// 游戏结束
	if frame.GameOver {
		ui.gameOverBox.SetText(fmt.Sprintf("\nScore: %d\n\n(Press ENTER or ESC to continue)", frame.Score))
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

// paintTetrisBlock 绘制方块
func paintTetrisBlock(blockType tetris.BlockType) string {
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
