package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Lox struct {
	HadError bool
}

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage : jlox [help]")
		os.Exit(1)
	}
	lox := Lox{}
	if len(args) == 2 {
		byt, err := os.ReadFile(args[1])
		if err != nil {
			panic(err)
		}

		source := string(byt)
		lox.run(source)
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
			lox.run(inp)
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

func (lox *Lox) run(source string) {
	scannner := CreateScanner(source, lox)
	tokens := scannner.scanTokens()
	// for _, tok := range tokens {
	// 	fmt.Println(tok.toString())
	// }
    if lox.HadError {
        os.Exit(69)
    }
	parser := CreateParser(tokens, lox)
	exprTree, err := parser.parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := exprTree.accept(&Interpreter{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res)
}
