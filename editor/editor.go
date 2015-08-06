// editor.go

// part of the pat package
package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var LogS func(string)

// One per file.
type Editor struct {
	file  []byte
	dot   [][]int
	saved [][]int
}

func (ed *Editor) SaveDot() {
	ed.saved = ed.dot // might not work.
}

func (ed *Editor) UnSaveDot() {
	ed.dot = ed.saved
}

func escapeSpace(str string) string {
	str = strings.Replace(str, "\\\\", "&#doubleslash;", -1)
	str = strings.Replace(str, "\\n", "\n", -1)
	str = strings.Replace(str, "\\t", "\t", -1)
	str = strings.Replace(str, "&#doubleslash;", "\\\\", -1)
	return str
}

// Because lines matter sometimes.
func escapeRegex(str string) string {
	str = strings.Replace(str, "\\\\", "&#doubleslash;", -1)
	str = strings.Replace(str, "\\N", ".*\n", -1)
	str = strings.Replace(str, "&#doubleslash;", "\\\\", -1)
	return str
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
	case "", "?":
		if args[0] == "" {
			return nil // when there's just a slash
		}
		re, err := regexp.Compile(escapeRegex(args[0]))
		if err != nil {
			return err
		}
		var newScope []int
		for _, scope := range ed.dot {
			if name == "" {
				newScope = re.FindIndex(ed.file[scope[0]:scope[1]])
			} else {
				indexes := re.FindAllIndex(ed.file[scope[0]:scope[1]], -1)
				if len(indexes) == 0 {
					newScope = nil
				} else {
					newScope = indexes[len(indexes)-1]
				}
			}
			if newScope == nil {
				continue
			}
			ed.dot = [][]int{[]int{newScope[0] + scope[0], newScope[1] + scope[0]}}
			LogS(fmt.Sprint(ed.dot))
			return nil
		}
		return errors.New("No match found.")
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
		re, err := regexp.Compile(escapeRegex(args[0]))
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.xCommand(s, re) }, ed.dot)
	case "a":
		ed.dot = ed.insertCommand(ed.dot, []byte(escapeSpace(args[0])), false)
	case "c":
		ed.dot = ed.cCommand(ed.dot, []byte(escapeSpace(args[0])))
	case "i":
		ed.dot = ed.insertCommand(ed.dot, []byte(escapeSpace(args[0])), true)
	case "d":
		ed.dot = ed.dCommand(ed.dot)
	case "g":
		re, err := regexp.Compile(escapeRegex(args[0]))
		if err != nil {
			return err
		}
		ed.dot = apply(func(s []int) [][]int { return ed.matchCommand(s, re, true) }, ed.dot)
	case "y":
		re, err := regexp.Compile(escapeRegex(args[0]))
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
	return &Editor{file, [][]int{[]int{0, len(file)}}, [][]int{}}
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
// func (ed *Editor) insertCommand(scopes [][]int, addition []byte, beginning bool) [][]int {
// 	var offset, off, index int
// 	outScopes := make([][]int, 0)
// 	for _, scope := range scopes {
// 		off = len(addition)
// 		if beginning {
// 			index = scope[0] + offset
// 		} else {
// 			index = scope[1] + offset
// 		}
//
// 		ed.file = ed.file[:len(ed.file)+off+1]
// 		ed.file[len(ed.file)-1] = 0
// 		copy(ed.file[index+off:], ed.file[index:])
// 		copy(ed.file[index:index+off], addition)
//
// 		outScopes = append(outScopes, []int{scope[0] + offset, index + off})
//
// 		offset += off
// 	}
// 	return outScopes
// }

func (ed *Editor) insertCommand(scopes [][]int, addition []byte, beginning bool) [][]int {
	// No places do do anything, don't.
	if len(scopes) == 0 {
		return [][]int{}
	}

	var finalSum int
	addLength := len(addition)
	for range scopes {
		finalSum += addLength
	}
	new_file := make([]byte, len(ed.file)+finalSum)

	outscopes := make([][]int, len(scopes))
	var outscopeI int

	var startOrEnd int
	if beginning {
		startOrEnd = 0
	} else {
		startOrEnd = 1
	}

	var j, scopeIndex int
	currentscope := scopes[scopeIndex]
	for k := range ed.file {
		// Iterate over the scopes as we go.
		if k > currentscope[startOrEnd] {
			scopeIndex++
			if len(scopes) > scopeIndex {
				currentscope = scopes[scopeIndex]
			}
		}

		// add the insertion
		if k == currentscope[startOrEnd] {
			copy(new_file[j:j+addLength], addition)
			j += addLength

			if beginning {
				outscopes[outscopeI] = []int{j - addLength, j}
			} else {
				outscopes[outscopeI] = []int{j - addLength, j}
			}
			outscopeI++
			// LogS(fmt.Sprint(outscopes))
		}

		// continue copying
		new_file[j] = ed.file[k]
		j++
	}
	ed.file = new_file
	return outscopes
}

// fast because doesn't allocate memory
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
	// No places do do anything, don't.
	if len(scopes) == 0 {
		return [][]int{}
	}

	var finalSum int
	addLength := len(addition)
	for _, scope := range scopes {
		finalSum += addLength - (scope[1] - scope[0])
	}
	new_file := make([]byte, len(ed.file)+finalSum)

	//outscopes := make([][]int, len(scopes))
	var outscopes [][]int

	var j, scopeIndex int
	currentscope := scopes[scopeIndex]
	for k := 0; k < len(ed.file); k++ {
		scopediff := currentscope[1] - currentscope[0]

		// Iterate over the scopes as we go.
		if k > currentscope[1] {
			scopeIndex++
			if len(scopes) > scopeIndex {
				currentscope = scopes[scopeIndex]
			}
		}

		// add the insertion
		if k == currentscope[0] {
			copy(new_file[j:j+addLength], addition)
			j += addLength
			if addLength > 0 {
				outscopes = append(outscopes, []int{j - addLength, j})
			} else {
				outscopes = append(outscopes, []int{j + addLength, j})
			}
			// Ignore those characters...
			k += scopediff
		}

		// continue copying
		if len(ed.file) > k && len(new_file) > j {
			new_file[j] = ed.file[k]
		}
		j++
	}
	ed.file = new_file
	return outscopes
}
