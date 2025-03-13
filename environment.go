package main

import "errors"

type Environment struct {
	Values  map[string]any
	PrevEnv *Environment
}

func (env *Environment) Get(name string) (any, error) {
	res, found := env.Values[name]
	if !found {
		if env.PrevEnv == nil {
			return nil, errors.New(name + " is not defined")
		}
		val, err := env.PrevEnv.Get(name)
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

func CreateEnvironment(prevEnv Environment, values map[string]any) *Environment {
	return &Environment{
		PrevEnv: &prevEnv,
		Values:  values,
	}
}
