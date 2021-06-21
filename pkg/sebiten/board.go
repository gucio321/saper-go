package sebiten

import (
	"bytes"
	"image"
	_ "image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/gucio321/saper-go/pkg/board"
	"github.com/gucio321/saper-go/pkg/sebiten/assets"
)

func newBoard(board *board.Board) (*gameBoard, error) {
	result := &gameBoard{
		Board: board,
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

	img, _, err := image.Decode(bytes.NewReader(assets.Flag))
	if err != nil {
		return nil, err
	}

	result.flag = ebiten.NewImageFromImage(img)

	return result, nil
}

type gameBoard struct {
	*board.Board
	font        font.Face
	buttonPress bool
	flag        *ebiten.Image
}

func (b *gameBoard) Update() {
	x, y := ebiten.CursorPosition()
	idxX := x / fieldSize
	idxY := y / fieldSize

	switch {
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
		if !b.buttonPress {
			b.LeftClick(idxY, idxX)
			b.buttonPress = true
		}
	case ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight):
		if !b.buttonPress {
			b.RightClick(idxY, idxX)
			b.buttonPress = true
		}
	default:
		b.buttonPress = false
	}
}

func (b *gameBoard) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(fieldSize*int(b.Width()), fieldSize*int(b.Height()))

	for y := 0; y < int(b.Height())*fieldSize; y++ {
		if y%fieldSize == 0 {
			continue
		}

		for x := 0; x < int(b.Width())*fieldSize; x++ {
			if x%fieldSize == 0 {
				continue
			}

			// b.Board.LeftClick(y/fieldSize, x/fieldSize)

			idxX, idxY := x/fieldSize, y/fieldSize
			field := b.Field(idxY, idxX)
			_, bgColor := field.GetColors()

			img.Set(x, y, bgColor)

		}

	}

	screen.DrawImage(img, &ebiten.DrawImageOptions{})

	// render labels
	for y := 0; y < int(b.Height()); y++ {
		for x := 0; x < int(b.Width()); x++ {
			// base position of the field
			posX, posY := x*fieldSize, fieldSize+y*fieldSize

			field := b.Field(y, x)

			switch field.State() {
			case board.MarkedBomb:
				// r := image.Rect(0, 0, posX, posY)
				screen.DrawImage(b.flag, &ebiten.DrawImageOptions{})
			default:
				s := field.String()

				textColor, _ := field.GetColors()

				labelSize := text.BoundString(b.font, s)
				labelW := labelSize.Dx()
				labelH := labelSize.Dy()

				labelX := posX + (fieldSize-labelW)/2
				labelY := posY - (fieldSize-labelH)/2

				text.Draw(screen, field.String(), b.font, labelX, labelY, textColor)
			}
		}
	}
}
