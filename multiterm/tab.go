package multiterm

import (
	"io"
	"os/exec"
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
	visibleBuffer *Matrix
	cmd           *exec.Cmd
	stdin         pipe
	stdout        pipe
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
		t.visibleBuffer = newMatrix()
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
	if matrix.Height > windowHeight {

		startOfWindow := matrix.Height - t.scrollHeight - windowHeight
		endOfWindow := startOfWindow + windowHeight - 1

		if t.scrollHeight < 0 {
			endOfWindow += t.scrollHeight
		}

		startOutOfBounds := startOfWindow < 0 || startOfWindow > matrix.Height-1
		endOutOfBounds := endOfWindow < 0 || endOfWindow > matrix.Height-1

		if endOutOfBounds || startOutOfBounds {
			return
		}

		matrix.Content = matrix.Content[startOfWindow:endOfWindow]

	}

	t.manager.updateBuffer(matrix)
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
	if t.scrollHeight < t.visibleBuffer.Height-t.manager.height {
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

	if matrix, err := makeMatrix(t.buffer, width, -1, t.manager.bg, false); err != "" {
		t.manager.Print(err)
	} else {
		matrix.X = t.startX
		matrix.Y = 0
		t.visibleBuffer = matrix
	}

	/*
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
	*/

}

//RunCommand blah blah blah
func (t *Tab) RunCommand(s []string) error {
	cmd := exec.Command(s[0], s[1:]...)

	cmd.Stdout = t.stdout.w
	cmd.Stdin = t.stdin.r

	return cmd.Run()
}
