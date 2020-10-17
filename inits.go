package board

var display = make(map[int]int)
var knightAttacks = make(map[int]bitBoard)
var kingAttacks = make(map[int]bitBoard)
var squareName = make(map[int]string)
var squareNameToInt = make(map[string]int)
var pieceITOA = make(map[int]string)
var colors = make(map[color]string)

func initDisplay() {
	const (
		whiteKing   = '♔'
		whiteQueen  = '♕'
		whiteRook   = '♖'
		whiteBishop = '♗'
		whiteKnight = '♘'
		whitePawn   = '♙'
		blackKing   = '♚'
		blackQueen  = '♛'
		blackRook   = '♜'
		blackBishop = '♝'
		blackKnight = '♞'
		blackPawn   = '♟'
	)
	display[WPawn] = whitePawn
	display[WKnight] = whiteKnight
	display[WBishop] = whiteBishop
	display[WRook] = whiteRook
	display[WQueen] = whiteQueen
	display[WKing] = whiteKing

	display[Empty] = ' '

	display[BPawn] = blackPawn
	display[BKnight] = blackKnight
	display[BBishop] = blackBishop
	display[BRook] = blackRook
	display[BQueen] = blackQueen
	display[BKing] = blackKing

}

func initKingAttacks() {
	mv := bitBoard(0)
	a := 0
	for i := A1; i <= H8; i++ {
		rank, file := i/8, i%8
		a = i + N
		if a >= A1 && a <= H8 && rank <= 6 {
			mv = mv | (1 << a)
		}
		a = i + NE
		if a >= A1 && a <= H8 && rank <= 6 && file <= 6 {
			mv = mv | (1 << a)
		}
		a = i + NW
		if a >= A1 && a <= H8 && rank <= 6 && file >= 1 {
			mv = mv | (1 << a)
		}
		a = i + E
		if a >= A1 && a <= H8 && file <= 6 {
			mv = mv | (1 << a)
		}
		a = i + W
		if a >= A1 && a <= H8 && file >= 1 {
			mv = mv | (1 << a)
		}
		a = i + S
		if a >= A1 && a <= H8 && rank >= 1 {
			mv = mv | (1 << a)
		}
		a = i + SE
		if a >= A1 && a <= H8 && rank >= 1 && file <= 6 {
			mv = mv | (1 << a)
		}
		a = i + SW
		if a >= A1 && a <= H8 && rank >= 1 && file >= 1 {
			mv = mv | (1 << a)
		}

		kingAttacks[i] = bitBoard(mv)
		// fmt.Printf("--------%s----------\n", squareName[i])
		// drawBitboard(mv)
		mv = 0
	}

}

func initKnightAttacks() {
	mv := bitBoard(0)
	a := 0
	for i := A1; i <= H8; i++ {
		rank, file := i/8, i%8
		a = i + NEN
		if a >= A1 && a <= H8 && rank <= 5 && file <= 6 {
			mv = mv | (1 << a)
		}
		a = i + NEE
		if a >= A1 && a <= H8 && rank <= 6 && file <= 5 {
			mv = mv | (1 << a)
		}
		a = i + NWN
		if a >= A1 && a <= H8 && rank <= 5 && file >= 1 {
			mv = mv | (1 << a)
		}
		a = i + NWW
		if a >= A1 && a <= H8 && rank <= 6 && file >= 2 {
			mv = mv | (1 << a)
		}
		a = i + SES
		if a >= A1 && a <= H8 && rank >= 2 && file <= 6 {
			mv = mv | (1 << a)
		}
		a = i + SEE
		if a >= A1 && a <= H8 && rank >= 1 && file <= 5 {
			mv = mv | (1 << a)
		}
		a = i + SWS
		if a >= A1 && a <= H8 && rank >= 2 && file >= 1 {
			mv = mv | (1 << a)
		}
		a = i + SWW
		if a >= A1 && a <= H8 && rank >= 1 && file >= 2 {
			mv = mv | (1 << a)
		}

		knightAttacks[i] = bitBoard(mv)
		// fmt.Printf("--------%s----------\n", squareName[i])
		// drawBitboard(mv)
		mv = 0
	}
}

