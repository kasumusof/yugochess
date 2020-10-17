package board

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Board that holds the chess position
type Board struct {
	square    [64]int
	colors    [2]bitBoard
	pieces    [6]bitBoard
	turn      color
	castlings castlings
	rule50    int
	moveNo    int
	enpassant int
	history   moves
	state     states
}

// NewBoard function creates an empty board
func NewBoard() *Board {
	b := Board{}
	b.Clear()
	return &b
}

// NewStdBoard function creates board with std chess position in place
func NewStdBoard() *Board {
	b := fenToBoard(StartPos)
	return b
}

// Clear make a chess board empty
func (b *Board) Clear() {
	for i := A1; i <= H8; i++ {
		b.square[i] = Empty
	}
	for i := 0; i < 2; i++ {
		b.colors[i].clear()
	}
	for i := 0; i < 6; i++ {
		b.pieces[i].clear()
	}
	b.turn = WHITE
	b.castlings = 0x0
	b.rule50 = 0
	b.moveNo = 1
	b.enpassant = 0
	b.history = moves{}
	b.state = states{}
}

func (b *Board) set(piece int, sq int) {
	col := WHITE
	p := abs(piece) - 1
	if piece < 0 {
		col = BLACK
	}
	b.unset(sq)
	b.square[sq] = piece
	b.colors[col].set(sq)
	b.pieces[p].set(sq)
	return
}

func (b *Board) unset(sq int) {
	p12 := b.square[sq]
	if p12 == Empty {
		return
	}
	col := WHITE
	p6 := abs(p12) - 1
	if p12 < 0 {
		col = BLACK
	}
	b.square[sq] = Empty
	b.colors[col].unset(sq)
	b.pieces[p6].unset(sq)
}

func (b *Board) allPieces() bitBoard {
	return b.colors[WHITE] | b.colors[BLACK]
}

func (b *Board) bitBoard(col color, piece int) bitBoard {
	return b.colors[col] & b.pieces[piece]
}

func (b *Board) strtAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	var dMove bitBoard
	// to go up
	for i := sq + N; sq/8 != 7 && i <= H8; i += N {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}

	// to go down
	for i := sq + S; sq/8 != 0 && i >= A1; i += S {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}

	// to go right
	for i := sq + E; sq%8 != 7 && i <= H8 && i%8 != 0; i += E {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}

	// to go left
	for i := sq + W; sq%8 != 0 && i >= A1 && i%8 != 7; i += W {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}

	return mv
}
func (b *Board) diagAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	var dMove bitBoard
	// to go north-east conditions are, in order; not from the last rank, not from the last file, not back to the first file(for west<->east movements) ,cant go beyond the last square
	for i := sq + NE; sq/8 != 7 && sq%8 != 7 && i%8 != 0 && i <= H8; i += NE {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}
	// to go south-west conditions are, in order; not from the first rank, not from the first file, not back to the last file(for west<->east movements) ,cant go beyond the first square
	for i := sq + SW; sq/8 != 0 && sq%8 != 0 && i%8 != 7 && i >= A1; i += SW {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}
	// to go north-west conditions are, in order; not from the last rank, not from the first file, not back to the last file(for west<->east movements) ,cant go beyond the last square
	for i := sq + NW; sq/8 != 7 && sq%8 != 0 && i%8 != 7 && i <= H8; i += NW {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}
	// to go south-east conditions are, in order; not from the first rank, not from the first file, not back to the last file(for west<->east movements) ,cant go beyond the first square
	for i := sq + SE; sq/8 != 0 && sq%8 != 7 && i%8 != 0 && i >= A1; i += SE {
		dMove = bitBoard(1 << i)
		if b.allPieces()&dMove != 0 {
			mv |= dMove
			break
		}
		mv |= dMove
	}
	return mv
}

// gen attacks
func (b *Board) kingAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	mv |= kingAttacks[sq]
	return mv
}
func (b *Board) knightAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	mv |= knightAttacks[sq]
	return mv
}
func (b *Board) rookAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	mv |= b.strtAttacks(sq)
	return mv
}
func (b *Board) bishopAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	mv |= b.diagAttacks(sq)
	return mv
}
func (b *Board) queenAttacks(sq int) bitBoard {
	mv := bitBoard(0)
	mv |= b.diagAttacks(sq) | b.strtAttacks(sq)
	return mv
}
func (b *Board) pawnAttacks(sq int, col color) bitBoard {
	var mv, lAttack, rAttack bitBoard
	if col == WHITE {
		lAttack = (bitBoard(1<<sq) &^ FileA) << NW
		rAttack = (bitBoard(1<<sq) &^ FileH) << NE
	} else {
		lAttack = (bitBoard(1<<sq) &^ FileH) >> NW
		rAttack = (bitBoard(1<<sq) &^ FileA) >> NE
	}

	mv |= lAttack | rAttack
	return mv
}

func (b *Board) genAttacks(col color) bitBoard {
	attack := bitBoard(0)
	for p := 0; p < 6; p++ {
		piece := b.bitBoard(col, p)
		for n := piece; n.countSet() != 0; n.nextSet() {
			sq := n.firstSet()
			switch p {
			case Pawn:
				attack |= b.pawnAttacks(sq, col)
			case Knight:
				attack |= b.knightAttacks(sq)
			case Bishop:
				attack |= b.bishopAttacks(sq)
			case Rook:
				attack |= b.rookAttacks(sq)
			case Queen:
				attack |= b.queenAttacks(sq)
			case King:
				attack |= b.kingAttacks(sq)
			}

		}
	}
	return attack
}

// to check if king is in check
func (b *Board) kingInCheck(col color) bool {
	return b.bitBoard(col, King)&b.genAttacks(col.opp()) != 0
}

