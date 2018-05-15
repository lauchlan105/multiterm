package multiterm

import (
	"fmt"
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
}

//Init returns
func Init() (terminal Terminal, tab Tab) {

	defaultSplitRune := '|'

	terminal = Terminal{
		fg:        termbox.ColorDefault,
		bg:        termbox.ColorDefault,
		splitRune: defaultSplitRune,
		splitCell: termbox.Cell{
			Ch: '|',
			Fg: termbox.ColorDefault,
			Bg: termbox.ColorDefault,
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

	termbox.Init()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	t.resize()

	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			if e.Key == 3 { //Ctrl+C
				return
			}
		case termbox.EventResize:
			fmt.Println("Resized")
			t.resize()
		case termbox.EventMouse:
			if e.Key == termbox.MouseLeft {
				newTab := t.NewTab()
				newTab.Open()
				fmt.Println("Added tab")
				fmt.Println(len(t.activeTabs))
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

func (t *Terminal) resize() {
	t.width, t.height = termbox.Size()
}

func (t *Terminal) printTabs() {

	numOfSeps := len(t.activeTabs) - 1

	if numOfSeps > 0 {
		sepDist := t.width / len(t.activeTabs)

		for h := 0; h < t.height; h += t.width {
			for sep := 0; sep < numOfSeps; sep++ {
				t.buffer[h+(sepDist*sep)] = t.splitCell
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
