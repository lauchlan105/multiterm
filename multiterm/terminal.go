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
	splitCell termbox.Cell

	tabs       map[string]Tab
	activeTabs []*Tab
	focus      *Tab
	buffer     []termbox.Cell
	sepIndexes []int
	tabWidth   int
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

	newTab := terminal.NewTab()

	return terminal, newTab

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
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			if e.Key == 3 {
				return
			}
		case termbox.EventResize:
			t.printAll()
		case termbox.EventMouse:
			if e.Key == termbox.MouseLeft {
				newTab := t.NewTab()
				newTab.Open()
				t.focus = &newTab
			} else if e.Key == termbox.MouseMiddle {
				t.focus.Terminate()
				t.focus = t.activeTabs[0]
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

	//Save new tab width (this includes the sep)
	if len(t.activeTabs) <= 0 {
		t.tabWidth = t.width
	} else {
		t.tabWidth = t.width / len(t.activeTabs)
	}

	//clear the saved seperator indexes
	t.sepIndexes = make([]int, 0)

	//Push new seperator indexes to sepIndexes slice
	numSeps := len(t.activeTabs) - 1
	for sep := 1; sep <= numSeps; sep++ {

		// '-1' offsets to start printing from 0 (not 1)
		x := (sep * t.tabWidth) - 1
		t.sepIndexes = append(t.sepIndexes, x)

	}

}

func (t *Terminal) printAll() {

	//Update printing-related variables
	t.updateSize()

	//Clear buffer
	t.buffer = make([]termbox.Cell, t.width*t.height)

	//Print actual chars
	t.printTabs()
	t.printSeps()

	//Copy buffer to termbox
	t.updateBuffer()

}

//print all tabs
func (t *Terminal) printTabs() {
	for indexOfTab, tab := range t.activeTabs {
		tab.print(indexOfTab)
	}
}

//print all tab seperators
func (t *Terminal) printSeps() {

	//Print seperators
	for _, x := range t.sepIndexes {

		// '-1' offsets to start printing from 0 (not 1)
		for h := 0; h < t.height; h++ {
			row := h * t.width
			dest := row + x
			if dest < len(t.buffer) {
				t.buffer[row+x] = t.splitCell
			}
		}

	}
}

//copy t.buffer to termbox buffer
func (t *Terminal) updateBuffer() {
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
		buffer:  make([]string, 0),
	}

	tab.buffer = append(tab.buffer, []string{
		"ID: " + tab.id + "\n",
		"Title: " + tab.name + "\n",
	}...)

	//Add to terminal
	t.tabs[tab.id] = tab

	//return
	return t.tabs[tab.id]

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
