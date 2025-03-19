package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Lox struct {
	HadError bool
	*Interpreter
	*Resolver
}

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage : jlox [help]")
		os.Exit(1)
	}
	interpreter := CreateAndSetupInterpreter()
	lox := Lox{
		Interpreter: interpreter,
		Resolver:    CreateResolver(interpreter, make([]map[string]bool, 0)),
	}
	if len(args) == 2 {
		byt, err := os.ReadFile(args[1])
		if err != nil {
			panic(err)
		}
		source := string(byt)
		lox.run(source, false)
	} else {
		lox.showPrompt()
	}
}

func (lox *Lox) showPrompt() {
	for {
		fmt.Printf("> ")
		reader := bufio.NewReader(os.Stdin)
		inp, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		inp = strings.TrimSpace(inp)
		if inp == "" {
			break
		} else {
			lox.run(inp, true)
			lox.HadError = false
		}
	}
}

func (lox *Lox) error(token Token, msg string) {
	if token.Type == EOF {
		lox.printError(token.Line, "at "+"end", msg)
	} else {
		lox.printError(token.Line, "at "+token.Lexeme, msg)
	}
}

func (lox *Lox) printError(line int, where string, msg string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, msg)
	lox.HadError = true
}

func (lox *Lox) run(source string, replMode bool) {
	scannner := CreateScanner(source, lox)
	tokens := scannner.scanTokens()
	// for _, tok := range tokens {
	// 	fmt.Println(tok.toString())
	// }
	if lox.HadError {
		os.Exit(69)
	}
	parser := CreateParser(tokens, lox)
	statements, err := parser.parse()
	lox.Interpreter.interpret(statements, replMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
