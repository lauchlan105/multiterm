package multiterm

import (
	"io"
	"os/exec"

	termbox "github.com/nsf/termbox-go"
)

//Tab asdf
type Tab struct {
	active  bool
	id      string
	manager *Terminal
	Name    string

	endX   int
	startX int

	buffer        string
	scrollHeight  int
	visibleBuffer [][]termbox.Cell

	cmd    *exec.Cmd
	stdin  pipe
	stdout pipe
}

type pipe struct {
	r *io.PipeReader
	w *io.PipeWriter
}

//Terminate kills the current tab
//via the Terminal objects removeTab(id) func
func (t *Tab) Terminate() {

	//Close tabs
	t.Close()

	//Close io pipes
	t.stdin.w.Close()
	t.stdin.r.Close()
	t.stdout.w.Close()
	t.stdout.r.Close()

	//Remove tab from manager
	t.manager.removeTab(t.id)

}

//Open sets the tab's active state to true
// and adds it to the managers activeTabs slice
func (t *Tab) Open() {
	if t.active {
		return
	}
	t.active = true
	t.scrollHeight = 0

	//Write any input to the input pipe
	/*
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				b, _ := reader.ReadBytes('\n')
				t.stdin.w.Write(b)
			}
		}()
	*/

	if t.visibleBuffer == nil {
		t.visibleBuffer = make([][]termbox.Cell, 0)
	}

	t.manager.activeTabs = append(t.manager.activeTabs, t)
	t.manager.printAll()
}

//Close sets the tab's active state to false
// and removes it from the managers activeTabs slice
func (t *Tab) Close() {
	if !t.active {
		return
	}
	t.active = false

	//Remove t from active tabs slice
	for i, tab := range t.manager.activeTabs {
		if tab.id == t.id {
			t.manager.activeTabs = append(t.manager.activeTabs[:i],
				t.manager.activeTabs[i+1:]...)
			break
		}
	}
	t.manager.printAll()
}

func (t *Tab) printTab() {

	if !t.active {
		return
	}

	windowHeight := t.manager.height

	t.setVisibleBuffer()
	matrix := t.visibleBuffer

	//if the matrix has more lines than
	//the window height allows -> remove all lines
	//outside the frame (based off scrollheight)
	if len(matrix) > windowHeight {

		startOfWindow := len(matrix) - t.scrollHeight - windowHeight
		endOfWindow := startOfWindow + windowHeight - 1

		if t.scrollHeight < 0 {
			endOfWindow += t.scrollHeight
		}

		startOutOfBounds := startOfWindow < 0 || startOfWindow > len(matrix)-1
		endOutOfBounds := endOfWindow < 0 || endOfWindow > len(matrix)-1

		if endOutOfBounds || startOutOfBounds {
			return
		}

		matrix = matrix[startOfWindow:endOfWindow]

	}

	for rowIndex, row := range matrix {
		for colIndex, Ch := range row {
			y := rowIndex * t.manager.width
			x := colIndex + t.startX
			dest := y + x
			t.manager.buffer[dest] = Ch
		}
	}

}

//Print prints to the terminal without appending a new line char
func (t *Tab) Print(str string) {
	t.buffer += str
	if t.active {
		t.manager.printAll()
	}
}

//Println prints a line to the terminal
func (t *Tab) Println(str string) {
	t.Print(str + "\n")
}

//ScrollUp asdf
func (t *Tab) ScrollUp() {
	if t.scrollHeight < len(t.visibleBuffer)-t.manager.height {
		t.scrollHeight++
	}
}

//ScrollDown asdf
func (t *Tab) ScrollDown() {
	if t.manager.height-(t.scrollHeight*-1) > 2 {
		t.scrollHeight--
	}
}

//adjustVisibileBuffer shifts characters
//around to match the tab's width
func (t *Tab) setVisibleBuffer() {

	width := t.endX - t.startX
	matrix := make([][]termbox.Cell, 0)
	matrix = append(matrix, make([]termbox.Cell, 0))

	for _, ch := range t.buffer {

		//if last row doesn't exist -> make it
		if len(matrix) == 0 {
			matrix = append(matrix, make([]termbox.Cell, 0))
		}

		//if last row is full
		//create a new row
		if len(matrix[len(matrix)-1]) == width {
			if ch == '\n' {
				return
			}
			matrix = append(matrix, make([]termbox.Cell, 0))
		}

		//Get index of last row
		rowInd := len(matrix) - 1

		//Append char to current matrix row
		matrix[rowInd] = append(matrix[rowInd], termbox.Cell{Ch: ch})

		//If the current char is a newline,
		//append whitespaces to the rest of the row
		if ch == '\n' {
			numberOfWhiteSpaces := width - len(matrix[rowInd])
			for numberOfWhiteSpaces > 0 {
				matrix[rowInd] = append(matrix[rowInd], termbox.Cell{Ch: ' '})
				numberOfWhiteSpaces--
			}
		}

	}

	t.visibleBuffer = matrix
}

//RunCommand blah blah blah
func (t *Tab) RunCommand(s []string) error {
	cmd := exec.Command(s[0], s[1:]...)

	cmd.Stdout = t.stdout.w
	cmd.Stdin = t.stdin.r

	return cmd.Run()
}
