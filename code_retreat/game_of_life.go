package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

func initBoard(rows int, cols int) [][]int {
	// Create capacity
	board := make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
	}

	// Randomize board
	for row := range board {
		for col := range board[row] {
			board[row][col] = rand.Intn(2)
		}
	}
	return board
}

func updateBoard(board [][]int) [][]int {
	// Create new board
	newBoard := make([][]int, len(board))
	for i := range newBoard {
		newBoard[i] = make([]int, len(board[i]))
	}

	for row := range board {
		for col := range board[row] {
			// Count of the amount of neighbors
			count := 0

			// Checking neighbors
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					// Skip self cell
					if i == 0 && j == 0 {
						continue
					}
					// Checking bounds and counting live neighbors
					if r, c := row+i, col+j; r >= 0 && r < len(board) && c >= 0 && c < len(board[row]) && board[r][c] == 1 {
						count++
					}
				}
			}

			// Game of Life rules
			if (board[row][col] == 1 && (count == 2 || count == 3)) || (board[row][col] == 0 && count == 3) {
				newBoard[row][col] = 1
			} else {
				newBoard[row][col] = 0
			}
		}
	}
	return newBoard
}

func printBoard(board [][]int) {
	// Clear previous output
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	var output strings.Builder

	// Print new output
	for row := range board {
		for col := range board[row] {
			if board[row][col] == 1 {
				output.WriteString(color.New(color.BgWhite).Sprint(" "))
			} else {
				output.WriteString(color.New(color.BgBlack).Sprint(" "))
			}
		}
		output.WriteString("\n")
	}

	fmt.Print(output.String())
}

func main() {
	rows := 40
	cols := 80
	board := initBoard(rows, cols)

	// Set generations
	for i := 0; i <= 200; i++ {
		updatedBoard := updateBoard(board)
		printBoard(updatedBoard)
		time.Sleep(100 * time.Millisecond)
		board = updatedBoard
	}
}
