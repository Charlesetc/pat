// pat.go

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charlesetc/pat/display"
	"github.com/charlesetc/pat/editor"
	"github.com/charlesetc/pat/input"
	"github.com/nsf/termbox-go"
)

var Log func([]rune)
var files []string // filenames
// var editors []*editor // editors for the filenames
var ed *editor.Editor // the current editor.

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
		switch {
		// Parse no-argument ones.
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
			command = []string{c, lines[i+1], lines[i+2]}
			i++
		// Everything else gets one argument.
		default:
			command = []string{c, lines[i+1]}
		}
		i += 2
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

		case e.Key == 1: // Control-A
			input.CursorAtBeginning()
		case e.Key == 3: // Control-C
			Exit()
		case e.Key == 5: // Control-E
			input.CursorAtEnd()

		// Scrolling
		case e.Key == 14 || e.Key == 65516:
			display.ScrollDown()
		case e.Key == 16 || e.Key == 65517:
			display.ScrollUp()

		case e.Key == 27:
			input.Reset()
		case e.Key == 32: //Space
			input.AddRune(' ')
			input.Draw()
		case e.Key == 13: // Return
			runes := input.Runes()
			input.Reset()
			if len(runes) == 0 {
				break
			}
			// err := ed.Command(string(runes[0]), string(runes[1:]))
			// if err != nil {
			// 	panic(err) // for the time being..
			// }

			commands := parseLine(string(runes))
			for _, command := range commands {
				ed.Command(command[0], command[1:])
			}

			display.ShowFile([]rune(ed.String()))
			display.Draw()

		// Cursor Movement
		case e.Key == 65515:
			input.CursorLeft()
		case e.Key == 65514:
			input.CursorRight()

		// Delete Key
		case e.Key == 127 || e.Key == 8:
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

	Log = makeLog()
	display.Log = Log
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
