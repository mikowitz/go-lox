package lox

import (
	"fmt"
	"strconv"
)

// Scanner performs lexical analysis on Lox source code.
type Scanner struct {
	source []rune  // The source code as runes for Unicode support
	tokens []Token // The scanned tokens
	lox    *Lox    // Reference to the interpreter for error reporting

	start, current, line int // Position tracking in the source
}

// NewScanner creates a new Scanner for the given source code.
func NewScanner(source string, lox *Lox) Scanner {
	return Scanner{
		source: []rune(source),
		lox:    lox,
		line:   1,
	}
}

// ScanTokens scans the entire source code and produces a list of tokens.
// The returned slice includes all scanned tokens plus a final EOF token.
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if isDigit(r) {
			s.scanNumber()
		} else if isAlpha(r) {
			s.scanIdentifier()
		} else {
			s.lox.error(s.line, fmt.Sprintf("Unexpected character %q", r))
		}
	}
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lox.error(s.line, "Unterminated string")
		return
	}

	s.advance() // The closing "

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(String, string(value))
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	floatStr := string(s.source[s.start:s.current])
	n, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		// NOTE: We should never get here because of the scanner's pre-validation
		s.lox.error(s.line, fmt.Sprintf("Unparseable float: %s", floatStr))
		return
	}
	s.addTokenWithLiteral(Number, n)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	identifier := string(s.source[s.start:s.current])
	tokenType := Identifier
	if tType, ok := keywords[identifier]; ok {
		tokenType = tType
	}
	s.addToken(tokenType)
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	var r rune // Zero value is '\0'
	if s.isAtEnd() {
		return r
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	var r rune // Zero value is '\0'
	if s.current+1 >= len(s.source) {
		return r
	}
	return s.source[s.current+1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, string(lexeme), literal, s.line))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// isDigit returns true if r is a decimal digit (0-9).
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// isAlpha returns true if r is an alphabetic character or underscore.
// Valid identifier start characters.
func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

// isAlphaNumeric returns true if r is alphanumeric or underscore.
// Valid identifier characters.
func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
