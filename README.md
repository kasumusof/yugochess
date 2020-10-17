# yugochess

  

yugochess is a chess board written in golang with the go standard library. It can be used to read fen strings, make fen strings from position, play moves and of course generate moves.

  

## Getting Started

To create a board with the new standard position

  

```
package main

  

import (

board "github.com/kasumusof/yugochess"

)

  

func main() {

b := board.NewBoardFromFEN(board.StartPos)

board.DrawBoard(b)

}

```

  
  

### Prerequisite