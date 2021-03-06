package multiterm

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nsf/termbox-go"
)

//Terminal asdf
type Terminal struct {

	//Terminal Attributes
	height    int
	width     int
	fg        termbox.Attribute
	bg        termbox.Attribute
	splitCell termbox.Cell

	//Tab Attributes
	tabs       map[string]Tab
	activeTabs []*Tab
	focus      *Tab
	buffer     []termbox.Cell
	stopChan   chan bool

	//Popup Attributes
	PopupDefaultColor termbox.Attribute
	PopupErrorColor   termbox.Attribute
	PopupPosition     Position
	PopupWidth        int //percentage of terminal
	PopupTime         int //seconds
	popups            map[string]*Popup
}

//Init returns
func Init() (terminal Terminal) {

	return Terminal{
		fg: termbox.ColorDefault,
		bg: termbox.ColorDefault,
		splitCell: termbox.Cell{
			Bg: termbox.ColorWhite,
		},

		tabs:       make(map[string]Tab, 0),
		activeTabs: make([]*Tab, 0),
		stopChan:   make(chan bool, 1),

		PopupDefaultColor: termbox.ColorCyan,
		PopupErrorColor:   termbox.ColorRed,
		PopupPosition:     bottomRight,
		PopupWidth:        70,
		PopupTime:         5,
		popups:            make(map[string]*Popup),
	}

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

	//Begin event listeners
	go func() {
		for {

			switch e := termbox.PollEvent(); e.Type {

			case termbox.EventKey:

				//Ctrl + C
				if e.Key == 3 {
					t.Stop()
				}

			case termbox.EventMouse:

				//Switch depending on which mouse event happened
				switch key := e.Key; key {

				case termbox.MouseLeft:
					newTab := t.NewTab()
					newTab.Open()
					t.focus = newTab

				case termbox.MouseRight:

				case termbox.MouseMiddle:
					if tab := t.getMouseFocus(e.MouseX); tab != nil &&
						len(t.activeTabs) > 1 {
						tab.Close()
						t.focus = t.activeTabs[len(t.activeTabs)-1]
					}

				case termbox.MouseWheelUp:
					if tab := t.getMouseFocus(e.MouseX); tab != nil {
						tab.ScrollUp()
					}

				case termbox.MouseWheelDown:
					if tab := t.getMouseFocus(e.MouseX); tab != nil {
						tab.ScrollDown()
					}

				}

				//Update all
				t.printAll()

			}

		}
	}()

	//Listen for system signals
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		t.stopChan <- true
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		for {
			if len(t.popups) == 0 {
				t.activeTabs[0].Println("was 0")
			} else {
				t.activeTabs[0].Println(strconv.Itoa(len(t.popups)))
			}
			time.Sleep(2 * time.Second)
		}
	}()

}

//Wait for Stop() to be called
func (t *Terminal) Wait() {
	for <-t.stopChan {
		return
	}
}

//Stop executes the necessary functions to
//shut down the Multi-terminal properly
func (t *Terminal) Stop() {
	termbox.Close()
	for _, tab := range t.tabs {
		tab.Terminate()
	}
	t.stopChan <- true
}

//Updates width and height vars and
//tab widths
func (t *Terminal) updateSize() {

	t.width, t.height = termbox.Size()
	numTabs := len(t.activeTabs)

	if numTabs > 0 {

		//divides width of terminal excluding spacers
		avgTabWidth := (t.width - (numTabs - 1)) / numTabs
		leftOverSpace := (t.width - (numTabs - 1)) % numTabs

		startX := 0
		endX := -1
		for i, tab := range t.activeTabs {
			startX = endX + 1
			endX = startX + avgTabWidth

			//Distribute the left over space
			if i < leftOverSpace {
				endX++
			}

			tab.startX = startX
			tab.endX = endX

		}

	}

}

//Call update size and then print tabs and seperators
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

	//Asynchronously print tabs
	blockChan := make(chan bool, 1)
	waitingOn := len(t.activeTabs)
	for _, tab := range t.activeTabs {
		go func(tab *Tab) {
			tab.printTab()
			blockChan <- true
		}(tab)
	}

	//wait for all to print
	//before exiting func
	for <-blockChan {
		waitingOn--
		if waitingOn == 0 {
			return
		}
	}
}

//print all tab seperators
func (t *Terminal) printSeps() {

	if len(t.activeTabs) > 1 {

		//tabs that start after a seperator
		relativeTabs := t.activeTabs[1:]

		for _, tab := range relativeTabs {
			for h := 0; h < t.height; h++ {

				/*
				 * row is the index within buffer
				 * where the row starts. Not the actual row
				 */
				row := h * t.width
				col := tab.startX - 1
				dest := row + col

				//if out of bounds
				if dest < 0 || dest >= len(t.buffer) {
					continue
				}

				t.buffer[dest] = t.splitCell

			}
		}
	}

}

////////////////////////////////
//  BUFFER RELATED FUNCTIONS  //
////////////////////////////////

//copy t.buffer to termbox buffer
func (t *Terminal) updateBuffer() {
	copy(termbox.CellBuffer(), t.buffer)
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
}

/////////////////////////////////////
//  TAB CREATION AND MANIPULATION  //
/////////////////////////////////////

//NewTab Generates and creates a new tab
func (t *Terminal) NewTab() *Tab {

	//Create stdin/stdout pipes
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()

	//Create new tab
	tab := Tab{
		manager: t,
		id:      t.generateTabID(),
		Name:    "Untitled",
		stdin: pipe{
			r: inR,
			w: inW,
		},
		stdout: pipe{
			r: outR,
			w: outW,
		},
	}

	//Print anything printed to the output pipe
	go func() {
		scanner := bufio.NewScanner(tab.stdout.r)
		for scanner.Scan() {
			tab.Println("Output: " + scanner.Text())
		}
	}()

	tab.Println("ID: " + tab.id)
	tab.Println("Title: " + tab.Name)

	//Add to terminal
	t.tabs[tab.id] = tab

	//return
	return &tab

}

//removeTab removes a tab from the terminal
//	@Param tab id (int)
func (t *Terminal) removeTab(id string) {
	delete(t.tabs, id)
}

//Returns an available id based off existing tab ids
func (t *Terminal) generateTabID() string {

	nextID := 0
	for range t.tabs {
		if _, found := t.tabs[strconv.Itoa(nextID)]; !found {
			return strconv.Itoa(nextID)
		}
		nextID++
	}
	return strconv.Itoa(nextID)

}

//Return which tab is under the mouse
func (t *Terminal) getMouseFocus(mouseX int) *Tab {

	for _, tab := range t.activeTabs {
		if mouseX < tab.endX {
			return tab
		}
	}
	return nil
}
