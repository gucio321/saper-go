package main

import (
	"github.com/AllenDang/giu"

	"github.com/gucio321/saper-go/pkg/sgiu"
)

const (
	windowW, windowH         = 800, 600
	boardW, boardH, numMines = 5, 5, 3
)

func main() {
	wnd := giu.NewMasterWindow("Saper-go", windowW, windowH, 0)
	wnd.Run(func() {
		giu.SingleWindow("game").Layout(
			sgiu.Create("examplegame", boardW, boardH, numMines),
		)
	})
}
