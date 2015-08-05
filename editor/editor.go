// editor.go

// part of the pat package
package editor

import (
	"errors"
	"regexp"
)

var LogS func(string)

// One per file.
type Editor struct {
	file []byte
	dot  [][]int
}

func apply(f func([]int) [][]int, ints [][]int) [][]int {
	outScopes := make([][]int, 0)
	for _, scope := range ints {
		outScopes = append(outScopes, f(scope)...)
	}
	return outScopes
}

func (ed *Editor) Command(name string, args []string) error {
	switch name {
	case "x":
		if len(args) != 1 {
			return errors.New("Wrong number of arguments for the 'x' command")
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.xCommand(s, re) }, ed.dot)
	case "a":
		if len(args) != 1 {
			return errors.New("Wrong number of arguments for the 'a' command")
		}
		ed.dot = ed.aCommand(ed.dot, []byte(args[0]))
	case "g":
		if len(args) != 1 {
			return errors.New("Wrong number of arguments for the 'g' command")
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.gCommand(s, re) }, ed.dot)
	}

	return nil
}

func (ed *Editor) String() string {
	return string(ed.file)
}

func NewEditor(file []byte) *Editor {
	return &Editor{file, [][]int{[]int{0, len(file)}}}
}

func (ed *Editor) xCommand(scope []int, re *regexp.Regexp) [][]int {
	scopes := re.FindAllIndex(ed.file[scope[0]:scope[1]], -1)
	for _, s := range scopes {
		s[0] += scope[0]
		s[1] += scope[0]
	}
	return scopes
}

func (ed *Editor) gCommand(scope []int, re *regexp.Regexp) [][]int {
	if re.Match(ed.file[scope[0]:scope[1]]) {
		return [][]int{scope}
	}
	return [][]int{}
}

// Assumes increasing indicies
func (ed *Editor) aCommand(scopes [][]int, addition []byte) [][]int {
	var offset, off, index int
	outScopes := make([][]int, 0)
	for _, scope := range scopes {
		off = len(addition)
		index = scope[1] + offset

		ed.file = ed.file[:len(ed.file)+off]
		copy(ed.file[index+off:], ed.file[index:])
		copy(ed.file[index:index+off], addition)

		outScopes = append(outScopes, []int{scope[0] + offset, index + off})

		offset += off
	}
	return outScopes
}