// moves for the board
func (b *Board) knightMoves(sq int, col color) bitBoard {
	return b.knightAttacks(sq) & ^b.colors[col]
}
func (b *Board) bishopMoves(sq int, col color) bitBoard {
	return b.bishopAttacks(sq) & ^b.colors[col]
}
func (b *Board) rookMoves(sq int, col color) bitBoard {
	return b.rookAttacks(sq) & ^b.colors[col]
}
func (b *Board) queenMoves(sq int, col color) bitBoard {
	return b.queenAttacks(sq) & ^b.colors[col]
}
func (b *Board) kingMoves(sq int, col color) bitBoard {
	cast := bitBoard(0)

	if col == WHITE {
		if b.castlings&shortW != 0 && b.castlingsCond(col, "S") && bitBoard(1<<sq)&Rank1 != 0 {
			cast.set(G1)
		}
		if b.castlings&longW != 0 && b.castlingsCond(col, "L") && bitBoard(1<<sq)&Rank1 != 0 {
			cast.set(C1)
		}
	}
	if col == BLACK {
		if b.castlings&shortB != 0 && b.castlingsCond(col, "S") && bitBoard(1<<sq)&Rank8 != 0 {
			cast.set(G8)
		}
		if b.castlings&longB != 0 && b.castlingsCond(col, "L") && bitBoard(1<<sq)&Rank8 != 0 {
			cast.set(C8)
		}
	}
	return (b.kingAttacks(sq) | cast) & ^b.colors[col]
}
func (b *Board) castlingsCond(col color, t string) bool {
	if b.kingInCheck(col) {
		return false
	}
	allPieces := b.allPieces()
	attacks := b.genAttacks(col.opp())
	if col == WHITE {
		switch t {
		case "L":
			if b.castlings&longW == 0 {
				return false
			}
			if allPieces&bitBoard(1<<C1) != 0 {
				return false
			}
			if attacks&bitBoard(1<<C1) != 0 {
				return false
			}
			if allPieces&bitBoard(1<<D1) != 0 {
				return false
			}
			if attacks&bitBoard(1<<D1) != 0 {
				return false
			}
		case "S":
			if b.castlings&0b1000 == 0 {
				return false
			}
			if allPieces&bitBoard(1<<G1) != 0 {
				return false
			}
			if attacks&bitBoard(1<<G1) != 0 {
				return false
			}
			if allPieces&bitBoard(1<<F1) != 0 {
				return false
			}
			if attacks&bitBoard(1<<F1) != 0 {
				return false
			}
		}
	} else {
		switch t {
		case "L":
			if b.castlings&0b0001 == 0 {
				return false
			}
			if allPieces&bitBoard(1<<C8) != 0 {
				return false
			}
			if attacks&bitBoard(1<<C8) != 0 {
				return false
			}
			if allPieces&bitBoard(1<<D8) != 0 {
				return false
			}
			if attacks&bitBoard(1<<D8) != 0 {
				return false
			}
		case "S":
			if b.castlings&0b0010 == 0 {
				return false
			}
			if allPieces&bitBoard(1<<G8) != 0 {
				return false
			}
			if attacks&bitBoard(1<<G8) != 0 {
				return false
			}
			if allPieces&bitBoard(1<<F8) != 0 {
				return false
			}
			if attacks&bitBoard(1<<F8) != 0 {
				return false
			}
		}
	}
	return true
}
func (b *Board) pawnEnp(sq int, col color) bitBoard {
	mv := bitBoard(0)
	if b.enpassant != 0 {
		if col == WHITE {
			if (sq+NE == b.enpassant && sq%8 != 7) || (sq+NW == b.enpassant && sq%8 != 0) {
				mv |= bitBoard(1 << b.enpassant)
			}
		} else {
			if (sq-NE == b.enpassant && sq%8 != 0) || (sq-NW == b.enpassant && sq%8 != 7) {
				mv |= bitBoard(1 << b.enpassant)
			}
		}
	}
	return mv
}
func (b *Board) pawnMoves(sq int, col color) bitBoard {
	var mv, one, two, p bitBoard
	p = bitBoard(1 << sq)
	allPieces := b.allPieces()

	if col == WHITE {
		one = (p << N) &^ allPieces
		two = ((one << N) & Rank4) &^ allPieces
	} else {
		one = (p >> N) &^ allPieces
		two = ((one >> N) & Rank5) &^ allPieces
	}
	mv = one | two
	// TODO: implement enpassant moves for real
	mv |= b.pawnEnp(sq, col)
	// if b.enpassant != 0 {
	// 	if col == WHITE {
	// 		if (sq+NE == b.enpassant && sq%8 != 7) || (sq+NW == b.enpassant && sq%8 != 0) { // && b.pawnAttacks(sq, col)&bitBoard(1<<b.enpassant) != 0 {
	// 			mv |= bitBoard(1 << b.enpassant)
	// 		}
	// 	} else {
	// 		if (sq-NE == b.enpassant && sq%8 != 0) || (sq-NW == b.enpassant && sq%8 != 7) { // && b.pawnAttacks(sq, col)&bitBoard(1<<b.enpassant) != 0 {
	// 			mv |= bitBoard(1 << b.enpassant)
	// 		}
	// 	}
	// }
	return mv | (b.pawnAttacks(sq, col) & b.colors[col.opp()])
}
func (b *Board) legalMoves(piece, sq int, col color) bitBoard {
	mvs := bitBoard(0)
	switch piece {
	case Pawn:
		mvs = b.pawnMoves(sq, col)
	case Knight:
		mvs = b.knightMoves(sq, col)
	case Bishop:
		mvs = b.bishopMoves(sq, col)
	case Rook:
		mvs = b.rookMoves(sq, col)
	case Queen:
		mvs = b.queenMoves(sq, col)
	case King:
		mvs = b.kingMoves(sq, col)
	}
	return mvs
}

