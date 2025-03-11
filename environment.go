package main

type Environment struct {
	values map[string]any
}

func (env *Environment) Get(name string) any {
	val, ok := env.values[name]
	if !ok {
		panic("TODO : handle error")
	}
	return val
}

func (env *Environment) Set(name string, value any) {
	env.values[name] = value
}
