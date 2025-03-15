package main

import "fmt"

type Callee interface {
	arity() int
	call(i *Interpreter, args *[]Expression) (any, error)
	toString() string
}

// TODO : maybe this should be Function struct and not FunctionDeclaration ?
func (f *FunctionDeclaration) call(interpreter *Interpreter, args *[]Expression) (any, error) {
	prevEnv := interpreter.Environment
	newEnv := CreateEnvironment(*prevEnv, make(map[string]any))
	interpreter.Environment = newEnv

	if f.arity() != len(*args) {
		return nil, CreateRuntimeError(f.Identifier, fmt.Sprintf("Expected %d arguments but got %d .", f.arity(), len(*args)))
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
			return nil, err
		}
	}

	interpreter.Environment = prevEnv
	return nil, nil
}

func (f *FunctionDeclaration) arity() int {
	return len(f.Params)
}

func (f *FunctionDeclaration) toString() string {
	return "<" + f.Identifier.Lexeme + ">"
}
