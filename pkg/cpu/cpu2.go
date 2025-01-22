package cpu

// import "fmt"

// // Define your initial instruction set
// const (
//     ADD = iota
//     SUB
//     MOV
//     PRINT
//     HALT
// )

// type Instruction struct {
//     Opcode int
//     Arg1   int
//     Arg2   int
//     Arg3   int
// }

// type CPU struct {
//     registers 4[]int
//     memory    []Instruction
//     pc        int // Program Counter
//     opcodes   map[int]func(*CPU, Instruction) // Map opcode to handler function
// }

// // New CPU initializes a CPU with a default set of instructions
// func NewCPU() *CPU {
//     cpu := &CPU{
//         opcodes: make(map[int]func(*CPU, Instruction)),
//     }

//     // Register default instructions
//     cpu.registerOpcode(ADD, cpu.add)
//     cpu.registerOpcode(SUB, cpu.sub)
//     cpu.registerOpcode(MOV, cpu.mov)
//     cpu.registerOpcode(PRINT, cpu.print)
//     cpu.registerOpcode(HALT, cpu.halt)

//     return cpu
// }

// // Register a new opcode (method to handle it)
// func (cpu *CPU) registerOpcode(opcode int, handler func(*CPU, Instruction)) {
//     cpu.opcodes[opcode] = handler
// }

// // Default Instructions
// func (cpu *CPU) add(instruction Instruction) {
//     cpu.registers[instruction.Arg3] = cpu.registers[instruction.Arg1] + cpu.registers[instruction.Arg2]
// }

// func (cpu *CPU) sub(instruction Instruction) {
//     cpu.registers[instruction.Arg3] = cpu.registers[instruction.Arg1] - cpu.registers[instruction.Arg2]
// }

// func (cpu *CPU) mov(instruction Instruction) {
//     cpu.registers[instruction.Arg3] = instruction.Arg1
// }

// func (cpu *CPU) print(instruction Instruction) {
//     fmt.Printf("R%d: %d\n", instruction.Arg1, cpu.registers[instruction.Arg1])
// }

// func (cpu *CPU) halt(instruction Instruction) {
//     fmt.Println("Halting CPU.")
// }

// // Run the CPU cycle
// func (cpu *CPU) run() {
//     for cpu.pc < len(cpu.memory) {
//         // Fetch
//         instruction := cpu.memory[cpu.pc]

//         // Decode & Execute
//         if handler, ok := cpu.opcodes[instruction.Opcode]; ok {
//             handler(cpu, instruction)
//         } else {
//             fmt.Println("Unknown opcode")
//             break
//         }

//         // Increment Program Counter
//         cpu.pc++
//     }
// }
