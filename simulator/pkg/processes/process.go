package processes

import (
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/metrics"
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

func (s State) String() string {
	switch s {
	case New:
		return "New"
	case Ready:
		return "Ready"
	case Running:
		return "Running"
	case Blocked:
		return "Blocked"
	case Terminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

type ProcessState struct {
	SP int // Address of the top of the stack

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
	Pid                 int               // id of the process
	Name                string            // Name of process, only for showing in list.
	PageTable           *memory.PageTable // index is the same as the virtual page number
	State               State             // New, Ready, Running, Blocked
	ProcessState        ProcessState      // State of the registers
	NextFreeCodeAddress uint32            // next address for the storing instructions
	Priority            int               // Priority of the process
	PageAmount          int               // Number of pages available to the process
	Metrics             metrics.Metrics   // Metrics for the process
	UpdateChan          chan bool         // Channel to notify the process has changed
}
