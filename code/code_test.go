package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		// Fetch the 65534th constant from the constant pool and load it to the stack
		// The operand needs 2 bytes to represent as we use uint16 namely 0xFF and 0xFE
		// Why not uint32? Then we use less bytes and thus instructions are smaller
		// We can also check if the most significant byte 0xFF comes first
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, tt := range tests {

		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want: %d, got: %d", len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want:%d, got: %d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	// Byte offset - Opcode (1 byte) - Operand (2 bytes)
	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := Instructions{} // Flatten the slice of slices into 1 slice

	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant: %q\ngot :%q",
			expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instructions := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandRead, n := ReadOperands(def, instructions[1:]) // Skip the opcode
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want: %d, got: %d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandRead[i] != want {
				t.Errorf("operand wrong. want: %d, got: %d", want, operandRead[i])
			}
		}
	}
}
