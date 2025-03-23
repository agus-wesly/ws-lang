package main

import (
	"errors"
)

type Environment struct {
	Identifiers      []*Identifier
	PrevEnv     *Environment
	Interpreter *Interpreter
}

type Identifier struct {
	Value any
	Name  string
}

func (env *Environment) lookUpVariable(name string, expr Expression) (any, error) {
	local, ok := env.Interpreter.Locals[expr]
	if !ok {
		return env.findInGlobal(name)
	}

	currEnv := env.GetAt(local.Distance)
	result := currEnv.Identifiers[local.Index]

	return result.Value, nil
}

func (env *Environment) GetAt(dist int) *Environment {
	curr := env
	for i := 0; i < dist; i += 1 {
		curr = curr.PrevEnv
	}
	return curr
}

func (env *Environment) findInGlobal(name string) (any, error) {
	for _, val := range env.Interpreter.Globals.Identifiers {
		if val.Name == name {
			return val.Value, nil
		}
	}
	return nil, errors.New("Not found in global")

}

func (env *Environment) findByName(name string) (*Identifier, bool) {
	for _, val := range env.Identifiers {
		if val.Name == name {
			return val, true
		}
	}
	return nil, false
}

func (env *Environment) GetCurrentBlock(name string) (any, error) {

	res, found := env.findByName(name)
	if !found {
		return nil, errors.New(name + " is not defined")
	}
	return res, nil
}

func (env *Environment) Set(name string, value any) {
	env.Identifiers = append(env.Identifiers, &Identifier{Value: value, Name: name})
}

func (env *Environment) AssignAt(distance int, token Token, value any) {
	targetEnv := env.GetAt(distance)

	val, found := targetEnv.findByName(token.Lexeme)
	if !found {
		panic("Unreachable")
	}

	val.Value = value
}

func CreateEnvironment(prevEnv *Environment, interpreter *Interpreter) *Environment {
	return &Environment{
		PrevEnv:     prevEnv,
		Identifiers:      make([]*Identifier, 0),
		Interpreter: interpreter,
	}
}
