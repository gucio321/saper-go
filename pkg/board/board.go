package board

import (
	"math/rand"
	"strconv"
)

type pos struct {
	row,
	index int
}

type Board struct {
	Fields [][]*Field
	frozen bool
	numMines,
	width,
	height uint
}

func NewBoard(w, h, numMines uint) *Board {
	result := &Board{
		Fields:   make([][]*Field, h),
		width:    w,
		height:   h,
		numMines: numMines,
	}

	result.Fill()

	return result
}

func (b *Board) Fill() {
	// reset board
	for n := range b.Fields {
		b.Fields[n] = make([]*Field, b.width)

		for i := range b.Fields[n] {
			b.Fields[n][i] = &Field{}
		}
	}

	// fill board
	for mine := 0; mine < int(b.numMines); {
		mineH := rand.Intn(int(b.height))
		mineW := rand.Intn(int(b.width))
		if !b.Fields[mineH][mineW].IsBomb() {
			b.Fields[mineH][mineW].value = Bomb

			neighbours := b.Neighbours(mineH, mineW)

			for _, n := range neighbours {
				if !b.Fields[n.row][n.index].IsBomb() {
					b.Fields[n.row][n.index].value++
				}
			}

			mine++
		}
	}
}

func (b *Board) Retry() {
	b.frozen = false
	b.Fill()
}

// Neighbours returns a list of fields connecting with given
func (b *Board) Neighbours(row, index int) []pos {
	result := make([]pos, 0)

	possible := []pos{
		{row - 1, index - 1},
		{row - 1, index},
		{row - 1, index + 1},
		{row, index - 1},
		{row, index + 1},
		{row + 1, index - 1},
		{row + 1, index},
		{row + 1, index + 1},
	}

	// check if indexes exists
	for _, p := range possible {
		if p.row >= 0 && p.row < int(b.height) &&
			p.index >= 0 && p.index < int(b.width) {
			result = append(result, p)
		}
	}

	return result
}

func (b *Board) Field(r, i int) *Field {
	return b.Fields[r][i]
}

func (b *Board) LeftClick(row, idx int) (IsLose bool) {
	if b.frozen {
		return
	}

	field := b.Field(row, idx)

	switch field.state {
	case Open:
		// https://github.com/gucio321/saper-go/issues/1
		return false
	case MarkedUncertain, Blank:
		field.state = Open
	}

	switch field.value {
	case Bomb:
		b.Lose()
		return true
	case 0:
		neighbours := b.Neighbours(row, idx)

		for _, c := range neighbours {
			b.LeftClick(c.row, c.index)
		}
	}

	return false
}

func (b *Board) RightClick(row, idx int) {
	if b.frozen {
		return
	}

	field := b.Field(row, idx)

	switch field.state {
	case Open:
		// noop
	case Blank:
		field.state = MarkedBomb
	case MarkedBomb:
		field.state = MarkedUncertain
	case MarkedUncertain:
		field.state = Blank
	}
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