func initSquareName() {

	colors[WHITE] = "white"
	colors[BLACK] = "black"

	squareName[A1] = "a1"
	squareName[B1] = "b1"
	squareName[C1] = "c1"
	squareName[D1] = "d1"
	squareName[E1] = "e1"
	squareName[F1] = "f1"
	squareName[G1] = "g1"
	squareName[H1] = "h1"

	squareName[A2] = "a2"
	squareName[B2] = "b2"
	squareName[C2] = "c2"
	squareName[D2] = "d2"
	squareName[E2] = "e2"
	squareName[F2] = "f2"
	squareName[G2] = "g2"
	squareName[H2] = "h2"

	squareName[A3] = "a3"
	squareName[B3] = "b3"
	squareName[C3] = "c3"
	squareName[D3] = "d3"
	squareName[E3] = "e3"
	squareName[F3] = "f3"
	squareName[G3] = "g3"
	squareName[H3] = "h3"

	squareName[A4] = "a4"
	squareName[B4] = "b4"
	squareName[C4] = "c4"
	squareName[D4] = "d4"
	squareName[E4] = "e4"
	squareName[F4] = "f4"
	squareName[G4] = "g4"
	squareName[H4] = "h4"

	squareName[A5] = "a5"
	squareName[B5] = "b5"
	squareName[C5] = "c5"
	squareName[D5] = "d5"
	squareName[E5] = "e5"
	squareName[F5] = "f5"
	squareName[G5] = "g5"
	squareName[H5] = "h5"

	squareName[A6] = "a6"
	squareName[B6] = "b6"
	squareName[C6] = "c6"
	squareName[D6] = "d6"
	squareName[E6] = "e6"
	squareName[F6] = "f6"
	squareName[G6] = "g6"
	squareName[H6] = "h6"

	squareName[A7] = "a7"
	squareName[B7] = "b7"
	squareName[C7] = "c7"
	squareName[D7] = "d7"
	squareName[E7] = "e7"
	squareName[F7] = "f7"
	squareName[G7] = "g7"
	squareName[H7] = "h7"

	squareName[A8] = "a8"
	squareName[B8] = "b8"
	squareName[C8] = "c8"
	squareName[D8] = "d8"
	squareName[E8] = "e8"
	squareName[F8] = "f8"
	squareName[G8] = "g8"
	squareName[H8] = "h8"

	squareNameToInt["a1"] = A1
	squareNameToInt["b1"] = B1
	squareNameToInt["c1"] = C1
	squareNameToInt["d1"] = D1
	squareNameToInt["e1"] = E1
	squareNameToInt["f1"] = F1
	squareNameToInt["g1"] = G1
	squareNameToInt["h1"] = H1

	squareNameToInt["a2"] = A2
	squareNameToInt["b2"] = B2
	squareNameToInt["c2"] = C2
	squareNameToInt["d2"] = D2
	squareNameToInt["e2"] = E2
	squareNameToInt["f2"] = F2
	squareNameToInt["g2"] = G2
	squareNameToInt["h2"] = H2

	squareNameToInt["a3"] = A3
	squareNameToInt["b3"] = B3
	squareNameToInt["c3"] = C3
	squareNameToInt["d3"] = D3
	squareNameToInt["e3"] = E3
	squareNameToInt["f3"] = F3
	squareNameToInt["g3"] = G3
	squareNameToInt["h3"] = H3

	squareNameToInt["a4"] = A4
	squareNameToInt["b4"] = B4
	squareNameToInt["c4"] = C4
	squareNameToInt["d4"] = D4
	squareNameToInt["e4"] = E4
	squareNameToInt["f4"] = F4
	squareNameToInt["g4"] = G4
	squareNameToInt["h4"] = H4

	squareNameToInt["a5"] = A5
	squareNameToInt["b5"] = B5
	squareNameToInt["c5"] = C5
	squareNameToInt["d5"] = D5
	squareNameToInt["e5"] = E5
	squareNameToInt["f5"] = F5
	squareNameToInt["g5"] = G5
	squareNameToInt["h5"] = H5

	squareNameToInt["a6"] = A6
	squareNameToInt["b6"] = B6
	squareNameToInt["c6"] = C6
	squareNameToInt["d6"] = D6
	squareNameToInt["e6"] = E6
	squareNameToInt["f6"] = F6
	squareNameToInt["g6"] = G6
	squareNameToInt["h6"] = H6

	squareNameToInt["a7"] = A7
	squareNameToInt["b7"] = B7
	squareNameToInt["c7"] = C7
	squareNameToInt["d7"] = D7
	squareNameToInt["e7"] = E7
	squareNameToInt["f7"] = F7
	squareNameToInt["g7"] = G7
	squareNameToInt["h7"] = H7

	squareNameToInt["a8"] = A8
	squareNameToInt["b8"] = B8
	squareNameToInt["c8"] = C8
	squareNameToInt["d8"] = D8
	squareNameToInt["e8"] = E8
	squareNameToInt["f8"] = F8
	squareNameToInt["g8"] = G8
	squareNameToInt["h8"] = H8

	pieceITOA[Pawn] = "P"
	pieceITOA[Knight] = "N"
	pieceITOA[Bishop] = "B"
	pieceITOA[Rook] = "R"
	pieceITOA[Queen] = "Q"
	pieceITOA[King] = "K"
}
