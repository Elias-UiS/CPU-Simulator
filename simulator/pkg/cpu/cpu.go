package cpu

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/settings"
	"CPU-Simulator/simulator/pkg/translator"
	"fmt"
	"sync"
	"time"
)

// Instruction Opcodes
const (
	ADD   = 1  // Adds value to the accumulator
	SUB   = 2  // Subtracts value from the accumulator
	STORE = 3  // Store value in accumulator to a memory address
	PRINT = 4  // Print the value of a register
	HALT  = 5  // Stop execution
	JUMP  = 6  // Sets the PC to a specific address
	LOAD  = 7  // Loads the data from an address to the accumulator
	CLEAR = 8  // Clears the accumulator		| Sets the accumulator to 0
	PUSH  = 9  // Pushes the value to the stack
	POP   = 10 // Pops the value from the stack
)

var OpcodeNames = map[int]string{
	ADD:   "ADD",
	SUB:   "SUB",
	STORE: "STORE",
	PRINT: "PRINT",
	HALT:  "HALT",
	JUMP:  "JUMP",
	LOAD:  "LOAD",
	CLEAR: "CLEAR",
	PUSH:  "PUSH",
	POP:   "POP",
}

var OpcodeValues = map[string]int{
	"ADD":   ADD,
	"SUB":   SUB,
	"STORE": STORE,
	"PRINT": PRINT,
	"HALT":  HALT,
	"JUMP":  JUMP,
	"LOAD":  LOAD,
	"CLEAR": CLEAR,
	"PUSH":  PUSH,
	"POP":   POP,
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

	SP int // Stack Pointer
	//SR int // Status Register/Flags
}

type CPU struct {
	Registers        Registers
	opcodes          map[int]func(*CPU) // Map opcode to a handler function
	mmu              *memory.MMU
	IsPaused         bool                          // Flag to check if the CPU is paused
	EventHandler     func(cpu *CPU)                // Event handler to notify the OS about the cycle
	InstructionCount int                           // Count of instructions executed for this process instance
	PageFaultHandler func(*CPU, uint32, int) error // Handler to let os decide when the pte is invalid.
	PauseWaitGroup   *sync.WaitGroup               // To let the cpu wait for the OnCycle function to finish.
}

// Register a new opcode
func (cpu *CPU) registerOpcode(opcode int, handler func(*CPU)) {
	cpu.opcodes[opcode] = handler
}

