package main

import (
	"errors"
)

type Environment struct {
	//Values  map[string]any
	Values      []*Value
	PrevEnv     *Environment
	Interpreter *Interpreter
}

type Value struct {
	Value any
	Name  string
}

func (env *Environment) lookUpVariable(name string, expr Expression) (any, error) {
	// TODO : selain distance harus juga mereturn index pada array
	local, ok := env.Interpreter.Locals[expr]
	if !ok {
		return env.findInGlobal(name)
	}

	currEnv := env.GetAt(local.Distance)
	result := currEnv.Values[local.Index]

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
	for _, val := range env.Interpreter.Globals.Values {
		if val.Name == name {
			return val.Value, nil
		}
	}
	return nil, errors.New("Not found in global")

}

// func (env *Environment) Get(name string) (any, error) {
// 	res, found := env.Values[name]
//
// 	if !found {
// 		if env.PrevEnv == nil {
// 			return nil, errors.New(name + " is not defined")
// 		}
// 		val, err := env.PrevEnv.lookUpVariable(name, nil)
// 		if err != nil {
// 			return nil, errors.New(name + " is not defined.")
// 		}
// 		res = val
// 	}
// 	return res, nil
// }

func (env *Environment) findByName(name string) (*Value, bool) {
	for _, val := range env.Values {
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
	env.Values = append(env.Values, &Value{Value: value, Name: name})
	// env.Values[name] = value
}

// func (env *Environment) Assign(name string, value any) error {
// 	_, found := env.Values[name]
// 	if found {
// 		env.Values[name] = value
// 	} else {
// 		if env.PrevEnv == nil {
// 			return errors.New(name + " is not defined")
// 		}
// 		err := env.PrevEnv.Assign(name, value)
// 		if err != nil {
// 			return errors.New(name + " is not defined.")
// 		}
// 	}
// 	return nil
// }

func (env *Environment) AssignAt(distance int, token Token, value any) {
	targetEnv := env.GetAt(distance)

	val, found := targetEnv.findByName(token.Lexeme)
	// _, found := targetEnv.Values[token.Lexeme]
	if !found {
		panic("Unreachable")
	}

	val.Value = value
}

func CreateEnvironment(prevEnv *Environment, interpreter *Interpreter) *Environment {
	return &Environment{
		PrevEnv:     prevEnv,
		Values:      make([]*Value, 0),
		Interpreter: interpreter,
	}
}
