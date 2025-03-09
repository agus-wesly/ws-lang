package main

import (
	"errors"
)

type Parser struct {
	*Lox
	Current int
	Tokens  []Token
}

func CreateParser(tokens []Token, lox *Lox) *Parser {
	return &Parser{
		Current: 0,
		Tokens:  tokens,
		Lox:     lox,
	}
}

func (p *Parser) parse() ([]Statement, error) {
	arr := []Statement{}
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		arr = append(arr, stmt)
	}
	return arr, nil
}

func (p *Parser) statement() (Statement, error) {
	if p.match(PRINT) {
		parsed, err := p.parsePrint()
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}
	parsed, err := p.parseExpressionStatement()
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (p *Parser) parsePrint() (Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

    err = p.consume(SEMICOLON, "Expected ; after expression")
	if err != nil {
		return nil, err
	}
	return CreatePrintStatement(expr), nil
}
func (p *Parser) parseExpressionStatement() (Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	err = p.consume(SEMICOLON,  "Expected ; after expression")
	if err != nil {
		return nil, err
	}
	return CreateExpressionStatement(expr), nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseComma()
}

// Comma operator evaluates left side, discards it and then
// evaluates and return right side.
func (p *Parser) parseComma() (Expression, error) {
	left, err := p.parseTernary()
	if err != nil {
		return nil, err
	}

	for p.match(COMMA) {
		p.previous()
		right, err := p.parseTernary()
		if err != nil {
			return nil, err
		}
		left = right
	}
	return left, nil
}

func (p *Parser) parseTernary() (Expression, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return nil, err
	}
	if p.match(QUESTION_MARK) {
		left, err := p.parseTernary()
		if err != nil {
			return nil, err
		}

		err = p.consume(COLON,  "Expected : inside ternary operator")
		if err != nil {
			return nil, err
		}

		right, err := p.parseTernary()
		if err != nil {
			return nil, err
		}

		expr = CreateTernary(expr, left, right)
	}
	return expr, nil
}

// (!=, ==)
func (p *Parser) parseEquality() (Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := p.previous()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = CreateBinary(left, operator, right)
	}
	return left, nil
}

// (>, <, <=, >=)
func (p *Parser) parseComparison() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = CreateBinary(left, operator, right)
	}
	return left, nil
}

// (+, -)
func (p *Parser) parseTerm() (Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = CreateBinary(left, operator, right)
	}
	return left, nil
}

// (/, *)
func (p *Parser) parseFactor() (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = CreateBinary(left, operator, right)
	}
	return left, nil
}

// (-1)
func (p *Parser) parseUnary() (Expression, error) {
	if p.match(BANG, MINUS) {
		operand := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return CreateUnary(right, &operand), nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expression, error) {
	if p.match(STRING, NUMBER, TRUE, FALSE, NIL, CHAR) {
		cur := p.previous()
		return CreateLiteral(cur.Literal), nil
	} else {
		if p.match(LEFT_PAREN) {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			err = p.consume(RIGHT_PAREN, "Expected )")
			if err != nil {
				return nil, err
			} else {
				return CreateGroup(expr), nil
			}
		}
	}
	return nil, CreateRuntimeError(p.peek(), "Unknown symbol")
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return &p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(expr TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.Tokens[p.Current].Type == expr
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current += 1
	}
	return p.previous()

}

func (p *Parser) consume(tokenType TokenType, msg string) error {
	if p.peek().Type == tokenType {
		p.advance()
		return nil
	}
	return CreateRuntimeError(p.peek(), msg)
}

func (p *Parser) error() error {
	p.Lox.error(*p.peek(), "Syntax error")
	return errors.New("Syntax error")
}

func (p *Parser) synchronize() {
	// WARNING : This might cause an error
	for !p.isAtEnd() {
		switch p.peek().Type {
		case SEMICOLON:
		case FOR:
		case WHILE:
		case IF:
		case VAR:
		case PRINT:
		case RETURN:
		case CLASS:
		case FUN:
			return
		}
		p.advance()
	}
}
