package main

import (
	"errors"
	"fmt"
	"testing"
)

func Do(testCase string, expect string) error {
	lox := Lox{
		Interpreter: CreateAndSetupInterpreter(),
	}
	scannner := CreateScanner(testCase, &lox)
	tokens := scannner.scanTokens()
	if lox.HadError {
		return errors.New("Compile error")
	}
	parser := CreateParser(tokens, &lox)
	statements, err := parser.parse()
	lox.Interpreter.interpret(statements, false)
	if err != nil {
		return err
	}
	return nil
}

func TestExpression(t *testing.T) {
	cases := [][2]string{
		{"1 + 2;", "3"},
		{"5 - 3;", "2"},
		{"4 * 6;", "24"},
		{"20 / 5;", "4"},
		{"(3 + 4) * 2;", "14"},
		{"10 - (2 + 3);", "5"},
		{"3 * (4 + 2) / 3;", "6"},
		{"true && false;", "false"},
		{"true || false;", "true"},
		{"!true;", "false;"},
		{"5 > 3;", "true;"},
		{"5 < 3;", "false;"},
		{"5 >= 5;", "true;"},
		{"4 <= 3;", "false;"},
		{"5 == 5;", "true;"},
		{"5 != 5;", "false;"},
	}

	for i, c := range cases {
		if err := Do(c[0], c[1]); err != nil {
			t.Error(fmt.Sprintf("Wrong on %d : %s\n", i, err.Error()))
		}
	}
}
