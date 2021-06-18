package main

import (
	"github.com/AllenDang/giu"
	"github.com/gucio321/saper-go/pkg/sgiu"
)

func main() {
	wnd := giu.NewMasterWindow("Saper-go", 800, 600, 0, nil)
	wnd.Run(func() {
		giu.SingleWindow("game").Layout(
			sgiu.Create(30, 16, 99),
		)
	})
}
