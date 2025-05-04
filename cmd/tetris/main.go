package main

import (
	"log"

	"github.com/yhlooo/go-tetris/pkg/ui/tty"
)

func main() {
	ui := tty.NewGameUI()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
