package multiterm

import (
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

//Popup blah
type Popup struct {
	manager *Terminal

	width  int
	height int

	X int
	Y int

	id      string
	active  bool
	content string
	matrix  [][]termbox.Cell
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

	str = str + "\n"
	matrix := make([][]termbox.Cell, 0)

	//Create first line
	matrix = append(matrix, make([]termbox.Cell, 0))
	matrix = append(matrix, make([]termbox.Cell, 0))
	longestRow := 0

	//Convert string into cell matrix
	for _, ch := range str {

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
			if len(*currentRow) > longestRow {
				longestRow = len(*currentRow)
			}

			matrix = append(matrix, make([]termbox.Cell, 0))
			continue
		}

		//Append current ch to end of the last row
		newCell := termbox.Cell{Ch: ch, Bg: bgColor}
		*currentRow = append(*currentRow, newCell)

	}

	//Fill out the popup
	for i, row := range matrix {

		toAppend := longestRow - len(row)

		for toAppend >= 0 {
			matrix[i] = append(matrix[i], termbox.Cell{Ch: ' ', Bg: bgColor})
			toAppend--
		}

	}

	//Create new popup
	newPopup := &Popup{
		manager: t,
		id:      strconv.Itoa(time.Now().Second()),
		width:   longestRow,
		height:  len(matrix),
		matrix:  matrix,
		content: str,
		active:  false,
	}

	//Set position based vars
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
		newPopup.Y = t.height - newPopup.height
		break
	default: //default is also bottomRight
		newPopup.X = t.width - newPopup.width - 2
		newPopup.Y = t.height - newPopup.height
		break
	}

	newPopup.print()

}

func (p *Popup) print() {

	defer p.kill()
	p.active = true

	//Shift all other popups
	for _, popup := range p.manager.popups {
		pos := p.manager.PopupPosition
		if pos == topLeft || pos == topRight {
			popup.Y += p.height + 1
		} else {
			popup.Y -= p.height + 1
		}
	}

	p.manager.popups[p.id] = p

	go func() {
		for p.active {

			for r, row := range p.matrix {

				destRow := p.manager.width * (p.Y + r)

				for c, cell := range row {

					destCol := p.X + c
					dest := destRow + destCol
					if dest < 0 || dest >= len(p.manager.buffer) {
						continue
					}
					p.manager.buffer[dest] = cell

				}

			}

			p.manager.updateBuffer()
			time.Sleep(250 * time.Millisecond)
		}
	}()

}

func (p *Popup) kill() {
	go func(p *Popup) {
		time.Sleep(time.Duration(p.manager.PopupTime) * time.Second)
		p.active = false
		delete(p.manager.popups, p.id)
		p.manager.printAll()
	}(p)
}
