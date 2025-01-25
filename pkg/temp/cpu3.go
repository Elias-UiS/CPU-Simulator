package temp

// import "fmt"

// // A basic CPU simulator

// type CPU struct {
// 	PC        int       // Program Counter
// 	Registers [8]int    // 8 General Purpose Registers
// 	Memory    [1024]int // Simple memory (1024 words)
// 	Cache     [2][4]int // Simple 2-level cache, each level has 4 slots
// }

// type Instruction struct {
// 	OpCode int // Operation code (e.g., ADD, SUB)
// 	Reg1   int // First operand register
// 	Reg2   int // Second operand register
// 	Reg3   int // Destination register
// }

// // Initialize the CPU with default values
// func NewCPU() *CPU {
// 	return &CPU{
// 		PC:        0,
// 		Registers: [8]int{0, 0, 0, 0, 0, 0, 0, 0},
// 		Memory:    [1024]int{},
// 		Cache:     [2][4]int{},
// 	}
// }

// const (
// 	ADD = iota
// 	LOAD
// 	STORE
// )

// // Fetch the instruction from memory
// func (cpu *CPU) Fetch() Instruction {
// 	// For simplicity, assume PC points directly to memory locations
// 	return Instruction{
// 		OpCode: cpu.Memory[cpu.PC],
// 		Reg1:   int(cpu.Memory[cpu.PC+1]),
// 		Reg2:   int(cpu.Memory[cpu.PC+2]),
// 		Reg3:   int(cpu.Memory[cpu.PC+3]),
// 	}
// }

// // Decode and Execute the instruction
// func (cpu *CPU) Execute(inst Instruction) {
// 	switch inst.OpCode {
// 	case ADD:
// 		cpu.Registers[inst.Reg3] = cpu.Registers[inst.Reg1] + cpu.Registers[inst.Reg2]
// 	case LOAD:
// 		cpu.Registers[inst.Reg1] = cpu.Memory[inst.Reg2]
// 	case STORE:
// 		cpu.Memory[inst.Reg1] = cpu.Registers[inst.Reg2]
// 	}
// }

// // Simulate a single cycle
// func (cpu *CPU) Cycle() {
// 	inst := cpu.Fetch()
// 	cpu.Execute(inst)
// 	cpu.PC += 4 // Move PC to the next instruction (assuming 4-byte instructions)
// }

// func Run() {
// 	// Create and initialize the CPU
// 	cpu := NewCPU()

// 	// Example program in memory (simple instructions)
// 	cpu.Memory[0] = ADD
// 	cpu.Memory[1] = 0 // Reg1
// 	cpu.Memory[2] = 1 // Reg2
// 	cpu.Memory[3] = 2 // Reg3

// 	// Simulate one cycle
// 	cpu.Cycle()

// 	// Output the result (value in Reg3)
// 	fmt.Println("Reg3 value:", cpu.Registers[2])
// }
