package processes

import (
	"CPU-Simulator/simulator/pkg/memory"
)

// Process State Enum
type State int

const (
	New State = iota
	Ready
	Running
	Blocked
	Terminated
)

type ProcessState struct {
	StackPointer uint32 // Address of the top of the stack

	PC  int // Program Pointer			| Holds address
	AC  int // Accumulator
	MAR int // Memory Address Registers | Holds address

	// in IR
	IROpType  int // 0: Direct, 1: Access memory
	IROpcode  int // Operation code
	IROperand int // Address in Memory

	// In MDR
	IsInstruction bool // Flag to indicate what type of data is stored
	MDROpType     int  // 0: Direct, 1: Access memory
	MDROpcode     int  // Operation code
	MDROperand    int  // Address in Memory
	Data          int  // If holding a data value
}

type PCB struct {
	Pid                 int              // id of the process
	Name                string           // Name of process, only for showing in list.
	PageTable           memory.PageTable // index is the same as the virtual page number
	State               State            // New, Ready, Running, Blocked
	ProcessState        ProcessState
	NextFreeCodeAddress uint32 // next address for the storing instructions
	Priority            int
}