// classifying move types for the move struct flags
func (b *Board) isTypeEnpassant(fr, to int, col color) bool {
	var attack, opp bitBoard
	var a = 0
	if b.enpassant != to {
		return false
	}
	if col == WHITE {
		a = to + S
	} else {
		a = to + N
	}
	// fmt.Println("from enpassant:", squareName[a])
	if a >= A1 && a <= H8 {
		opp = bitBoard(1 << a)
	}
	attack = opp & b.bitBoard(col.opp(), Pawn) &^ bitBoard(1<<to)
	return attack != 0
}

func (b *Board) isTypeCapture(piece, fr, to int, col color) bool {
	var attack bitBoard
	opp := b.colors[col.opp()] & bitBoard(1<<to)
	switch piece {
	case Pawn:
		attack = b.pawnMoves(fr, col)
		return b.isTypeEnpassant(fr, to, col) || attack&opp != 0
	default:
		attack = b.legalMoves(piece, fr, col)
	}
	return attack&opp != 0
}

func (b *Board) isTypeProm(fr, to int, col color) bool {

	if col == WHITE {
		return bitBoard(1<<to)&Rank8 != 0
	}
	return bitBoard(1<<to)&Rank1 != 0
}

func (b *Board) isDoublePawnMove(fr, to int, col color) bool {
	if col == WHITE {
		return fr/8 == 1 && fr+N+N == to
	}
	return fr/8 == 6 && fr+S+S == to
}

func (b *Board) isShortCaslte(fr, to int, col color) bool {
	if col == WHITE {
		return bitBoard(1<<fr)&Rank1 != 0 && bitBoard(1<<to)&FileG != 0
	}
	return bitBoard(1<<fr)&Rank8 != 0 && bitBoard(1<<to)&FileG != 0
}
func (b *Board) isLongCaslte(fr, to int, col color) bool {
	if col == WHITE {
		return bitBoard(1<<fr)&Rank1 != 0 && bitBoard(1<<to)&FileC != 0
	}
	return bitBoard(1<<fr)&Rank8 != 0 && bitBoard(1<<to)&FileC != 0
}

func (b *Board) packMove(piece, fr, to int, col color) move {
	mv := move{}
	mv.fr = fr
	mv.to = to
	mv.flag = mvQuiet
	cap := b.isTypeCapture(piece, fr, to, col)
	if cap {
		mv.flag = mvCapture
	}
	switch piece {
	case Pawn:
		if b.isTypeEnpassant(fr, to, col) {
			mv.flag |= mvEnp
		}
		if b.isDoublePawnMove(fr, to, col) {
			mv.flag |= mvDbPawn
		}
		if b.isTypeProm(fr, to, col) {
			mv.flag |= mvNProm
			mv.flag |= mvBProm
			mv.flag |= mvRProm
			mv.flag |= mvQProm
		}
		if b.isTypeProm(fr, to, col) && cap {
			mv.flag |= mvNPromCap
			mv.flag |= mvBPromCap
			mv.flag |= mvRPromCap
			mv.flag |= mvQPromCap
		}

	case King:
		if b.isShortCaslte(fr, to, col) {
			mv.flag = mvCastSh
		}
		if b.isLongCaslte(fr, to, col) {
			mv.flag = mvCastLn
		}
	}
	return mv
}

// get the square that is set for enpassant if a pawn move double
func (b *Board) getEpSqr(to int, col color) int {
	if col == WHITE {
		return to - N
	}
	return to + N
}

// get the square that the opponent pawn is on when enpassant was set
func (b *Board) getEpCapSqr(col color) int {
	if col == WHITE {
		return b.enpassant - N
	}
	return b.enpassant + N
}

