package systemState

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/os"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"time"
)

// CPU               []*cpu.CPU
// Memory            *memory.Memory
// MMU               *memory.MMU
// ProcessTable      *processes.ProcessTable
// FreeList          *memory.FreeList
// Scheduler         scheduler.SchedulerInterface

type State struct {
	CurrentTime  string
	Loop         int
	PubSub       *PubSub[State] `json:"-"`
	Registers    cpu.Registers
	Memory       memory.Memory
	ProcessTable processes.ProcessTable
	FreeList     memory.FreeList
	Scheduler    scheduler.SchedulerInterface
}

func CreateState() *State {
	state := &State{}
	pubSub := NewPubSub[State]()
	state.PubSub = pubSub
	return state
}

func (state State) UpdateState(os *os.OS) {
	logger.Log.Println("INFO: Starting UpdateState()")
	for {
		state.Loop += 1
		timeNow := time.Now()
		formattedTime := timeNow.Format("2006-01-02 | 15:04:05")
		state.CurrentTime = formattedTime
		logger.Log.Printf("INFO: Starting UpdateState() %d \n", state.Loop)
		//state.ProcessTable = state.DeepCopyProcessTable(os.ProcessTable)

		state.Registers = os.CPU[0].Registers
		logger.Log.Println(state.Registers.SP)
		state.PubSub.Publish(state)
		time.Sleep(1000 * time.Millisecond)
	}
}

func (state State) DeepCopyProcessTable(processTable *processes.ProcessTable) processes.ProcessTable {
	newProcessTable := processes.ProcessTable{}
	for i, pcb := range processTable.ProcessMap {
		newProcessTable.ProcessMap[i] = state.DeepCopyPCB(pcb)
	}
	return newProcessTable
}

func (state State) DeepCopyPCB(pcb *processes.PCB) *processes.PCB {
	newPCB := &processes.PCB{}
	pageTable := state.DeepCopyPageTable(pcb.PageTable)
	newPCB.Pid = pcb.Pid
	newPCB.Name = pcb.Name
	newPCB.State = pcb.State
	newPCB.PageTable = &pageTable
	return newPCB
}

func (state State) DeepCopyPageTable(pt *memory.PageTable) memory.PageTable {
	newPageTable := memory.PageTable{}
	for i, value := range pt.Entries {
		newPageTable.Entries[i] = state.DeepCopyPTE(value)
	}

	return newPageTable
}

func (state State) DeepCopyPTE(pte *memory.PTE) *memory.PTE {
	newPTE := &memory.PTE{}
	newPTE.FrameNumber = pte.FrameNumber
	newPTE.Type = pte.Type
	newPTE.Valid = pte.Valid
	return newPTE
}

func (state State) DeepCopyRegister(register cpu.Registers) cpu.Registers {
	newRegister := cpu.Registers{}
	newInstruction := cpu.Instruction{}

	newInstruction.OpType = register.IR.OpType
	newInstruction.Opcode = register.IR.Opcode
	newInstruction.Operand = register.IR.Operand

	newRegister.AC = register.AC
	newRegister.IR = newInstruction
	newRegister.MAR = register.MAR
	newRegister.MDR = register.MDR
	newRegister.PC = register.PC
	newRegister.SP = register.SP
	return newRegister
}
