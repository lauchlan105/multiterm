package multiterm

import (
	"github.com/lauchlan105/multiterm/tab"
	"github.com/nsf/termbox-go"
)

//Terminal asdf
type Terminal struct {
	height int
	width  int
	fg     termbox.Attribute
	bg     termbox.Attribute
	tabs   []tab.Tab
}

//Init returns
func Init(textColor, backgroundColor termbox.Attribute) (terminal Terminal, tab tab.Tab) {

	terminal = Terminal{
		height: 0,
		width:  0,
		fg:     textColor,
		bg:     backgroundColor,
	}
	tab = terminal.NewTab()

	return terminal, terminal.tabs[0]

}

//NewTab creates and returns a new tab. (Initialised but not open)
func (t *Terminal) NewTab() tab.Tab {
	return tab.Tab{}
}

//CloseTab closes a tab
func (t *Terminal) CloseTab() {

}