// move to make a move on the board
func (b *Board) move(mv move) error {
	var err error
	prev := state{
		square:    b.square,
		colors:    b.colors,
		pieces:    b.pieces,
		enpassant: b.enpassant,
		rule50:    b.rule50,
		castlings: b.castlings,
	}
	fr := mv.fr
	to := mv.to
	p12 := b.square[fr]
	col := WHITE

	// // check for stalemate or checkmate on previous move
	// if b.checkMate() {
	// 	err = fmt.Errorf("checkmate: %s is victorious", colors[col.opp()])
	// 	return err
	// }
	// // check for stalemate or checkmate on previous move
	// if b.staleMate() {
	// 	err = fmt.Errorf("stalemate: draw")
	// 	return err
	// }

	if p12 < 0 {
		col = BLACK
	}
	p6 := abs(p12) - 1
	//TODO: moves validation
	//check to see if the from square is not empty
	if p12 == Empty {
		err = fmt.Errorf("from square is empty")
		return err
	}
	// check to see if the fr piece is of the right color
	if col != b.turn {
		err = fmt.Errorf("piece color is not of right color: %s turn", colors[col])
		return err
	}
	//check to see if move is available in piece moves
	if bitBoard(1<<to)&b.legalMoves(p6, fr, col) == 0 {
		err = fmt.Errorf("move is not legal")
		return err
	}

	if mv.flag == mvCastLn {
		if !b.castlingsCond(col, "L") {
			err = fmt.Errorf("%s cant castle long", colors[col])
			return err
		}
		a := A1
		d := D1
		rook := WRook
		if col == BLACK {
			a = A8
			d = D8
			rook = BRook
		}
		b.unset(a)
		b.set(rook, d)
	}
	if mv.flag == mvCastSh {
		if !b.castlingsCond(col, "S") {
			err = fmt.Errorf("%s cant castle short", colors[col])
			return err
		}
		a := H1
		d := F1
		rook := WRook
		if col == BLACK {
			a = H8
			d = F8
			rook = BRook
		}
		b.unset(a)
		b.set(rook, d)
	}
	if p6 == King {
		if col == WHITE {
			b.castlings.unset(longW)
			b.castlings.unset(shortW)
		} else {
			b.castlings.unset(longB)
			b.castlings.unset(shortB)
		}
	}
	if p6 == Rook {
		if fr == A1 {
			b.castlings.unset(longW)
		}
		if fr == H1 {
			// fmt.Println("I have lost short caslte")
			b.castlings.unset(shortW)
		}
		if fr == A8 {
			b.castlings.unset(shortB)
		}
		if fr == A8 {
			b.castlings.unset(shortB)
		}
	}
	b.unset(fr)
	b.set(p12, to)

	if mv.flag == mvEnp {
		sq := b.getEpCapSqr( /*fr,*/ col)
		// fmt.Println(squareName[sq])
		b.unset(sq)
	}

	// handling promotions
	if mv.flag == mvNProm || mv.flag == mvNPromCap {
		p6 := Knight
		p12 := p6 + 1
		if col == BLACK {
			p12 = -p12
		}
		b.set(p12, to)
	}
	if mv.flag == mvBProm || mv.flag == mvBPromCap {
		p6 := Bishop
		p12 := p6 + 1
		if col == BLACK {
			p12 = -p12
		}
		b.set(p12, to)
	}
	if mv.flag == mvRProm || mv.flag == mvRPromCap {
		p6 := Rook
		p12 := p6 + 1
		if col == BLACK {
			p12 = -p12
		}
		b.set(p12, to)
	}
	if mv.flag == mvQProm || mv.flag == mvQPromCap {
		p6 := Queen
		p12 := p6 + 1
		if col == BLACK {
			p12 = -p12
		}
		b.set(p12, to)
	}

	b.enpassant = 0
	if mv.flag&mvCapture != 0 || p6 == Pawn {
		b.rule50 = 0
	} else {
		b.rule50++
	}
	if b.isDoublePawnMove(fr, to, col) {
		b.enpassant = b.getEpSqr(to, col)
	}
	b.state.push(prev)
	b.history.push(mv)
	// incremeting the move no when it is black turn
	if col == BLACK {
		b.moveNo++
	}
	b.turn.flip()
	if b.kingInCheck(col) {
		err = fmt.Errorf("%s king in check", colors[col])
		b.unmove()
		return err
	}
	return nil
}

// unmove taking back a move
func (b *Board) unmove() error {
	prev, err := b.state.pop()
	if err != nil {
		return err
	}

	_, err = b.history.pop()
	if err != nil {
		return err
	}

	b.square = prev.square
	b.colors = prev.colors
	b.pieces = prev.pieces
	b.enpassant = prev.enpassant
	b.rule50 = prev.rule50
	b.castlings = prev.castlings
	if b.turn == WHITE {
		b.moveNo--
	}
	b.turn.flip()
	return nil
}

func (b *Board) moveToString(mv move) []string {
	var moves []string
	mov := ""
	p12 := b.square[mv.fr]
	p6 := abs(p12) - 1
	name := pieceITOA[p6]
	// fmt.Println("from move to string", name, pieceITOA, p6, mv.fr)
	fr := squareName[mv.fr]
	to := squareName[mv.to]
	cap := ""
	prom := ""
	if mv.flag == mvCastLn {
		mov = "O-O-O"
		moves = append(moves, mov)
		return moves
	}
	if mv.flag == mvCastSh {
		mov = "O-O"
		moves = append(moves, mov)
		return moves
	}
	if mv.flag&mvCapture != 0 {
		cap = "x"
	}
	if mv.flag&0b1000 != 0 {
		if mv.flag&mvNProm != 0 {
			prom = "N"
			mov = name + ":" + fr + cap + to + prom
			moves = append(moves, mov)
		}
		if mv.flag&mvBProm != 0 {
			prom = "B"
			mov = name + ":" + fr + cap + to + prom
			moves = append(moves, mov)
		}
		if mv.flag&mvRProm != 0 {
			prom = "R"
			mov = name + ":" + fr + cap + to + prom
			moves = append(moves, mov)
		}
		if mv.flag&mvQProm != 0 {
			prom = "Q"
			mov = name + ":" + fr + cap + to + prom
			moves = append(moves, mov)
		}
		return moves
	}
	mov = name + ":" + fr + cap + to + prom
	moves = append(moves, mov)
	return moves
}

// availableMovesStr  visual representation of the available moves
func (b *Board) availableMovesStr(col color) []string {
	var moves []string
	var err error
	check := make(map[string]int)
	for p := 0; p < 6; p++ {
		piece := b.bitBoard(col, p)
		for j := piece; j.countSet() != 0; j.nextSet() {
			fr := j.firstSet()
			mv := b.legalMoves(p, fr, col)
			for j := mv; j.countSet() != 0; j.nextSet() {
				to := j.firstSet()
				mvs := b.moveToString(b.packMove(p, fr, to, col))
				for _, str := range mvs {
					if contains(str, ':') {
						str = string(str[2:])
					}
					// str = strings.Split(str, ":")[1]
					if mv, ok := b.stringToMove(str, col); ok {
						err = b.move(mv)
						if err == nil {
							b.unmove()
							for _, toUse := range b.moveToString(mv) {
								if _, ok := check[toUse]; ok {
									continue
								}
								check[toUse] = 1
								moves = append(moves, toUse)
							}
						}
					}
				}
			}
		}
	}
	return moves
}

// availableMovesMv generates legal moves on the board
func (b *Board) availableMovesMv(col color) moves {
	mvs := moves{}

	for _, strmv := range b.availableMovesStr(col) {
		if mv, ok := b.stringToMove(strmv, col); ok {
			mv.str = strmv
			mvs.push(mv)
		}
	}
	return mvs
}