func (cpu *CPU) fetch() {
	logger.Log.Println("INFO: CPU fetch() instruction")
	virtualAddr := uint32(cpu.Registers.PC)
	physicalAddr, errStruct := cpu.mmu.TranslateAddress(virtualAddr)
	logger.Log.Printf("INFO: CPU.Fetch() - PhysicalAddr: %d", physicalAddr)
	if physicalAddr == -1 {
		fmt.Println(errStruct.Text)
	}
	if errStruct != nil {
		cpu.Pause()
		go cpu.PageFaultHandler(cpu, virtualAddr, errStruct.VPN)
		return
	}
	cpu.Registers.MAR = physicalAddr
	bindings.MarBinding.Set(cpu.Registers.MAR)

	time.Sleep(100 * time.Millisecond)

	instructionBits, err := cpu.mmu.Read(uint32(cpu.Registers.MAR))
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	logger.Log.Println(instructionBits)

	instructionType := (instructionBits >> 31) & 0x1 // Extract the first bit
	opcode := instructionBits & 0x7FFFFFFF           // Extract the last 15 bits

	virtualAddr2 := uint32(cpu.Registers.PC + 1)
	physicalAddr2, errStruct := cpu.mmu.TranslateAddress(virtualAddr2)
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	if errStruct != nil {
		cpu.Pause()
		go cpu.PageFaultHandler(cpu, virtualAddr, errStruct.VPN)
		return
	}
	cpu.Registers.MAR = physicalAddr2
	bindings.MarBinding.Set(cpu.Registers.MAR)
	time.Sleep(100 * time.Millisecond)
	addressBits, err := cpu.mmu.Read(uint32(cpu.Registers.MAR))
	if physicalAddr == -1 {
		fmt.Println(err)
	}

	operand := addressBits
	instruction := Instruction{instructionType, opcode, operand}
	var mdr MDR = MDR{IsInstruction: true, Instruction: instruction}
	cpu.Registers.MDR = mdr

	// Update bindings
	bindings.MdrIsInstructionBinding.Set(cpu.Registers.MDR.IsInstruction)
	bindings.MdrInstructionOpTypeBinding.Set(cpu.Registers.MDR.Instruction.OpType)
	bindings.MdrInstructionOpCodeBinding.Set(cpu.Registers.MDR.Instruction.Opcode)
	bindings.MdrInstructionOperandBinding.Set(cpu.Registers.MDR.Instruction.Operand)
	bindings.MdrDataBinding.Set(cpu.Registers.MDR.Data)
	time.Sleep(100 * time.Millisecond)

	if cpu.Registers.MDR.IsInstruction {
		cpu.Registers.IR = cpu.Registers.MDR.Instruction
		bindings.InstructionOpTypeBinding.Set(cpu.Registers.IR.OpType)
		bindings.InstructionOpCodeBinding.Set(cpu.Registers.IR.Opcode)
		bindings.InstructionOperandBinding.Set(cpu.Registers.IR.Operand)
	}
	time.Sleep(100 * time.Millisecond)
	instructionAddress := translator.FindNextInstructionAddress(settings.MemType, cpu.Registers.PC+2)
	cpu.Registers.PC = instructionAddress
	bindings.PcBinding.Set(cpu.Registers.PC) // Update binding
	time.Sleep(100 * time.Millisecond)
	return
}

func (cpu *CPU) decode() {
	logger.Log.Println("INFO: CPU decode() instruction")
	if cpu.Registers.IR.OpType == 0 {
		var mdr MDR = MDR{IsInstruction: false, Data: cpu.Registers.MDR.Instruction.Operand}
		cpu.Registers.MDR = mdr

		// Update bindings
		bindings.MdrIsInstructionBinding.Set(cpu.Registers.MDR.IsInstruction)
		bindings.MdrInstructionOpTypeBinding.Set(cpu.Registers.MDR.Instruction.OpType)
		bindings.MdrInstructionOpCodeBinding.Set(cpu.Registers.MDR.Instruction.Opcode)
		bindings.MdrInstructionOperandBinding.Set(cpu.Registers.MDR.Instruction.Operand)
		bindings.MdrDataBinding.Set(cpu.Registers.MDR.Data)

	} else {
		cpu.Registers.MAR = cpu.Registers.IR.Operand
		bindings.MarBinding.Set(cpu.Registers.MAR)
	}
	time.Sleep(100 * time.Millisecond)

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
		// Update bindings
		bindings.MdrIsInstructionBinding.Set(cpu.Registers.MDR.IsInstruction)
		bindings.MdrInstructionOpTypeBinding.Set(cpu.Registers.MDR.Instruction.OpType)
		bindings.MdrInstructionOpCodeBinding.Set(cpu.Registers.MDR.Instruction.Opcode)
		bindings.MdrInstructionOperandBinding.Set(cpu.Registers.MDR.Instruction.Operand)
		bindings.MdrDataBinding.Set(cpu.Registers.MDR.Data)
		time.Sleep(100 * time.Millisecond)
	}

	var opcode int = cpu.Registers.IR.Opcode
	if handler, exists := cpu.opcodes[opcode]; exists {
		handler(cpu)
	}
	for i := range cpu.mmu.PageTable.Entries {
		logger.Log.Printf("INFO: CPU.execute() - PageTable Entry nr: %d -> %d", i, cpu.mmu.PageTable.Entries[i].FrameNumber)
	}
	return
}

