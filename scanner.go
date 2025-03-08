package main

import (
	"strconv"
)

type Scanner struct {
	*Lox

	Source string
	Tokens []Token

	start     int
	current   int
	lineCount int
}

func CreateScanner(src string, lox *Lox) *Scanner {
	return &Scanner{
		Source: src,
		Tokens: make([]Token, 0),
		Lox:    lox,

		start:     0,
		current:   0,
		lineCount: 0,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.processChar()
	}
	s.Tokens = append(s.Tokens, Token{Type: EOF, Literal: nil, Lexeme: "EOF", Line: s.lineCount})
	return s.Tokens
}

func (s *Scanner) advance() rune {
	cur := rune(s.Source[s.current])
	s.current += 1
	return cur
}

func (s *Scanner) processChar() {
	ch := s.advance()

	switch ch {
	case '(':
		s.addToken(LEFT_PAREN)
		break
	case ')':
		s.addToken(RIGHT_PAREN)
		break
	case '{':
		s.addToken(LEFT_BRACE)
		break
	case '}':
		s.addToken(RIGHT_BRACE)
	case '?':
		s.addToken(QUESTION_MARK)
		break
	case ':':
		s.addToken(COLON)
		break
	case '.':
		s.addToken(DOT)
		break
	case '+':
		s.addToken(PLUS)
		break
	case '-':
		s.addToken(MINUS)
	case ';':
		s.addToken(SEMICOLON)
		break
	case '*':
		s.addToken(STAR)
		break
	case ',':
		s.addToken(COMMA)
		break
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
		break
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
		break
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.current += 1
			}
			break
		} else if s.match('*') {
			s.multilineComment()
			break
		} else {
			s.addToken(SLASH)
		}
		break
	case '"':
		s.string()
		break
	case '\n':
		s.lineCount += 1
		break
	case ' ':
		break
	case '\t':
		break
	default:
		if s.isNumber(ch) {
			s.number()
			break
		}
		if s.isAlpha(ch) {
			s.alpha()
			break
		}
		// TODO : handle unknown input
		s.Lox.error(Token{Type: EQUAL, Line: s.lineCount, Lexeme: string(ch)}, "Invalid token")
		break
	}
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.Source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.Source) {
		return rune(0)
	}
	return rune(s.Source[s.current+1])
}

func (s *Scanner) match(char rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.Source[s.current]) != char {
		return false
	}
	s.current += 1
	return true
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.lineCount += 1
		}
		s.advance()
	}
	if s.peek() != '"' {
		s.Lox.error(Token{Type: STRING, Line: s.lineCount, Lexeme: string(s.peek())}, "Invalid string")
		return
	}
	s.advance()
	s.addTokenLiteral(STRING, s.Source[s.start+1:s.current-1])
}
func (s *Scanner) multilineComment() {
	for !s.isAtEnd() {
		if s.peek() == '*' && s.peekNext() == '/' {
			s.advance()
			s.advance()
			break
		}
		if s.peek() == '\n' {
			s.lineCount += 1
		}
		s.advance()
	}
}

func (s *Scanner) isNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char == '_')
}
func (s *Scanner) isAlphaNumeric(char rune) bool {
	return s.isAlpha(char) || s.isNumber(char)
}

// .123
func (s *Scanner) isValidDigit() bool {
	return s.isNumber(s.peek())
}
func (s *Scanner) number() {
	for s.isNumber(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isNumber(s.peekNext()) {
		s.advance()
		for s.isNumber(s.peek()) {
			s.advance()
		}
	}
	parsed, err := strconv.ParseFloat(s.Source[s.start:s.current], 64)
	if err != nil {
		panic("UNREACHABLE")
	}
	s.addTokenLiteral(NUMBER, parsed)
}

func (s *Scanner) alpha() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	var tokenType TokenType
	text := s.Source[s.start:s.current]
	reservedType, found := keywords[text]
	if !found {
		tokenType = IDENTIFIER
	} else {
		tokenType = reservedType
	}
	s.addTokenLiteral(tokenType, s.Source[s.start:s.current])
}

// x=1+1\n
func (s *Scanner) addToken(tokenType TokenType) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, Token{Type: tokenType, Literal: nil, Lexeme: text, Line: s.lineCount})
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
    // TODO : check literal. ex : nil then type is nil
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, Token{Type: tokenType, Literal: literal, Lexeme: text, Line: s.lineCount})
}
