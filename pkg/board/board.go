package board

import (
	"crypto/rand"
	"image/color"
	"log"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/image/colornames"
)

type pos struct {
	row,
	index int
}

// Board represents a minesweeper board
type Board struct {
	Fields [][]*Field
	frozen bool
	numMines,
	width,
	height uint

	time   time.Time
	ticker *time.Ticker
	timeCB func()
}

// NewBoard creates a new board
func NewBoard(w, h, numMines uint) *Board {
	result := &Board{
		Fields:   make([][]*Field, h),
		width:    w,
		height:   h,
		numMines: numMines,
		ticker:   time.NewTicker(1 * time.Second),
	}

	result.fill()

	go func() {
		for range result.ticker.C {
			result.time = result.time.Add(time.Second)
			if result.timeCB != nil {
				result.timeCB()
			}
		}
	}()

	return result
}

// fill fills a board with a random values
func (b *Board) fill() {
	// reset board
	for n := range b.Fields {
		b.Fields[n] = make([]*Field, b.width)

		for i := range b.Fields[n] {
			b.Fields[n][i] = &Field{}
		}
	}

	// fill board
	for mine := 0; mine < int(b.numMines); {
		mineHBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(b.height)))
		if err != nil {
			log.Panic(err)
		}

		mineH := int(mineHBigInt.Int64())

		mineWBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(b.width)))
		if err != nil {
			log.Panic(err)
		}

		mineW := int(mineWBigInt.Int64())

		if !b.Fields[mineH][mineW].IsBomb() {
			b.Fields[mineH][mineW].value = Bomb

			neighbors := b.Neighbours(mineH, mineW)

			for _, n := range neighbors {
				if !b.Fields[n.row][n.index].IsBomb() {
					b.Fields[n.row][n.index].value++
				}
			}

			mine++
		}
	}
}

// Retry resets a board
func (b *Board) Retry() {
	b.frozen = false
	b.fill()
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

// Field returns a specified field
func (b *Board) Field(r, i int) *Field {
	return b.Fields[r][i]
}

// LeftClick handles left click (should be called on left click on field)
// in some graphic wrapper of the board
func (b *Board) LeftClick(row, idx int) (isLose bool) {
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
		b.lose()
		return true
	case 0:
		neighbors := b.Neighbours(row, idx)

		for _, c := range neighbors {
			b.LeftClick(c.row, c.index)
		}
	}

	return false
}

// RightClick handles right click (should be called on right click on field)
// in some graphic wrapper of the board
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

func (b *Board) StopTimer() {
	b.ticker.Stop()
}

func (b *Board) ResumeTimer() {
	b.ticker.Reset(time.Second)
}

func (b *Board) UpdateTime(cb func()) {
	b.timeCB = cb
}

func (b *Board) Time() time.Time {
	return b.time
}

// lose handles a lose event
func (b *Board) lose() {
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

// Field represents a board field
type Field struct {
	value Value
	state State
}

// IsBomb returns true if field's value is "Bomb"
func (f *Field) IsBomb() bool {
	return f.value == Bomb
}

// Value returns a value of field
func (f *Field) Value() Value {
	return f.value
}

// State returns a field's state
func (f *Field) State() State {
	return f.state
}

// String returns a field string
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

func (f *Field) GetColors() (textColor, bgColor color.RGBA) {
	switch f.State() {
	case Open:
		bgColor = colornames.Black
	default:
		bgColor = colornames.Green
	}

	switch f.State() {
	case MarkedBomb:
		textColor = colornames.Black
	case MarkedUncertain:
		textColor = colornames.Orange
	case Open:
		// nolint:gomnd // obvious meaning - a value of fields
		switch f.Value() {
		case Bomb:
			textColor = colornames.Red
		case 1:
			textColor = colornames.Green
		case 2:
			textColor = colornames.White
		case 3:
			textColor = colornames.Aqua
		case 4:
			textColor = colornames.Yellow
		case 5:
			textColor = colornames.Blue
		case 6:
			textColor = colornames.Violet
		case 7, 8: // TODO - check this colors
			textColor = colornames.White
		}
	}

	return
}

// Value represents a value of vield
type Value int

const (
	// Bomb is a 'mine' value
	Bomb Value = -1
)

// String returns a value string
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

// State represents field state
type State byte

const (
	// Blank - field isn covered
	Blank State = iota
	// Open - field is opened (and save or not ;-) )
	Open
	// MarkedBomb flag is present
	MarkedBomb
	// MarkedUncertain - uncertain flag present
	MarkedUncertain
)
