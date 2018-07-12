package multiterm

import (
	"github.com/nsf/termbox-go"
)

//Position a
type Position int

const (
	topLeft Position = iota
	topRight
	bottomLeft
	bottomRight
)

//Print displays the content as a popup.
func (t *Terminal) Print(content interface{}) {

	t.print(content.(string), t.popupDefaultColor)

}

//Error displays the content as a popup.
func (t *Terminal) Error(content interface{}) {

	t.print(content.(string), t.popupErrorColor)

}

func (t *Terminal) print(str string, bgColor termbox.Attribute) {

	fmtStr := make([][]string, 0)

	//Create first line
	fmtStr[0] = make([]string, 0)

	for i, ch := range str {

		if ch == '\n' {

		}

	}

}
