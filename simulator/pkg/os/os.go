package os

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"fmt"
	"time"
)

type OS struct {
	CPU               []*cpu.CPU
	Memory            *memory.Memory
	MMU               *memory.MMU
	ProcessTable      *processes.ProcessTable
	FreeList          *memory.FreeList
	ProcessController *Controller
	Scheduler         scheduler.SchedulerInterface
	osIsRunning       bool
	cpuIsRunning      bool
	Test              int
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
		CPU:               []*cpu.CPU{cpuInstance},
		Memory:            mem,
		MMU:               mmu,
		ProcessTable:      processTableStruct,
		FreeList:          freeList,
		ProcessController: controller,
		osIsRunning:       false,
		cpuIsRunning:      false,
		Scheduler:         scheduler,
		Test:              10,
	}
}

func (os *OS) StartSimulation() {
	if os.osIsRunning {
		return
	}
	logger.Log.Println("Starting simulation...")
	os.osIsRunning = true

	pcb := os.ProcessController.MakeTestProcessBasic()
	pcb2 := os.ProcessController.MakeTestProcessBasic2()
	pcb3 := os.ProcessController.MakeTestProcessBasic()

	// os.AddProcessToProcessTable(pcb)
	// os.AddProcessToProcessTable(pcb2)
	// os.AddProcessToProcessTable(pcb3)

	os.AddProcessToSchedulerQueue(pcb)
	os.AddProcessToSchedulerQueue(pcb2)
	os.AddProcessToSchedulerQueue(pcb3)

	nextPcb := os.Scheduler.GetNextProcess()
	os.ProcessController.SetPageTabletoMMU(nextPcb)
	nextPcb.Metrics.CpuStartTime = time.Now()
	os.CPU[0].EventHandler = os.OnCPUCycle
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
	os.UpdateMetricsPause()
}

func (os *OS) ResumeSimulation() {
	if os.cpuIsRunning {
		return
	}
	logger.Log.Println("Testing bug here")
	fmt.Println("Resuming simulation...")
	for i := range len(os.CPU) {
		logger.Log.Println("Testing bug here 2")
		os.CPU[i].Resume()
	}
	logger.Log.Println("Testing bug here 3")
	if os.Scheduler.GetRunningProcess().Metrics.CpuStartTime.IsZero() {

		os.Scheduler.GetRunningProcess().Metrics.CpuStartTime = time.Now()
	}
	os.cpuIsRunning = true
	os.UpdateMetricsResume()
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
	logger.Log.Println("Performing context switch...")
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

	cpuTimeAdd := time.Now().Sub(currentProcess.Metrics.CpuStartTime)
	currentProcess.Metrics.CpuTime += cpuTimeAdd
	currentProcess.Metrics.WaitingStartTime = time.Now()

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
	os.ProcessController.SetPageTabletoMMU(nextProcess)
	for index := range len(nextProcess.PageTable.Entries) {
		logger.Log.Printf("Index: %d, Value: %d\n", index, nextProcess.PageTable.Entries[uint16(index)].FrameNumber)
	}
	os.SetNewProcessState(nextProcess, cpu) // Sets the process state of pcb to cpu.
	nextProcess.State = processes.Running
	if nextProcess.Metrics.ResponseTime == 0 {
		nextProcess.Metrics.ResponseTime = time.Now().Sub(nextProcess.Metrics.ArrivalTime)
	}
	if os.cpuIsRunning {
		nextProcess.Metrics.CpuStartTime = time.Now()
		nextProcess.Metrics.WaitingTime += time.Now().Sub(nextProcess.Metrics.WaitingStartTime)
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

	bindings.MdrIsInstructionBinding.Set(cpu.Registers.MDR.IsInstruction)
	bindings.MdrInstructionOpTypeBinding.Set(cpu.Registers.MDR.Instruction.OpType)
	bindings.MdrInstructionOpCodeBinding.Set(cpu.Registers.MDR.Instruction.Opcode)
	bindings.MdrInstructionOperandBinding.Set(cpu.Registers.MDR.Instruction.Operand)
	bindings.MdrDataBinding.Set(cpu.Registers.MDR.Data)

	bindings.InstructionOpTypeBinding.Set(cpu.Registers.IR.OpType)
	bindings.InstructionOpCodeBinding.Set(cpu.Registers.IR.Opcode)
	bindings.InstructionOperandBinding.Set(cpu.Registers.IR.Operand)

	bindings.MarBinding.Set(cpu.Registers.MAR)
	bindings.AcBinding.Set(cpu.Registers.AC)
	bindings.PcBinding.Set(cpu.Registers.PC)

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
	pcb.Metrics.ArrivalTime = time.Now()
}

func (os *OS) GetScheduler() scheduler.SchedulerInterface {
	scheduler := scheduler.NewScheduler()
	return scheduler
}

func (os *OS) TestNumer() {
	for {
		os.Test += 1
		bindings.SharedValue.Set(os.Test)
		time.Sleep(500 * time.Millisecond)
	}
}

func (os *OS) OnCPUCycle(cpu *cpu.CPU) {
	// This is run after every CPU cycle. The os can decide whether to intervene
	pcb := os.Scheduler.GetRunningProcess()
	pcb.Metrics.InstructionAmount += 1
	logger.Log.Println("OS: CPU cycle completed.")
	if cpu.InstructionCount >= 6 {
		logger.Log.Println("OS: Exceeded instruction count per process instance.")
		logger.Log.Println("OS: Performing context switch.")
		cpu.Pause()
		cpu.InstructionCount = 0
		go os.ContextSwitch(cpu)
	} else {
		logger.Log.Println("OS: Continuing execution.")
	}
}

func (os *OS) UpdateMetricsPause() {
	// This updates the metrics when the CPU is paused
	// Such that the metrics "pauses" as well
	runningPCB := os.Scheduler.GetRunningProcess()
	runningPCB.Metrics.CpuTime += time.Now().Sub(runningPCB.Metrics.CpuStartTime)
	runningPCB.Metrics.SimulationTime += time.Now().Sub(runningPCB.Metrics.SimulationStartTime)

	for _, pcb := range os.Scheduler.GetReadyQueue() {
		pcb.Metrics.WaitingTime += time.Now().Sub(pcb.Metrics.WaitingStartTime)
	}
}

func (os *OS) UpdateMetricsResume() {
	// This updates the metrics when the CPU is resumed
	// Such that the metrics "resume" as well

	runningPCB := os.Scheduler.GetRunningProcess()
	runningPCB.Metrics.CpuStartTime = time.Now()
	runningPCB.Metrics.SimulationStartTime = time.Now()

	for _, pcb := range os.Scheduler.GetReadyQueue() {
		pcb.Metrics.WaitingStartTime = time.Now()
	}
}
