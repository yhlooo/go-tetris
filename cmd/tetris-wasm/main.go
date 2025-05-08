package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"

	"github.com/yhlooo/go-tetris/pkg/ui/web"
)

var (
	genStatic    = false
	staticPath   = "web/static"
	staticPrefix = ""
	listenAddr   = ":8000"
)

func init() {
	flag.BoolVar(&genStatic, "gen-static", genStatic, "Generate static files")
	flag.StringVar(&staticPath, "static-path", staticPath, "Generate static files to specified path")
	flag.StringVar(&staticPrefix, "static-prefix", staticPrefix, "URI prefix for static files")
	flag.StringVar(&listenAddr, "listen", listenAddr, "Listen address")
}

func main() {
	flag.Parse()

	app.Route("/", func() app.Composer {
		return web.NewGameUI()
	})
	app.RunWhenOnBrowser()

	if genStatic {
		// 生成静态文件
		log.Printf("generate static files to %s", staticPath)
		h := &app.Handler{
			Name: "Tetris",
		}
		if staticPrefix != "" {
			h.Resources = app.PrefixedLocation(staticPrefix)
		}
		if err := app.GenerateStaticWebsite(staticPath, h); err != nil {
			log.Fatal(err)
		}
	} else {
		// 运行 Server
		log.Printf("serving http on %s", listenAddr)
		http.Handle("/", &app.Handler{
			Name: "Tetris",
		})
		if err := http.ListenAndServe(listenAddr, nil); err != nil {
			log.Fatal(err)
		}
	}
}
