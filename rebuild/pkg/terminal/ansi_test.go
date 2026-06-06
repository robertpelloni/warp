package terminal

import "testing"

func TestANSIParser(t *testing.T) {
	p := NewANSIParser()

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
