package main

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
    Line    int
}

func CreateToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
        Line: line,
	}
}

func (t *Token) toString() string {
	return fmt.Sprint(t.Type, " ", t.Lexeme, " ", t.Literal)
}
