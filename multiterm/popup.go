package multiterm

import (
	"github.com/nsf/termbox-go"
)

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
func (t *Terminal) Print(content interface{}) {
	t.print(content.(string), t.popupDefaultColor)
}

//Error displays the content as a popup.
//Background color defined by terminal.popupErrorColor
func (t *Terminal) Error(content interface{}) {
	t.print(content.(string), t.popupErrorColor)
}

func (t *Terminal) print(str string, bgColor termbox.Attribute) {

	fmtStr := make([][]rune, 0)

	//Create first line
	fmtStr = append(fmtStr, make([]rune, 0))
	longestRow := 0

	for i, ch := range str {

		if ch == '\n' {

			//Re-check for longest row
			currentRowLength := len(fmtStr[len(fmtStr)-1])
			if currentRowLength > len(fmtStr[longestRow]) {
				longestRow = i
			}

			fmtStr = append(fmtStr, make([]rune, 0))
			continue
		}

		//Append current ch to end of the last row
		fmtStr[len(fmtStr)-1] = append(fmtStr[len(fmtStr)-1], ch)

	}

}
