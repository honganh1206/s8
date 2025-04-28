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
		// We can also check the most significant byte 0xFF comes first
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
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
