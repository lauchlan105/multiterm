package multiterm

import (
	"github.com/nsf/termbox-go"
)

//Terminal asdf
type Terminal struct {
	height int
	width  int
	fg     termbox.Attribute
	bg     termbox.Attribute
	tabs   []Tab
}

//Init returns
func Init(textColor, backgroundColor termbox.Attribute) (terminal Terminal, tab Tab) {

	terminal = Terminal{
		height: 0,
		width:  0,
		fg:     textColor,
		bg:     backgroundColor,
	}
	tab = terminal.NewTab()

	return terminal, terminal.tabs[0]

}

//CloseTab closes a tab
func (t *Terminal) CloseTab() {

}
