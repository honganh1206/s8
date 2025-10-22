package compiler

import (
	"fmt"
	"sort"

	"s8/ast"
	"s8/code"
	"s8/object"
)

type Compiler struct {
	// The internal instruction slice
	instructions code.Instructions
	// The constant pool
	constants   []object.Object
	symbolTable *SymbolTable
	// A stack of compilation scopes. Each scope is for a compiled function
	scopes     []CompilationScope
	scopeIndex int
}

// Compiled bytecode
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// Keep track of previously emitted instructions
// for branch instructioins
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	// We create duplications here
	// by creating a new compiler instance
	// then overwrite the symbol table and constants.
	// Not a problem for Go's GC though.
	compiler := New()
	// Preserve the symbol table and constant pool
	// in each iteration in repl.go
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	// Similar structure like Evaf l()
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.InfixExpression:
		// Reordering the operands
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case "|":
			c.emit(code.OpPipe)
		case ">>":
			c.emit(code.OpRShift)
		case "<<":
			c.emit(code.OpLShift)
		case "^":
			c.emit(code.OpExponent)
		case "&":
			c.emit(code.OpAmpersand)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		case "~":
			c.emit(code.OpTilde)
		case "++":
			c.emit(code.OpPreInc)
		case "--":
			c.emit(code.OpPreDec)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}
	case *ast.PostfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "++":
			c.emit(code.OpPostInc)
		case "--":
			c.emit(code.OpPostDec)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		// We can later modify the operand of OpJumpNotTruthy
		// AFTER we compile node.Consequence,
		// that way we know how far the VM has to jump.
		// This is called back-patching
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		// node.Consequence will also be compiled as an expression statement,
		// thus there will be an additional OpPop but we need to retain the latest statement
		// We need to get rid of this
		// since Consequence and Alternative need to leave a value on the stack
		// if (true) {
		// 	3;
		// 	2;
		// 	1; // This must be on the stack
		// }
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// We need this whether we have the Alternative or not
		// to jump to the next instruction after If-Else
		jumpPos := c.emit(code.OpJump, 9999)

		// Handle scoped instructions and main instructions
		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}
		// If not truthy but we have Alternatiive, jump to statements outside of Else block
		// If not truthy but there is no Alternative, jump to OpNull
		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			// A compile-time error!
			// With our evaluator we cannot throw an error before we pass bytecode to the VM
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.loadSymbols(symbol)
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))
	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		// Sort before compilation.
		// Not entirely required, but without this the tests will break
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}
		// Doubled num of values on the stack since k-v
		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		// Could be either array or hash literal here
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)
	case *ast.FunctionLiteral:
		c.enterScope()

		// Define function arguments as local bindings
		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		// Handle both implicit and explicit returns.
		// We do this after compiling the function's body and BEFORE leaving the scope,
		// since we still have access to the just-emitted instruction
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		// Change where compiled instructions are stored
		// and this time they are not in the main scope
		compiledFn := &object.CompiledFunction{Instructions: instructions, NumLocals: numLocals, NumParameters: len(node.Parameters)}

		fnIndex := c.addConstant(compiledFn)
		// Turning all functions to closures
		c.emit(code.OpClosure, fnIndex, 0)
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments))
	}

	return nil
}

// Return compiled bytecode
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

// Add the result of evaluation to the constant pool
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	// Return the index of the object at the end of the pool
	// The index also works as the identifier
	return len(c.constants) - 1
}

// "emit" is a compiler's term for "generate" and "output".
// Generate an instruction and add it to the internal instructions slice
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	return pos
}

// emit invokes this
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	// Add function statements (now as instructions) to the local scope
	c.scopes[c.scopeIndex].instructions = updatedInstructions
	return posNewInstruction
}

// Type-safe way to check the latest emitted instruction
// without having to do casting from and to bytes
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	prev := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = prev
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// Retain a value on a stack
func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

// Replace an instruction at an arbitrary offset (pos)
// in the isntruction slice
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()
	for i := range newInstruction {
		ins[pos+i] = newInstruction[i]
	}
}

// Re-create the instruction with the new operand
// assuming we only replace instructions of the same type
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	// Enclose the newly created symbol table of the current scopeIndex to the global symbol table.
	// At this point our compiler uses a fresh symbol table, and we can still use symbols from the global table (closure, baby!)
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	// Remove the  top-of-the-stack scope
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer
	return instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) loadSymbols(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	}
}
