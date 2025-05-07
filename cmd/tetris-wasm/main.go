package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/ui/web"
)

var (
	static bool
)

func init() {
	flag.BoolVar(&static, "static", false, "generate static files")
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	flag.Parse()

	app.Route("/", func() app.Composer {
		return web.NewGameUI()
	})
	app.RunWhenOnBrowser()

	if static {
		log.Println("generate static files to ./web")
		if err := app.GenerateStaticWebsite("web", &app.Handler{
			Name: "Tetris",
		}); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("serving http on :8000")
		http.Handle("/", &app.Handler{
			Name: "Tetris",
		})
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatal(err)
		}
	}
}
