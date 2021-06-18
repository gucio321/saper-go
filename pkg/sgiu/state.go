package sgiu

import (
	"fmt"

	"github.com/AllenDang/giu"

	"github.com/gucio321/saper-go/pkg/board"
)

type widgetState struct {
	board *board.Board
}

func (s *widgetState) Dispose() {
	s.board = nil
}

func (w *widget) getStateID() string {
	return fmt.Sprintf("game_state")
}

func (w *widget) setState(s *widgetState) {
	giu.Context.SetState(w.getStateID(), s)
}

func (w *widget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(w.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		state = &widgetState{
			board: board.NewBoard(w.width, w.height, w.numMines),
		}

		w.setState(state)
	}

	return state
}
