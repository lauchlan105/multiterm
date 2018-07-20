package multiterm

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

//Matrix is a 2D termbox cell array
//used to parse through to the printer function
type Matrix struct {
	X       int
	Y       int
	Height  int
	Width   int
	Content [][]termbox.Cell
}

func newMatrix() *Matrix {
	return &Matrix{
		X:       0,
		Y:       0,
		Height:  0,
		Width:   0,
		Content: make([][]termbox.Cell, 0),
	}
}

func makeMatrix(str string, maxWidth int, maxHeight int, bgColor termbox.Attribute, fillColor bool) (matrix *Matrix, err string) {

	matrix = newMatrix()
	matrix.Content = append(matrix.Content, make([]termbox.Cell, 0))
	longestRowIndex := -1

	for _, ch := range str {

		//if last row doesn't exist -> make it
		if len(matrix.Content) == 0 {
			matrix.Content = append(matrix.Content, make([]termbox.Cell, 0))
		}

		//get length of latest row
		currRowLength := len(matrix.Content[len(matrix.Content)-1])
		if longestRowIndex == -1 {
			longestRowIndex = len(matrix.Content) - 1
		} else if currRowLength > len(matrix.Content[longestRowIndex]) {
			longestRowIndex = len(matrix.Content) - 1
		}

		//if last row is full
		//create a new row
		if len(matrix.Content[len(matrix.Content)-1]) == maxWidth {
			if ch == '\n' {
				break
			}
			matrix.Content = append(matrix.Content, make([]termbox.Cell, 0))
		}

		//Get index of last row
		rowInd := len(matrix.Content) - 1

		//Append char to current matrix row
		matrix.Content[rowInd] = append(matrix.Content[rowInd], termbox.Cell{Ch: ch, Bg: bgColor})

		//If the current char is a newline,
		//append whitespaces to the rest of the row
		if ch == '\n' {
			numberOfWhiteSpaces := maxWidth - len(matrix.Content[rowInd])
			for numberOfWhiteSpaces > 0 {

				if fillColor {
					matrix.Content[rowInd] = append(matrix.Content[rowInd], termbox.Cell{Ch: ' ', Bg: bgColor})
				} else {
					matrix.Content[rowInd] = append(matrix.Content[rowInd], termbox.Cell{Ch: ' '})
				}

				numberOfWhiteSpaces--
			}
		}

	}

	creationSuccess := true

	//Set height and width
	matrix.Height = len(matrix.Content)
	if longestRowIndex >= 0 && longestRowIndex < len(matrix.Content) {
		matrix.Width = len(matrix.Content[longestRowIndex])
	} else {
		creationSuccess = false
		err += "longestRowIndex (" + strconv.Itoa(longestRowIndex) +
			") out of bounds of matrix (" + strconv.Itoa(len(matrix.Content)) + ")\n"
	}

	if matrix.Height > maxHeight && maxHeight != -1 {
		creationSuccess = false
		err += "Matrix creation exceed height by " + strconv.Itoa(matrix.Height-maxHeight) + ".\n"
	}

	if matrix.Width > maxWidth && maxWidth != -1 {
		creationSuccess = false
		err += "Matrix creation exceed height by " + strconv.Itoa(matrix.Width-maxWidth) + ".\n"
	}

	if !creationSuccess {
		err = "Failed to create matrix.\n" + err
	}

	return

}
