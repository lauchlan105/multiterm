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
	matrix  *Matrix
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
	matrix, err := makeMatrix(str, t.PopupWidth, -1, t.bg, false)

	if err != "" {
		matrix = newMatrix()
	}

	//Create new popup
	newPopup := &Popup{
		manager: t,
		id:      strconv.Itoa(time.Now().Second()),
		width:   matrix.Width,
		height:  matrix.Height,
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

	newPopup.matrix.X = newPopup.X
	newPopup.matrix.Y = newPopup.Y

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
			popup.matrix.Y += p.height + 1
		} else {
			popup.Y -= p.height + 1
			popup.matrix.Y -= p.height + 1
		}
	}

	p.manager.popups[p.id] = p

	go func() {
		for p.active {
			p.manager.updateBuffer(p.matrix)
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
