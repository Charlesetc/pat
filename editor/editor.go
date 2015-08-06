// editor.go

// part of the pat package
package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var LogS func(string)

// One per file.
type Editor struct {
	file []byte
	dot  [][]int
}

func (ed *Editor) nthLine(n int) (int, bool) {
	var newlines, index int
	var r byte
	for index, r = range ed.file {
		if newlines == n {
			return index, true
		}
		if r == '\n' {
			newlines++
		}
	}
	return index, false
}

func (ed *Editor) Highlights() [][]int {
	return ed.dot
}

func apply(f func([]int) [][]int, ints [][]int) [][]int {
	outScopes := make([][]int, 0)
	for _, scope := range ints {
		outScopes = append(outScopes, f(scope)...)
	}
	return outScopes
}

func (ed *Editor) multiLineSelect(l1, l2 int) error {
	index1, finished := ed.nthLine(l1)
	if !finished {
		return errors.New(fmt.Sprintf("No such line %d.", l1))
	}
	index2, _ := ed.nthLine(l2)
	ed.dot = [][]int{[]int{index1, index2}}
	return nil
}

func (ed *Editor) Command(name string, args []string) error {
	switch name {
	case "line":
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		ed.multiLineSelect(n, n+1)
	case ",":
		var err error
		var n1, n2 int
		if args[0] == "" {
			n1 = 0
		} else {
			n1, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}
		if args[1] == "" || args[1] == "$" {
			n2 = -1
		} else {
			n2, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
		}
		ed.multiLineSelect(n1, n2)
	case "x":
		re, err := regexp.Compile(args[0])
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.xCommand(s, re) }, ed.dot)
	case "a":
		ed.dot = ed.insertCommand(ed.dot, []byte(args[0]), false)
	case "c":
		ed.dot = ed.cCommand(ed.dot, []byte(args[0]))
	case "i":
		ed.dot = ed.insertCommand(ed.dot, []byte(args[0]), true)
	case "d":
		LogS(fmt.Sprint(ed.dot))
		ed.dot = ed.dCommand(ed.dot)
		LogS(fmt.Sprint(ed.dot))
	case "g":
		re, err := regexp.Compile(args[0])
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.matchCommand(s, re, true) }, ed.dot)
	case "y":
		re, err := regexp.Compile(args[0])
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.matchCommand(s, re, false) }, ed.dot)
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

func (ed *Editor) matchCommand(scope []int, re *regexp.Regexp, keep bool) [][]int {
	// Strange xor simplification
	if re.Match(ed.file[scope[0]:scope[1]]) == keep {
		return [][]int{scope}
	}
	return [][]int{}
}

// Assumes increasing indicies
func (ed *Editor) insertCommand(scopes [][]int, addition []byte, beginning bool) [][]int {
	var offset, off, index int
	outScopes := make([][]int, 0)
	for _, scope := range scopes {
		off = len(addition)
		if beginning {
			index = scope[0] + offset
		} else {
			index = scope[1] + offset
		}

		ed.file = ed.file[:len(ed.file)+off+1]
		ed.file[len(ed.file)-1] = 0
		copy(ed.file[index+off:], ed.file[index:])
		copy(ed.file[index:index+off], addition)

		outScopes = append(outScopes, []int{scope[0] + offset, index + off})

		offset += off
	}
	return outScopes
}

func (ed *Editor) dCommand(scopes [][]int) [][]int {
	var offset, off, index int
	outScopes := make([][]int, 0)
	for _, scope := range scopes {
		off = scope[1] - scope[0]
		index = scope[0] - offset
		ed.file = append(ed.file[:index], ed.file[index+off:]...)
		offset += off
		outScopes = append(outScopes, []int{index, index})
	}
	return outScopes
}

func (ed *Editor) cCommand(scopes [][]int, addition []byte) [][]int {
	var offset, off, index int
	outScopes := make([][]int, 0)
	for _, scope := range scopes {
		off = len(addition) - (scope[1] - scope[0])
		insert := len(addition)
		index = scope[0] + offset
		LogS(fmt.Sprint(len(ed.file), off, index))
		ed.file = ed.file[:len(ed.file)+off]
		copy(ed.file[index+insert:], ed.file[index+(scope[1]-scope[0]):])
		copy(ed.file[index:index+insert], addition)
		offset += off
		outScopes = append(outScopes, []int{index, index})
	}
	return outScopes
}
