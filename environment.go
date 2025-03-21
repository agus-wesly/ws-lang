package main

import (
	"errors"
)

type Environment struct {
	Values  map[string]any
	PrevEnv *Environment
	*Interpreter
}

func (env *Environment) lookUpVariable(name string, expr Expression) (any, error) {
	distance, ok := env.Interpreter.Locals[expr]
	if !ok {
		return env.findInGlobal(name)
	}

	currEnv := env.GetAt(distance)
	result, found := currEnv.Values[name]
	if !found {
		panic("Unreachable")
	}

	return result, nil
}

func (env *Environment) GetAt(dist int) *Environment {
	curr := env
	for i := 0; i < dist; i += 1 {
		curr = env.PrevEnv
	}
	return curr
}

func (env *Environment) findInGlobal(name string) (any, error) {
	curr := env
	for curr.PrevEnv != nil {
		curr = curr.PrevEnv
	}
	val, ok := curr.Values[name]
	if !ok {
        return nil, errors.New("Not found in global")
	}
	return val, nil

}

func (env *Environment) Get(name string) (any, error) {
	res, found := env.Values[name]

	if !found {
		if env.PrevEnv == nil {
			return nil, errors.New(name + " is not defined")
		}
		val, err := env.PrevEnv.lookUpVariable(name, nil)
		if err != nil {
			return nil, errors.New(name + " is not defined.")
		}
		res = val
	}
	return res, nil
}

func (env *Environment) GetCurrentBlock(name string) (any, error) {
	res, found := env.Values[name]
	if !found {
		return nil, errors.New(name + " is not defined")
	}
	return res, nil
}

func (env *Environment) Set(name string, value any) {
	env.Values[name] = value
}

func (env *Environment) Assign(name string, value any) error {
	_, found := env.Values[name]
	if found {
		env.Values[name] = value
	} else {
		if env.PrevEnv == nil {
			return errors.New(name + " is not defined")
		}
		err := env.PrevEnv.Assign(name, value)
		if err != nil {
			return errors.New(name + " is not defined.")
		}
	}
	return nil
}

func CreateEnvironment(prevEnv *Environment, values map[string]any, interpreter *Interpreter) *Environment {
	return &Environment{
		PrevEnv:     prevEnv,
		Values:      values,
		Interpreter: interpreter,
	}
}
