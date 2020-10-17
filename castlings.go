package board

// castlings type to represent the castlings
type castlings uint8

func (c *castlings) set(cast castlings) {
	*c |= cast
}
func (c *castlings) unset(cast castlings) {
	*c &= ^cast
}

// the different castlings
const (
	shortW = castlings(0b1000)
	longW  = castlings(0b0100)
	shortB = castlings(0b0010)
	longB  = castlings(0b0001)
)
