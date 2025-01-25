package cpu

import "fmt"

// Instruction Opcodes
const (
	ADD   = iota // Add two registers
	SUB          // Subtract two registers
	MOV          // Move a value to a register
	PRINT        // Print the value of a register
	HALT         // Stop execution
)

var opcodeNames = map[int]string{
	ADD:   "ADD",
	SUB:   "SUB",
	MOV:   "MOV",
	PRINT: "PRINT",
	HALT:  "HALT",
}

// Instruction represents a single CPU instruction.
type Instruction struct {
	Opcode  int // Operation code
	Operand int // Address in Memory

}

type MDR struct {
	IsInstruction bool        // Flag to indicate what type of data is stored
	Instruction   Instruction // If holding an instruction
	Data          int         // If holding a data value
}

type Registers struct {
	R0 int // General Purpose Register 1
	R1 int // General Purpose Register 2
	R2 int // General Purpose Register 3
	R3 int // General Purpose Register 4

	PC int         // Program Pointer			| Holds address
	IR Instruction // Instruction Register		| Holds instruction
	AC int         // Accumulator

	MAR int // Memory Address Registers | Holds address
	MDR MDR // Memory Data Registers	| Holds instruction

	SP int // Stack Pointer
	SR int // Status Register/Flags
}

type CPU struct {
	Registers       Registers
	Instructions    Instruction
	Memory          [1024]int
	Cache           [16]int
	InstructionList []Instruction
	opcodes         map[int]func(*CPU) // Map opcode to a handler function
}

// Register a new opcode
func (cpu *CPU) registerOpcode(opcode int, handler func(*CPU)) {
	cpu.opcodes[opcode] = handler
}

func (cpu *CPU) fetch() {
	cpu.Registers.MAR = cpu.Registers.PC
	var mdr MDR = MDR{IsInstruction: true, Instruction: cpu.InstructionList[cpu.Registers.MAR]}
	cpu.Registers.MDR = mdr
	if cpu.Registers.MDR.IsInstruction {
		cpu.Registers.IR = cpu.Registers.MDR.Instruction
	}
	cpu.Registers.PC += 1
}

func (cpu *CPU) decode() {
	cpu.Registers.MAR = cpu.Registers.IR.Operand
	var mdr MDR = MDR{IsInstruction: false, Data: cpu.Memory[cpu.Registers.MAR]}
	cpu.Registers.MDR = mdr
}

func (cpu *CPU) execute() {
	var opcode int = cpu.Registers.IR.Opcode
	if handler, exists := cpu.opcodes[opcode]; exists {
		handler(cpu)
	}
}

func getOpcodeName(opcode int) string {
	if name, exists := opcodeNames[opcode]; exists {
		return name
	}
	return "Unknown Opcode"
}

// Instructions
func add(cpu *CPU) {
	cpu.Registers.AC += cpu.Registers.MDR.Data
}

func sub(cpu *CPU) {
	cpu.Registers.AC -= cpu.Registers.MDR.Data
}

func print(cpu *CPU) {
	var name string = getOpcodeName(cpu.Registers.IR.Opcode)
	fmt.Printf("%s: %d\n", name, cpu.Registers.AC)
}

func halt(cpu *CPU) {
	fmt.Println("Halting CPU.")
}

// Initialize the CPU with default values
func NewCPU() *CPU {
	cpu := &CPU{
		opcodes: make(map[int]func(*CPU)),
	}

	// Adds default instructions to opcodes
	cpu.registerOpcode(ADD, add)
	cpu.registerOpcode(SUB, sub)
	cpu.registerOpcode(PRINT, print)
	cpu.registerOpcode(HALT, halt)

	return cpu
}

// Run executes the CPU simulation.
func Run() {
	// Simulated memory containing instructions
	var instructions []Instruction = []Instruction{
		{Opcode: ADD, Operand: 1},
		{Opcode: PRINT},
		{Opcode: ADD, Operand: 2},
		{Opcode: PRINT},
		{Opcode: SUB, Operand: 3},
		{Opcode: PRINT},
	}

	cpu := NewCPU()
	cpu.InstructionList = instructions

	cpu.Memory[0] = 100
	cpu.Memory[1] = 1
	cpu.Memory[2] = 2
	cpu.Memory[3] = 3
	for cpu.Registers.PC < len(instructions) {
		cpu.fetch()
		cpu.decode()
		cpu.execute()
	}
}
