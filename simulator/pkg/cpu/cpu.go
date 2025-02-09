package cpu

import (
	"CPU-Simulator/simulator/pkg/memory"
	"fmt"
)

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
	OpType  int // 0: Direct, 1: Access memory
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
	Registers Registers
	opcodes   map[int]func(*CPU) // Map opcode to a handler function
	mmu       *memory.MMU
	pcb       *memory.PCB
}

// Register a new opcode
func (cpu *CPU) registerOpcode(opcode int, handler func(*CPU)) {
	cpu.opcodes[opcode] = handler
}

func (cpu *CPU) fetch() {
	physicalAddr, err := cpu.mmu.TranslateAddress(uint32(cpu.Registers.PC))
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	cpu.Registers.MAR = physicalAddr
	typeAndOpcode, err := cpu.mmu.Read(uint32(cpu.Registers.MAR))
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	instructionType := (typeAndOpcode >> 7) & 0x1 // Extract the first bit
	opcode := typeAndOpcode & 0x7F                // Extract the last 7 bits

	operand, err := cpu.mmu.Read(uint32(cpu.Registers.MAR + 1))
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	instruction := Instruction{instructionType, opcode, operand}
	var mdr MDR = MDR{IsInstruction: true, Instruction: instruction}
	cpu.Registers.MDR = mdr
	if cpu.Registers.MDR.IsInstruction {
		cpu.Registers.IR = cpu.Registers.MDR.Instruction
	}
	cpu.Registers.PC += 2
}

func (cpu *CPU) decode() {
	if cpu.Registers.MDR.Instruction.OpType == 0 {
		var mdr MDR = MDR{IsInstruction: false, Data: cpu.Registers.MDR.Instruction.Operand}
		cpu.Registers.MDR = mdr
	} else {
		cpu.Registers.MAR = cpu.Registers.IR.Operand
	}

}

func (cpu *CPU) execute() {
	if cpu.Registers.MDR.IsInstruction {
		value, err := cpu.mmu.Read(uint32(cpu.Registers.MAR))
		if err != nil {
			fmt.Println(err)
		}
		var mdr MDR = MDR{IsInstruction: false, Data: value}
		cpu.Registers.MDR = mdr
	}

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
	var oldAcc int = cpu.Registers.AC
	cpu.Registers.AC += cpu.Registers.MDR.Data
	fmt.Printf("Prev ACC: %d + Data: %d = New ACC: %d \n", oldAcc, cpu.Registers.MDR.Data, cpu.Registers.AC)
}

func sub(cpu *CPU) {
	var oldAcc int = cpu.Registers.AC
	cpu.Registers.AC -= cpu.Registers.MDR.Data
	fmt.Printf("Prev ACC: %d - Data: %d = New ACC: %d \n", oldAcc, cpu.Registers.MDR.Data, cpu.Registers.AC)

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
func Run(cpu *CPU) {

	// cpu.Memory[0] = 100
	// cpu.Memory[1] = 10
	// cpu.Memory[2] = 20
	// cpu.Memory[3] = 5

	// var instructions []Instruction = []Instruction{
	// 	{Opcode: ADD, Operand: 1},
	// 	{Opcode: ADD, Operand: 2},
	// 	{Opcode: SUB, Operand: 3},
	// }
	// cpu.InstructionList = instructions

	// for i := 0; i < 100; i++ {
	// 	// Simulated memory containing instructions

	// 	for cpu.Registers.PC < len(instructions) {
	// 		cpu.fetch()
	// 		cpu.decode()
	// 		cpu.execute()
	// 		time.Sleep(time.Second * time.Duration(settings.UpdateTimer))
	// 	}
	// 	cpu.Registers.PC = 0
	// 	fmt.Println("New loop")

	// }

}
