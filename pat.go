// pat.go

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/charlesetc/pat/display"
	"github.com/charlesetc/pat/input"
	"github.com/nsf/termbox-go"
)

var Log func([]rune)
var files []string // filenames

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
			Log(input.Runes())
			input.Reset()

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

	if len(files) == 0 {
		fmt.Println("usage: pat [file]")
		os.Exit(0)
	}

	display.Init(!(contains(flags, "--bottom") || contains(flags, "-b"))) // topbar

	Log = makeLog()
	display.Log = Log
	input.Log = Log
	input.LogS = LogS
}

func main() {
	defer display.Reset()

	bytes, err := ioutil.ReadFile(files[0])
	if err != nil {
		panic(err)
	}

	display.Show([]rune(string(bytes)), []rune{})
	display.Draw()
	Poll()
}
