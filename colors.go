package board

// color type to represent black and white
type color int

func (c *color) opp() color {
	if *c == WHITE {
		return BLACK
	}
	return WHITE
}

func (c *color) flip() {
	*c = c.opp()
}
