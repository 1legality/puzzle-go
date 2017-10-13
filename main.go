package main

import (
	"bufio"
	"fmt"
	"os"
)

var puzzleState [7][7]int;

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
			// fmt.Printf("%d: %d %c\n", line, n, rune)
			puzzleState[line][n] = int(r - '0') //convert rune to int
		}
		line++;
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

// TODO : findNextMove, undoLastMove, RecursiveFunction "resolve" which returns true or false

func main() {
	file2lines("./test.puzzle")

	fmt.Print(puzzleState)
}
