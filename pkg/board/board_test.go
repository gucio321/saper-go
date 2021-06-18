package board

import (
	"fmt"
	"testing"
)

func Test_String(t *testing.T) {
	b := NewBoard(30, 16, 99)
	fmt.Println(b.String())
	t.Fail()
}
