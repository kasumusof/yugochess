package board

import (
	"fmt"
	"time"
)

var now time.Time

func init() {
	initDisplay()
	initKingAttacks()
	initKnightAttacks()
	initSquareName()
}

func main() {
	fen := "k7/8/8/8/8/8/7R/1RK5 w - - 0 1"
	fen = "k5r1/7P/8/8/8/8/8/K7 w - - 0 1"
	// fen = "1rk5/7r/8/8/8/8/8/K7 b - - 0 1"
	// fen = " 2bqkbn1/2pppp2/np2N3/r3P1p1/p2N2B1/5Q2/PPPPKPP1/RNB2r2 w KQkq - 0 1"
	// fen = " k7/8/8/8/8/pppppppp/PPPPPPPP/R3K2R w KQkq - 0 1"
	// fen = startPos
	// fen = "rnbqkbnr/ppp3pp/4p3/3pPp2/3P4/8/PPP2PPP/RNBQKBNR w KQkq f6 0 4"

	b := fenToBoard(fen)
	// b.listen(InputChan())
	// b.turn = BLACK
	// c := genMovesAndEval(b)
	// fmt.Println(c)
	// fmt.Println(search(b, 1))
	// b.listen(InputChan())
	DrawBoard(b)
	mv, _ := b.uciStrToMv("e8g8")
	fmt.Println(mv, mv.flag)
	fmt.Println(colors[b.turn], "to play")
	fmt.Println(b.AvailableMoves(b.turn))
}
