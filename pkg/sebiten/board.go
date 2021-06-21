package sebiten

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/gucio321/saper-go/pkg/board"
)

func newBoard(board *board.Board) (*gameBoard, error) {
	result := &gameBoard{
		Board: board,
		image: ebiten.NewImage(fieldSize*int(board.Width()), fieldSize*int(board.Height())),
	}

	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	result.font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fieldSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	result.rebuildImage()

	return result, nil
}

type gameBoard struct {
	*board.Board
	font        font.Face
	buttonPress bool
	image       *ebiten.Image
}

func (b *gameBoard) Update() {
	x, y := ebiten.CursorPosition()
	idxX := x / fieldSize
	idxY := y / fieldSize

	switch {
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
		if !b.buttonPress {
			b.LeftClick(idxY, idxX)
			b.rebuildImage()
			b.buttonPress = true
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight):
		if !b.buttonPress {
			b.RightClick(idxY, idxX)
			b.rebuildImage()
			b.buttonPress = true
		}
	default:
		b.buttonPress = false
	}
}

func (b *gameBoard) Draw(screen *ebiten.Image) {
	screen.DrawImage(b.image, &ebiten.DrawImageOptions{})
}

func (b *gameBoard) rebuildImage() {
	// render labels
	for y := 0; y < int(b.Height()); y++ {
		for x := 0; x < int(b.Width()); x++ {
			// base position of the field
			posX, posY := x*fieldSize, y*fieldSize

			field := b.Field(y, x)

			textColor, bgColor := field.GetColors()

			for imgX := 1; imgX < fieldSize; imgX++ {
				for imgY := 1; imgY < fieldSize; imgY++ {
					b.image.Set(posX+imgX, posY+imgY, bgColor)
				}
			}

			s := field.String()

			labelSize := text.BoundString(b.font, s)
			labelW := labelSize.Dx()
			labelH := labelSize.Dy()

			labelX := posX + (fieldSize-labelW)/2
			labelY := fieldSize + posY - (fieldSize-labelH)/2

			text.Draw(b.image, field.String(), b.font, labelX, labelY, textColor)
		}
	}
}