func (b *Board) availableMovesSet(col color) map[string]int {
	var check = make(map[string]int)
	for _, str := range b.availableMovesStr(col) {
		check[str] = 1
	}
	return check
}

func (b *Board) stringToMove(str string, col color) (move, bool) {
	if contains(str, ':') {
		str = string(str[2:])
	}
	str = trimSpaces(str)
	var mv move
	var ok bool
	mv.flag = 0
	if upper(str) == "O-O" {
		mv.flag = mvCastSh
		mv.fr = b.getKingPos(col)
		if col == WHITE {
			mv.to = G1
		} else {
			mv.to = G8
		}
		return mv, true
	} else if upper(str) == "O-O-O" {
		mv.flag = mvCastLn
		mv.fr = b.getKingPos(col)
		if col == WHITE {
			mv.to = C1
		} else {
			mv.to = C8
		}
		return mv, true
	} else if len(str) < 4 {
		return move{}, false
	}
	fr := string(str[:2])
	to := string(str[2:4])
	isCap := strings.Contains(str, "x")
	isProm := strings.Contains(upper(str), "Q") || strings.Contains(upper(str), "R") || strings.Contains(upper(str), "B") || strings.Contains(upper(str), "N")
	if isCap {
		mv.flag = mvCapture
		to = string(str[3:5])
	}
	if isProm {
		s := string(str[len(str)-1])
		switch upper(s) {
		case "Q":
			mv.flag = mvQProm
			if isCap && isProm {
				mv.flag = mvQPromCap
			}
		case "R":
			mv.flag = mvRProm
			if isCap && isProm {
				mv.flag = mvRPromCap
			}
		case "B":
			mv.flag = mvBProm
			if isCap && isProm {
				mv.flag = mvBPromCap
			}
		case "N":
			mv.flag = mvNProm
			if isCap && isProm {
				mv.flag = mvNPromCap
			}
		}
	}
	mv.fr, ok = squareNameToInt[fr]
	if !ok {
		return move{}, false
	}
	mv.to, ok = squareNameToInt[to]
	if !ok {
		return move{}, false
	}

	if b.isTypeEnpassant(mv.fr, mv.to, col) {
		mv.flag = mvEnp
	}
	if b.isDoublePawnMove(mv.fr, mv.to, col) {
		mv.flag = mvDbPawn
	}
	// mv.flag = 0
	return mv, true
}

func (b *Board) getKingPos(col color) int {
	v := b.bitBoard(col, King)
	return v.firstSet()
}

// listen function to listen for input on the board
func (b *Board) listen(inputChan chan string) {
	listenStart := time.Now()
	quit := false
	var input, cmd, others string
	// var err error
	for quit == false {
		select {
		case input = <-inputChan:
			input = trimSpaces(input)
			d := strings.Split(input, " ")
			cmd = d[0]
			others = strings.Join(d[1:], " ")
			// fmt.Println("\033[H\033[2J")
		}
		switch lower(cmd) {
		case "history", "h":
			count := 1
			for idx, mv := range b.history {
				if idx%2 == 1 {
					fmt.Printf("%s ", mv.str)
					count++
					continue
				}
				fmt.Printf("%v. %s ", count, mv.str)
			}
			fmt.Println()

		case "moveno":
			fmt.Println("move", b.moveNo)
		case "rule50":
			fmt.Println(b.rule50, "moves played without capturees or pawn moves")
		case "print", "p":
			// handlePrint(b, others)
			DrawBoard(b)
		case "new", "n":
			c := b
			b = fenToBoard(others)
			if b == nil {
				fmt.Println("invalid fen")
				b = c
			}
		case "startpos", "start":
			b = fenToBoard(StartPos)

		case "state":
			if b.StaleMate() {
				fmt.Printf("Stalemate: draw\n")
			} else if b.CheckMate() {
				fmt.Printf("checkmate: %s is victorious\n", colors[b.turn.opp()])
			} else {
				fmt.Printf("game in play %s to play\n", colors[b.turn])
			}
		case "unmove", "um":
			err := b.unmove()
			if err != nil {
				fmt.Println(err)
			}
		case "move", "m", "mv":
			var err error
			others = prepMoves(others)
			mvs := strings.Split(others, " ")
			for i := 0; i < len(mvs) && err == nil; i++ {
				if b.CheckMate() {
					err = fmt.Errorf("checkmate: %s is victorious", colors[b.turn.opp()])

				}
				// check for stalemate or checkmate on previous move
				if b.StaleMate() {
					err = fmt.Errorf("stalemate: draw")
				}
				str := mvs[i]
				a, err := b.validateStd(str, b.turn)
				if err != nil {
					fmt.Println(err, str)
				}
				mv, ok := b.stringToMove(a, b.turn)
				if ok {
					mv.str = str
					err := b.move(mv)
					if err != nil {
						fmt.Println(err, str)

						break
					}
					if b.CheckMate() {
						err = fmt.Errorf("checkmate: %s is victorious", colors[b.turn.opp()])

					}
					// check for stalemate or checkmate on previous move
					if b.StaleMate() {
						err = fmt.Errorf("stalemate: draw")
					}

				} else {
					fmt.Println("error understanding move", str)
					break
				}

			}
			if err != nil {
				fmt.Println(err)
			}

			// handleMove(b, others)

		case "fen", "f":
			fmt.Println(b.ToFEN())
		case "castlings", "castling", "castle", "castles":
			if b.castlings&longW != 0 {
				fmt.Print("white:long ")
			}
			if b.castlings&shortW != 0 {
				fmt.Print("white:short ")
			}
			if b.castlings&longB != 0 {
				fmt.Print("black:long ")
			}
			if b.castlings&shortB != 0 {
				fmt.Print("black:short ")
			}
			fmt.Println()
		case "check":
			if b.kingInCheck(WHITE) {
				fmt.Println("white king in check")
			} else if b.kingInCheck(BLACK) {
				fmt.Println("black king in check")
			} else {
				fmt.Println("No checks on board")
			}

		case "turn":
			fmt.Println(colors[b.turn], "to play")
		case "moves", "ms", "mvs":
			q := b.availableMovesStr(b.turn)
			fmt.Println("Legal moves for", colors[b.turn], ":", len(q))
			fmt.Println(q)
		case "quit", "q":
			quit = true
			fmt.Println("Quitting program")
		default:
			fmt.Println("error:", cmd, "is not a valid command")
		}
	}
	fmt.Println("\nfunc listen ran for:", time.Since(listenStart))
}

