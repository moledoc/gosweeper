package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type tile struct {
	ind, row, col int
	hint          int // -1 indicates mine
	hintStr       string
	isOpened      bool
	isFlagged     bool
}

type Tiles []tile

type Game struct {
	tiles *Tiles
	start time.Time
}

var closedTile = "#"
var flaggedTile = "X"

var nRows int = 3
var nCols int = 3
var nMines int = 3

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

func makeMines() []int {
	perm := rand.Perm(nRows * nCols)[:nMines]
	sort.Ints(perm)
	return perm
}

func makeTiles() *Tiles {
	tiles := make(Tiles, nRows*nCols)
	mines := makeMines()
	var j int
	for i := range tiles {
		tiles[i] = tile{
			ind:     i,
			row:     (i / nRows),
			col:     i % nCols,
			hintStr: "#",
		}
		if j < nMines && i == mines[j] {
			tiles[i].hint--
			j++
		}
	}
	return (&tiles)
}

func (tiles *Tiles) addHints() *Tiles {
	for i, elem := range *tiles {
		for d := 0; elem.hint == -1 && d < len(dirs); d += 2 {
			if (i/nRows == 0 && dirs[d] == -1) || // top row
				(i/nRows == nRows-1 && dirs[d] == 1) || // bottom row
				(i%nCols == 0 && dirs[d+1] == -1) || // left column
				(i%nCols == nCols-1 && dirs[d+1] == 1) { // right column
				continue
			}
			if x := i + nCols*dirs[d] + dirs[d+1]; (*tiles)[x].hint != -1 {
				(*tiles)[x].hint++
			}
		}
	}
	return tiles
}

func (tiles *Tiles) printer(debug bool) {
	fmt.Print("r\\c|")
	for i := 0; i < nCols; i++ {
		fmt.Printf("%3v", i)
	}
	fmt.Print("\n---+")
	for i := 0; i < nCols; i++ {
		fmt.Printf("---")
	}
	fmt.Println()
	for i, elem := range *tiles {
		if i%nCols == 0 && i != 0 {
			fmt.Println()
		}
		if i%nCols == 0 {
			fmt.Printf("%3v|", i/nRows)
		}
		if debug {
			fmt.Printf("%3v", elem.hint)
			continue
		}
		fmt.Printf("%3v", elem.hintStr)
	}
	fmt.Println()
}

func (tiles *Tiles) open(loc int) error {
	if (*tiles)[loc].isOpened {
		return nil
	}
	(*tiles)[loc].hintStr = fmt.Sprintf("%v", (*tiles)[loc].hint)
	(*tiles)[loc].isOpened = true
	if (*tiles)[loc].hint == -1 {
		return errors.New("Game Over! You hit a mine.")
	}
	if (*tiles)[loc].hint > 0 {
		return nil
	}
	for d := 0; d < len(dirs); d += 2 {
		if (loc/nRows == 0 && dirs[d] == -1) || // top row
			(loc/nRows == nRows-1 && dirs[d] == 1) || // bottom row
			(loc%nCols == 0 && dirs[d+1] == -1) || // left column
			(loc%nCols == nCols-1 && dirs[d+1] == 1) { // right column
			continue
		}
		if next := loc + nCols*dirs[d] + dirs[d+1]; (*tiles)[next].hint != -1 {
			tiles.open(next)
		}
	}
	return nil
}

func (tiles *Tiles) onlyMines() bool {
	for _, elem := range *tiles {
		if (!elem.isOpened || elem.isFlagged) && elem.hint != -1 {
			return false
		}
	}
	return true
}

func (_ *Game) userInput() (int, string, error) {
	var input string
	fmt.Println("Enter <row><col><action> [o: open; x: flag; c: clear flag]:")
	fmt.Scanln(&input)
	x, err := strconv.Atoi(string(input[0]))
	if err != nil || x >= nRows {
		return 0, "", errors.New("Incorrect row value")
	}
	y, err := strconv.Atoi(string(input[1]))
	if err != nil || y >= nCols {
		return 0, "", errors.New("Incorrect col value")
	}
	return x*nRows + y, string(input[2]), nil
}

func (game *Game) play() {
	var err error

	game.tiles.printer(true)
	fmt.Println("........................................\n")
	game.tiles.printer(false)
	fmt.Println("========================================\n")

	for {
		if err != nil {
			log.Fatalln(err)
		}
		loc, action, userErr := game.userInput()
		if userErr != nil {
			fmt.Println(userErr)
			continue
		}
		switch action {
		case "o":
			err = game.tiles.open(loc)
		case "x":
			if !(*game.tiles)[loc].isFlagged && !(*game.tiles)[loc].isOpened {
				(*game.tiles)[loc].hintStr = flaggedTile
				(*game.tiles)[loc].isFlagged = true
			}
		case "c":
			if (*game.tiles)[loc].isFlagged {
				(*game.tiles)[loc].hintStr = closedTile
				(*game.tiles)[loc].isFlagged = false
			}
		default:
			fmt.Println("Unknown action")
			continue
		}
		game.tiles.printer(false)
		fmt.Println("========================================\n")
		if (*game.tiles).onlyMines() {
			fmt.Println("Congratulations! You found all the mines")
			fmt.Printf("Time: %v\n", time.Since((*game).start))
			return
		}
	}
}

func main() {
	// rand.Seed(time.Now().Unix())
	game := Game{
		tiles: makeTiles().addHints(),
		start: time.Now(),
	}
	game.play()
}
