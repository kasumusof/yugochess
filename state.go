package board

import "fmt"

type state struct {
	square    [64]int
	colors    [2]bitBoard
	pieces    [6]bitBoard
	enpassant int
	rule50    int
	castlings castlings
}

type states []state

func (s *states) push(st state) {
	*s = append(*s, st)
}

// TODO: make this return errors
func (s *states) pop() (*state, error) {
	var idx int
	if len(*s) == 0 {
		return nil, fmt.Errorf("state is empty")
	}
	idx = len(*s) - 1
	st := (*s)[idx]
	*s = (*s)[:idx]
	return &st, nil
}


