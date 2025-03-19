package main

import (
	"errors"
	"fmt"
)

type Environment struct {
	Values  map[string]any
	PrevEnv *Environment
	*Interpreter
}

func (env *Environment) Get(name string, expr Expression) (any, error) {
	distance, ok := env.Interpreter.Locals[expr]
	if !ok {
		res := env.findInGlobal(name)
		return res, nil
	}

	fmt.Println(env)
	values := env.GetValuesFromDistance(distance)
    fmt.Println(values, distance)
	res, found := values[name]
	if !found {
        fmt.Println("err : ", distance)
		panic("Unreachable")
	}

	return res, nil
}

func (env *Environment) GetValuesFromDistance(dist int) map[string]any {
	curr := env
	for i := 0; i < dist; i += 1 {
		curr = env.PrevEnv
	}
	return curr.Values
}

func (env *Environment) findInGlobal(name string) any {
	curr := env
	for curr.PrevEnv != nil {
		curr = curr.PrevEnv
	}
	val, ok := curr.Values[name]
	if !ok {
		panic("Not found in global")
	}
	return val

}

func (env *Environment) _Get(name string) (any, error) {
	if true {
		panic("Disabled")
	}
	res, found := env.Values[name]

	// Get the distance

	// Go through the target map based on the distance

	if !found {
		if env.PrevEnv == nil {
			return nil, errors.New(name + " is not defined")
		}
		val, err := env.PrevEnv.Get(name, nil)
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
