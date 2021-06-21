package sebiten

import (
	"log"

	"github.com/gucio321/saper-go/pkg/board"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	fieldSize = 20
)

func Create(w, h, m uint) *Game {
	var err error

	result := &Game{
		width:    w,
		height:   h,
		numMines: m,
	}

	result.board, err = newBoard(board.NewBoard(w, h, m))
	if err != nil {
		log.Print(err)
	}

	return result
}

type Game struct {
	board *gameBoard
	width,
	height,
	numMines uint
}

func (g *Game) Update() error {
	g.board.Update()
	return nil
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return outsideW, outsideH
}

func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(colornames.Blue)

	g.board.Draw(screen)
}
