package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// MakeGrid TODO:
func MakeGrid[V int | string](r, c int, t V) [][]V {
	grid := make([][]V, r)
	for i := 0; i < r; i++ {
		tmpRow := make([]V, c)
		grid[i] = tmpRow
	}
	// if reflect.TypeOf(t) == "string" {
	// nRows := len(grid)
	// nCols := len(grid[0])
	// for r := 0; r < nRows; r++ {
	// for c := 0; c < nCols; c++ {
	// grid[r][c] = "#"
	// }
	// }
	// }
	return grid
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
func AddMines(grid [][]int, mines []int) [][]int {
	var i int
	var j int
	for r := 0; r < len(grid) && j < len(mines); r++ {
		for c := 0; c < len(grid[0]) && j < len(mines); c++ {
			if mines[j] == i {
				grid[r][c]--
				j++
			}
			i++
		}
	}
	return grid
}

// AddHints TODO:
func AddHints(grid [][]int) [][]int {
	dirs := []int{
		-1, 0, // n
		1, 0, // s
		0, -1, // w
		0, 1, // e
		-1, -1, // nw
		1, -1, // sw
		-1, 1, // ne
		1, 1, // se
	}
	nRows := len(grid)
	nCols := len(grid[0])
	for r := 0; r < nRows; r++ {
		for c := 0; c < nCols; c++ {
			for i := 0; grid[r][c] == -1 && i < len(dirs); i += 2 {
				tmpR := r + dirs[i]
				tmpC := c + dirs[i+1]
				if tmpR < 0 || tmpR >= nRows ||
					tmpC < 0 || tmpC >= nCols ||
					grid[tmpR][tmpC] == -1 {
					continue
				}
				grid[tmpR][tmpC]++
			}
		}
	}
	return grid
}

// printer
func printer[V int | string](grid [][]V) {
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			fmt.Printf("%3v", grid[r][c])
		}
		fmt.Println()
	}
}

// Round
func Round(grid *[][]int, board *[][]string, x int, y int, action string) { // ([][]int,[][]string){
	switch action {
	case "o":
		// TODO:
		break
	case "f":
		(*board)[x][y] = "F"
	case "x":
		(*board)[x][y] = "X"
	case "c":
		(*board)[x][y] = ""
	default:
		break
	}
	// return grid,board
}

func main() {
	cols := flag.Int("c", 10, "The number of columns in the grid")
	rows := flag.Int("r", 10, "The number of rows in the grid")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	rand.Seed(time.Now().Unix())
	mines := MakeMines(*rows, *cols, 10)
	fmt.Println(mines)
	grid := MakeGrid(*rows, *cols, 0)
	grid = AddMines(grid, mines)
	printer(grid)
	fmt.Println("-----------------------------------------")
	grid = AddHints(grid)
	printer(grid)
	fmt.Println("-----------------------------------------")
	board := MakeGrid(*rows, *cols, "")
	printer(board)
	fmt.Println("-----------------------------------------")
	Round(&grid, &board, 4, 4, "f")
	printer(board)
	fmt.Println("-----------------------------------------")
}
