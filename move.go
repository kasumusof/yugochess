package board

import "fmt"

type move struct {
	fr   int
	to   int
	flag int
	str  string
}

// moves struct holding a list of move
type moves []move

func (m *moves) push(mv move) {
	*m = append(*m, mv)
}

// TODO: make this return errors
func (m *moves) pop() (*move, error) {
	var idx int
	if len(*m) == 0 {
		return nil, fmt.Errorf("moves empty already")
	}
	idx = len(*m) - 1
	mv := (*m)[idx]
	*m = (*m)[:idx]
	return &mv, nil
}

// TODO might improve move flags
const (
	mvQuiet   = 0b0000
	mvDbPawn  = 0b0001
	mvCastSh  = 0b0010
	mvCastLn  = 0b0011
	mvCapture = 0b0100
	mvEnp     = 0b0101

	mvNProm = 0b1000
	mvBProm = 0b1001
	mvRProm = 0b1010
	mvQProm = 0b1011

	mvNPromCap = 0b1100
	mvBPromCap = 0b1101
	mvRPromCap = 0b1110
	mvQPromCap = 0b1111
)
