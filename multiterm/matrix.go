package multiterm

import (
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

func makeMatrix(str string, maxWidth int, maxHeight int) *Matrix {

	matrix := newMatrix()
	matrix.Content = append(matrix.Content, make([]termbox.Cell, 0))

	for _, ch := range str {

		//if maxHeight -> break
		if len(matrix.Content) == maxHeight {
			break
		}

		//if last row doesn't exist -> make it
		if len(matrix.Content) == 0 {
			matrix.Content = append(matrix.Content, make([]termbox.Cell, 0))
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
		matrix.Content[rowInd] = append(matrix.Content[rowInd], termbox.Cell{Ch: ch})

		//If the current char is a newline,
		//append whitespaces to the rest of the row
		if ch == '\n' {
			numberOfWhiteSpaces := maxWidth - len(matrix.Content[rowInd])
			for numberOfWhiteSpaces > 0 {
				matrix.Content[rowInd] = append(matrix.Content[rowInd], termbox.Cell{Ch: ' '})
				numberOfWhiteSpaces--
			}
		}

	}

	return matrix

}
