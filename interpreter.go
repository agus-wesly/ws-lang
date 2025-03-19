package main

import (
	"errors"
	"fmt"
	"strconv"
)

type ExpressionVisitor interface {
	VisitLiteral(l *Literal) any
	VisitIdentifier(i *Identifier) (any, error)
	VisitUnary(u *Unary) (any, error)
	VisitBinary(b *Binary) (any, error)
	VisitTernary(t *Ternary) (any, error)
	VisitLogicalOperator(l *LogicalOperator) (any, error)
	VisitGrouping(g *Grouping) (any, error)
	VisitVarAssignment(v *VarAssignment) (any, error)
	VisitFunction(f *Function) (any, error)
}

type Interpreter struct {
	*Environment
	Locals map[Expression]int
}

func CreateAndSetupInterpreter() *Interpreter {
	interpreter := &Interpreter{
		Environment: CreateEnvironment(nil, map[string]any{}, nil),
		Locals:      make(map[Expression]int),
	}
	interpreter.Environment.Interpreter = interpreter
	SetupInterpreter(interpreter)

	return interpreter
}

func (i *Interpreter) interpret(statements []Statement, replMode bool) {
	for _, stmt := range statements {
		if stmt != nil {
			res, err := stmt.accept(i)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if replMode && res != nil {
				fmt.Println(res)
			}
		}
	}
}

func (i *Interpreter) VisitVarAssignment(v *VarAssignment) (any, error) {
	name := *(&v.Token.Lexeme)
	_, err := i.Get(name)
	if err != nil {
		return nil, CreateRuntimeError(v.Token, "Unknown variable: "+name)
	}

	val, err := v.Expr.accept(i)
	if err != nil {
		return nil, err
	}
	err = i.Assign(name, val)
	if err != nil {
		return nil, err
	}

	return (i.Values[name]), nil
}

func (i *Interpreter) VisitFunction(f *Function) (any, error) {
	_calle, err := f.Identifier.accept(i)
	if err != nil {
		return nil, err
	}

	calle, ok := _calle.(Callee)
	if !ok {
		return nil, CreateRuntimeError(f.Token, "Identifier `"+f.Token.Lexeme+"` is not a function")
	}
	val, err := calle.call(i, f.Token, f.Args)
	if err != nil {
		return nil, err
	}

	// Todo : when _calle is another function, try to comes up with correct return type
	return val, nil
}

func (i *Interpreter) VisitFunctionDeclaration(f *FunctionDeclaration) (any, error) {
	f.Env = i.Environment
	// Save into current scope
	name := f.Identifier.Lexeme
	i.Environment.Set(name, f)
	return nil, nil
}

func (i *Interpreter) VisitReturnStatement(r *ReturnStatement) error {
	// TODO : fix this so that it will evaluate in HERE not in the function call
	return r
}

func (i *Interpreter) VisitVarDeclaration(v *VarDeclaration) (any, error) {
	name := v.Identifier.Lexeme
	var value any = v.Expr
	if value != nil {
		exprValue, err := v.Expr.accept(i)
		if err != nil {
			fmt.Println("Unreachable in var declr")
			return nil, err
		}
		value = exprValue
	}
	_, err := i.Environment.GetCurrentBlock(name)
	if err == nil {
		// Variable is redeclared
		return nil, CreateRuntimeError(v.Identifier, "Redeclaration of name "+name)
	}
	i.Environment.Set(name, value)
	return value, nil
}

// TODO : this should return value, not just error
func (i *Interpreter) VisitBlockStatement(b *BlockStatement) error {
	prevEnv := i.Environment
	defer func() {
		i.Environment = prevEnv
	}()

	newEnv := CreateEnvironment(prevEnv, make(map[string]any), i)

	i.Environment = newEnv
	for _, stmt := range b.Statements {
		_, err := stmt.accept(i)
		if err != nil {
			// i.Environment = prevEnv
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitIfStatement(ifs *IfStatement) error {
	stmt, err := ifs.Expr.accept(i)
	if err != nil {
		return err
	}
	if i.isTruthy(stmt) {
		_, err := ifs.IfStmt.accept(i)
		if err != nil {
			return err
		}
	} else if ifs.ElseStmt != nil {
		_, err := ifs.ElseStmt.accept(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitWhileStatement(w *WhileStatement) error {
	for {
		val, err := w.Expr.accept(i)
		if err != nil {
			return err
		}

		if !i.isTruthy(val) {
			break
		}

		_, err = w.Stmt.accept(i)
		if err != nil {
			if err == BreakStmtErr {
				break
			}
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitPrintStatement(p *PrintStatement) error {
	expr, err := i.evaluate(p.Expr)
	if err != nil {
		return err
	}
	// TODO : need to have `.toString()` method to be called
	fmt.Println(expr)
	return nil
}

func (i *Interpreter) VisitExpressionStatement(p *ExpressionStatement) (any, error) {
	val, err := i.evaluate(p.Expr)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) VisitGrouping(g *Grouping) (any, error) {
	return i.evaluate(g.Expression)
}

func (i *Interpreter) VisitLogicalOperator(l *LogicalOperator) (any, error) {
	left, err := l.Left.accept(i)
	if err != nil {
		return nil, err
	}

	if l.Name.Type == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}
	right, err := l.Right.accept(i)
	if err != nil {
		return nil, err
	}

	return right, nil
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
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil

	case PLUS:
		if err := i.checkExprNumber(b.Operator, left, right); err == nil {
			return left.(float64) + right.(float64), nil
		}
		if err := i.checkExprString(b.Operator, left, right); err == nil {
			return left.(string) + right.(string), nil
		}
		if err := i.checkExprString(b.Operator, left); err == nil {
			if err := i.checkExprNumber(b.Operator, right); err == nil {
				return left.(string) + fmt.Sprint(right.(float64)), nil
			}
		}
		if err := i.checkExprNumber(b.Operator, left); err == nil {
			if err := i.checkExprString(b.Operator, right); err == nil {
				return fmt.Sprint(left.(float64)) + right.(string), nil
			}
		}
		return nil, CreateRuntimeError(b.Operator, "Addition not supported")

	case SLASH:
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		// Divide by 0
		if right == 0.0 {
			return nil, CreateRuntimeError(b.Operator, "Cannot divide by 0")
		}
		return left.(float64) / right.(float64), nil

	case STAR:
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	case GREATER:
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case GREATER_EQUAL:
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		r := (left.(float64)) >= (right.(float64))
		return r, nil

	case LESS:
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil

	case LESS_EQUAL:
		if err := i.checkExprNumber(b.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil

	case BANG_EQUAL:
		return (left != right), nil

	case EQUAL_EQUAL:
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

func (i *Interpreter) VisitIdentifier(identifier *Identifier) (any, error) {
	val, err := i.Environment.Get(identifier.name.Lexeme)
	if err != nil {
		return nil, err
	}
	_, isUninitialized := val.(Nil)
	if isUninitialized {
		return nil, CreateRuntimeError(identifier.name, "Variable must be initialized before used")
	}
	return val, nil
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
	return errors.New(fmt.Sprintf("[line %d] Runtime Error : %s\n", token.Line, msg))
}
