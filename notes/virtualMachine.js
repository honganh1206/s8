const PUSH = "PUSH";
const ADD = "ADD";
const MINUS = "MINUS";

let virtualMachine = function (program) {
  let programCounter = 0;
  let stack = [];
  let stackPointer = 0;

  while (programCounter < program.length) {
    let currentInstruction = program[programCounter];
    switch (currentInstruction) {
      case PUSH:
        // Put the instruction at the top of the stack
        stack[stackPointer] = program[programCounter + 1];
        stackPointer++;
        programCounter++;
        break;
      case ADD:
        // Get the right and left operand on the 2nd and 3rd element of the stack
        right = stack[stackPointer - 1];
        stackPointer--;
        left = stack[stackPointer - 1];
        stackPointer--;

        stack[stackPointer] = left + right;
        stackPointer++;
        break;

      case MINUS:
        right = stack[stackPointer - 1];
        stackPointer--;
        left = stack[stackPointer - 1];
        stackPointer--;

        stack[stackPointer] = left - right;
        stackPointer++;
        break;
    }
    programCounter++;
  }
  console.log("stacktop: ", stack[stackPointer - 1]);
};

let program = [PUSH, 3, PUSH, 4, ADD, PUSH, 5, MINUS];

virtualMachine(program);
