package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gucio321/saper-go/pkg/sebiten"
)

const (
	windowW, windowH      = 800, 600
	boardW, boardH, mines = 30, 16, 99
)

func main() {
	ebiten.SetWindowSize(windowW, windowH)
	ebiten.SetWindowTitle("Saper-Go")
	if err := ebiten.RunGame(sebiten.Create(boardW, boardH, mines)); err != nil {
		log.Fatal(err)
	}
}
