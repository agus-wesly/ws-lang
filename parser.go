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
		stmt, err := p.parseDeclaration()
		if err != nil {
			arr = append(arr, nil)
			p.synchronize()
			continue
		}
		arr = append(arr, stmt)
	}
	return arr, nil
}

func (p *Parser) parseDeclaration() (Statement, error) {
	if p.match(LET) {
		return p.parseVarDeclaration()
	}
	if p.match(FUN) {
		return p.parseFunctionDeclaration()
	}
	return p.parseStatement()
}

func (p *Parser) parseStatement() (Statement, error) {
	if p.match(PRINT) {
		parsed, err := p.parsePrint()
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}
	if p.match(WHILE) {
		return p.parseWhile()
	}
	if p.match(LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return CreateBlock(statements), nil
	}
	if p.match(IF) {
		return p.parseIf()
	}
	if p.match(FOR) {
		return p.parseFor()
	}
	if p.match(BREAK) {
		return p.parseBreak()
	}

	parsed, err := p.parseExpressionStatement()
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (p *Parser) parseIf() (Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	ifStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var elseStmt Statement = nil
	if p.match(ELSE) {
		elseStmt, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return CreateIfStatement(expr, ifStmt, elseStmt), nil
}

// while (expr) stmt
func (p *Parser) parseWhile() (Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return CreateWhileStatement(expr, stmt), nil
}

func (p *Parser) parseFor() (Statement, error) {
	_, err := p.consume(LEFT_PAREN, "Expected left parentheses ')' after for")
	if err != nil {
		return nil, err
	}

	var declr Statement
	if p.match(SEMICOLON) {
		declr = nil
	} else if p.check(LET) {
		declr, err = p.parseDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		declr, err = p.parseExpressionStatement()
	}

	var condition Expression
	if p.check(SEMICOLON) {
		condition = nil
	} else {
		condition, err = p.parseExpression()
	}
	p.consume(SEMICOLON, "Expected semicolon ';'")

	var incrementer Expression
	if p.check(RIGHT_PAREN) {
		incrementer = nil
	} else {
		stmt, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		incrementer = stmt
	}
	_, err = p.consume(RIGHT_PAREN, "Expected closing parentheses ')'")
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// Construct
	arrs := []Statement{body}
	if incrementer != nil {
		arrs = append(arrs, CreateExpressionStatement(incrementer))
	}
	body = CreateBlock(arrs)

	var expr Expression = nil
	if condition == nil {
		condition = CreateLiteral(true)
	}
	expr = condition

	res := []Statement{CreateWhileStatement(expr, body)}
	if declr != nil {
		res = append([]Statement{declr}, res...)
	}

	return CreateBlock(res), nil
}

func (p *Parser) block() ([]Statement, error) {
	statements := make([]Statement, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := p.consume(RIGHT_BRACE, "Expected closing bracket '}'")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) parseVarDeclaration() (Statement, error) {
	identifier, err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}
	var initValue Expression = CreateLiteral(Nil{})
	if p.match(EQUAL) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		initValue = expr
	}

	if _, err := p.consume(SEMICOLON, "Missing semicolon ; at the end of statement"); err != nil {
		return nil, err
	}
	return CreateVarDeclaration(initValue, identifier), nil
}

func (p *Parser) parseFunctionDeclaration() (Statement, error) {
	identifier, err := p.consume(IDENTIFIER, "Expected identifier after function declaration")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, "Expected opening parentheses '(' in function declaration")
	if err != nil {
		return nil, err
	}

	params := []*Token{}
	if p.match(IDENTIFIER) {
		params = append(params, p.previous())

		for p.match(COMMA) {
			param, err := p.consume(IDENTIFIER, "Expected identifier after ','")
			if err != nil {
				return nil, err
			}
			params = append(params, param)
		}
	}

	if len(params) >= 255 {
		return nil, CreateRuntimeError(identifier, "Can't have more than 255 args")
	}

	_, err = p.consume(RIGHT_PAREN, "Expected closing parentheses ')'")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expected opening braces '{' in function declaration")
	if err != nil {
		return nil, err
	}

	stmts, err := p.block()
	if err != nil {
		return nil, err
	}

	return CreateFunctionDeclaration(identifier, params, stmts), nil
}

