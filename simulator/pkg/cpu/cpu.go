package cpu

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"fmt"
	"time"
)

// Instruction Opcodes
const (
	ADD   = 1 // Adds value to the accumulator
	SUB   = 2 // Subtracts value to the accumulator
	STORE = 3 // Store value in accumulator to a memory address
	PRINT = 4 // Print the value of a register
	HALT  = 5 // Stop execution
	JUMP  = 6 // Sets the PC to a specific address
	LOAD  = 7 // Loads the data from an address to the accumulator
	CLEAR = 8 // Clears the accumulator		| Sets the accumulator to 0
)

var opcodeNames = map[int]string{
	ADD:   "ADD",
	SUB:   "SUB",
	STORE: "STORE",
	PRINT: "PRINT",
	HALT:  "HALT",
	JUMP:  "JUMP",
	LOAD:  "LOAD",
	CLEAR: "CLEAR",
}

// Instruction represents a single CPU instruction.
type Instruction struct {
	OpType  int // 0: Direct, 1: Access memory
	Opcode  int // Operation code
	Operand int // Address in Memory

}

func (instruction *Instruction) ToInt() uint64 {
	// Ensure the instruction parts are within the correct bit ranges:
	// - OpType should be 1 bit (0 or 1)
	// - Opcode should be 31 bits (0 to 2147483647)
	// - Operand should be 32 bits (0 to 4294967295)

	var result uint64
	result |= uint64(instruction.OpType) << 63 // OpType takes the first bit (bit 63)
	result |= uint64(instruction.Opcode) << 32 // Opcode takes the next 31 bits (bits 32-62)
	result |= uint64(instruction.Operand)      // Operand takes the last 32 bits (bits 0-31)

	return result
}

type MDR struct {
	IsInstruction bool        // Flag to indicate what type of data is stored
	Instruction   Instruction // If holding an instruction
	Data          int         // If holding a data value
}

type Registers struct {
	//R0 int // General Purpose Register 1

	PC int         // Program Pointer			| Holds address
	IR Instruction // Instruction Register		| Holds instruction
	AC int         // Accumulator

	MAR int // Memory Address Registers | Holds address
	MDR MDR // Memory Data Registers	| Holds instruction

	//SP int // Stack Pointer
	//SR int // Status Register/Flags
}

type CPU struct {
	Registers  Registers
	opcodes    map[int]func(*CPU) // Map opcode to a handler function
	mmu        *memory.MMU
	IsPaused   bool      // Flag to check if the CPU is paused
	pauseChan  chan bool // Channel to signal when to pause
	resumeChan chan bool // Channel to signal when to resume
}

// Register a new opcode
func (cpu *CPU) registerOpcode(opcode int, handler func(*CPU)) {
	cpu.opcodes[opcode] = handler
}

func (cpu *CPU) fetch() {
	logger.Log.Println("INFO: CPU fetch() instruction")
	virtualAddr := uint32(cpu.Registers.PC)
	physicalAddr, err := cpu.mmu.TranslateAddress(virtualAddr)
	logger.Log.Printf("INFO: CPU.Fetch() - PhysicalAddr: %d", physicalAddr)
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	cpu.Registers.MAR = physicalAddr
	instructionBits, err := cpu.mmu.Read(uint32(cpu.Registers.MAR))
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	logger.Log.Println(instructionBits)

	instructionType := (instructionBits >> 31) & 0x1 // Extract the first bit
	opcode := instructionBits & 0x7FFFFFFF           // Extract the last 15 bits

	virtualAddr2 := uint32(cpu.Registers.PC + 1)
	physicalAddr2, err := cpu.mmu.TranslateAddress(virtualAddr2)
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	cpu.Registers.MAR = physicalAddr2
	addressBits, err := cpu.mmu.Read(uint32(cpu.Registers.MAR))
	if physicalAddr == -1 {
		fmt.Println(err)
	}

	operand := addressBits
	instruction := Instruction{instructionType, opcode, operand}
	var mdr MDR = MDR{IsInstruction: true, Instruction: instruction}
	cpu.Registers.MDR = mdr
	if cpu.Registers.MDR.IsInstruction {
		cpu.Registers.IR = cpu.Registers.MDR.Instruction
	}
	cpu.Registers.PC += 2
}

