package board

import (
	"math/bits"
)

// BitBoard uint64 based type to manipulate board
type bitBoard uint64

func (bb *bitBoard) clear() {
	*bb = 0
}
func (bb *bitBoard) set(sq int) {
	*bb |= bitBoard(1 << sq)
}
func (bb *bitBoard) unset(sq int) {
	*bb &= ^bitBoard(1 << sq)
}
func (bb *bitBoard) countSet() int {
	return bits.OnesCount(uint(*bb))
}
func (bb *bitBoard) firstSet() int {
	return bits.TrailingZeros(uint(*bb))
}
func (bb *bitBoard) nextSet() int {
	a := bb.firstSet()
	bb.unset(a)
	return a
}
