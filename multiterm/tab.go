package multiterm

import (
	"strings"

	termbox "github.com/nsf/termbox-go"
)

//Tab asdf
type Tab struct {
	manager *Terminal
	id      string
	name    string
	active  bool
	buffer  []string
}

//Terminate kills the current tab
//via the Terminal objects removeTab(id) func
func (t *Tab) Terminate() {
	t.Close()
	t.manager.removeTab(t.id)
}

//Open sets the tab's active state to true
// and adds it to the managers activeTabs slice
func (t *Tab) Open() {
	if t.active {
		return
	}
	t.active = true

	t.manager.activeTabs = append(t.manager.activeTabs, t)
	t.manager.printAll()
}

//Close sets the tab's active state to false
// and removes it from the managers activeTabs slice
func (t *Tab) Close() {
	if !t.active {
		return
	}
	t.active = false

	//Iterate tabs and break when tab is found
	for i, tab := range t.manager.activeTabs {
		if tab.id == t.id {
			t.manager.activeTabs = append(t.manager.activeTabs[:i],
				t.manager.activeTabs[i+1:]...)
			break
		}
	}
	t.manager.printAll()
}

func (t *Tab) print(indexOfTab int) {

	width := t.manager.tabWidth - 1 //width minus seperator
	startingX := indexOfTab * t.manager.tabWidth

	combinedBuffer := strings.Join(t.buffer, "")

	matrix := make([][]termbox.Cell, 0)
	matrix = append(matrix, make([]termbox.Cell, 0))

	for _, ch := range combinedBuffer {

		appendChar := func(ch rune) {

			//if last row doesn't exist -> make it
			if len(matrix) == 0 {
				matrix = append(matrix, make([]termbox.Cell, 0))
			}

			//if last row is full
			//create a new row
			if len(matrix[len(matrix)-1]) == width {
				matrix = append(matrix, make([]termbox.Cell, 0))
			}

			//Get last row
			rowInd := len(matrix) - 1

			//Append char to current row
			matrix[rowInd] = append(matrix[rowInd], termbox.Cell{Ch: ch})

			//If the current char is a newline,
			//append whitespaces to the rest of the row
			if ch == '\n' {
				numberOfWhiteSpaces := width - len(matrix[rowInd])
				for numberOfWhiteSpaces > 0 {
					matrix[rowInd] = append(matrix[rowInd], termbox.Cell{Ch: ' '})
					numberOfWhiteSpaces--
				}
			}
		}
		appendChar(ch)

	}

	for rowIndex, row := range matrix {
		for colIndex, Ch := range row {
			y := rowIndex * t.manager.width
			x := colIndex + startingX
			dest := y + x
			t.manager.buffer[dest] = Ch
		}
	}

}