func (cpu *CPU) decode() {
	logger.Log.Println("INFO: CPU decode() instruction")
	if cpu.Registers.IR.OpType == 0 {
		var mdr MDR = MDR{IsInstruction: false, Data: cpu.Registers.MDR.Instruction.Operand}
		cpu.Registers.MDR = mdr
	} else {
		cpu.Registers.MAR = cpu.Registers.IR.Operand
	}

}

func (cpu *CPU) execute() {
	logger.Log.Println("INFO: CPU execute() instruction")
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
	for i := range cpu.mmu.PageTable.Entries {
		logger.Log.Printf("INFO: CPU.execute() - PageTable Entry nr: %d -> %d", i, cpu.mmu.PageTable.Entries[i].FrameNumber)
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
	logger.Log.Println("INFO: CPU add()")
	var oldAcc int = cpu.Registers.AC
	cpu.Registers.AC += cpu.Registers.MDR.Data
	fmt.Printf("Prev ACC: %d + Data: %d = New ACC: %d \n", oldAcc, cpu.Registers.MDR.Data, cpu.Registers.AC)
}

func sub(cpu *CPU) {
	logger.Log.Println("INFO: CPU sub()")
	var oldAcc int = cpu.Registers.AC
	cpu.Registers.AC -= cpu.Registers.MDR.Data
	fmt.Printf("Prev ACC: %d - Data: %d = New ACC: %d \n", oldAcc, cpu.Registers.MDR.Data, cpu.Registers.AC)

}

func print(cpu *CPU) {
	var name string = getOpcodeName(cpu.Registers.IR.Opcode)
	fmt.Printf("%s: %d\n", name, cpu.Registers.AC)
}

func store(cpu *CPU) {
	logger.Log.Println("INFO: CPU store()")
	value := cpu.Registers.AC
	destination := cpu.Registers.MDR.Data
	logger.Log.Printf("DEBUG: store() address %d", destination)
	physAddr, err := cpu.mmu.TranslateAddress(uint32(destination))
	if err != nil {
		logger.Log.Printf("ERROR: Store() %s", err)
		return
	}
	cpu.mmu.Write(uint32(physAddr), uint32(value))
	logger.Log.Println(destination)
	logger.Log.Println(value)
}

func halt(cpu *CPU) {
	logger.Log.Println("INFO: CPU halt()")
}

func jump(cpu *CPU) {
	logger.Log.Println("INFO: CPU jump()")
	cpu.Registers.PC = cpu.Registers.MDR.Data
}

func clear(cpu *CPU) {
	logger.Log.Println("INFO: CPU clear()")
	cpu.Registers.AC = 0
}

// Initialize the CPU with default values
func NewCPU(mmu *memory.MMU) *CPU {
	logger.Log.Println("INFO: CPU New()")
	cpu := &CPU{
		opcodes: make(map[int]func(*CPU)),
		mmu:     mmu,
	}

	// Adds default instructions to opcodes
	cpu.registerOpcode(ADD, add)
	cpu.registerOpcode(SUB, sub)
	cpu.registerOpcode(PRINT, print)
	cpu.registerOpcode(HALT, halt)
	cpu.registerOpcode(STORE, store)
	cpu.registerOpcode(JUMP, jump)
	cpu.registerOpcode(CLEAR, clear)

	return cpu
}

func (cpu *CPU) Run() {
	cpu.IsPaused = false
	cpu.pauseChan = make(chan bool)  // To pause the CPU
	cpu.resumeChan = make(chan bool) // To resume the CPU
	logger.Log.Println("INFO: CPU Run()")

	for {
		if len(cpu.mmu.PageTable.Entries) == 0 {
			logger.Log.Println("No page table found")
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		// Handle pause/resume logic
		select {
		case <-cpu.pauseChan:
			cpu.IsPaused = true
			logger.Log.Println("INFO: CPU paused")
			// Waits until resume signal is received
			<-cpu.resumeChan
			cpu.IsPaused = false
			logger.Log.Println("INFO: CPU resumed")
		default:
			// Proceed with CPU operations
			if !cpu.IsPaused {
				cpu.fetch()
				time.Sleep(500 * time.Millisecond)
				cpu.decode()
				time.Sleep(500 * time.Millisecond)
				cpu.execute()
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

func (cpu *CPU) Pause() {
	cpu.pauseChan <- true
}

func (cpu *CPU) Resume() {
	cpu.resumeChan <- true
}

// func Run(cpu *CPU) {
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
