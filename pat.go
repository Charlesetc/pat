// pat.go

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/charlesetc/pat/display"
	"github.com/charlesetc/pat/editor"
	"github.com/charlesetc/pat/input"
	"github.com/charlesetc/pat/stack"
	"github.com/nsf/termbox-go"
)

var Log func([]rune)
var files []string // filenames
// var editors []*editor // editors for the filenames
var ed *editor.Editor // the current editor.
var commandHistory *stack.Stack

func lineCommand(command string) (bool, []string, int) {
	var parsed bool
	var output []string
	var numberParsed int
	switch {
	case len(command) == 0:
		parsed = false
	case strings.ContainsRune(command, ','):
		strs := strings.SplitN(command, ",", 2)
		_, err1 := strconv.Atoi(strs[0])
		n, err2 := strconv.Atoi(strs[1])
		if (err1 != nil && strs[0] != "") || (err2 != nil && strs[1] != "$" && strs[1] != "") {
			parsed = false
			break
		}
		parsed = true
		output = []string{",", strs[0], strs[1]}
		numberParsed = len(strs[0]) + 1 + len(strconv.Itoa(n))
		if len(strs[1]) == 0 {
			numberParsed-- // because "" translates to 0
		}
	default:
		n, err := strconv.Atoi(command)
		parsed = err == nil
		numberParsed = len(strconv.Itoa(n))
		output = []string{"line", command}
	}

	return parsed, output, numberParsed
}

func parseLine(line string) [][]string {
	commands := make([][]string, 0)

	// this works, it might be ugly, but it works.
	line = strings.Replace(line, "\\/", "&#sslash;", -1)
	line = strings.Replace(line, "/", "&#slash;", -1)
	line = strings.Replace(line, "&#sslash;", "/", -1)
	lines := strings.Split(line, "&#slash;")
	var i int
	var command []string
	var c string

Lines:
	for i < len(lines) {
		c = lines[i]
		c = strings.Replace(c, " ", "", -1) // Get rid of space

		isLineCommand, lineResult, numberParsed := lineCommand(c)

		switch {

		case isLineCommand:
			command = lineResult
			lines[i] = lines[i][:numberParsed]
			if numberParsed < len(c) {
				i--
			}
		case c == "color":
			display.RandomColor()
			return [][]string{}
		case c == "d":
			command = []string{c}

		// Parse ?re?
		case len(c) > 0 && c[0] == '?':
			command = []string{"?", c[1 : len(c)-1]}
			break

		// Do nothing at the end.
		case i == len(lines)-1:
			break Lines

		// Parse S
		case c == "s":
			command1 := []string{"x", lines[i+1]}
			commands = append(commands, command1)
			command = []string{"c", lines[i+2]}
			i += 2
		// Everything else gets one argument.
		default:
			command = []string{c, lines[i+1]}
			i++ // One argument
		}
		i++ // for itself
		commands = append(commands, command)
	}

	LogS(fmt.Sprint(commands))

	return commands
}

func makeLog() func([]rune) {
	file, err := os.Create(".log")
	if err != nil {
		panic(err)
	}
	return func(input []rune) {
		file.WriteString(string(input))
		if input[len(input)-1] != '\n' {
			file.WriteString("\n")
		}
		file.Sync()
	}
}

func LogS(str string) {
	Log([]rune(str))
}

func Exit() {
	LogS("Exiting now")
	display.Reset()
	os.Exit(0)
}

func Poll() {
	for {
		e := termbox.PollEvent()

		// LogS(fmt.Sprintf("%d", e.Key))

		switch {
		case e.Type == termbox.EventResize:
			display.Resize()
			display.Draw()

		case e.Key == termbox.KeyCtrlA: // Control-A
			input.CursorAtBeginning()
		case e.Key == termbox.KeyCtrlE: // Control-E
			input.CursorAtEnd()
		case e.Key == termbox.KeyCtrlC: // Control-C
			Exit()

			// Scrolling
		case e.Key == termbox.KeyCtrlN:
			display.ScrollDown()
		case e.Key == termbox.KeyCtrlP:
			display.ScrollUp()

		case e.Key == termbox.KeyEsc:
			input.Reset()
		case e.Key == termbox.KeySpace: //Space
			input.AddRune(' ')
			input.Draw()
		case e.Key == termbox.KeyEnter: // Return
			runes := input.Runes()
			commandHistory.Add(runes)
			input.Reset()
			if len(runes) == 0 {
				break
			}

			commands := parseLine(string(runes))
			for _, command := range commands {
				ed.Command(command[0], command[1:])
			}

			display.ShowFile([]rune(ed.String()))
			display.Highlight(ed.Highlights())
			display.Draw()

		// // Arrow keys
		// Cursor Movement
		case e.Key == termbox.KeyArrowLeft:
			input.CursorLeft()
		case e.Key == termbox.KeyArrowRight:
			input.CursorRight()
		case e.Key == termbox.KeyArrowUp:
			past, err := commandHistory.Pop()
			if err != nil {
				// Empty stack.
				// show alert later.
				break
			}
			input.Reset()
			input.SetRunes(past.([]rune))
			input.Draw()
			break
		case e.Key == termbox.KeyArrowDown:
			input.Reset()
			past, err := commandHistory.UnPop()
			if err != nil {
				// Empty stack.
				// show alert later.
				break
			}
			input.SetRunes(past.([]rune))
			input.Draw()
			break

		// Delete Key
		case e.Key == termbox.KeyBackspace || e.Key == termbox.KeyBackspace2 || e.Key == termbox.KeyDelete:
			input.Backspace()
			input.Draw()
		case e.Key == 0: // All other normal chars
			input.AddRune(e.Ch)
			input.Draw()
		}
	}
}

func contains(strings []string, match string) bool {
	for _, str := range strings {
		if str == match {
			return true
		}
	}
	return false
}

func init() {
	flags := make([]string, 0)
	files = make([]string, 0)

	for _, arg := range os.Args[1:] {
		if arg[0] == '-' {
			flags = append(flags, arg)
			continue
		}
		files = append(files, arg)
	}

	if contains(flags, "-v") || contains(flags, "--version") {
		fmt.Println("The Glorious Pat Text Editor : v0.0.1")
		os.Exit(0)
	} else if contains(flags, "-rc") {
		fmt.Println("Yay Recurse Center!")
		os.Exit(0)
	}

	if len(files) == 0 {
		fmt.Println("usage: pat [file]")
		os.Exit(0)
	}

	display.Init(!(contains(flags, "--bottom") || contains(flags, "-b"))) // topbar
	commandHistory = stack.New()

	Log = makeLog()
	display.Log = Log
	display.LogS = LogS
	input.Log = Log
	input.LogS = LogS
	editor.LogS = LogS
}

func main() {
	defer display.Reset()

	bytes, err := ioutil.ReadFile(files[0])
	if err != nil {
		panic(err)
	}

	ed = editor.NewEditor(bytes)

	display.Show([]rune(string(bytes)), []rune{})
	display.Draw()
	Poll()
}
