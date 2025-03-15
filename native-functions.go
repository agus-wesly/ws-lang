package main

import "time"

type Clock struct{}

func (c *Clock) call(*Interpreter, *[]Expression) (any, error) {
	return time.Now().UnixMilli(), nil
}

func (c *Clock) arity() int {
	return 0
}

func (c *Clock) toString() string {
	return "<native fn>"
}

func SetupInterpreter(i *Interpreter) {
	i.Environment.Set("clock", &Clock{})
	i.Environment.Set("bar", &FunctionDeclaration{})
}
