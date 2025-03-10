package os

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"fmt"
	"time"
)

type OS struct {
	CPU           []*cpu.CPU
	Memory        *memory.Memory
	MMU           *memory.MMU
	ProcessTable  *processes.ProcessTable
	FreeList      *memory.FreeList
	CpuController *Controller
	Scheduler     scheduler.SchedulerInterface
	osIsRunning   bool
	cpuIsRunning  bool
	Test          int
}

func NewOS() *OS {
	// Initialize memory
	mem := memory.NewMemory() // Example: 512 frames, 4KB pages

	// Initialize MMU
	mmu := memory.NewMMU(mem)

	// Initialize CPU
	cpuInstance := cpu.NewCPU(mmu)

	// Initialize free list
	freeList := memory.NewFreeList()

	processTableStruct := processes.CreateProcessTable()

	scheduler := scheduler.NewScheduler()

	// Initialize processController
	controller := createController(mmu, freeList, processTableStruct)

	return &OS{
		CPU:           []*cpu.CPU{cpuInstance},
		Memory:        mem,
		MMU:           mmu,
		ProcessTable:  processTableStruct,
		FreeList:      freeList,
		CpuController: controller,
		osIsRunning:   false,
		cpuIsRunning:  false,
		Scheduler:     scheduler,
		Test:          10,
	}
}

func (os *OS) StartSimulation() {
	if os.osIsRunning {
		return
	}
	logger.Log.Println("Starting simulation...")
	os.osIsRunning = true

	pcb := os.CpuController.MakeTestProcessBasic()
	pcb2 := os.CpuController.MakeTestProcessBasic2()
	pcb3 := os.CpuController.MakeTestProcessBasic()

	// os.AddProcessToProcessTable(pcb)
	// os.AddProcessToProcessTable(pcb2)
	// os.AddProcessToProcessTable(pcb3)

	os.AddProcessToSchedulerQueue(pcb)
	os.AddProcessToSchedulerQueue(pcb2)
	os.AddProcessToSchedulerQueue(pcb3)

	nextPcb := os.Scheduler.GetNextProcess()
	os.CpuController.SetPageTabletoMMU(nextPcb)
	go os.CPU[0].Run() // Run CPU in a separate goroutine
	os.cpuIsRunning = true

	// For testing REMOVE later
	logger.Log.Println("Number of free Fames", os.FreeList.NumberOfFreeFrames)
	os.Memory.Frames[15][0] = 50
	os.Memory.Frames[15][6] = 50
	queue := os.Scheduler.GetReadyQueue()
	logger.Log.Println("Info: Checking Ready Queue")
	for i := range len(queue) {
		logger.Log.Println(queue[i])
	}
	/////

}

func (os *OS) PauseSimulation() {
	if !os.cpuIsRunning {
		return
	}
	fmt.Println("Pausing simulation...")
	for i := range len(os.CPU) {
		os.CPU[i].Pause()
	}
	os.cpuIsRunning = false
}

func (os *OS) ResumeSimulation() {
	if os.osIsRunning {
		return
	}
	logger.Log.Println("Testing bug here")
	fmt.Println("Resuming simulation...")
	for i := range len(os.CPU) {
		logger.Log.Println("Testing bug here 2")
		os.CPU[i].Resume()
	}
	logger.Log.Println("Testing bug here 3")
	os.cpuIsRunning = true
}

// TODO
func (os *OS) StopSimulation() {

	fmt.Println("Stopping simulation...")
	//os.CPU.Stop()
}

func (os *OS) Reset() {
	fmt.Println("Resetting OS...")
	// Reset memory, MMU, CPU, and free list
	os.Memory = memory.NewMemory()
	os.MMU = memory.NewMMU(os.Memory)
	os.CPU = []*cpu.CPU{cpu.NewCPU(os.MMU)}
	os.FreeList = memory.NewFreeList()
}

