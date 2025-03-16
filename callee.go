package main

import "fmt"

type Callee interface {
	arity() int
	call(i *Interpreter, token *Token, args *[]Expression) (any, error)
	toString() string
}

func (f *FunctionDeclaration) call(interpreter *Interpreter, token *Token, args *[]Expression) (any, error) {
	prevEnv := interpreter.Environment
	defer func() {
		interpreter.Environment = prevEnv
	}()

	newEnv := CreateEnvironment(*prevEnv, make(map[string]any))
	interpreter.Environment = newEnv

	if f.arity() != len(*args) {
		return nil, CreateRuntimeError(token, fmt.Sprintf("Expected %d arguments but got %d .", f.arity(), len(*args)))
	}

	for i, param := range f.Params {
		val, err := (*args)[i].accept(interpreter)
		if err != nil {
			return nil, err
		}
		interpreter.Environment.Set(param.Lexeme, val)
	}

	for _, stmt := range f.Stmts {
		_, err := stmt.accept(interpreter)
		if err != nil {
			if returnStmt, ok := err.(*ReturnStatement); ok {
				// Encountered return keyword
				return f.handleReturn(interpreter, returnStmt)
			}

			return nil, err
		}
	}

	return nil, nil
}

func (f *FunctionDeclaration) arity() int {
	return len(f.Params)
}

func (f *FunctionDeclaration) handleReturn(interpreter *Interpreter, ret *ReturnStatement) (any, error) {
	if ret.Expr != nil {
		val, err := ret.Expr.accept(interpreter)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, nil
}

func (f *FunctionDeclaration) toString() string {
	return "<" + f.Identifier.Lexeme + ">"
}
