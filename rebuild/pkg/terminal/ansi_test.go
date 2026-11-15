package terminal

import "testing"

func TestANSIParser(t *testing.T) {
	p := NewANSIParser(80, 24)

	input := []byte("Hello \x1b[31mRed\x1b[0m Text")
	expected := "Hello Red Text"

	result := p.Parse(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}

	input = []byte("Command\x1b[1;32mExecution\x1b[m")
	expected = "CommandExecution"
	result = p.Parse(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func BenchmarkANSIParser(b *testing.B) {
	p := NewANSIParser(80, 24)
	input := []byte("Hello \x1b[31mRed\x1b[0m Text with \x1b[1;32mExecution\x1b[m codes")
	for i := 0; i < b.N; i++ {
		p.Parse(input)
	}
}
