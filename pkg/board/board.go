package board

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Board struct {
	Fields [][]*Field
	frozen bool
}

func NewBoard(w, h, numMines uint) *Board {
	result := &Board{
		Fields: make([][]*Field, h),
	}

	for n := range result.Fields {
		result.Fields[n] = make([]*Field, w)

		for i := range result.Fields[n] {
			result.Fields[n][i] = &Field{}
		}
	}

	for mine := 0; mine < int(numMines); {
		mineH := rand.Intn(int(h))
		mineW := rand.Intn(int(w))
		if !result.Fields[mineH][mineW].IsBomb() {
			w, h := int(w), int(h)
			result.Fields[mineH][mineW].value = Bomb

			if iw, ih := mineW-1, mineH-1; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW, mineH-1; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW+1, mineH-1; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW-1, mineH; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW+1, mineH; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW-1, mineH+1; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW, mineH+1; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			if iw, ih := mineW+1, mineH+1; iw >= 0 && ih >= 0 &&
				iw < w && ih < h &&
				!result.Fields[ih][iw].IsBomb() {
				result.Fields[ih][iw].value++
			}

			mine++
		}
	}

	return result
}

func (b *Board) Field(r, i int) *Field {
	return b.Fields[r][i]
}

func (b *Board) String() string {
	board := *b
	for h := range board.Fields {
		for w := range board.Fields[h] {
			if board.Fields[h][w].value == Bomb {
				board.Fields[h][w].value = 9
			}
		}
	}

	boardS := fmt.Sprintln(board.Fields)
	boardS = boardS[2:]
	return strings.ReplaceAll(boardS, "] [", "\n")
}

func (b *Board) LeftClick(row, idx int) (IsLose bool) {
	if b.frozen {
		return
	}

	if handled := b.Fields[row][idx].LeftClick(); !handled {
		return
	}

	switch field := b.Fields[row][idx]; field.value {
	case Bomb:
		b.Lose()
		return true
	case 0:
		neighbours := []struct{ r, i int }{
			{row - 1, idx - 1},
			{row - 1, idx},
			{row - 1, idx + 1},
			{row, idx - 1},
			{row, idx + 1},
			{row + 1, idx - 1},
			{row + 1, idx},
			{row + 1, idx + 1},
		}

		for _, c := range neighbours {
			if c.r >= 0 && c.i >= 0 &&
				c.r < len(b.Fields) && c.i < len(b.Fields[0]) {
				b.LeftClick(c.r, c.i)
			}
		}
	}

	return false
}

func (b *Board) RightClick(row, idx int) {
	if b.frozen {
		return
	}

	b.Field(row, idx).RightClick()
}

func (b *Board) Lose() {
	b.frozen = true

	for _, row := range b.Fields {
		for _, f := range row {
			if f.value != Bomb {
				continue
			}

			if f.state == MarkedBomb {
				continue
			}

			f.state = Open
		}
	}
}

type Field struct {
	value Value
	state State
}

func (f *Field) IsBomb() bool {
	return f.value == Bomb
}

// LeftClick handler
func (f *Field) LeftClick() (handled bool) {
	// no action when bomb marked
	switch f.state {
	case MarkedBomb, Open:
		return false
	}

	f.state = Open

	return true
}

// RightClick handler
func (f *Field) RightClick() {
	switch f.state {
	case Open:
		return
	case Blank:
		f.state = MarkedBomb
	case MarkedBomb:
		f.state = MarkedUncertain
	case MarkedUncertain:
		f.state = Blank
	}
}

func (f *Field) Value() Value {
	return f.value
}

func (f *Field) State() State {
	return f.state
}

func (f Field) String() string {
	lookup := map[State]string{
		Blank:           "",
		Open:            f.value.String(),
		MarkedBomb:      "P",
		MarkedUncertain: "?",
	}

	s, ok := lookup[f.state]
	if !ok {
		return "Err"
	}

	return s
}

type Value int

const (
	Bomb Value = -1
)

func (v Value) String() string {
	lookup := map[Value]string{
		Bomb: "X",
	}

	s, ok := lookup[v]
	if !ok {
		return strconv.Itoa(int(v))
	}

	return s
}

type State byte

const (
	Blank State = iota
	Open
	MarkedBomb
	MarkedUncertain
)
