package cpu

// import "fmt"

// // Instruction Opcodes
// const (
// 	ADD   = iota // Add two registers
// 	SUB          // Subtract two registers
// 	MOV          // Move a value to a register
// 	PRINT        // Print the value of a register
// 	HALT         // Stop execution
// )

// // Instruction represents a single CPU instruction.
// type Instruction struct {
// 	Opcode int // Operation code
// 	Arg1   int // First argument (register index or value)
// 	Arg2   int // Second argument (register index or value)
// 	Arg3   int // Destination register index
// }

// // Run executes a simple CPU simulation.
// func Run() {
// 	// Simulated memory containing instructions
// 	var memory []Instruction = []Instruction{
// 		{Opcode: MOV, Arg1: 5, Arg3: 0},          // MOV 5 -> R0
// 		{Opcode: MOV, Arg1: 10, Arg3: 1},         // MOV 10 -> R1
// 		{Opcode: ADD, Arg1: 0, Arg2: 1, Arg3: 2}, // ADD R0 + R1 -> R2
// 		{Opcode: SUB, Arg1: 1, Arg2: 0, Arg3: 3}, // SUB R1 - R0 -> R3
// 		{Opcode: PRINT, Arg1: 2},                 // PRINT R2
// 		{Opcode: PRINT, Arg1: 3},                 // PRINT R3
// 		{Opcode: HALT},                           // HALT
// 	}

// 	// Registers to hold data
// 	var registers [4]int = [4]int{}

// 	// Program Counter (PC) to track the current instruction
// 	var pc int = 0

// 	// Simulate the Fetch-Decode-Execute cycle
// 	for pc < len(memory) {
// 		// Fetch the next instruction
// 		var instruction Instruction = memory[pc]

// 		// Decode and Execute the instruction
// 		switch instruction.Opcode {
// 		case MOV:
// 			// Move a value to a register
// 			registers[instruction.Arg3] = instruction.Arg1
// 		case ADD:
// 			// Add values from two registers and store the result
// 			registers[instruction.Arg3] = registers[instruction.Arg1] + registers[instruction.Arg2]
// 		case SUB:
// 			// Subtract the value in one register from another
// 			registers[instruction.Arg3] = registers[instruction.Arg1] - registers[instruction.Arg2]
// 		case PRINT:
// 			// Print the value of a register
// 			fmt.Printf("R%d: %d\n", instruction.Arg1, registers[instruction.Arg1])
// 		case HALT:
// 			// Stop the execution
// 			fmt.Println("Halting CPU.")
// 			return
// 		default:
// 			// Handle unknown instructions
// 			fmt.Println("Unknown instruction. Halting.")
// 			return
// 		}

// 		// Increment the program counter to move to the next instruction
// 		pc++
// 	}
// }
