package terminal

import (
	"bytes"
)

type ANSIParser struct {
	Buffer bytes.Buffer
	State  int
	Cols   int
	Rows   int
	X      int
	Y      int
}

const (
	StateNormal = iota
	StateESC
	StateCSI
)

func NewANSIParser(cols, rows int) *ANSIParser {
	return &ANSIParser{State: StateNormal, Cols: cols, Rows: rows}
}

// Parse processes raw bytes and extracts human-readable text,
// stripping or handling basic ANSI escape sequences.
// This is a foundational state machine for terminal emulation.
func (p *ANSIParser) Parse(data []byte) string {
	var result bytes.Buffer
	for _, b := range data {
		switch p.State {
		case StateNormal:
			if b == 0x1b {
				p.State = StateESC
			} else {
				result.WriteByte(b)
				if b == '\n' {
					p.X = 0
					p.Y++
				} else if b != '\r' {
					p.X++
					if p.X >= p.Cols {
						p.X = 0
						p.Y++
					}
				}
				if p.Y >= p.Rows {
					p.Y = p.Rows - 1
				}
			}
		case StateESC:
			if b == '[' {
				p.State = StateCSI
			} else {
				// Non-CSI escape sequence, just return to normal for now
				p.State = StateNormal
			}
		case StateCSI:
			if b >= 0x40 && b <= 0x7e {
				// End of CSI sequence
				p.State = StateNormal
			}
			// Intermediate bytes are ignored for stripping
		}
	}
	return result.String()
}
