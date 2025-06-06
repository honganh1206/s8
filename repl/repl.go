package repl

import (
	"bufio"
	"fmt"
	"io"
	"s8/compiler"
	"s8/evaluator"
	"s8/lexer"
	"s8/object"
	"s8/parser"
	"s8/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// env persists between calls to Eval()
	// env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		// expanded := evaluator.ExpandMacros(program, macroEnv)

		// evaluated := evaluator.Eval(expanded, env)

		// if evaluated != nil {
		// 	// We return a string representation of the obj here
		// 	io.WriteString(out, evaluated.Inspect())
		// 	io.WriteString(out, "\n")
		// }

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "executing bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := machine.LastPoppedStackElement()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")

	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
