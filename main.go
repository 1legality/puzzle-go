package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var puzzleState [7][7]int
var moves []move
var bannedMoves []move
var pegsOnBoard = 0

type move struct {
	line        int
	column      int
	direction   string
	puzzleState [7][7]int
}

func file2lines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	var line int = 0;
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())

		for n, r := range scanner.Text() {

			if int(r - '0') == 1 {
				pegsOnBoard++
			}

			puzzleState[line][n] = int(r - '0') //convert rune to int
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

// TODO : findNextMove, undoLastMove, RecursiveFunction "resolve" which returns true or false

func findNextMove(line int, column int) (bool, string) {
	// Use Konami order up down left right
	if line >= 2 &&
		puzzleState[line-1][column] == 1 &&
		puzzleState[line-2][column] == 2 &&
		!isMoveBanned(move{line, column, "↑", puzzleState}) {
		// up
		moves = append(moves, move{line, column, "↑", puzzleState})

		puzzleState[line][column] = 2
		puzzleState[line-1][column] = 2
		puzzleState[line-2][column] = 1

		return true, "↑"
	} else if line <= 4 &&
		puzzleState[line+1][column] == 1 &&
		puzzleState[line+2][column] == 2 &&
		!isMoveBanned(move{line, column, "↓", puzzleState}) {
		// down

		moves = append(moves, move{line, column, "↓", puzzleState})

		puzzleState[line][column] = 2
		puzzleState[line+1][column] = 2
		puzzleState[line+2][column] = 1

		return true, "↓"
	} else if column >= 2 &&
		puzzleState[line][column-1] == 1 &&
		puzzleState[line][column-2] == 2 &&
		!isMoveBanned(move{line, column, "←", puzzleState}) {
		// left

		moves = append(moves, move{line, column, "←", puzzleState})
		puzzleState[line][column] = 2
		puzzleState[line][column-1] = 2
		puzzleState[line][column-2] = 1

		return true, "←"
	} else if column <= 4 &&
		puzzleState[line][column+1] == 1 &&
		puzzleState[line][column+2] == 2 &&
		!isMoveBanned(move{line, column, "→", puzzleState}) {
		// right
		moves = append(moves, move{line, column, "→", puzzleState})

		puzzleState[line][column] = 2
		puzzleState[line][column+1] = 2
		puzzleState[line][column+2] = 1

		return true, "→"
	}
	return false, "x"
}

func isMoveBanned(newMove move) bool {
	for _, v := range bannedMoves {
		if v == newMove {
			return true
		}
	}
	return false
}

var iteration = 0

func resolve() bool {
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			if puzzleState[line][column] == 1 {
				foundNextMove, direction := findNextMove(line, column)
				if foundNextMove {
					iteration++
					//fmt.Println(iteration)
					fmt.Println("moving ", line+1, ":", column+1, ":", direction, ", pegs left : ", pegsOnBoard)
					//printPuzzle()

					pegsOnBoard--
					if pegsOnBoard == 1 {
						return true
					}

					return resolve()
				}
			}

			if len(moves) > 1 && line == 6 && column == 6 {
				bannedMoves = append(bannedMoves, moves[len(moves)-1])
				puzzleState = moves[len(moves)-1].puzzleState
				moves = moves[:len(moves)-1]

				pegsOnBoard++

				return resolve()
			}
		}
	}
	return false
}

func printPuzzle() {
	fmt.Println("Printing puzzle")
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			fmt.Print(puzzleState[line][column], " ")
		}
		fmt.Println()
	}
}

func verifyIfWin() bool {
	counter := 0
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			if puzzleState[line][column] == 1 {
				counter++
			}
		}
	}

	if counter == 1 {
		return true
	} else {
		return false
	}
}

func main() {
	file2lines("./english.test.puzzle")

	startTime := time.Now()

	if resolve() {
		fmt.Println("Won!")
	} else {
		fmt.Println("Found no solution")
	}
	// printPuzzle()
	fmt.Println("Done in ", time.Now().Sub(startTime))
}
