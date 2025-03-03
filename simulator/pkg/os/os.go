package os

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"fmt"
)

type OS struct {
	CPU           []*cpu.CPU
	Memory        *memory.Memory
	MMU           *memory.MMU
	ProcessList   map[uint32]*processes.PCB
	FreeList      *memory.FreeList
	CpuController *Controller
	Scheduler     scheduler.SchedulerInterface
	osIsRunning   bool
	cpuIsRunning  bool
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

	// Initialize processController
	controller := createController(mmu, freeList)

	return &OS{
		CPU:           []*cpu.CPU{cpuInstance},
		Memory:        mem,
		MMU:           mmu,
		ProcessList:   make(map[uint32]*processes.PCB),
		FreeList:      freeList,
		CpuController: controller,
		osIsRunning:   false,
		cpuIsRunning:  false,
	}
}

func (os *OS) StartSimulation() {
	if os.osIsRunning {
		return
	}
	logger.Log.Println("Starting simulation...")
	os.osIsRunning = true
	scheduler := os.GetScheduler()
	os.Scheduler = scheduler
	pcb := os.CpuController.MakeTestProcessBasic()
	pcb2 := os.CpuController.MakeTestProcessBasic2()
	pcb3 := os.CpuController.MakeTestProcessBasic()
	os.AddProcessToScheduler(pcb)
	os.AddProcessToScheduler(pcb2)
	os.AddProcessToScheduler(pcb3)
	nextPcb := os.Scheduler.GetNextProcess()
	os.CpuController.SetPageTabletoMMU(nextPcb)
	go os.CPU[0].Run() // Run CPU in a separate goroutine
	os.cpuIsRunning = true
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
	if os.cpuIsRunning || !os.osIsRunning {
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

func (os *OS) ContextSwitch(cpu cpu.CPU) {
	currentProcess := os.Scheduler.GetRunningProcess()
	logger.Log.Println("currentProcess.Pid")
	logger.Log.Println(currentProcess.Pid)
	if currentProcess == nil {
		logger.Log.Println("No processes currently running.")
		return
	}
	nextProcess := os.Scheduler.GetNextProcess()
	if nextProcess == nil {
		logger.Log.Println("No processes in the ready queue.")
		return
	}
	logger.Log.Println("nextProcess.Pid")
	logger.Log.Println(nextProcess.Pid)

	logger.Log.Println("Context switching from process", currentProcess.Pid)
	newPCB := os.Scheduler.GetNextProcess()
	os.SaveProcessState(currentProcess, cpu) // Saves the process state of pcb from cpu to pcb..
	os.CpuController.SetPageTabletoMMU(newPCB)
	os.SetNewProcessState(newPCB, cpu) // Sets the process state of pcb to cpu.
	cpu.Resume()
}

// Function to save the state of the current process
func (os *OS) SaveProcessState(pcb *processes.PCB, cpu cpu.CPU) {
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
func (os *OS) SetNewProcessState(pcb *processes.PCB, cpu cpu.CPU) {
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

func (os *OS) AddProcessToMap(pcb *processes.PCB) {
	id := pcb.Pid
	os.ProcessList[uint32(id)] = pcb
	return
}
func (os *OS) AddProcessToScheduler(pcb *processes.PCB) {
	os.Scheduler.AddProcess(pcb)
	os.AddProcessToMap(pcb)
}

func (os *OS) GetScheduler() scheduler.SchedulerInterface {
	scheduler := scheduler.NewScheduler()
	return scheduler
}
