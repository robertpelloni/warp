package terminal

import (
	"strings"
	"sync"
)

type Cell struct {
	Char     rune
	FG, BG   string
}

type ScreenBuffer struct {
	Grid [][]Cell
	Cols int
	Rows int
	mu   sync.RWMutex
}

func NewScreenBuffer(cols, rows int) *ScreenBuffer {
	grid := make([][]Cell, rows)
	for i := range grid {
		grid[i] = make([]Cell, cols)
		for j := range grid[i] {
			grid[i][j] = Cell{Char: ' '}
		}
	}
	return &ScreenBuffer{Grid: grid, Cols: cols, Rows: rows}
}

func (b *ScreenBuffer) SetCell(x, y int, char rune) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if x >= 0 && x < b.Cols && y >= 0 && y < b.Rows {
		b.Grid[y][x].Char = char
	}
}

func (b *ScreenBuffer) GetString() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var sb strings.Builder
	for _, row := range b.Grid {
		for _, cell := range row {
			sb.WriteRune(cell.Char)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
