// display.go

package display

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var currentFile, currentLine []rune
var width, height, currentScroll, currentCursor int
var topBar bool
var Log func([]rune)

// Ranges over the currentFile and
// returns the positions of each rune.
func overLines(str []rune, f func(int, int, int, rune)) {
	var x, y int
	for i, r := range str {
		f(i, x, y, r)
		if r == '\n' || x >= width {
			y++
			x = 0
		} else { // not a newline
			x++
		}
	}
}

func viewHeight() int {
	return height - 3 // right now, just one extra line.
}

func Draw() {
	termbox.Clear(0, 0)
	DrawBox()
	var offset int
	if topBar {
		offset = 4
	}
	overLines(viewSection(), func(i, x, y int, r rune) {
		if r != '\n' {
			termbox.SetCell(x, y+offset, r, 0, 0)
		}
	})
	termbox.Flush()
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func DrawBox() {
	var offset int
	if topBar {
		offset = height - 4
	}
	edit_box_width := width - 4
	midx := 2
	midy := height - 2 - offset
	termbox.SetCell(midx-1, midy, '│', 0, 0)
	termbox.SetCell(midx+edit_box_width, midy, '│', 0, 0)
	termbox.SetCell(midx-1, midy-1, '┌', 0, 0)
	termbox.SetCell(midx-1, midy+1, '└', 0, 0)
	termbox.SetCell(midx+edit_box_width, midy-1, '┐', 0, 0)
	termbox.SetCell(midx+edit_box_width, midy+1, '┘', 0, 0)
	fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '─'})
	fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '─'})
	DrawLine()
}

func ScrollDown() {
	if currentScroll < numberOfLines() {
		currentScroll++
	}
	Draw()
}

func ScrollUp() {
	if currentScroll != 0 {
		currentScroll--
	}
	Draw()
}

func viewSection() []rune {
	var start, end int
	overLines(currentFile, func(i, x, y int, r rune) {
		if y == currentScroll && x == 0 {
			start = i
		}
		if y == currentScroll+viewHeight() && x == 0 {
			end = i - 1
		}
	})
	if end == 0 { // Didn't reach the end.
		end = len(currentFile) - 2
	}
	Log([]rune(fmt.Sprintf("%d - %d, %d", start, end, currentScroll)))
	return currentFile[start:end]
}

func numberOfLines() int {
	var a int
	overLines(currentFile, func(i, x, y int, r rune) {
		a = y
	})
	return a
}

func Reset() {
	termbox.Close()
}

func Resize() {
	termbox.Flush()
	width, height = termbox.Size()
}

func Init(bar bool) {
	termbox.Init()
	Resize()
	topBar = bar
}

func ShowLine(line []rune) {
	currentLine = line
}

func LineWidth() int {
	return width - 6 - 1
}

func DrawLine() {
	x := 3
	var offset int
	if topBar {
		offset = height - 4 //w/e
	}
	y := height - 2 - offset

	termbox.SetCursor(x+currentCursor, y)

	// ClearLine
	for i := x; i < width-2; i++ {
		termbox.SetCell(i, y, ' ', 0, 0)
	}

	// Fill it in
	for _, r := range currentLine {
		termbox.SetCell(x, y, r, 0, 0)
		x++
	}
	termbox.Flush()
}

func SetCursor(x int) {
	currentCursor = x
}

func Flush() {
	termbox.Flush()
}

func Show(main []rune, line []rune) {
	currentFile = main
	currentLine = line
}
