package processes

import (
	"CPU-Simulator/simulator/pkg/memory"
	"fmt"
)

// make processes

func MakeTestProcess(controller *memory.Controller, mmu *memory.MMU) {
	pcb := controller.MakeProcess()
	fmt.Println(pcb)
}