func (os *OS) ContextSwitch(cpu *cpu.CPU) {
	if !cpu.IsPaused {
		cpu.Pause()
	}
	for i := range os.Scheduler.GetReadyQueue() {
		logger.Log.Println("Ready Queue list 1: ", os.Scheduler.GetReadyQueue()[i].Pid)
	}
	currentProcess := os.Scheduler.GetRunningProcess()
	logger.Log.Printf("currentProcess.Pid: %d", currentProcess.Pid)
	if currentProcess == nil {
		logger.Log.Println("No processes currently running.")
		return
	}
	for i := range os.Scheduler.GetReadyQueue() {
		logger.Log.Println("Ready Queue list 2: ", os.Scheduler.GetReadyQueue()[i].Pid)
	}
	nextProcess := os.Scheduler.GetNextProcess()
	for i := range os.Scheduler.GetReadyQueue() {
		logger.Log.Println("Ready Queue list 3: ", os.Scheduler.GetReadyQueue()[i].Pid)
	}
	logger.Log.Printf("Info: IsTerminated Check: %d", currentProcess.State)
	if currentProcess.State != processes.Terminated {
		currentProcess.State = processes.Ready
		os.AddProcessToSchedulerQueue(currentProcess)
	}

	if nextProcess == nil {
		logger.Log.Println("No processes in the ready queue.")
		return
	}
	logger.Log.Printf("nextProcess.Pid: %d", nextProcess.Pid)
	logger.Log.Println("Context switching from process", currentProcess.Pid)
	for i := range os.Scheduler.GetReadyQueue() {
		logger.Log.Println("Ready Queue list 4: ", os.Scheduler.GetReadyQueue()[i].Pid)
	}
	logger.Log.Println("Context switching to process", nextProcess.Pid)
	os.SaveProcessState(currentProcess, cpu) // Saves the process state of pcb from cpu to pcb..
	os.CpuController.SetPageTabletoMMU(nextProcess)
	for index := range len(nextProcess.PageTable.Entries) {
		logger.Log.Printf("Index: %d, Value: %d\n", index, nextProcess.PageTable.Entries[uint16(index)].FrameNumber)
	}
	os.SetNewProcessState(nextProcess, cpu) // Sets the process state of pcb to cpu.
	nextProcess.State = processes.Running
	if os.cpuIsRunning {
		cpu.Resume()
	}

}

// Function to save the state of the current process
func (os *OS) SaveProcessState(pcb *processes.PCB, cpu *cpu.CPU) {
	logger.Log.Println("Saving state for process", pcb.Pid)

	// Simulate saving program counter and stack pointer (this is an abstract example)
	pcb.ProcessState.AC = cpu.Registers.AC
	pcb.ProcessState.Data = cpu.Registers.MDR.Data
	pcb.ProcessState.IROpType = cpu.Registers.IR.OpType
	pcb.ProcessState.IROpcode = cpu.Registers.IR.Opcode
	pcb.ProcessState.IROperand = cpu.Registers.IR.Operand
	pcb.ProcessState.IsInstruction = cpu.Registers.MDR.IsInstruction
	pcb.ProcessState.MAR = cpu.Registers.MAR
	pcb.ProcessState.MDROpType = cpu.Registers.MDR.Instruction.OpType
	pcb.ProcessState.MDROpcode = cpu.Registers.MDR.Instruction.Opcode
	pcb.ProcessState.MDROperand = cpu.Registers.MDR.Instruction.Operand
	pcb.ProcessState.PC = cpu.Registers.PC
}

// Function to set the process state of the new scheduled process to the cpu.
func (os *OS) SetNewProcessState(pcb *processes.PCB, cpu *cpu.CPU) {
	logger.Log.Println("SetNewProcessState()", pcb.Pid)

	// Simulate saving program counter and stack pointer (this is an abstract example)
	cpu.Registers.AC = pcb.ProcessState.AC
	cpu.Registers.MDR.Data = pcb.ProcessState.Data
	cpu.Registers.IR.OpType = pcb.ProcessState.IROpType
	cpu.Registers.IR.Opcode = pcb.ProcessState.IROpcode
	cpu.Registers.IR.Operand = pcb.ProcessState.IROperand
	cpu.Registers.MDR.IsInstruction = pcb.ProcessState.IsInstruction
	cpu.Registers.MAR = pcb.ProcessState.MAR
	cpu.Registers.MDR.Instruction.OpType = pcb.ProcessState.MDROpType
	cpu.Registers.MDR.Instruction.Opcode = pcb.ProcessState.MDROpcode
	cpu.Registers.MDR.Instruction.Operand = pcb.ProcessState.MDROperand
	cpu.Registers.PC = pcb.ProcessState.PC
}

func (os *OS) GetCpu() *cpu.CPU {
	if len(os.CPU) == 0 {
		return nil
	}

	return os.CPU[0]
}

func (os *OS) AddProcessToProcessTable(pcb *processes.PCB) {
	os.ProcessTable.AddProcessToTable(pcb)
}
func (os *OS) AddProcessToSchedulerQueue(pcb *processes.PCB) {
	os.Scheduler.AddProcess(pcb)
}

func (os *OS) GetScheduler() scheduler.SchedulerInterface {
	scheduler := scheduler.NewScheduler()
	return scheduler
}

func (os *OS) TestNumer() {

	for {
		os.Test += 1
		time.Sleep(500 * time.Millisecond)
	}
}
