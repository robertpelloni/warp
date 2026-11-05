package editor

import (
	"unicode"
)

// TokenType represents a lexical token class for code syntax highlighting.
type TokenType int

const (
	TokenText TokenType = iota
	TokenKeyword
	TokenString
	TokenComment
	TokenNumber
	TokenOperator
	TokenIdentifier
)

// Token represents a highlighted token span within a source file.
type Token struct {
	Type  TokenType
	Text  string
	Start int
	End   int
}

var goKeywords = map[string]bool{
	"package": true, "import": true, "func": true, "return": true,
	"var": true, "const": true, "type": true, "struct": true, "interface": true,
	"if": true, "else": true, "for": true, "range": true, "switch": true,
	"case": true, "default": true, "select": true, "go": true, "defer": true,
	"chan": true, "map": true, "break": true, "continue": true,
}

var rustKeywords = map[string]bool{
	"fn": true, "let": true, "mut": true, "struct": true, "enum": true,
	"impl": true, "trait": true, "pub": true, "use": true, "mod": true,
	"match": true, "if": true, "else": true, "loop": true, "while": true,
	"for": true, "return": true, "self": true, "Self": true, "type": true,
}

var pythonKeywords = map[string]bool{
	"def": true, "class": true, "import": true, "from": true, "return": true,
	"if": true, "elif": true, "else": true, "for": true, "while": true,
	"try": true, "except": true, "finally": true, "with": true, "as": true,
	"pass": true, "break": true, "continue": true, "lambda": true, "yield": true,
}

// Tokenize processes content and returns categorized tokens based on language.
func Tokenize(content string, lang string) []Token {
	var tokens []Token
	runes := []rune(content)
	n := len(runes)
	i := 0

	for i < n {
		r := runes[i]

		// Line comments
		if (r == '/' && i+1 < n && runes[i+1] == '/') || (lang == "python" && r == '#') {
			start := i
			for i < n && runes[i] != '\n' {
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Text:  string(runes[start:i]),
				Start: start,
				End:   i,
			})
			continue
		}

		// Block comments (/* ... */)
		if r == '/' && i+1 < n && runes[i+1] == '*' {
			start := i
			i += 2
			for i+1 < n && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			if i+1 < n {
				i += 2
			} else {
				i = n
			}
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Text:  string(runes[start:i]),
				Start: start,
				End:   i,
			})
			continue
		}

		// String literals ("..." or '...' or `...`)
		if r == '"' || r == '\'' || r == '`' {
			quote := r
			start := i
			i++
			for i < n {
				if runes[i] == '\\' && i+1 < n && quote != '`' {
					i += 2
					continue
				}
				if runes[i] == quote {
					i++
					break
				}
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenString,
				Text:  string(runes[start:i]),
				Start: start,
				End:   i,
			})
			continue
		}

		// Numeric literals
		if unicode.IsDigit(r) {
			start := i
			for i < n && (unicode.IsDigit(runes[i]) || runes[i] == '.' || runes[i] == 'x' || runes[i] == 'b') {
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenNumber,
				Text:  string(runes[start:i]),
				Start: start,
				End:   i,
			})
			continue
		}

		// Identifiers and Keywords
		if unicode.IsLetter(r) || r == '_' {
			start := i
			for i < n && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
			}
			word := string(runes[start:i])
			tokType := TokenIdentifier
			switch lang {
			case "go":
				if goKeywords[word] {
					tokType = TokenKeyword
				}
			case "rust":
				if rustKeywords[word] {
					tokType = TokenKeyword
				}
			case "python":
				if pythonKeywords[word] {
					tokType = TokenKeyword
				}
			default:
				if goKeywords[word] || rustKeywords[word] || pythonKeywords[word] {
					tokType = TokenKeyword
				}
			}
			tokens = append(tokens, Token{
				Type:  tokType,
				Text:  word,
				Start: start,
				End:   i,
			})
			continue
		}

		// Operators & Punctuation
		if isOperator(r) {
			tokens = append(tokens, Token{
				Type:  TokenOperator,
				Text:  string(r),
				Start: i,
				End:   i + 1,
			})
			i++
			continue
		}

		// Plain text / whitespace
		start := i
		for i < n && !unicode.IsLetter(runes[i]) && !unicode.IsDigit(runes[i]) && !isOperator(runes[i]) &&
			runes[i] != '"' && runes[i] != '\'' && runes[i] != '`' && runes[i] != '/' && runes[i] != '#' {
			i++
		}
		if i == start {
			i++
		}
		tokens = append(tokens, Token{
			Type:  TokenText,
			Text:  string(runes[start:i]),
			Start: start,
			End:   i,
		})
	}

	return tokens
}

func isOperator(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '%', '=', '!', '<', '>', '&', '|', '^', '~', ':', ';', ',', '.', '(', ')', '[', ']', '{', '}':
		return true
	default:
		return false
	}
}
