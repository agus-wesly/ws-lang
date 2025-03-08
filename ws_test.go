package main

import (
	"testing"
)

func TestScanner(t *testing.T) {
	input := "1 + -1 * 3"
	expect := []*Token{
		CreateToken(NUMBER, "1", "1", 1),
		CreateToken(PLUS, "+", nil, 1),
		CreateToken(MINUS, "-", nil, 1),
		CreateToken(NUMBER, "1", "1", 1),
		CreateToken(STAR, "*", nil, 1),
		CreateToken(NUMBER, "3", "3", 1),
		CreateToken(EOF, "EOF", nil, 1),
	}

	l := Lox{}
	scanner := CreateScanner(input, &l)
	tokens := scanner.scanTokens()

	for i := 0; i < len(tokens); i += 1 {
		if tokens[i].toString() != expect[i].toString() {
			t.Errorf("Not match %s -> %s\n", tokens[i].toString(), expect[i].toString())
		}
	}
}

func TestScannerParserAndInterpreter(t *testing.T) {
	input := "1 + -1 * 3"

	right := CreateBinary(CreateUnary(CreateLiteral(1.0), CreateToken(MINUS, "-", nil, 1)), *CreateToken(STAR, "*", nil, 1), CreateLiteral(3.0))
	expect := CreateBinary(CreateLiteral(1.0), *CreateToken(PLUS, "+", nil, 1), right)
	expectedValue, err := expect.accept(&Interpreter{})
	if err != nil {
		t.Error(err)
	}

	lox := Lox{}
	scanner := CreateScanner(input, &lox)
	tokens := scanner.scanTokens()
	parser := CreateParser(tokens, &lox)
	inputTree, err := parser.parse()
	if err != nil {
		t.Error(err)
		return
	}

	inputValue, err := inputTree.accept(&Interpreter{})
	if err != nil {
		t.Error(err)
		return
	}

	if inputValue != expectedValue {
		t.Errorf("NOT EQUAL %s WITH %s\n", inputValue, expectedValue)
		return
	}
}
