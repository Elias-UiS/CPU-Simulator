package os

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"CPU-Simulator/simulator/pkg/scheduler/schedulerFiles"
	"CPU-Simulator/simulator/pkg/settings"
	"fmt"
	"sync"
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
	OsIsRunning       bool
	cpuIsRunning      bool
	StepMode          bool
	MidContextSwitch  bool
	OSMutex           sync.RWMutex
	MemoryController  *memory.MemoryController
	Stopping          bool
}

func NewOS() *OS {
	// Initialize memory
	mem := memory.NewMemory() // Example: 512 frames, 4KB pages

	// Initialize MMU
	mmu := memory.NewMMU(mem)

	// Initialize MemoryController
	memoryController := memory.NewMemoryController(mem)

	// Initialize CPU
	cpuInstance := cpu.NewCPU(mmu, memoryController)

	// Initialize free list
	freeList := memory.NewFreeList()

	processTableStruct := processes.CreateProcessTable()

	scheduler := schedulerFiles.NewScheduler()

	// Initialize processController
	controller := createController(mmu, freeList, processTableStruct, memoryController)

	return &OS{
		CPU:               []*cpu.CPU{cpuInstance},
		Memory:            mem,
		MMU:               mmu,
		ProcessTable:      processTableStruct,
		FreeList:          freeList,
		ProcessController: controller,
		OsIsRunning:       false,
		cpuIsRunning:      false,
		Scheduler:         scheduler,
		StepMode:          false,
		MemoryController:  memoryController,
	}
}

func (os *OS) StartSimulation() {
	if !os.Stopping {
		for _, value := range os.ProcessTable.ProcessMap {
			os.Scheduler.AddProcess(value)
			value.Metrics.ArrivalTime = time.Now()
		}
	}
	os.Stopping = false
	if os.OsIsRunning {
		return
	}
	logger.Log.Println("Starting simulation...")
	os.OsIsRunning = true

	os.CPU[0].PageFaultHandler = os.PageFaultHandler
	os.CPU[0].InterruptHandler = os.InterruptHandler

	logger.Log.Println(os.Scheduler)
	nextPcb := os.Scheduler.GetNextProcess()
	os.ProcessController.SetPageTabletoMMU(nextPcb)
	os.SetNewProcessState(nextPcb, os.CPU[0])
	nextPcb.Metrics.CpuStartTime = time.Now()
	os.CPU[0].EventHandler = os.OnCPUCycle
	bindings.NameBinding.Set(nextPcb.Name)
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
	os.UpdateMetricsPause()
}

func (os *OS) ResumeSimulation() {

	if os.cpuIsRunning {
		return
	}

	fmt.Println("Resuming simulation...")
	for i := range len(os.CPU) {
		os.CPU[i].Resume()
	}
	if os.Scheduler.GetRunningProcess().Metrics.CpuStartTime.IsZero() {

		os.Scheduler.GetRunningProcess().Metrics.CpuStartTime = time.Now()
	}
	os.cpuIsRunning = true
	os.UpdateMetricsResume()
}

func (os *OS) StopSimulation() {
	fmt.Println("Stopping simulation...")
	os.Stopping = true
	done := os.CPU[0].Stop()
	logger.Log.Println("CPU stopped:", done)
	os.Reset()

}

func (os *OS) Reset() {
	fmt.Println("Resetting OS...")
	// Reset memory, MMU, CPU, and free list
	os.Memory = memory.NewMemory()
	os.MMU = memory.NewMMU(os.Memory)
	newMemoryController := memory.NewMemoryController(os.Memory)
	newCpuInstance := cpu.NewCPU(os.MMU, newMemoryController)
	os.Scheduler = schedulerFiles.NewScheduler()
	os.ProcessTable = processes.CreateProcessTable()
	os.FreeList = memory.NewFreeList()
	os.ProcessController = createController(os.MMU, os.FreeList, os.ProcessTable, newMemoryController)
	os.CPU = []*cpu.CPU{newCpuInstance}
	os.OsIsRunning = false
	fmt.Println("Finished resetting OS")
}

