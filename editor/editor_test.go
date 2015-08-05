// editor_test.go

package editor

import (
	"bytes"
	"regexp"
	"testing"
)

func Equal(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Error(a, "!=", b)
	}
}

func TestCommandFunc(t *testing.T) {
	ed := NewEditor([]byte("aabcc"))
	err := ed.Command("x", "b.")
	if err != nil {
		panic(err)
	}
	err = ed.Command("x", "c")
	if err != nil {
		panic(err)
	}
	Equal(t, len(ed.dot), 1)
	Equal(t, ed.dot[0][0], 3)
	Equal(t, ed.dot[0][1], 4)
}

func TestACommand(t *testing.T) {
	ed := NewEditor([]byte("hello there hello"))
	scopes := ed.aCommand([][]int{[]int{0, 2}, []int{4, 6}}, []byte("wow"))
	if bytes.Compare(ed.file, []byte("hewowllo wowthere hello")) != 0 {
		t.Error("Got ed file: %s", string(ed.file))
	}
	Equal(t, scopes[0][0], 0)
	Equal(t, scopes[0][1], 5)
	Equal(t, scopes[1][0], 7)
	Equal(t, scopes[1][1], 12)
}

func TestXCommand(t *testing.T) {
	re, err := regexp.Compile("h.")
	if err != nil {
		panic(err)
	}
	ed := NewEditor([]byte("hh hack"))
	scopes := ed.xCommand([]int{0, 7}, re)
	Equal(t, scopes[0][0], 0)
	Equal(t, scopes[0][1], 2)
	Equal(t, scopes[1][0], 3)
	Equal(t, scopes[1][1], 5)
	scopes = ed.xCommand([]int{1, 7}, re)
	Equal(t, scopes[0][0], 1)
	Equal(t, scopes[0][1], 3)
	Equal(t, scopes[1][0], 3)
	Equal(t, scopes[1][1], 5)
}
