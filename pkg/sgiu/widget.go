package sgiu

import (
	"strconv"

	"golang.org/x/image/colornames"

	"github.com/AllenDang/giu"
)

const btnSize = 20

type widget struct {
	width, height, numMines uint
	id                      string
}

// Create creates a widget
func Create(id string, w, h, m uint) giu.Widget {
	return &widget{
		width:    w,
		height:   h,
		numMines: m,
		id:       id,
	}
}

// Build builds a widget (giu.Widget implementation)
func (w *widget) Build() {
	state := w.getState()

	giu.Layout{
		// timer
		giu.Custom(func() {
			t := state.board.Time()
			label := strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
			giu.Label(label).Build()
		}),
		// board
		giu.Custom(func() {
			for r := 0; r < int(w.height); r++ {
				row := []giu.Widget{}
				for idx := 0; idx < int(w.width); idx++ {
					idx := idx

					field := state.board.Field(r, idx)

					textColor, bgColor := field.GetColors()

					row = append(row,
						giu.Style().
							SetColor(giu.StyleColorText, textColor).
							SetColor(giu.StyleColorButton, bgColor).
							SetColor(giu.StyleColorButtonHovered, bgColor).
							SetColor(giu.StyleColorButtonActive, colornames.Black).
							SetStyle(giu.StyleVarItemSpacing, 50, 0).To(
							giu.Button(field.String()+"##"+w.id+"boarditem"+strconv.Itoa(r)+strconv.Itoa(idx)).
								Size(btnSize, btnSize),
						),
						giu.Custom(func() {
							if !giu.IsItemHovered() {
								return
							}

							switch {
							case giu.IsMouseClicked(giu.MouseButtonLeft):
								state.board.LeftClick(r, idx)
							case giu.IsMouseClicked(giu.MouseButtonRight):
								state.board.RightClick(r, idx)
							}
						}),
					)
				}

				giu.Style().
					SetStyle(giu.StyleVarItemSpacing, 0, 0).
					To(giu.Row(row...)).Build()
			}
		}),
		giu.Button("Reset##" + w.id + "resetButton").OnClick(func() {
			state.board.Retry()
		}),
	}.Build()
}