func (os *OS) ContextSwitch(cpu *cpu.CPU) {
	logger.Log.Println("Performing context switch...")
	os.MidContextSwitch = true
	if !cpu.IsPaused {
		cpu.Pause()
	}

	currentProcess := os.Scheduler.GetRunningProcess()

	logger.Log.Printf("currentProcess.Pid: %d", currentProcess.Pid)
	if currentProcess == nil {
		logger.Log.Println("No processes currently running.")
		os.MidContextSwitch = false
		return
	}

	cpuTimeAdd := time.Now().Sub(currentProcess.Metrics.CpuStartTime)
	currentProcess.Metrics.CpuTime += cpuTimeAdd
	if currentProcess.State == processes.Terminated {
		currentProcess.Metrics.TurnaroundTime = time.Now().Sub(currentProcess.Metrics.ArrivalTime)
	}
	currentProcess.Metrics.WaitingStartTime = time.Now()
	if currentProcess.State != processes.Terminated {
		currentProcess.State = processes.Ready
	}
	nextProcess := os.Scheduler.GetNextProcess()

	if nextProcess == nil {
		logger.Log.Println("No processes in the ready queue.")
		os.MidContextSwitch = false
		return
	}
	logger.Log.Println("Context switching from process", currentProcess.Pid)
	logger.Log.Println("Context switching to process", nextProcess.Pid)
	os.SaveProcessState(currentProcess, cpu) // Saves the process state of pcb from cpu to pcb..
	os.ProcessController.SetPageTabletoMMU(nextProcess)
	for index := range len(nextProcess.PageTable.Entries) {
		logger.Log.Printf("Index: %d, Value: %d\n", index, nextProcess.PageTable.Entries[index].FrameNumber)
	}
	os.SetNewProcessState(nextProcess, cpu) // Sets the process state of pcb to cpu.
	nextProcess.State = processes.Running
	if nextProcess.Metrics.ResponseTime == 0 {
		nextProcess.Metrics.ResponseTime = time.Now().Sub(nextProcess.Metrics.ArrivalTime)
	}
	if os.cpuIsRunning {
		nextProcess.Metrics.CpuStartTime = time.Now()
		nextProcess.Metrics.WaitingTime += time.Now().Sub(nextProcess.Metrics.WaitingStartTime)
		cpu.InstructionCount = 0
		if !os.StepMode {
			cpu.Resume()
		}
	}
	os.MidContextSwitch = false

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
	pcb.ProcessState.SP = cpu.Registers.SP
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
	cpu.Registers.SP = pcb.ProcessState.SP

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
	bindings.SpBinding.Set(cpu.Registers.SP)

	bindings.NameBinding.Set(pcb.Name)
	logger.Log.Println("SetNewProcessState END")
	return
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

func (os *OS) OnCPUCycle(cpu *cpu.CPU) {
	// This is run after every CPU cycle. The os can decide whether to intervene
	if os.Stopping {
		logger.Log.Println("Stop OnCPUCycle, simulation stopping.")
		return
	}
	pcb := os.Scheduler.GetRunningProcess()
	pcb.Metrics.InstructionAmount += 1
	logger.Log.Println("OS: CPU cycle completed.")

	pc := cpu.Registers.PC
	phys, errStruct := os.MMU.TranslateAddress(uint32(pc))
	if errStruct != nil {
		logger.Log.Println("Error: OS | OnCPUCycle | TranslateAddress\n", errStruct.Text)
		pcb.State = processes.Terminated
		cpu.Pause()
		os.ContextSwitch(cpu)
	}
	instruction, err := os.MemoryController.Read(uint32(phys))

	pc2 := cpu.Registers.PC + 1
	phys2, errStruct := os.MMU.TranslateAddress(uint32(pc2))
	if errStruct != nil {
		logger.Log.Println("Error: OS | OnCPUCycle | TranslateAddress\n", errStruct.Text)
		pcb.State = processes.Terminated
		cpu.Pause()
		os.ContextSwitch(cpu)
	}
	instruction2, err := os.MemoryController.Read(uint32(phys2))
	if err != nil {
		logger.Log.Println("Error: OS | OnCPUCycle | Read\n", err)
		pcb.State = processes.Terminated
		logger.Log.Println("OS: Process terminated.")
		cpu.InstructionCount = 0
		cpu.Pause()
		os.ContextSwitch(cpu)
		return
	}
	if instruction == 0 && instruction2 == 0 {
		pcb.State = processes.Terminated
		logger.Log.Println("OS: Process terminated.")
		cpu.InstructionCount = 0
		cpu.Pause()
		os.ContextSwitch(cpu)
	}
	if cpu.InstructionCount >= settings.InstructionLimitPerRun {
		logger.Log.Println(cpu.InstructionCount)
		logger.Log.Println(settings.InstructionLimitPerRun)
		logger.Log.Println("OS: Exceeded instruction count per process instance.")
		logger.Log.Println("OS: Performing context switch.")
		cpu.Pause()
		cpu.InstructionCount = 0
		bindings.InstructionCount.Set(cpu.InstructionCount)
		os.ContextSwitch(cpu)
	} else {
		logger.Log.Println("OS: Continuing execution.")
		if os.StepMode {
			os.cpuIsRunning = false
			cpu.Pause()
		}
	}
	cpu.Resume()

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

func (os *OS) PageFaultHandler(cpu *cpu.CPU, address uint32, vpn int) error {
	logger.Log.Println("ERROR: PageFaultHandler()")
	pcb := os.Scheduler.GetRunningProcess()
	if pcb.Limits.CodeStart <= address && address <= pcb.Limits.CodeEnd {
		logger.Log.Println("ERROR: PageFaultHandler() 2")
		err := os.ProcessController.AllocateFrameToPage(pcb, vpn, 1)
		pcb.State = processes.Terminated
		os.ContextSwitch(cpu)
		return err
		// This should not happen as code is mapped wholly at creation. The code segment is valid, and won't cause a pagefault.
	}
	logger.Log.Println("ERROR: PageFaultHandler() 3")
	if pcb.Limits.HeapStart <= address && address <= pcb.Limits.HeapEnd {
		logger.Log.Println("ERROR: PageFaultHandler() 4")
		stackPointer := cpu.Registers.SP
		if uint32(stackPointer) == address {
			pcb.State = processes.Terminated
			logger.Log.Println("Pausing +1")
			cpu.Pause()
			os.ContextSwitch(cpu)
			return fmt.Errorf("ERROR: Stackpointer pointing outside stack")
		}
		err := os.ProcessController.AllocateFrameToPage(pcb, vpn, 1)
		if err != nil {
			return err
		}
	}
	if pcb.Limits.StackStart <= address && address <= pcb.Limits.StackEnd {
		logger.Log.Println("ERROR: PageFaultHandler() 5")
		err := os.ProcessController.AllocateFrameToPage(pcb, vpn, 1)
		if err != nil {
			return err
		}
		logger.Log.Println("ERROR: PageFaultHandler() 6")
		cpu.Resume()
		return nil
	}
	logger.Log.Println("ERROR: PageFaultHandler() 7")
	pcb.State = processes.Terminated
	os.ContextSwitch(cpu)
	return nil
}

func (os *OS) InterruptHandler(cpu *cpu.CPU) error {
	logger.Log.Println("INFO: InterruptHandler()")
	code := cpu.InterruptCode
	switch code {
	case 1:
		// Handle pageFault interrupt
		address := cpu.Registers.FAR
		vpn := cpu.Registers.FVR
		err := os.PageFaultHandler(cpu, uint32(address), vpn)
		if err != nil {
			logger.Log.Println(err)
		}

	case 2:

	default:
		logger.Log.Println("ERROR: InterruptHandler() - Unknown interrupt code")
		process := os.Scheduler.GetRunningProcess()
		process.State = processes.Terminated
		os.ContextSwitch(cpu)
	}

	return nil

}
