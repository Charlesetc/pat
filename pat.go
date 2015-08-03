package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var Log func(string)
var d *Display

type Display struct {
	current string
	height  int
	width   int
	scroll  int
}

func makeLog() func(string) {
	file, err := os.Create(".pat.log")
	if err != nil {
		panic(err)
	}
	return func(input string) {
		file.WriteString(input)
		if input[len(input)-1] != '\n' {
			file.WriteString("\n")
		}
		file.Sync()
	}
}

func (d *Display) viewPort() string {
	for _, r := range string {

	}
}

func (d *Display) Draw() {
	d.Clear()
	Log(d.current)
	var i, j int
	for _, r := range d.current {
		if r == '\n' || i >= d.width {
			j++
			i = 0
			continue
		}
		termbox.SetCell(i, j, r, 0, 0)
		i++
	}
	termbox.Flush()
}

func (d *Display) Show(str string) {
	d.current = str
}

func NewDisplay() *Display {
	d := new(Display)
	d.current = ""
	return d
}

func (d *Display) Clear() {
}

func Reset() {
	termbox.Close()
}

func Exit() {
	Log("Exiting now")
	Reset()
	os.Exit(0)
}

func (d *Display) Poll() {
	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventInterrupt:
			Exit()
		case termbox.EventKey:
			switch e.Ch {
			case 'q':
				Exit()
			}
		}
	}
}

func (d *Display) Escape(code string) {
	fmt.Printf("\033%s", code)
}

func init() {
	termbox.Init()
	d = NewDisplay()
	d.width, d.height = termbox.Size()
	Log = makeLog()
}

func main() {
	defer Reset()
	d.Show("lllo,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Hello,Helo,Hello,Hellohhhhhh,\n世a界!")
	d.Draw()
	d.Poll()
}