func GetOpcodeName(opcode int) string {
	if name, exists := OpcodeNames[opcode]; exists {
		return name
	}
	return "Unknown Opcode"
}

// Instructions
func add(cpu *CPU) {
	logger.Log.Println("INFO: CPU add()")
	var oldAcc int = cpu.Registers.AC
	cpu.Registers.AC += cpu.Registers.MDR.Data
	bindings.AcBinding.Set(cpu.Registers.AC)
	fmt.Printf("Prev ACC: %d + Data: %d = New ACC: %d \n", oldAcc, cpu.Registers.MDR.Data, cpu.Registers.AC)
	return
}

func sub(cpu *CPU) {
	logger.Log.Println("INFO: CPU sub()")
	var oldAcc int = cpu.Registers.AC
	cpu.Registers.AC -= cpu.Registers.MDR.Data
	bindings.AcBinding.Set(cpu.Registers.AC)
	fmt.Printf("Prev ACC: %d - Data: %d = New ACC: %d \n", oldAcc, cpu.Registers.MDR.Data, cpu.Registers.AC)
	return
}

func print(cpu *CPU) {
	var name string = GetOpcodeName(cpu.Registers.IR.Opcode)
	fmt.Printf("%s: %d\n", name, cpu.Registers.AC)
	return
}

func store(cpu *CPU) {
	logger.Log.Println("INFO: CPU store()")
	value := cpu.Registers.AC
	destination := cpu.Registers.MDR.Data
	logger.Log.Printf("DEBUG: store() address %d", destination)
	physAddr, errStruct := cpu.mmu.TranslateAddress(uint32(destination))
	if errStruct != nil {
		logger.Log.Printf("ERROR: push() %s", errStruct.Text)
		cpu.Registers.PC -= 2
		cpu.Pause()
		go cpu.PageFaultHandler(cpu, uint32(destination), errStruct.VPN)
		return
	}
	cpu.mmu.Write(uint32(physAddr), uint32(value))
	logger.Log.Println(destination)
	logger.Log.Println(value)
	return
}

func halt(cpu *CPU) {
	logger.Log.Println("INFO: CPU halt()")
	return
}

func jump(cpu *CPU) {
	logger.Log.Println("INFO: CPU jump()")
	cpu.Registers.PC = cpu.Registers.MDR.Data
	bindings.PcBinding.Set(cpu.Registers.PC)
	return
}

func clear(cpu *CPU) {
	logger.Log.Println("INFO: CPU clear()")
	cpu.Registers.AC = 0
	bindings.AcBinding.Set(cpu.Registers.AC)
	return
}

func push(cpu *CPU) {
	logger.Log.Println("INFO: CPU push()")
	destination := cpu.Registers.SP - 1
	logger.Log.Printf("INFO: CPU push() - destination: %d", destination)
	vpn, offset := translator.TranslateAddressToVPNandOffset(cpu.Registers.SP - 1)
	logger.Log.Printf("INFO: push to VPN: %d, Offset: %d", vpn, offset)
	stackPointer := translator.TranslateVPNandOffsetToAddress(vpn-1, settings.PageSize-1)
	logger.Log.Printf("INFO: CPU push() - stackPointer: %d", stackPointer)
	if offset >= settings.PageSize {
		destination = int(translator.TranslateVPNandOffsetToAddress(vpn, settings.PageSize-1))
		logger.Log.Println("INFO: CPU push() - stackPointer: %d", stackPointer)
	}
	physAddr, errStruct := cpu.mmu.TranslateAddress(uint32(destination))
	logger.Log.Println("INFO: CPU push() - physAddr: %d", physAddr)
	if errStruct != nil {
		logger.Log.Printf("ERROR: push() %s", errStruct.Text)
		cpu.Registers.PC -= 2
		cpu.Pause()
		go cpu.PageFaultHandler(cpu, uint32(destination), errStruct.VPN)
		return
	}

	cpu.Registers.SP = destination
	bindings.SpBinding.Set(cpu.Registers.SP)

	value := cpu.Registers.MDR.Data
	cpu.mmu.Write(uint32(physAddr), uint32(value))

	logger.Log.Println(physAddr)
	logger.Log.Println(value)
	return
}

