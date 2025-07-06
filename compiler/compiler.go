package compiler

import (
	"fmt"
	"s8/ast"
	"s8/code"
	"s8/object"
)

type Compiler struct {
	instructions code.Instructions
	// The internal instruction slice
	constants []object.Object
	// The constant pool
	lastInstruction EmittedInstruction
	// The one we have just emitted
	previousInstruction EmittedInstruction
	// The one before the recently emitted instruction
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

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
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
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		if node.Alternative == nil {
			// Jump to statements outside of If block
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		} else {
			// Part of node.Consequence to skip over the else branch
			// in case the condition is truthy
			jumpPos := c.emit(code.OpJump, 9999)

			// Jump to node.Alternative
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

			// Jump to statements outside of If-Else block
			afterAlternativePos := len(c.instructions)
			c.changeOperand(jumpPos, afterAlternativePos)
		}
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Return compiled bytecode
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
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

// "emit" is a compiler's term for "generate" and "output"
// Generate an instruction and add it to the results (could be a file, memory collection, etc.)
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

// Type-safe way to check the latest emitted instruction
// without having to do casting from and to bytes
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	prev := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = prev
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

// Replace an instruction at an arbitrary offset (pos)
// in the isntruction slice
func (c *Compiler) replaceInstructions(pos int, newInstruction []byte) {
	for i := range newInstruction {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// Re-create the instruction with the new operand
// assuming we only replace instructions of the same type
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstructions(opPos, newInstruction)
}