// TODO: make this function work
// also make it return all possible moves string and give errors if it is not explicit
func (b *Board) standardNotationToMvStr(str string, col color) ([]string, error) {
	str = trimSpaces(str)
	if endsWith(str, '+') {
		str = str[:len(str)-1]
	}
	if endsWith(str, '#') {
		str = str[:len(str)-1]
	}

	result := ""
	var ret []string
	var name, frSqr, cap, toSqr, prom string
	var err error

	//castles remain the same
	if lower(str) == "o-o" || lower(str) == "o-o-o" {
		result = upper(str)
		ret = append(ret, result)
		return ret, nil
	}

	// if len == 2 then move is a pawn move(e4)
	if len(str) == 2 {
		name = "P"
		to, ok := squareNameToInt[lower(str)]
		if !ok {
			err = fmt.Errorf("not a valid square name")
			return ret, err
		}
		// if "to" exists then it must be empty (not really) itself and there must be an empty square one move below it or 2 above it
		switch col {
		case WHITE:
			// checking for 2
			a, _ := strconv.Atoi(string(str[1]))
			if a == 8 {
				err = fmt.Errorf("please input a promotion piece")
				return ret, err
			}
			if a < 3 || a > 7 {
				err = fmt.Errorf("illegal move %s ", str)
				return ret, err
			}
			to1 := to + S
			// check for to1 not being out of range

			frSqr = squareName[to1]
			toSqr = str
			result = name + ":" + frSqr + cap + toSqr + prom
			ret = append(ret, result)

			// to2 exists
			if endsWith(str, '4') {
				to2 := to1 + S
				frSqr = squareName[to2]
				toSqr = str
				result = name + ":" + frSqr + cap + toSqr + prom
				ret = append(ret, result)
			}

		case BLACK:
			// checking for 2
			a, _ := strconv.Atoi(string(str[1]))
			if a < 2 || a > 6 {
				err = fmt.Errorf("illegal move %s ", str)
				return ret, err
			}
			to1 := to + N
			// check for to1 not being out of range

			// to2 exists
			if endsWith(str, '5') {
				to2 := to1 + N
				frSqr = squareName[to2]
				toSqr = str
				result = name + ":" + frSqr + cap + toSqr + prom
				ret = append(ret, result)
			}

			frSqr = squareName[to1]
			toSqr = str
			result = name + ":" + frSqr + cap + toSqr + prom
			ret = append(ret, result)

		}
	}

	// if len == 3 then move is a pawn promotion(a8q) or a piece move (Nf3)
	//TODO add one for the fancy pawn captures (ab3)
	if len(str) == 3 {
		begin := startsWithAny(str, "NBKRQ")
		end := endsWithAny(str, "NBKRQ")
		// we have a piece move
		if begin {
			name = string(str[0])
			toName := string(str[1:])
			toSqr = toName
			to, ok := squareNameToInt[toName]
			if !ok {
				err = fmt.Errorf("not a valid square name: %s", toName)
				return ret, err
			}
			bb := bitBoard(0)
			p12 := 0
			switch name {
			case "N":
				bb = b.knightAttacks(to)
				p12 = Knight + 1
				if col == BLACK {
					p12 = -p12
				}
			case "B":
				bb = b.bishopAttacks(to)
				p12 = Bishop + 1
				if col == BLACK {
					p12 = -p12
				}
			case "R":
				bb = b.rookAttacks(to)
				p12 = Rook + 1
				if col == BLACK {
					p12 = -p12
				}
			case "Q":
				bb = b.queenAttacks(to)
				p12 = Queen + 1
				if col == BLACK {
					p12 = -p12
				}
			case "K":
				bb = b.kingAttacks(to)
				p12 = King + 1
				if col == BLACK {
					p12 = -p12
				}
			default:
				err = fmt.Errorf("not a valid piece name: %s", name)
				return ret, err
			}
			// bb := (b.knightAttacks(to))
			for i := bb; i.countSet() != 0; i.nextSet() {
				if b.square[i.firstSet()] == p12 {

					frSqr = squareName[i.firstSet()]
					result = name + ":" + frSqr + toSqr
					ret = append(ret, result)
				}
			}

		}
		// we have a pawn promotion
		if end {
			name = "P"
			toName := trimSpaces(string(str[:2]))
			// pawn promotions must be to the first or last rank
			to, ok := squareNameToInt[lower(toName)]
			if !ok {
				err = fmt.Errorf("not a valid square name: %s", toName)
				return ret, err
			}
			if col == WHITE && string(toName[1]) != "8" {
				err = fmt.Errorf("not valid promotion square for %s", colors[col])
				return ret, err
			}
			if col == BLACK && string(toName[1]) != "1" {
				err = fmt.Errorf("not valid promotion square for %s", colors[col])
				return ret, err
			}
			prom = "Q"
			frSqr = squareName[to+S]
			toSqr = toName
			if col == BLACK {
				frSqr = squareName[to+N]
			}
			if endsWith(str, 'N') {
				prom = "N"
			}
			if endsWith(str, 'B') {
				prom = "B"
			}
			if endsWith(str, 'R') {
				prom = "R"
			}
			result = name + ":" + frSqr + cap + toSqr + prom
			ret = append(ret, result)
			// return ret, nil

		}
	}

	// if len == 4 then move is a pawn capture move(axb7) or piece specific move(N1c3, Nbc3) or piece capture (Nxb7)
	if len(str) == 4 {
		begin := startsWithAny(str, "NBKRQ")
		captr := contains(str, 'x')
		// we have a piece specific move
		if begin && !captr {
			name = string(str[0])
			toName := string(str[2:])
			flag := string(str[1])
			toSqr = toName
			to, ok := squareNameToInt[toName]
			if !ok {
				err = fmt.Errorf("not a valid square name: %s", toName)
				return ret, err
			}
			bb := bitBoard(0)
			p12 := 0
			intersect := bitBoard(0)

			switch upper(flag) {
			case "1":
				intersect = Rank1
			case "2":
				intersect = Rank2
			case "3":
				intersect = Rank3
			case "4":
				intersect = Rank4
			case "5":
				intersect = Rank5
			case "6":
				intersect = Rank6
			case "7":
				intersect = Rank7
			case "8":
				intersect = Rank8
			case "A":
				intersect = FileA
			case "B":
				intersect = FileB
			case "C":
				intersect = FileC
			case "D":
				intersect = FileD
			case "E":
				intersect = FileE
			case "F":
				intersect = FileF
			case "G":
				intersect = FileG
			case "H":
				intersect = FileH
			default:
				err = fmt.Errorf("not a valid flag: %s in %s", flag, str)
				return ret, err
			}
			switch name {
			case "N":
				bb = b.knightAttacks(to)
				p12 = Knight + 1
				if col == BLACK {
					p12 = -p12
				}
			case "B":
				bb = b.bishopAttacks(to)
				p12 = Bishop + 1
				if col == BLACK {
					p12 = -p12
				}
			case "R":
				bb = b.rookAttacks(to)
				p12 = Rook + 1
				if col == BLACK {
					p12 = -p12
				}
			case "Q":
				bb = b.queenAttacks(to)
				p12 = Queen + 1
				if col == BLACK {
					p12 = -p12
				}
			case "K":
				bb = b.kingAttacks(to)
				p12 = King + 1
				if col == BLACK {
					p12 = -p12
				}
			default:
				err = fmt.Errorf("not a valid piece name: %s", name)
				return ret, err
			}
			bb &= intersect
			// drawBitBoard(bb)
			// bb := (b.knightAttacks(to))
			for i := bb; i.countSet() != 0; i.nextSet() {
				if b.square[i.firstSet()] == p12 {

					frSqr = squareName[i.firstSet()]
					result = name + ":" + frSqr + toSqr
					ret = append(ret, result)
				}
			}

		} else if begin && captr { // we have a piece capture move
			cap = "x"
			name = string(str[0])

			toName := string(str[2:])
			redo := name + toName
			//TODO try making this better by including an x
			return b.standardNotationToMvStr(redo, col)
			// to, ok := squareNameToInt[toName]
			// if !ok {
			// 	err = fmt.Errorf("not a valid square name: %s", toName)
			// 	return ret, err
			// }
		} else if !begin && captr { // then we have a pawn capture
			name = "P"
			flag := string(str[0])
			toName := string(str[2:])
			// fmt.Println(flag, toName)
			// fmt.Println(flag[0], toName[0], flag[0]-toName[0])
			if abs(int(flag[0])-int(toName[0])) != 1 {
				err = fmt.Errorf("not a valid pawn capture: %s", str)
				return ret, err
			}
			cap = "x"
			to, ok := squareNameToInt[toName]
			if !ok {
				err = fmt.Errorf("not a valid square name: %s", toName)
				return ret, err
			}
			fr := 0
			if col == WHITE {
				if int(flag[0]) < int(toName[0]) {
					if to%8 != 0 {
						fr = to - NE
					}
				} else {
					if to%8 != 7 {
						fr = to - NW
					}
				}
			} else {
				if int(flag[0]) < int(toName[0]) {
					if to%8 != 0 {
						fr = to + NW
					}
				} else {
					if to%8 != 7 {
						fr = to + NE
					}
				}
			}
			frSqr = squareName[fr]
			toSqr = squareName[to]
			result = name + ":" + frSqr + cap + toSqr + prom
			ret = append(ret, result)
		}
	}

	// if len == 5 then move is a pawn capture promotion (axb8q) or piece specific capture
	if len(str) == 5 {
		capture := contains(str, 'x')
		pro := endsWithAny(str, "QNBR")
		// we have a pawn capture promotion
		if pro && capture {
			name = "P"
			prom = string(str[4])
			flag := string(str[0])
			toName := string(str[2:4])
			// fmt.Println(flag, toName)
			// fmt.Println(flag[0], toName[0], flag[0]-toName[0])
			if abs(int(flag[0])-int(toName[0])) != 1 {
				err = fmt.Errorf("not a valid pawn capture: %s", str)
				return ret, err
			}
			if !endsWith(toName, '8') && col == WHITE {
				err = fmt.Errorf("not a valid pawn promotion move: %s", str)
				return ret, err
			}
			if !endsWith(toName, '1') && col == BLACK {
				err = fmt.Errorf("not a valid pawn promotion move: %s", str)
				return ret, err
			}
			cap = "x"
			to, ok := squareNameToInt[toName]
			if !ok {
				err = fmt.Errorf("not a valid square name: %s", toName)
				return ret, err
			}
			fr := 0
			if col == WHITE {
				if int(flag[0]) < int(toName[0]) {
					if to%8 != 0 {
						fr = to - NE
					}
				} else {
					if to%8 != 7 {
						fr = to - NW
					}
				}
			} else {
				if int(flag[0]) < int(toName[0]) {
					if to%8 != 0 {
						fr = to + NW
					}
				} else {
					if to%8 != 7 {
						fr = to + NE
					}
				}
			}
			frSqr = squareName[fr]
			toSqr = squareName[to]
			result = name + ":" + frSqr + cap + toSqr + prom
			ret = append(ret, result)
		} else if capture { // piece specific capture
			redo := string(str[:2]) + string(str[3:])
			return b.standardNotationToMvStr(redo, col)
		}
	}
	// result = name + ":" + frSqr + cap + toSqr + prom
	return ret, nil
}

