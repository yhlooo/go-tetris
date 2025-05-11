package web

import (
	"fmt"
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// renderMain 渲染主要内容
func (ui *GameUI) renderMain() app.UI {
	return app.Div().Class("tetris-main").Body(
		ui.renderGame(),
		ui.renderHelp(),
		ui.renderAbout(),
	)
}

// renderGame 渲染游戏界面
func (ui *GameUI) renderGame() app.UI {
	fieldWidth := 100
	fieldHeight := 100
	if ui.field != nil {
		fieldWidth, fieldHeight = ui.field.Size()
	}
	return app.Div().Class("tetris-game").Body(
		app.Div().Class("tetris-game-sidebar tetris-game-sidebar-left").Body(
			app.Div().Class("tetris-block-booth").Body(
				app.Div().Class("tetris-game-sub-title").Text("HOLD"),
				app.Div().Class("tetris-block").Body(ui.hold),
			),
			app.Div().Class("tetris-score-box").Body(
				app.Div().Body(
					app.Div().Class("tetris-game-sub-title").Text("SCORE"),
					app.Div().Text(strconv.Itoa(ui.score)),
				),
				app.Div().Body(
					app.Div().Class("tetris-game-sub-title").Text("LEVEL"),
					app.Div().Text(strconv.Itoa(ui.level)),
				),
				app.Div().Body(
					app.Div().Class("tetris-game-sub-title").Text("LINES"),
					app.Div().Text(strconv.Itoa(ui.clearLines)),
				),
			),
		),
		app.Div().Class("tetris-game-field").Body(
			app.If(ui.page == "", func() app.UI {
				return app.Div().Class("tetris-game-menu").Body(
					app.Button().Text("Start").OnClick(func(ctx app.Context, _ app.Event) { ui.toGame(ctx) }),
				)
			}).ElseIf(ui.page == "paused", func() app.UI {
				return app.Div().Class("tetris-game-menu").Body(
					app.Button().Text("Resume").OnClick(func(ctx app.Context, _ app.Event) { ui.toGame(ctx) }),
					app.Button().Text("Quit").OnClick(func(ctx app.Context, _ app.Event) { ui.toStartMenu(ctx) }),
				)
			}).ElseIf(ui.page == "over", func() app.UI {
				return app.Div().Class("tetris-game-menu").Body(
					app.Div().Class("tetris-game-sub-title").Text("Game Over"),
					app.Div().Text(fmt.Sprintf("Score: %d", ui.score)).Style("margin-bottom", "15px"),
					app.Button().Text("Ok").OnClick(func(ctx app.Context, _ app.Event) { ui.toStartMenu(ctx) }),
				)
			}).Else(func() app.UI {
				return ui.field
			}),
		).Styles(map[string]string{
			"width":  strconv.Itoa(fieldWidth) + "px",
			"height": strconv.Itoa(fieldHeight) + "px",
		}),
		app.Div().Class("tetris-game-sidebar tetris-game-sidebar-right").Body(
			app.Div().Class("tetris-block-booth").Body(
				app.Div().Class("tetris-game-sub-title").Text("NEXT"),
				app.Range(ui.next).Slice(func(i int) app.UI {
					return app.Div().Class("tetris-block").Body(ui.next[i])
				}),
			),
		),
	)
}

// renderHelp 渲染帮助信息
func (ui *GameUI) renderHelp() app.UI {
	return app.Div().Class("tetris-help tetris-page").Body(
		app.H2().Text("Help"),
		app.P().Body(
			app.Text("Up / w / i : Rotate Right"), app.Br(),
			app.Text("Left / a / j : Move Left"), app.Br(),
			app.Text("Right / d / l : Move Right"), app.Br(),
			app.Text("Down / s / k : Soft Drop"), app.Br(),
			app.Text("z : Rotate Left"), app.Br(),
			app.Text("c : Hold"), app.Br(),
			app.Text("Space : Hard Drop"), app.Br(),
			app.Text("ESC : Pause"),
		),
		app.Button().Text("Ok"),
	)
}

// renderAbout 渲染关于信息
func (ui *GameUI) renderAbout() app.UI {
	return app.Div().Class("tetris-about tetris-page").Body(
		app.H2().Text("About"),
		app.P().Text(`Tetris is the addictive puzzle game created by Alexey Pajitnov in 1984.
In the decades to follow, Tetris became one of the most successful and recognizable video games,
appearing on nearly every gaming platform available.`),
		app.P().Body(
			app.Text("This version is an open source implementation of Tetris, created by yhlooo in 2025, see "),
			app.A().Text("https://github.com/yhlooo/go-tetris").Href("https://github.com/yhlooo/go-tetris"),
			app.Text(" ."),
		),
		app.Button().Text("Ok"),
	)
}