func (p *Parser) parsePrint() (Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMICOLON, "Expected ; after expression")
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

	_, err = p.consume(SEMICOLON, "Expected ; after expression: "+p.peek().Lexeme)
	if err != nil {
		return nil, err
	}
	return CreateExpressionStatement(expr), nil
}

func (p *Parser) parseBreak() (Statement, error) {
	breakStmt := CreateBreakStatement()
	_, err := p.consume(SEMICOLON, "Expected semicolon")
	if err != nil {
		return nil, err
	}
	return breakStmt, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAssignment()
}

// a = 4;
func (p *Parser) parseAssignment() (Expression, error) {
	expr, err := p.parseComma()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		// Check if expr is var
		exprVar, ok := expr.(*Identifier)
		if !ok {
			// If not then it must return error
			return nil, CreateRuntimeError(p.peek(), "Invalid identifier")
		}
		value, err := p.parseComma()
		if err != nil {
			return nil, err
		}
		return CreateVarAssignment(exprVar.name, value), nil
	}
	return expr, nil
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
	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.match(QUESTION_MARK) {
		left, err := p.parseTernary()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(COLON, "Expected : inside ternary operator")
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

func (p *Parser) parseOr() (Expression, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	if p.match(OR) {
		name := p.previous()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		return CreateLogicalOperator(left, right, name), nil
	}
	return left, nil
}

func (p *Parser) parseAnd() (Expression, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	if p.match(AND) {
		name := p.previous()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		return CreateLogicalOperator(left, right, name), nil
	}
	return left, nil
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
		return CreateUnary(right, operand), nil
	}
	return p.parseFunction()
}

func (p *Parser) parseFunction() (Expression, error) {
	identifier, err := p.parsePrimary()
	token := p.previous()
	if err != nil {
		return nil, err
	}
	for p.match(LEFT_PAREN) {
		args := []Expression{}
		if !p.check(RIGHT_PAREN) {
			expr, err := p.parseTernary()
			if err != nil {
				return nil, err
			}

			args = append(args, expr)
			for p.match(COMMA) {
				expr, err = p.parseTernary()
				if err != nil {
					return nil, err
				}
				args = append(args, expr)
			}
		}
		_, err = p.consume(RIGHT_PAREN, "Expected closing parentheses ')'")
		if err != nil {
			return nil, err
		}
		identifier = CreateFunctionExpression(identifier, &args, token)
	}

	return identifier, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	if p.match(STRING, NUMBER, TRUE, FALSE, NIL, CHAR) {
		cur := p.previous()
		return CreateLiteral(cur.Literal), nil
	} else if p.match(IDENTIFIER) {
		cur := p.previous()
		return CreateIdentifier(cur), nil
	} else {
		if p.match(LEFT_PAREN) {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			_, err = p.consume(RIGHT_PAREN, "Expected )")
			if err != nil {
				return nil, err
			} else {
				return CreateGroup(expr), nil
			}
		}
	}
	return nil, CreateRuntimeError(p.peek(), "Unknown symbol"+p.peek().Lexeme)
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return &p.Tokens[p.Current]
}

func (p *Parser) previous() *Token {
	return &p.Tokens[p.Current-1]
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

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.Current += 1
	}
	return p.previous()

}

func (p *Parser) consume(tokenType TokenType, msg string) (*Token, error) {
	if p.peek().Type == tokenType {
		p.advance()
		return p.previous(), nil
	}
	return &Token{}, CreateRuntimeError(p.peek(), msg)
}

func (p *Parser) error() error {
	p.Lox.error(*p.peek(), "Syntax error")
	return errors.New("Syntax error")
}

func (p *Parser) synchronize() {
	// WARNING : This might cause an error
	for !p.isAtEnd() {
		switch p.peek().Type {
		case FOR:
			fallthrough
		case WHILE:
			fallthrough
		case SEMICOLON:
            p.advance()
			fallthrough
		case IF:
			fallthrough
		case LET:
			fallthrough
		case PRINT:
			fallthrough
		case RETURN:
			fallthrough
		case CLASS:
			fallthrough
		case FUN:
			return
		}
		p.advance()
	}
}