func (b *Board) validateStd(str string, col color) (string, error) {
	valid := ""
	var ids []int
	var err error
	mvs, err := b.standardNotationToMvStr(str, col)
	if err != nil {
		return "", err
	}
	if len(mvs) == 0 {
		err = fmt.Errorf("that is not validated valid move")
		return "", err
	}
	if len(mvs) == 1 {
		valid = mvs[0]
	}
	if len(mvs) > 1 {
		canDo := 0
		for idx, mv := range mvs {
			m, ok := b.stringToMove(mv, col)
			if !ok {
				continue
			}
			err = b.move(m)
			if err == nil {
				b.unmove()
				canDo++
				ids = append(ids, idx)
			}
		}
		if canDo == 0 {
			err = fmt.Errorf("that is not validated valid move")
			return "", err
		}
		if canDo == 1 {
			valid = mvs[ids[0]]
		}
		if canDo > 1 {
			err = fmt.Errorf("ambigious move, please be specific")
			return "", err
		}
	}
	return valid, nil
}

// moveToSTDNot converts move struct to standard notation
func (b *Board) moveToSTDNot(move move) string {
	result := ""
	var name, frSqr, toSqr, cap, prom string
	p12 := b.square[move.fr]
	p6 := abs(p12) - 1
	frSqr = squareName[move.fr]

	switch move.flag {
	case mvCapture, mvEnp:
		cap = "x"
	case mvCastLn:
		result = "O-O-O"
		return result
	case mvCastSh:
		result = "O-O"
		return result
	case mvNPromCap:
		cap = "x"
		prom = "N"
	case mvRPromCap:
		cap = "x"
		prom = "R"
	case mvBPromCap:
		cap = "x"
		prom = "B"
	case mvQPromCap:
		cap = "x"
		prom = "Q"
	case mvNProm:
		prom = "N"
	case mvRProm:
		prom = "R"
	case mvBProm:
		prom = "B"
	case mvQProm:
		prom = "Q"
	}

	switch p6 {
	case Pawn:
		name = ""
		result = name + string(frSqr[0]) + cap + toSqr + prom
		return result
	case Knight:
		name = "N"
	case Bishop:
		name = "B"
	case Rook:
		name = "R"
	case Queen:
		name = "Q"
	case King:
		name = "K"
	}

	result = name + frSqr + cap + toSqr

	return result
}

