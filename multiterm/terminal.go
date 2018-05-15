package multiterm

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

//Terminal asdf
type Terminal struct {
	height    int
	width     int
	fg        termbox.Attribute
	bg        termbox.Attribute
	splitRune rune
	splitCell termbox.Cell

	tabs       map[string]Tab
	activeTabs []*Tab
	focus      *Tab
	buffer     []termbox.Cell
	sepIndexes []int
}

//Init returns
func Init() (terminal Terminal, tab Tab) {

	terminal = Terminal{
		fg: termbox.ColorDefault,
		bg: termbox.ColorDefault,
		splitCell: termbox.Cell{
			Bg: termbox.ColorWhite,
		},

		tabs:       make(map[string]Tab, 0),
		activeTabs: make([]*Tab, 0),
	}
	tab = terminal.NewTab()
	tab.Open()

	return terminal, terminal.tabs[tab.id]

}

////////////////////////////////
// TERMINAL MANAGER FUNCTIONS //
////////////////////////////////

//Start asdf
func (t *Terminal) Start() {

	//Init termbox
	termbox.Init()

	//Set the terminal objects that rely
	//on the initialised termbox
	t.updateSize()
	t.buffer = make([]termbox.Cell, t.height*t.width)

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			if e.Key == 3 { //Ctrl+C
				return
			}
		case termbox.EventResize:
			t.updateSize()
			t.printTabs()
		case termbox.EventMouse:
			if e.Key == termbox.MouseLeft {
				newTab := t.NewTab()
				newTab.Open()
				t.printTabs()
			}
		}
	}

}

//Stop executes the necessary functions to
//shut down the Multi-terminal properly
func (t *Terminal) Stop() {
	termbox.Close()
}

func (t *Terminal) updateSize() {
	t.width, t.height = termbox.Size()
	t.height += 5
}

func (t *Terminal) printTabs() {

	t.buffer = make([]termbox.Cell, t.width*t.height)

	numOfSeps := len(t.activeTabs) - 1

	if numOfSeps > 0 {
		sepDist := (t.width) / len(t.activeTabs)

		for h := 0; h < t.height; h++ {
			for sep := 0; sep < numOfSeps; sep++ {
				row := h * t.width
				col := (sep + 1) * sepDist
				t.buffer[row+col] = t.splitCell
			}
		}

	}

	//Update and flush termbox buffer
	copy(termbox.CellBuffer(), t.buffer)
	if err := termbox.Flush(); err != nil {

	}

}

/////////////////////////////////////
//  TAB CREATION AND MANIPULATION  //
/////////////////////////////////////

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

//NumTabs asdf
func (t *Terminal) NumTabs() string {
	return strconv.Itoa(len(t.tabs))
}
