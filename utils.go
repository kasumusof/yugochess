package board

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var trimSpaces = func(text string) string {
	text = strings.TrimSpace(text)
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")
	return text
}
var upper = func(text string) string {
	text = strings.ToUpper(text)
	return text
}
var lower = func(text string) string {
	text = strings.ToLower(text)
	return text
}

var prepMoves = func(text string) string {
	text = strings.TrimSpace(text)
	space := regexp.MustCompile(`\d+\.`)
	text = space.ReplaceAllString(text, " ")
	equals := regexp.MustCompile(`=`)
	text = equals.ReplaceAllString(text, "")
	text = trimSpaces(text)
	return text
}

// function for parsing notations of chess
var startsWith = func(text string, char byte) bool {
	if len(text) == 0 {
		return false
	}
	return text[0] == char
}
var startsWithAny = func(text string, chars string) bool {
	if len(text) == 0 {
		return false
	}
	for _, char := range chars {
		if text[0] == byte(char) {
			return true
		}
	}
	return false
}
var endsWithAny = func(text string, chars string) bool {
	if len(text) == 0 {
		return false
	}
	for _, char := range chars {
		if text[len(text)-1] == byte(char) {
			return true
		}
	}
	return false
}
var endsWith = func(text string, char byte) bool {
	if len(text) == 0 {
		return false
	}
	return text[len(text)-1] == char
}
var contains = func(text string, char byte) bool {
	return strings.Contains(text, string(char))
}

func drawBitBoard(bb bitBoard) {
	// start := time.Now()

	bitString := fmt.Sprintf("%064b", bb)
	var c, d string
	// fmt.Println("______Drawing_BitBoard_Started______")

	fmt.Println("    A   B   C   D   E   F   G   H")

	fmt.Println("  +---|---|---|---|---|---|---|---+")
	for i := 0; i < 8; i++ {
		// c = string(bitString[i*8 : (i+1)*8])
		c = string(bitString[i*8 : (i+1)*8])
		for j := 7; j >= 0; j-- {
			d = string(c[j])
			if j == 0 {
				fmt.Printf("%2s |%2v", d, 8-i)
				continue
			}
			if j == 7 {
				fmt.Printf("%1v |%2s |", 8-i, d)
				continue
			}
			fmt.Printf("%2s |", d)
		}
		fmt.Println()
		fmt.Println("  +---|---|---|---|---|---|---|---+")

	}
	fmt.Println("    A   B   C   D   E   F   G   H")
	// fmt.Println("_______Drawing_BitBoard_Ended_______\n", time.Since(start))

}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// DrawBoard function to draw a board to the  terminal
func DrawBoard(b *Board) {
	// start := time.Now()
	var c []int
	var d string
	// fmt.Println("_______Drawing_Board_Started_______")

	fmt.Println("    A   B   C   D   E   F   G   H")

	fmt.Println("  +---|---|---|---|---|---|---|---+")
	for i := 7; i >= 0; i-- {
		c = b.square[i*8 : (i+1)*8]
		for j := 0; j < 8; j++ {
			d = string(display[c[j]])
			// d = string(c[i])
			// d = "p"

			if j == 7 {
				fmt.Printf("%2s |%2v", d, i+1)
				continue
			}
			if j == 0 {
				fmt.Printf("%1v |%2s |", i+1, d)
				continue
			}
			fmt.Printf("%2s |", d)
		}
		fmt.Println()
		fmt.Println("  +---|---|---|---|---|---|---|---+")

	}
	fmt.Println("    A   B   C   D   E   F   G   H")
	// fmt.Println("_______Drawing_Board_Ended_______\n", time.Since(start))

}

// InputChan this function returns a string channel and creates an anonymous function that listens for input on that channel
func InputChan() chan string {
	input := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != io.EOF && len(text) > 0 {
				text = trimSpaces(text)
				input <- text
			}
		}

	}()
	return input
}
