package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const gridSize = 9

type SudokuSolver struct {
	board /*2d arr*/ [gridSize][gridSize]int
}

func (x *SudokuSolver) isValid(row, col, num int) bool { 
	for i := 0; i < gridSize; i++ {
		if x.board[row][i] == num || x.board[i][col] == num {
			return false
		}
	}
	startRow, startCol := row/3*3, col/3*3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if x.board[startRow+i][startCol+j] == num {
				return false
			}
		}
	}
	return true
}

func (x *SudokuSolver) solve() bool { 
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if x.board[row][col] == 0 {
				for num := 1; num <= gridSize; num++ {
					if x.isValid(row, col, num) {
						x.board[row][col] = num
						if x.solve() { 
							return true
						}
						x.board[row][col] = 0 
					}
				}
				return false
			}
		}
	}
	return true
}

func main() {
	a := app.New()
	w := a.NewWindow("Sudoku Solver")
	w.Resize(fyne.NewSize(450, 500))

	var entries [gridSize][gridSize]*widget.Entry
	grid := container.NewGridWithColumns(gridSize)
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			entry := widget.NewEntry()
			entry.SetPlaceHolder("")
			entries[i][j] = entry

			entryBox := container.NewPadded(entry)

			if j%3 == 2 && j < gridSize-1 { 
				entryBox = container.NewBorder(nil, nil, nil, canvas.NewRectangle(color.NRGBA{100, 0, 0, 255}), entryBox)
			}
			if i%3 == 2 && i < gridSize-1 {
				entryBox = container.NewBorder(nil, canvas.NewRectangle(color.NRGBA{100, 0, 0, 255}), nil, nil, entryBox)
			}
			grid.Add(entryBox)
		}
	}

	solveButton := widget.NewButton("Solve", func() {
		var board [gridSize][gridSize]int
		for i := 0; i < gridSize; i++ {
			for j := 0; j < gridSize; j++ {
				value, err := strconv.Atoi(entries[i][j].Text)
				if err != nil || value < 0 || value > 9 {
					entries[i][j].SetText("0")
					value = 0
				}
				board[i][j] = value
			}
		}

		x := SudokuSolver{board: board}
		if x.solve() {
			for i := 0; i < gridSize; i++ {
				for j := 0; j < gridSize; j++ {
					entries[i][j].SetText(strconv.Itoa(x.board[i][j]))
				}
			}
		} else {
			fmt.Println("No solution found.")
		}
	})

	clearButton := widget.NewButton("Clear", func() {
		for i := 0; i < gridSize; i++ {
			for j := 0; j < gridSize; j++ {
				entries[i][j].SetText("")
			}
		}
	})

	buttons := container.NewHBox(solveButton, clearButton)
	content := container.NewVBox(grid, layout.NewSpacer(), buttons)

	w.SetContent(content)
	w.ShowAndRun()
}
