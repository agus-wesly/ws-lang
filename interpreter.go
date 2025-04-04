package main

import (
	"errors"
	"fmt"
	"strconv"
)

type ExpressionVisitor interface {
	VisitLiteral(l *Literal) any
	VisitIdentifier(i *IdentifierExpr) (any, error)
	VisitUnary(u *Unary) (any, error)
	VisitBinary(b *Binary) (any, error)
	VisitTernary(t *Ternary) (any, error)
	VisitLogicalOperator(l *LogicalOperator) (any, error)
	VisitGrouping(g *Grouping) (any, error)
	VisitVarAssignment(v *VarAssignment) (any, error)
	VisitFunction(f *Function) (any, error)
}

type Interpreter struct {
	Environment *Environment
	Locals      map[Expression]*Local
	Globals     *Environment
}

type Local struct {
	Distance int
	Index    int
}

func CreateAndSetupInterpreter() *Interpreter {
	globalInterpreter := CreateEnvironment(nil, nil)
	interpreter := &Interpreter{
		Environment: globalInterpreter,
		Locals:      make(map[Expression]*Local),
		Globals:     globalInterpreter,
	}
	globalInterpreter.Interpreter = interpreter
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
	newValue, err := v.Expr.accept(i)
	if err != nil {
		return nil, err
	}

	local, found := i.Locals[v]
	if found {
		i.Environment.AssignAt(local.Distance, *v.Token, newValue)
	} else {
		// search in global
		for _, val := range i.Globals.Identifiers {
			if val.Name == v.Token.Lexeme {
				val.Value = newValue
			}
		}
	}

	return newValue, nil
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
	name := f.Identifier.Lexeme
	_, err := i.Environment.GetCurrentBlock(name)
	if err == nil {
		// Function is Redeclarated
        // TODO : maybe we can make this compile time ?
		return nil, CreateRuntimeError(f.Identifier, "Redeclaration of name "+name)
	}

	i.Environment.Set(name, f)
	return nil, nil
}

func (i *Interpreter) VisitReturnStatement(r *ReturnStatement) (any, error) {
	if r.Expr == nil {
		return nil, r
	}
	val, err := r.Expr.accept(i)
	if err != nil {
		return nil, err
	}
	return val, r
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

func (i *Interpreter) VisitBlockStatement(b *BlockStatement) (any, error) {
	prevEnv := i.Environment
	defer func() {
		i.Environment = prevEnv
	}()

	newEnv := CreateEnvironment(prevEnv, i)

	i.Environment = newEnv
	for _, stmt := range b.Statements {
		val, err := stmt.accept(i)
		if err != nil {
			return val, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitIfStatement(ifs *IfStatement) (any, error) {
	stmt, err := ifs.Expr.accept(i)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(stmt) {
		val, err := ifs.IfStmt.accept(i)
		if err != nil {
			return val, err
		}
	} else if ifs.ElseStmt != nil {
		val, err := ifs.ElseStmt.accept(i)
		if err != nil {
			return val, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitWhileStatement(w *WhileStatement) (any, error) {
	for {
		val, err := w.Expr.accept(i)
		if err != nil {
			return val, err
		}

		if !i.isTruthy(val) {
			break
		}

		val, err = w.Stmt.accept(i)
		if err != nil {
			if err == BreakStmtErr {
				break
			}
			return val, err
		}
	}
	return nil, nil
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

func (i *Interpreter) VisitIdentifier(identifier *IdentifierExpr) (any, error) {
	val, err := i.Environment.lookUpVariable(identifier.name.Lexeme, identifier)
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
