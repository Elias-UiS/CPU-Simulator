package os

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"fmt"
)

type OS struct {
	CPU           []*cpu.CPU
	Memory        *memory.Memory
	MMU           *memory.MMU
	ProcessList   []*PCB
	FreeList      *memory.FreeList
	CpuController *Controller
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

	controller := MakeController(mmu)

	return &OS{
		CPU:           []*cpu.CPU{cpuInstance},
		Memory:        mem,
		MMU:           mmu,
		ProcessList:   []*PCB{},
		FreeList:      freeList,
		CpuController: controller,
	}
}

func (os *OS) StartSimulation() {
	logger.Log.Println("Starting simulation...")

	pcb := os.CpuController.MakeTestProcessBasic()
	os.CpuController.ScheduleProcess(pcb)

	go os.CPU[0].Run() // Run CPU in a separate goroutine
}

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

func (os *OS) ContextSwitch(pcb *PCB) {
	logger.Log.Println("Context switching to process", pcb.Pid)
	os.CpuController.ScheduleProcess(pcb)
}

func (os *OS) GetCpu() *cpu.CPU {
	if len(os.CPU) == 0 {
		return nil
	}

	return os.CPU[0]
}
