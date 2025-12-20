package game

type BotAI struct {
	board *Board
	botSymbol int
	opponentSymbol int
}

// GetBestMove returns the optimal column using a priority-based strategy
func (b *BotAI) GetBestMove() int {
	if col := b.findWinningMove(b.botSymbol); col != -1 {
		return col
	}

	if col := b.findWinningMove(b.opponentSymbol); col != -1 {
		return col
	}

	if col := b.findThreateningMove(b.botSymbol); col != -1 {
		return col
	}

	if col := b.findThreateningMove(b.opponentSymbol); col != -1 {
		return col
	}

	centerCols := []int{3, 2, 4, 1, 5, 0, 6}
	for _, col := range centerCols {
		if b.isValidMove(col) {
			return col
		}
	}

	for c := 0; c < Cols; c++ {
		if b.isValidMove(c) {
			return c
		}
	}

	return -1
}

// findWinningMove returns a column that would result in an immediate win
func (b *BotAI) findWinningMove(symbol int) int {
	for col := 0; col < Cols; col++ {
		if !b.isValidMove(col) {
			continue
		}

		row := b.getNextRow(col)
		if row == -1 {
			continue
		}

		b.board.Grid[row][col] = symbol

		if b.board.CheckWin(row, col, symbol) {
			b.board.Grid[row][col] = 0
			return col
		}

		b.board.Grid[row][col] = 0
	}
	return -1
}

// findThreateningMove returns a column that would create three connected pieces
func (b *BotAI) findThreateningMove(symbol int) int {
	for col := 0; col < Cols; col++ {
		if !b.isValidMove(col) {
			continue
		}

		row := b.getNextRow(col)
		if row == -1 {
			continue
		}

		b.board.Grid[row][col] = symbol

		if b.countConnected(row, col, symbol) >= 3 {
			b.board.Grid[row][col] = 0
			return col
		}

		b.board.Grid[row][col] = 0
	}
	return -1
}

// countConnected counts the maximum connected pieces for a position
func (b *BotAI) countConnected(row, col, symbol int) int {
	maxCount := 1

	// Check horizontal
	count := 1
	// Left
	for c := col - 1; c >= 0 && b.board.Grid[row][c] == symbol; c-- {
		count++
	}
	// Right
	for c := col + 1; c < Cols && b.board.Grid[row][c] == symbol; c++ {
		count++
	}
	if count > maxCount {
		maxCount = count
	}

	// Check vertical
	count = 1
	// Down only (can't go up in Connect Four)
	for r := row + 1; r < Rows && b.board.Grid[r][col] == symbol; r++ {
		count++
	}
	if count > maxCount {
		maxCount = count
	}

	// Check diagonal (top-left to bottom-right)
	count = 1
	// Up-left
	for r, c := row-1, col-1; r >= 0 && c >= 0 && b.board.Grid[r][c] == symbol; r, c = r-1, c-1 {
		count++
	}
	// Down-right
	for r, c := row+1, col+1; r < Rows && c < Cols && b.board.Grid[r][c] == symbol; r, c = r+1, c+1 {
		count++
	}
	if count > maxCount {
		maxCount = count
	}

	// Check diagonal (top-right to bottom-left)
	count = 1
	// Up-right
	for r, c := row-1, col+1; r >= 0 && c < Cols && b.board.Grid[r][c] == symbol; r, c = r-1, c+1 {
		count++
	}
	// Down-left
	for r, c := row+1, col-1; r < Rows && c >= 0 && b.board.Grid[r][c] == symbol; r, c = r+1, c-1 {
		count++
	}
	if count > maxCount {
		maxCount = count
	}

	return maxCount
}

// isValidMove checks if a column is not full
func (b *BotAI) isValidMove(col int) bool {
	if col < 0 || col >= Cols {
		return false
	}
	return b.board.Grid[0][col] == 0
}

// getNextRow returns the next available row in a column
func (b *BotAI) getNextRow(col int) int {
	for r := Rows - 1; r >= 0; r-- {
		if b.board.Grid[r][col] == 0 {
			return r
		}
	}
	return -1
}
