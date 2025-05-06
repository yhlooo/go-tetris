package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/ui/web"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	app.Route("/", func() app.Composer {
		return web.NewGameUI()
	})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name: "Tetris",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
