package set

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLBrace
	TokenRBrace
	TokenLParen
	TokenRParen
	TokenWord
	TokenString
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenWord:
		return "word"
	case TokenString:
		return "string"
	default:
		return fmt.Sprintf("Token(%d)", t)
	}
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Lexer struct {
	input []rune

	pos int

	line   int
	column int
}

func NewLexer(data []byte) *Lexer {
	return &Lexer{
		input:  []rune(string(data)),
		line:   1,
		column: 1,
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}

	return l.input[l.pos]
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.input) {
		return 0
	}

	r := l.input[l.pos]
	l.pos++

	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return r
}

func (l *Lexer) Next() (Token, error) {

	for {
		for unicode.IsSpace(l.peek()) {
			l.advance()
		}

		if l.peek() == 0 {
			return Token{
				Type:   TokenEOF,
				Line:   l.line,
				Column: l.column,
			}, nil
		}

		if l.peek() == ';' {
			for {
				r := l.advance()
				if r == '\n' || r == 0 {
					break
				}
			}

			continue
		}

		break
	}

	line := l.line
	column := l.column

	switch l.peek() {

	case '{':
		l.advance()
		return Token{
			Type:   TokenLBrace,
			Line:   line,
			Column: column,
		}, nil

	case '}':
		l.advance()
		return Token{
			Type:   TokenRBrace,
			Line:   line,
			Column: column,
		}, nil

	case '(':
		l.advance()
		return Token{
			Type:   TokenLParen,
			Line:   line,
			Column: column,
		}, nil

	case ')':
		l.advance()
		return Token{
			Type:   TokenRParen,
			Line:   line,
			Column: column,
		}, nil

	case '"':
		l.advance()
		var value []rune
		for {
			r := l.peek()
			if r == 0 {
				return Token{}, fmt.Errorf(
					"unterminated string at %d:%d",
					line,
					column,
				)
			}

			if r == '"' {
				l.advance()
				break
			}

			if r == '\\' {
				l.advance()
				switch l.peek() {
				case '"':
					value = append(value, '"')
					l.advance()
				case '\\':
					value = append(value, '\\')
					l.advance()
				case 'n':
					value = append(value, '\n')
					l.advance()
				case 't':
					value = append(value, '\t')
					l.advance()
				default:
					value = append(value, '\\')
				}
				continue
			}

			value = append(value, r)
			l.advance()
		}

		return Token{
			Type:   TokenString,
			Value:  string(value),
			Line:   line,
			Column: column,
		}, nil
	}

	var value []rune

	depth := 0
	for {
		r := l.peek()
		if r == 0 {
			break
		}

		if depth == 0 {
			if unicode.IsSpace(r) || r == '{' || r == '}' || r == ';' {
				break
			}
			// Only ')' terminates a word if it isn't part of a balanced (...)
			if r == ')' {
				break
			}
		}

		if r == '(' {
			depth++
		} else if r == ')' {
			depth--
		}

		value = append(value, r)
		l.advance()
	}

	return Token{
		Type:   TokenWord,
		Value:  string(value),
		Line:   line,
		Column: column,
	}, nil
}
