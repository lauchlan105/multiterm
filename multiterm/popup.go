package multiterm

import (
	"time"

	"github.com/nsf/termbox-go"
)

//Popup blah
type Popup struct {
	width  int
	height int

	X int
	Y int

	content     string
	matrix      [][]termbox.Cell
	popupTime   int
	timeoutChan chan bool
}

//Position determines which corner the popups spawn
type Position int

const (
	topLeft Position = iota
	topRight
	bottomLeft
	bottomRight
)

//Print displays the content as a popup.
//Background color defined by terminal.popupDefaultColor
func (t *Terminal) Print(content string) {
	t.print(content, t.PopupDefaultColor)
}

//Error displays the content as a popup.
//Background color defined by terminal.popupErrorColor
func (t *Terminal) Error(content string) {
	t.print(content, t.PopupErrorColor)
}

func (t *Terminal) print(str string, bgColor termbox.Attribute) {

	matrix := make([][]termbox.Cell, 0)

	//Create first line
	matrix = append(matrix, make([]termbox.Cell, 0))
	matrix = append(matrix, make([]termbox.Cell, 0))
	longestRow := 0

	//Convert string into cell matrix
	for i, ch := range str {

		currentRow := &matrix[len(matrix)-1]

		//if a new row => insert first padding rune
		if len(*currentRow) == 0 {
			newCell := termbox.Cell{Ch: ' ', Bg: bgColor}
			*currentRow = append(*currentRow, newCell)
		}

		//Move to next line on \n char
		//OR force newline if the current row is larger than
		//the max width defined under terminal.popupWidth
		if ch == '\n' || len(*currentRow) > (t.width/100)*t.PopupWidth {

			//Re-check for longest row
			currentRowLength := len(*currentRow)
			if currentRowLength > len(matrix[longestRow]) {
				longestRow = i
			}

			matrix = append(matrix, make([]termbox.Cell, 0))
			continue
		}

		//Append current ch to end of the last row
		newCell := termbox.Cell{Ch: ch, Bg: bgColor}
		*currentRow = append(*currentRow, newCell)

	}

	//Add last line of padding
	matrix = append(matrix, make([]termbox.Cell, 0))

	//Fill out the popup
	for _, row := range matrix {

		toAppend := longestRow - len(row)

		for toAppend > 0 {
			row = append(row, termbox.Cell{Ch: ' ', Bg: bgColor})
			toAppend--
		}

	}

	//Create new popup
	newPopup := &Popup{
		width:       longestRow,
		height:      len(matrix),
		matrix:      matrix,
		popupTime:   t.PopupTime,
		timeoutChan: make(chan bool, 1),
	}

	switch t.PopupPosition {
	case topLeft:
		newPopup.X = 1
		newPopup.Y = 1
		break
	case topRight:
		newPopup.X = t.width - newPopup.width - 2
		newPopup.Y = 1
		break
	case bottomLeft:
		newPopup.X = 1
		newPopup.Y = t.height - newPopup.height - 2
		break
	case bottomRight:
		newPopup.X = t.width - newPopup.width - 2
		newPopup.Y = t.height - newPopup.height - 2
		break
	}

	//append newPopup to start of popups slice
	t.popups = append([]*Popup{newPopup}, t.popups...)

}

func (p *Popup) kill() {
	go func(p *Popup) {
		for p.popupTime > 0 {
			time.Sleep(time.Second)
			p.popupTime--
		}
		p.timeoutChan <- false
	}(p)
}

func (p *Popup) print() {
	go func() {
		for <-p.timeoutChan {

		}
	}()
}
