// input.go

package input

import (
	"fmt"

	"github.com/charlesetc/pat/display"
)

var currentLine []rune
var cursor int
var currentScroll int
var Log func([]rune)
var LogS func(string)

func AddRune(r rune) {
	Log([]rune(fmt.Sprintf("%d, %d", len(currentLine), cursor)))
	if len(currentLine) <= cursor {
		currentLine = append(currentLine, r)
		cursor++
		checkScrollRight()
		return
	}
	copyLine := make([]rune, len(currentLine))
	copy(copyLine, currentLine)
	currentLine = append(append(copyLine[:cursor], r),
		currentLine[cursor:]...)
	if cursor-currentScroll >= display.LineWidth() {
		checkScrollRight()
	}
	cursor++
}

func checkScrollRight() {
	width := display.LineWidth()
	if len(currentLine)-currentScroll > width {
		currentScroll++
	}
}

func checkScrollLeft() {
	if len(currentLine)-currentScroll < 0 {
		currentScroll--
	}
}

func visibleLine() []rune {
	width := display.LineWidth()

	if len(currentLine) <= width {
		return currentLine
	}
	if len(currentLine) <= width+currentScroll {
		return currentLine[currentScroll:]
	}
	return currentLine[currentScroll : width+currentScroll]
}

func updateCursor() {
	display.SetCursor(cursor - currentScroll)
	display.DrawLine()
}

func CursorLeft() {
	if cursor != 0 {
		cursor--
	}
	if cursor < currentScroll {
		currentScroll--
		display.ShowLine(visibleLine())
	}
	updateCursor()
}

func CursorRight() {
	width := display.LineWidth()
	if cursor < len(currentLine) {
		cursor++
	}
	if cursor > width+currentScroll {
		LogS(fmt.Sprintf("-- %d", cursor))
		currentScroll++
		display.ShowLine(visibleLine())
	}
	updateCursor()
}

func Backspace() {
	if cursor == 0 {
		return
	}
	width := display.LineWidth()
	if cursor-currentScroll <= 0 && currentScroll != 0 {
		currentScroll -= width
		updateCursor()
	}
	if len(currentLine) <= cursor {
		currentLine = currentLine[:cursor-1]
		cursor--
	} else {
		currentLine = append(currentLine[:cursor-1], currentLine[cursor:]...)
		cursor--
	}

	if cursor == width {
		currentScroll = 0
		updateCursor()
	}
}

func Runes() []rune {
	return currentLine
}

func Reset() {
	currentLine = []rune{}
	currentScroll = 0
	cursor = 0
	updateCursor()
	display.ShowLine(visibleLine())
	display.DrawLine()
}

func CursorAtBeginning() {
	cursor = 0
	currentScroll = 0
	display.ShowLine(visibleLine())
	updateCursor()
}
func CursorAtEnd() {
	width := display.LineWidth()
	cursor = len(currentLine)
	if len(currentLine) < width {
		currentScroll = 0
	} else {
		currentScroll = len(currentLine) - width
	}
	display.ShowLine(visibleLine())
	updateCursor()
}

func Draw() {
	display.ShowLine(visibleLine())
	display.SetCursor(cursor - currentScroll)
	display.DrawLine()
	display.Flush()
}
