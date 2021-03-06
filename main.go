package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"github.com/cnf/structhash"
	"runtime/debug"
)

var puzzleState [7][7]int
var moves []move
var bannedPuzzleState = make(map[string] byte)
var pegsOnBoard = 0
var startTime = time.Now()

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

func findNextMove() (bool, move) {
	// Use Konami order up down left right
	var initialPuzzleState = puzzleState
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			if puzzleState[line][column] == 1 {
				if line >= 2 &&
					puzzleState[line-1][column] == 1 &&
					puzzleState[line-2][column] == 2 {

					puzzleStatePrototype := puzzleState
					puzzleStatePrototype[line][column] = 2
					puzzleStatePrototype[line-1][column] = 2
					puzzleStatePrototype[line-2][column] = 1

					if !isBannedBoard(puzzleStatePrototype) {
						// up
						puzzleState[line][column] = 2
						puzzleState[line-1][column] = 2
						puzzleState[line-2][column] = 1

						return true, move{line, column, "↑", initialPuzzleState}
					}
				}

				if line <= 4 &&
					puzzleState[line+1][column] == 1 &&
					puzzleState[line+2][column] == 2 {

					puzzleStatePrototype := puzzleState
					puzzleStatePrototype[line][column] = 2
					puzzleStatePrototype[line+1][column] = 2
					puzzleStatePrototype[line+2][column] = 1

					if !isBannedBoard(puzzleStatePrototype) {
						// down
						puzzleState[line][column] = 2
						puzzleState[line+1][column] = 2
						puzzleState[line+2][column] = 1

						return true, move{line, column, "↓", initialPuzzleState}
					}
				}

				if column >= 2 &&
					puzzleState[line][column-1] == 1 &&
					puzzleState[line][column-2] == 2 {

					puzzleStatePrototype := puzzleState
					puzzleStatePrototype[line][column] = 2
					puzzleStatePrototype[line][column-1] = 2
					puzzleStatePrototype[line][column-2] = 1

					if !isBannedBoard(puzzleStatePrototype) {
						// left
						puzzleState[line][column] = 2
						puzzleState[line][column-1] = 2
						puzzleState[line][column-2] = 1

						return true, move{line, column, "←", initialPuzzleState}
					}
				}

				if column <= 4 &&
					puzzleState[line][column+1] == 1 &&
					puzzleState[line][column+2] == 2 {

					puzzleStatePrototype := puzzleState
					puzzleStatePrototype[line][column] = 2
					puzzleStatePrototype[line][column+1] = 2
					puzzleStatePrototype[line][column+2] = 1

					if !isBannedBoard(puzzleStatePrototype) {
						// right
						puzzleState[line][column] = 2
						puzzleState[line][column+1] = 2
						puzzleState[line][column+2] = 1

						return true, move{line, column, "→", initialPuzzleState}
					}
				}
			}
		}
	}
	return false, move{}
}

func undo() bool {

	if len(moves) == 0 {
		return false
	}

	bannedPuzzleState[string(structhash.Md5(puzzleState, 1))] = 1

	puzzleState = moves[len(moves)-1].puzzleState
	moves = moves[:len(moves)-1]

	pegsOnBoard++

	return true
}

func isBannedBoard(puzzleStatePrototype [7][7]int) bool {
	if _, ok := bannedPuzzleState[string(structhash.Md5(puzzleStatePrototype, 1))]; ok {
		return true
	}
	return false
}

var iteration = 0

func resolve() bool {
	foundNextMove, newMove := findNextMove()
	if foundNextMove {
		iteration++
		pegsOnBoard--
		moves = append(moves, newMove)

		fmt.Printf("\rmoving %d:%d:%s, pegs left : %02d, banned moves : %06d, moves : %06d, timer : %s",
			newMove.column, newMove.line, newMove.direction, pegsOnBoard, len(bannedPuzzleState), len(moves), time.Now().Sub(startTime))

		if pegsOnBoard == 1 {
			return true
		}

		return resolve()
	} else {
		if !undo() {
			return false
		}

		return resolve()
	}
	return false
}

func printPuzzle(puzzle [7][7]int) {
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			fmt.Print(puzzle[line][column], " ")
		}
		fmt.Println()
	}
}

func main() {

	debug.SetMaxStack(14000000000)
	arg := os.Args[1]

	fmt.Println("Opening file", arg)

	file2lines(arg)

	startTime := time.Now()
	fmt.Println("Puzzle initial state")
	printPuzzle(puzzleState)

	win := false
	if resolve() {
		win = true
	}


	fmt.Println("\n_____________________")
	fmt.Println("Moves")
	for n, r := range moves {
		fmt.Printf("\n%d. %d:%d:%s\n", n + 1, r.column +1 , r.line + 1, r.direction)
		//printPuzzle(r.puzzleState)
	}
	fmt.Println("_____________________")

	if win {
		fmt.Println("\nComputer wins! :)")
	} else {
		fmt.Println("\nHumans win! :(")
	}
	fmt.Println("_____________________")
	fmt.Println("Stastistics")
	fmt.Println(iteration, "iterations")
	fmt.Println("Banned", len(bannedPuzzleState), "boards")
	fmt.Println("Solution uses", len(moves), "moves")
	fmt.Println("Done in ", time.Now().Sub(startTime))

	fmt.Println("Puzzle final state")
	printPuzzle(puzzleState)
}
