package multiterm

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

//Terminal asdf
type Terminal struct {
	height     int
	width      int
	fg         termbox.Attribute
	bg         termbox.Attribute
	tabs       map[string]Tab
	activeTabs []*Tab
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
	tab.Open()

	return terminal, terminal.tabs[tab.id]

}

//Start asdf
func (t *Terminal) Start() {

	t.resize()

	go func() {
		for {
			switch e := termbox.PollEvent(); e.Type {
			case termbox.EventResize:

			}

		}
	}()

}

func (t *Terminal) resize() {
	t.width, t.height = termbox.Size()
}

//Stop asdf
func (t *Terminal) Stop() {

}

////////////////////////////////////
// TAB CREATION AND MANIPULATION  //
////////////////////////////////////

//NewTab Generates and creates a new tab
func (t *Terminal) NewTab() Tab {

	//Create new tab
	tab := Tab{
		manager: t,
		id:      t.generateTabID(),
		name:    "Untitled",
	}

	//Add to terminal
	t.tabs[tab.id] = tab

	//return
	return tab

}

//removeTab removes a tab from the terminal
//	@Param tab id (int)
func (t *Terminal) removeTab(id string) {
	delete(t.tabs, id)
}

func (t *Terminal) generateTabID() string {

	nextID := 0

	for range t.tabs {
		if _, ok := t.tabs[strconv.Itoa(nextID)]; ok {
			return strconv.Itoa(nextID)
		}
		nextID++
	}

	return strconv.Itoa(nextID)
}
