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

// TODO : findNextMove, undoLastMove, RecursiveFunction "resolve" which returns true or false

func findNextMove() (bool, move) {
	// Use Konami order up down left right
	var initialPuzzleState = puzzleState
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			if puzzleState[line][column] == 1 {
				if line >= 2 &&
					puzzleState[line-1][column] == 1 &&
					puzzleState[line-2][column] == 2 &&
					!isMoveBanned(move{line, column, "↑", puzzleState}) {
					// up
					puzzleState[line][column] = 2
					puzzleState[line-1][column] = 2
					puzzleState[line-2][column] = 1

					return true, move{line, column, "↑", initialPuzzleState}
				} else if line <= 4 &&
					puzzleState[line+1][column] == 1 &&
					puzzleState[line+2][column] == 2 &&
					!isMoveBanned(move{line, column, "↓", puzzleState}) {
					// down
					puzzleState[line][column] = 2
					puzzleState[line+1][column] = 2
					puzzleState[line+2][column] = 1

					return true, move{line, column, "↓", initialPuzzleState}
				} else if column >= 2 &&
					puzzleState[line][column-1] == 1 &&
					puzzleState[line][column-2] == 2 &&
					!isMoveBanned(move{line, column, "←", puzzleState}) {
					// left
					puzzleState[line][column] = 2
					puzzleState[line][column-1] = 2
					puzzleState[line][column-2] = 1

					return true, move{line, column, "←", initialPuzzleState}
				} else if column <= 4 &&
					puzzleState[line][column+1] == 1 &&
					puzzleState[line][column+2] == 2 &&
					!isMoveBanned(move{line, column, "→", puzzleState}) {
					// right
					puzzleState[line][column] = 2
					puzzleState[line][column+1] = 2
					puzzleState[line][column+2] = 1

					return true, move{line, column, "→", initialPuzzleState}
				}
			}
		}
	}
	return false, move{}
}

func undo() {
	bannedMoves = append(bannedMoves, moves[len(moves)-1])
	puzzleState = moves[len(moves)-1].puzzleState
	moves = moves[:len(moves)-1]

	pegsOnBoard++
}

func isMoveBanned(newMove move) bool {
	// iterate in reverse to save time
	for i := len(bannedMoves)-1; i >= 0; i-- {
		if bannedMoves[i] == newMove {
			return true
		}
	}
	return false
}

var iteration = 0

func resolve() bool {
	foundNextMove, newMove := findNextMove()
	if foundNextMove {
		iteration++
		moves = append(moves, newMove)
		fmt.Printf("\rmoving %d:%d:%s, pegs left : %02d, banned moves : %06d, moves : %06d, timer : %s",
			newMove.column, newMove.line, newMove.direction, pegsOnBoard, len(bannedMoves), len(moves), time.Now().Sub(startTime))

		pegsOnBoard--
		if pegsOnBoard == 1 {
			return true
		}

		return resolve()
	} else {
		undo()
		return resolve()
	}
	return false
}

func printPuzzle() {
	for line := 0; line < 7; line++ {
		for column := 0; column < 7; column++ {
			fmt.Print(puzzleState[line][column], " ")
		}
		fmt.Println()
	}
}

func main() {
	file2lines("./english.test.puzzle")

	startTime := time.Now()
	fmt.Println("Puzzle initial state")
	printPuzzle()

	if resolve() {
		fmt.Println("\nWon!")


	} else {
		fmt.Println("\nFound no solution")
	}

	fmt.Println("Stastistics")
	fmt.Println("_____________________")
	fmt.Println(iteration, "iterations")
	fmt.Println("Banned", len(bannedMoves), "moves")
	fmt.Println("Solution uses", len(moves), "moves")
	fmt.Println("Done in ", time.Now().Sub(startTime))

	fmt.Println("Puzzle final state")
	printPuzzle()
}
