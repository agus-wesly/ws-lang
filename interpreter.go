package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Visitor interface {
	VisitLiteral(l *Literal) any
	VisitUnary(u *Unary) (any, error)
	VisitBinary(b *Binary) (any, error)
	VisitTernary(t *Ternary) (any, error)
	VisitGrouping(g *Grouping) (any, error)
}

type Interpreter struct{}

func (i *Interpreter) VisitGrouping(g *Grouping) (any, error) {
	return i.evaluate(g.Expression)
}

func (i *Interpreter) VisitTernary(t *Ternary) (any, error) {
	left, err := i.evaluate(t.Left)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(left) {
		return i.evaluate(t.Center)
	} else {
		return i.evaluate(t.Right)
	}
}

func (i *Interpreter) VisitBinary(b *Binary) (any, error) {
	left, err := i.evaluate(b.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(b.Right)
	if err != nil {
		return nil, err
	}

	switch b.Operator.Type {
	case MINUS:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil

	case PLUS:
		if err := i.checkExprNumber(&b.Operator, left, right); err == nil {
			return left.(float64) + right.(float64), nil
		}
		if err := i.checkExprString(&b.Operator, left, right); err == nil {
			return left.(string) + right.(string), nil
		}
		if err := i.checkExprString(&b.Operator, left); err == nil {
			if err := i.checkExprNumber(&b.Operator, right); err == nil {
				return left.(string) + fmt.Sprint(right.(float64)), nil
			}
		}
		if err := i.checkExprNumber(&b.Operator, left); err == nil {
			if err := i.checkExprString(&b.Operator, right); err == nil {
				return fmt.Sprint(left.(float64)) + right.(string), nil
			}
		}
		return nil, CreateRuntimeError(&b.Operator, "Addition not supported")

	case SLASH:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil

	case STAR:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	case GREATER:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case GREATER_EQUAL:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		r := (left.(float64)) >= (right.(float64))
		return r, nil

	case LESS:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil

	case LESS_EQUAL:
		if err := i.checkExprNumber(&b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil

	case BANG_EQUAL:
		return (left != right), nil

	case EQUAL_EQUAL:
		fmt.Println(left, right)
		return (left == right), nil
	}
	return nil, errors.New("Unreachable")
}

// -1
func (i *Interpreter) VisitUnary(u *Unary) (any, error) {
	val, err := i.evaluate(u.Right)
	if err != nil {
		return nil, err
	}
	switch u.Operand.Type {
	case MINUS:
		val, ok := val.(float64)
		if !ok {
			return nil, CreateRuntimeError(u.Operand, "Conversion error")
		}
		return -(val), nil
	case BANG:
		return !i.isTruthy(val), nil
	}
	panic("Unreachable")
}

func (i *Interpreter) VisitLiteral(l *Literal) any {
	return l.Value
}

func (i *Interpreter) evaluate(exp Expression) (any, error) {
	return exp.accept(i)
}

func (i *Interpreter) isTruthy(exp any) bool {
	// TODO: handle nil
	if exp == nil || exp == 0 {
		return false
	}
	switch exp.(type) {
	case bool:
		return exp.(bool)
	}
	return true
}

func (i *Interpreter) isString(exp any) bool {
	_, ok := exp.(string)
	return ok
}

func (i *Interpreter) isNumber(exp any) bool {
	_, ok := exp.(float64)
	return ok
}

func parseFloat(x any) float64 {
	res, err := strconv.ParseFloat(x.(string), 32)
	if err != nil {
		panic(err)
	}
	return res
}

func (i *Interpreter) checkExprNumber(tok *Token, expressions ...any) error {
	for _, expr := range expressions {
		_, ok := expr.(float64)
		if !ok {
			return CreateRuntimeError(tok, "Parsing error")
		}
	}
	return nil
}

func (i *Interpreter) checkExprString(tok *Token, expressions ...any) error {
	for _, expr := range expressions {
		_, ok := expr.(string)
		if !ok {
			return CreateRuntimeError(tok, "Parsing error")
		}
	}
	return nil
}

func CreateRuntimeError(token *Token, msg string) error {
	return errors.New(fmt.Sprintf("[line %d] Compile Error : %s\n", token.Line, msg))
}
