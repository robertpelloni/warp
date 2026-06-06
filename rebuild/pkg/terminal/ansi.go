package terminal

import (
	"bytes"
)

type ANSIParser struct {
	Buffer bytes.Buffer
}

func NewANSIParser() *ANSIParser {
	return &ANSIParser{}
}

// Parse processes raw bytes and extracts human-readable text,
// stripping or handling basic ANSI escape sequences.
// This is a foundational state machine for terminal emulation.
func (p *ANSIParser) Parse(data []byte) string {
	var result bytes.Buffer
	i := 0
	for i < len(data) {
		if data[i] == 0x1b { // ESC
			if i+1 < len(data) && data[i+1] == '[' {
				// CSI sequence
				j := i + 2
				for j < len(data) && (data[j] >= 0x30 && data[j] <= 0x3f) {
					j++
				}
				for j < len(data) && (data[j] >= 0x20 && data[j] <= 0x2f) {
					j++
				}
				if j < len(data) && (data[j] >= 0x40 && data[j] <= 0x7e) {
					// End of sequence
					i = j + 1
					continue
				}
			}
		}
		result.WriteByte(data[i])
		i++
	}
	return result.String()
}
