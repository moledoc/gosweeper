package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var closedTile = "#"

var dirs []int = []int{
	-1, 0, // n
	1, 0, // s
	0, -1, // w
	0, 1, // e
	-1, -1, // nw
	1, -1, // sw
	-1, 1, // ne
	1, 1, // se
}

var nRows int
var nCols int

// MakeGrid TODO:
// func MakeGrid[V int | string](r, c int, t V) [][]V {
func MakeGrid(r, c int) [][]int {
	grid := make([][]int, r)
	for i := 0; i < r; i++ {
		tmpRow := make([]int, c)
		grid[i] = tmpRow
	}
	return grid
}

// MakeBoard TODO:
func MakeBoard(r, c int) [][]string {
	board := make([][]string, r)
	for i := 0; i < r; i++ {
		tmpRow := make([]string, c)
		for i := range tmpRow {
			tmpRow[i] = closedTile
		}
		board[i] = tmpRow
	}
	return board
}

// func MakeGrid(r, c int, _ int) [][]int {
// grid := make([][]int, r)
// for i := 0; i < r; i++ {
// tmpRow := make([]int, c)
// grid[i] = tmpRow
// }
// return grid
// }

// MakeMines TODO:
func MakeMines(r, c, m int) []int {
	perm := rand.Perm(r * c)[:m]
	sort.Ints(perm)
	return perm
}

// AddMines TODO:
func AddMines(grid *[][]int, mines *[]int) {
	var i int
	var j int
	for r := 0; r < len((*grid)) && j < len((*mines)); r++ {
		for c := 0; c < len((*grid)[0]) && j < len((*mines)); c++ {
			if (*mines)[j] == i {
				(*grid)[r][c]--
				j++
			}
			i++
		}
	}
}

// AddHints TODO:
func AddHints(grid *[][]int) {
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			for i := 0; (*grid)[r][c] == -1 && i < len(dirs); i += 2 {
				tmpR := r + dirs[i]
				tmpC := c + dirs[i+1]
				if tmpR < 0 || tmpR >= nRows ||
					tmpC < 0 || tmpC >= nCols ||
					(*grid)[tmpR][tmpC] == -1 {
					continue
				}
				(*grid)[tmpR][tmpC]++
			}
		}
	}
}

// printer
func printer[V int | string](grid *[][]V) {
	for r := 0; r < len((*grid)); r++ {
		for c := 0; c < len((*grid)[0]); c++ {
			fmt.Printf("%3v", (*grid)[r][c])
		}
		fmt.Println()
	}
}

// Round
func Round(grid *[][]int, board *[][]string, x int, y int, action string) error { // ([][]int,[][]string){
	switch action {
	case "o":
		// TODO:
		if (*grid)[x][y] == -1 {
			(*board)[x][y] = fmt.Sprintf("%v", (*grid)[x][y])
			return errors.New("Hit a mine")
		}
		Open(grid, board, x, y)
	case "x":
		if (*board)[x][y] == closedTile {
			(*board)[x][y] = "X"
		}
	case "c":
		if (*board)[x][y] == "X" {
			(*board)[x][y] = closedTile
		}
	default:
		break
	}
	// return grid,board
	return nil
}

func Open(grid *[][]int, board *[][]string, x int, y int) {
	(*board)[x][y] = fmt.Sprintf("%v", (*grid)[x][y])
	if (*grid)[x][y] > 0 {
		return
	}
	for i := 0; i < len(dirs); i += 2 {
		tmpR := x + dirs[i]
		tmpC := y + dirs[i+1]
		if tmpR < 0 || tmpR >= nRows ||
			tmpC < 0 || tmpC >= nCols ||
			(*grid)[tmpR][tmpC] == -1 ||
			(*board)[tmpR][tmpC] != closedTile {
			continue
		}
		Open(grid, board, tmpR, tmpC)
	}
}

func OnlyMines(grid *[][]int, board *[][]string) bool {
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			if ((*board)[r][c] == closedTile || (*board)[r][c] == "X") && (*grid)[r][c] > -1 {
				return false
			}
		}
	}
	return true
}

func main() {
	cols := flag.Int("c", 3, "The number of columns in the grid")
	rows := flag.Int("r", 3, "The number of rows in the grid")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	nRows = *rows
	nCols = *cols
	// rand.Seed(time.Now().Unix())
	mines := MakeMines(*rows, *cols, 3)
	grid := MakeGrid(*rows, *cols)
	AddMines(&grid, &mines)
	printer(&grid)
	fmt.Println("-----------------------------------------")
	AddHints(&grid)
	printer(&grid)
	fmt.Println("-----------------------------------------")
	board := MakeBoard(*rows, *cols)
	printer(&board)
	fmt.Println("-----------------------------------------")
	start := time.Now()
	var input string
	for {
		fmt.Println("Enter (row)(col)(action)")
		fmt.Scanln(&input)
		x, err := strconv.Atoi(string(input[0]))
		if err != nil {
			fmt.Println("Incorrect row value")
			err = nil
			continue
		}
		y, err := strconv.Atoi(string(input[1]))
		if err != nil {
			fmt.Println("Incorrect col value")
			err = nil
			continue
		}
		err = Round(&grid, &board, x, y, string(input[2]))
		printer(&board)
		fmt.Println("-----------------------------------------")
		if err != nil {
			log.Fatalln(err)
		}
		if OnlyMines(&grid, &board) {
			fmt.Println("Congratulations! All mines found")
			fmt.Printf("Time: %v\n", time.Since(start))
			return
		}
	}

}