// Show to display the board
func (b *Board) Show() {
	DrawBoard(b)
}

// CheckMate checks if the current board posiition is a checkmate
func (b *Board) CheckMate() bool {
	return b.kingInCheck(b.turn) && len(b.availableMovesMv(b.turn)) == 0
}

// StaleMate checks if the current board posiition is a checkmate
func (b *Board) StaleMate() bool {
	return !b.kingInCheck(b.turn) && len(b.availableMovesMv(b.turn)) == 0
}

// AvailableMoves  visual representation of the available moves
func (b *Board) AvailableMoves(col color) []string {
	var moves []string
	for _, mv := range b.availableMovesMv(col) {
		prom := ""
		switch mv.flag {
		case mvNProm, mvNPromCap:
			prom = "N"
		case mvBProm, mvBPromCap:
			prom = "B"
		case mvRProm, mvRPromCap:
			prom = "R"
		case mvQProm, mvQPromCap:
			prom = "Q"
		}
		mov := squareName[mv.fr] + squareName[mv.to] + prom
		moves = append(moves, mov)
	}
	return moves
}

func (b *Board) uciStrToMv(str string) (move, error) {
	mv := move{}
	var fr, to, flag int
	var err error

	if len(str) < 4 {
		err = fmt.Errorf("invalid move string")
		return mv, err
	}
	fr = squareNameToInt[string(str[:2])]
	to = squareNameToInt[string(str[2:4])]
	flag = mvQuiet

	col := WHITE
	p12 := b.square[fr]
	if p12 < 0 {
		col = BLACK
	}
	fmt.Println(colors[col])
	p6 := abs(p12) - 1

	if b.isDoublePawnMove(fr, to, col) {
		flag = mvDbPawn
	}
	if b.isTypeCapture(p6, fr, to, col) {
		flag = mvCapture
	}
	if b.isShortCaslte(fr, to, col) {
		flag = mvCastSh
	}
	if b.isLongCaslte(fr, to, col) {
		flag = mvCastLn
	}
	if b.isTypeEnpassant(fr, to, col) {
		flag = mvEnp
	}

	if len(str) == 5 {
		switch string(str[4]) {
		case "N":
			mv.flag |= mvNProm
		case "B":
			mv.flag |= mvBProm
		case "R":
			mv.flag |= mvRProm
		case "Q":
			mv.flag |= mvQProm
		}
	}
	mv.fr = fr
	mv.to = to
	mv.flag = flag
	return mv, nil
}

// Move to make a move on the board
func (b *Board) Move(mov move) error {
	return b.move(mov)
}

// Unmove to make a move on the board
func (b *Board) Unmove() error {
	return b.unmove()
}

// NewBoardFromFEN to create a board from a fen string
func NewBoardFromFEN(fen string) *Board {
	b := fenToBoard(fen)
	return b
}
