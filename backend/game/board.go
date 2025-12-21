package game

import "errors"

const (
	Rows = 6
	Cols = 7
)

type Board struct {
	Grid [Rows][Cols]int
}

func NewBoard() *Board {
	return &Board{}
}

// DropDisc places a disc in the specified column and returns the landing row
func (b *Board) DropDisc(col int, player int) (int, error) {
	if col < 0 || col >= Cols {
		return -1, errors.New("invalid column")
	}

	for r := Rows - 1; r >= 0; r-- {
		if b.Grid[r][col] == 0 {
			b.Grid[r][col] = player
			return r, nil
		}
	}

	return -1, errors.New("column full")
}

// CheckWin returns true if the last move resulted in a win
func (b *Board) CheckWin(row, col, player int) bool {
	directions := [][2]int{
		{0, 1},
		{1, 0},
		{1, 1},
		{1, -1},
	}

	for _, d := range directions {
		count := 1

		// Check forward
		for i := 1; i < 4; i++ {
			r, c := row+d[0]*i, col+d[1]*i
			if r < 0 || r >= Rows || c < 0 || c >= Cols || b.Grid[r][c] != player {
				break
			}
			count++
		}

		// Check backward
		for i := 1; i < 4; i++ {
			r, c := row-d[0]*i, col-d[1]*i
			if r < 0 || r >= Rows || c < 0 || c >= Cols || b.Grid[r][c] != player {
				break
			}
			count++
		}

		if count >= 4 {
			return true
		}
	}

	return false
}

// IsFull returns true if the board has no empty cells
func (b *Board) IsFull() bool {
	for c := 0; c < Cols; c++ {
		if b.Grid[0][c] == 0 {
			return false
		}
	}
	return true
}
