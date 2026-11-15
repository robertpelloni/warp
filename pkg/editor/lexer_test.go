package editor

import (
	"testing"
)

func TestTokenizeGo(t *testing.T) {
	src := "package main\n// greeting\nfunc main() {\n\tconst x = 42\n}"
	tokens := Tokenize(src, "go")

	hasPkg := false
	hasComment := false
	hasFunc := false
	hasNum := false

	for _, tok := range tokens {
		if tok.Type == TokenKeyword && tok.Text == "package" {
			hasPkg = true
		}
		if tok.Type == TokenComment && tok.Text == "// greeting" {
			hasComment = true
		}
		if tok.Type == TokenKeyword && tok.Text == "func" {
			hasFunc = true
		}
		if tok.Type == TokenNumber && tok.Text == "42" {
			hasNum = true
		}
	}

	if !hasPkg {
		t.Error("expected keyword 'package'")
	}
	if !hasComment {
		t.Error("expected comment '// greeting'")
	}
	if !hasFunc {
		t.Error("expected keyword 'func'")
	}
	if !hasNum {
		t.Error("expected number '42'")
	}
}

func TestTokenizeRust(t *testing.T) {
	src := "fn add(a: i32) -> i32 { let mut x = 10; x }"
	tokens := Tokenize(src, "rust")

	hasFn := false
	hasMut := false

	for _, tok := range tokens {
		if tok.Type == TokenKeyword && tok.Text == "fn" {
			hasFn = true
		}
		if tok.Type == TokenKeyword && tok.Text == "mut" {
			hasMut = true
		}
	}

	if !hasFn {
		t.Error("expected keyword 'fn'")
	}
	if !hasMut {
		t.Error("expected keyword 'mut'")
	}
}

func TestTokenizePython(t *testing.T) {
	src := "# Python script\ndef hello():\n\treturn \"world\""
	tokens := Tokenize(src, "python")

	hasDef := false
	hasComment := false
	hasStr := false

	for _, tok := range tokens {
		if tok.Type == TokenKeyword && tok.Text == "def" {
			hasDef = true
		}
		if tok.Type == TokenComment && tok.Text == "# Python script" {
			hasComment = true
		}
		if tok.Type == TokenString && tok.Text == "\"world\"" {
			hasStr = true
		}
	}

	if !hasDef {
		t.Error("expected keyword 'def'")
	}
	if !hasComment {
		t.Error("expected comment '# Python script'")
	}
	if !hasStr {
		t.Error("expected string '\"world\"'")
	}
}