func pop(cpu *CPU) {
	logger.Log.Println("INFO: CPU pop()")
	destination := cpu.Registers.SP
	physAddr, errStruct := cpu.mmu.TranslateAddress(uint32(destination))
	if errStruct != nil {
		logger.Log.Printf("ERROR: pop() %s", errStruct.Text)
		cpu.Registers.PC -= 2
		cpu.Pause()
		go cpu.PageFaultHandler(cpu, uint32(destination), errStruct.VPN)
		return
	}
	value, err := cpu.mmu.Read(uint32(physAddr))
	if err != nil {
		logger.Log.Printf("ERROR: pop() %s", err)
	}
	logger.Log.Println(value)
	cpu.Registers.AC = value
	bindings.AcBinding.Set(cpu.Registers.AC)

	newSP := cpu.Registers.SP + 1
	logger.Log.Printf("INFO: CPU push() - newSP: %d", newSP)
	vpn, offset := translator.TranslateAddressToVPNandOffset(newSP)

	if offset >= settings.PageSize {
		newSP = int(translator.TranslateVPNandOffsetToAddress(vpn+1, 0))
	}

	cpu.Registers.SP = newSP
	bindings.SpBinding.Set(cpu.Registers.SP)
	return
}

// Initialize the CPU with default values
func NewCPU(mmu *memory.MMU) *CPU {
	logger.Log.Println("INFO: CPU New()")
	cpu := &CPU{
		opcodes:        make(map[int]func(*CPU)),
		mmu:            mmu,
		PauseWaitGroup: new(sync.WaitGroup),
	}

	// Adds default instructions to opcodes
	cpu.registerOpcode(ADD, add)
	cpu.registerOpcode(SUB, sub)
	cpu.registerOpcode(PRINT, print)
	cpu.registerOpcode(HALT, halt)
	cpu.registerOpcode(STORE, store)
	cpu.registerOpcode(JUMP, jump)
	cpu.registerOpcode(CLEAR, clear)
	cpu.registerOpcode(PUSH, push)
	cpu.registerOpcode(POP, pop)

	return cpu
}

func (cpu *CPU) Run() {
	cpu.IsPaused = false
	logger.Log.Println("INFO: CPU Run()")

	for {
		if len(cpu.mmu.PageTable.Entries) == 0 {
			logger.Log.Println("No page table found")
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		cpu.PauseWaitGroup.Wait() // If paused it wait, else it runs.
		cpu.fetch()
		if cpu.IsPaused {
			continue
		}
		time.Sleep(settings.CpuFetchDecodeExecuteDelay * time.Millisecond)
		cpu.decode()
		time.Sleep(settings.CpuFetchDecodeExecuteDelay * time.Millisecond)
		cpu.execute()
		if cpu.IsPaused {
			continue
		}
		time.Sleep(settings.CpuFetchDecodeExecuteDelay * time.Millisecond)
		cpu.InstructionCount += 1
		bindings.InstructionCount.Set(cpu.InstructionCount)
		if cpu.EventHandler != nil {
			logger.Log.Println("INFO: CPU EventHandler()")
			cpu.Pause()
			go cpu.EventHandler(cpu) // Notify the OS about the cycle
		}
	}
}

func (cpu *CPU) Pause() {
	logger.Log.Println("Pausing +1")
	cpu.IsPaused = true
	cpu.PauseWaitGroup.Add(1)
}

func (cpu *CPU) Resume() {
	logger.Log.Println("Resuming +1")
	cpu.IsPaused = false
	cpu.PauseWaitGroup.Done()
}
