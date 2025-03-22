package main

import "fmt"

type Callee interface {
	arity() int
	call(i *Interpreter, token *Token, args *[]Expression) (any, error)
	toString() string
}

func (f *FunctionDeclaration) call(interpreter *Interpreter, token *Token, args *[]Expression) (any, error) {
	argsVal := make([]any, 0)
	for _, arg := range *args {
		val, err := arg.accept(interpreter)
		if err != nil {
			return nil, err
		}
		argsVal = append(argsVal, val)
	}

	prevEnv := interpreter.Environment
	defer func() {
		interpreter.Environment = prevEnv
	}()

	newEnv := CreateEnvironment(prevEnv, make(map[string]any), interpreter)
	interpreter.Environment = newEnv

	if f.arity() != len(*args) {
		return nil, CreateRuntimeError(token, fmt.Sprintf("Expected %d arguments but got %d .", f.arity(), len(*args)))
	}

	for i, param := range f.Params {
		interpreter.Environment.Set(param.Lexeme, argsVal[i])
	}

	for _, stmt := range f.Stmts {
		val, err := stmt.accept(interpreter)
		if err != nil {
			if _, ok := err.(*ReturnStatement); ok {
				// Encountered return keyword
				return val, nil
			}

			return nil, err
		}
	}

	return nil, nil
}

func (f *FunctionDeclaration) arity() int {
	return len(f.Params)
}

func (f *FunctionDeclaration) toString() string {
	return "<" + f.Identifier.Lexeme + ">"
}
