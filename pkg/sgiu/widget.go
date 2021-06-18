package sgiu

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/giu"
)

const btnSize = 30

type widget struct {
	width, height, numMines uint
}

func Create(w, h, m uint) giu.Widget {
	return &widget{
		width:    w,
		height:   h,
		numMines: m,
	}
}

func (w *widget) Build() {
	state := w.getState()

	giu.Layout{
		// board
		giu.Custom(func() {
			for r := 0; r < int(w.height); r++ {
				row := []giu.Widget{}
				for idx := 0; idx < int(w.width); idx++ {
					idx := idx
					row = append(row,
						giu.Button(state.board.Field(r, idx).String()+"##boarditem"+strconv.Itoa(r)+strconv.Itoa(idx)).
							Size(btnSize, btnSize),
						giu.Custom(func() {
							if !giu.IsItemHovered() {
								return
							}

							fmt.Println(r, idx)
							switch {
							case giu.IsMouseClicked(giu.MouseButtonLeft):
								state.board.LeftClick(r, idx)
							case giu.IsMouseClicked(giu.MouseButtonRight):
								state.board.RightClick(r, idx)
							}
						}),
					)
				}

				giu.Row(row...).Build()
			}
		}),
	}.Build()
}
