package main

import (
	"os"
	"testing"
)

func TestScannerParserAndInterpreter(t *testing.T) {
	byt, err := os.ReadFile("test.ws")
	if err != nil {
		panic(err)
	}
	source := string(byt)
	lox := Lox{
		Interpreter: CreateInterpreter(),
	}
	scannner := CreateScanner(source, &lox)
	tokens := scannner.scanTokens()
	if lox.HadError {
		t.Error("Error while scan")
	}
	parser := CreateParser(tokens, &lox)
	statements, err := parser.parse()
	lox.Interpreter.interpret(statements, false)
	if err != nil {
		t.Error(err)
		return
	}
}
