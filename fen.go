package board

import (
	"fmt"
	"strconv"
	"strings"
)

// fenToBoard function to create a new board fro a FEN string
func fenToBoard(fen string) *Board {
	fen = trimSpaces(fen)
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return nil
	}
	b := parseBoardFEN(parts[0])
	turn := parseTurnFEN(parts[1])
	castlings := parseCastlingFEN(parts[2])
	enpassant := parseEnPFEN(parts[3])
	rule50 := parseRule50FEN(parts[4])
	moveNo := parseMoveNoFEN(parts[5])

	b.turn = turn
	b.castlings = castlings
	b.enpassant = enpassant
	b.rule50 = rule50
	b.moveNo = moveNo
	return b
}

func parseBoardFEN(part string) *Board {
	b := &Board{}
	b.Clear()
	subs := strings.Split(part, "/")
	var sub string
	var sq int
	if len(subs) != 8 {
		fmt.Println("incorrect fen segments")
	}
	for i := 7; i >= 0; i-- {
		sub = subs[i]
		sq = (7 - i) * 8
		for j := 0; j < len(sub); j++ {
			c := string(sub[j])
			switch c {
			case "K":
				b.set(WKing, sq)
				sq++
			case "Q":
				b.set(WQueen, sq)
				sq++
			case "R":
				b.set(WRook, sq)
				sq++
			case "B":
				b.set(WBishop, sq)
				sq++
			case "N":
				b.set(WKnight, sq)
				sq++
			case "P":
				b.set(WPawn, sq)
				sq++
			case "k":
				b.set(BKing, sq)
				sq++
			case "q":
				b.set(BQueen, sq)
				sq++
			case "r":
				b.set(BRook, sq)
				sq++
			case "b":
				b.set(BBishop, sq)
				sq++
			case "n":
				b.set(BKnight, sq)
				sq++
			case "p":
				b.set(BPawn, sq)
				sq++
			default:
				a, err := strconv.Atoi(c)
				if err != nil {
					fmt.Println("unexpected input in fen string")
					return nil
				}
				sq += a
			}
		}
	}
	return b
}
func parseTurnFEN(part string) color {
	switch lower(part) {
	case "w":
		return WHITE
	case "b":
		return BLACK
	default:
		fmt.Println("got an unexpected input for fen color: choosing white")
	}
	return WHITE
}
func parseCastlingFEN(part string) castlings {
	c := castlings(0)
	// fmt.Printf("from fen.go parsing castlings: %s\n", part)

	for _, str := range part {
		if str == 'K' {
			c.set(shortW)
		}
		if str == 'Q' {
			c.set(longW)
		}
		if str == 'k' {
			c.set(shortB)
		}
		if str == 'q' {
			c.set(longB)
		}
	}
	// fmt.Printf("from fen.go parsing castlings: %04b\n", c)
	return c
}
func parseEnPFEN(part string) int {
	if enp, ok := squareNameToInt[lower(part)]; ok {
		if enp/8 == 2 || enp/8 == 5 {

			return enp
		}
	}
	return 0
}
func parseRule50FEN(part string) int {
	a, err := strconv.Atoi(part)
	if err != nil {
		return 0
	}
	return a
}
func parseMoveNoFEN(part string) int {
	a, err := strconv.Atoi(part)
	if err != nil {
		return 0
	}
	return a
}

// ToFEN this returns the FEN of the current board position
func (b *Board) ToFEN() string {
	fen := ""
	counter := 0
	anon := func(a bool) {
		if a {
			if counter != 0 {
				fen = fen + fmt.Sprintf("%v/", counter)
				counter = 0
			} else {
				fen = fen + "/"
			}
			return
		}
		if counter != 0 {
			fen = fen + fmt.Sprintf("%v", counter)
			counter = 0
		}
	}
	for row := 7; row >= 0; row-- {
		c := b.square[row*8 : (row+1)*8]
		for col := 0; col <= 7; col++ {
			switch c[col] {
			case WKing:
				anon(false)
				fen = fen + "K"
			case WQueen:
				anon(false)
				fen = fen + "Q"
			case WRook:
				anon(false)
				fen = fen + "R"
			case WBishop:
				anon(false)
				fen = fen + "B"
			case WKnight:
				anon(false)
				fen = fen + "N"
			case WPawn:
				anon(false)
				fen = fen + "P"
			case BKing:
				anon(false)
				fen = fen + "k"
			case BQueen:
				anon(false)
				fen = fen + "q"
			case BRook:
				anon(false)
				fen = fen + "r"
			case BBishop:
				anon(false)
				fen = fen + "b"
			case BKnight:
				anon(false)
				fen = fen + "n"
			case BPawn:
				anon(false)
				fen = fen + "p"

			case Empty:
				// fen = fen + "K"
				counter++
			}
		}
		if row != 0 {
			anon(true)
		}
		if row == 0 && counter != 0 {
			fen += strconv.Itoa(counter)
		}
	}

	//add color to move
	fen += " "
	if b.turn == WHITE {
		fen += "w"
	} else {
		fen += "b"
	}
	// add castling prev
	fen += " "
	cast := b.castlings
	if cast != 0 {
		if (cast & shortW) != 0 {
			fen += "K"
		}
		if (cast & longW) != 0 {
			fen += "Q"
		}

		if (cast & shortB) != 0 {
			fen += "k"
		}

		if (cast & longB) != 0 {
			fen += "q"
		}

	} else {
		fen += "-"
	}

	fen += " "
	if b.enpassant/8 == 2 || b.enpassant/8 == 5 {
		fen += squareName[b.enpassant]
	} else {
		fen += "-"
	}

	// add 50 move rule
	fen += " "
	fif := strconv.Itoa(b.rule50)
	fen += fif
	// add move number
	fen += " "
	moveNum := strconv.Itoa(b.moveNo)
	fen += moveNum

	return fen
}
