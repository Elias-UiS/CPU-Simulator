package os

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/processes"
	"fmt"
)

// Main struct
type Controller struct {
	nextFreeID   int
	ProcessTable []*processes.PCB
	mmu          *memory.MMU
	freelist     *memory.FreeList
}

func (controller *Controller) MakeProcess() (*processes.PCB, error) {
	logger.Log.Println("INFO: MakeProcess()")
	pageNum := 3
	pcb := &processes.PCB{
		Pid:          controller.nextFreeID,
		Name:         "Default",
		State:        processes.New,
		ProcessState: processes.ProcessState{},
		PageTable: memory.PageTable{
			Entries: make(map[uint16]*memory.PTE),
		},
		NextFreeCodeAddress: 0,
	}
	logger.Log.Println("INFO: MakeProcess() 2")
	list, err := controller.freelist.AllocateFrame(pageNum)
	logger.Log.Println("INFO: MakeProcess() 3")
	if err != nil {
		logger.Log.Println("INFO: MakeProcess() 4")
		fmt.Println("Error: ProcessController | addPages\n", err)
		return nil, fmt.Errorf("makeProcess failed\n", err)
	}
	logger.Log.Println("INFO: MakeProcess() 5")
	for i := range len(list) - 1 {
		segmentType := memory.PageType(i)
		pte := &memory.PTE{
			FrameNumber: uint32(list[i]),
			Type:        segmentType,
		}
		pcb.PageTable.Entries[pcb.PageTable.NextFreeIndex] = pte
		pcb.PageTable.NextFreeIndex += 1
	}
	controller.ProcessTable = append(controller.ProcessTable, pcb)
	controller.nextFreeID += 1
	return pcb, nil
}

func (controller *Controller) AddPages(pcb *processes.PCB, pageNum int, segmentType int) {
	list, err := controller.freelist.AllocateFrame(pageNum)
	if err != nil {
		fmt.Println("Error: ProcessController | addPages\n", err)
		return
	}
	for i := range len(list) {
		pte := &memory.PTE{
			FrameNumber: uint32(i),
			Type:        memory.PageType(segmentType),
		}
		pcb.PageTable.Entries[pcb.PageTable.NextFreeIndex] = pte
		pcb.PageTable.NextFreeIndex += 1
	}
}

func (controller *Controller) DeallocateFrameForProcess(pcb *processes.PCB) {
	list := []int{}
	for index := range pcb.PageTable.Entries {
		num := pcb.PageTable.Entries[index].FrameNumber
		list = append(list, int(num))
	}
	controller.freelist.DeallocateFrame(list)
}

func (controller *Controller) FindPCB(id int) (*processes.PCB, error) {
	for i := range len(controller.ProcessTable) {
		if controller.ProcessTable[i].Pid == id {
			return controller.ProcessTable[i], nil
		}
	}
	return nil, fmt.Errorf("could not find PCB, with Pid: %v", id)
}

func (controller *Controller) StoreInstruction(instruction uint64, id int) {
	pcb, err := controller.FindPCB(id)
	if err != nil {
		fmt.Println(err)
	}
	opcodeType := uint32(instruction >> 32)
	operand := uint32(instruction & 0xFFFFFFFF)
	controller.mmu.Write(pcb.NextFreeCodeAddress, opcodeType) // need to change this to use PTE.type instead, so it can go to other pages.
	if err != nil {
		logger.Log.Println(err)
	}
	pcb.NextFreeCodeAddress += 1
	controller.mmu.Write(pcb.NextFreeCodeAddress, operand) // need to change this to use PTE.type instead, so it can go to other pages.
	if err != nil {
		logger.Log.Println(err)
	}
	pcb.NextFreeCodeAddress += 1

}

func createController(mmu *memory.MMU, freelist *memory.FreeList) *Controller {
	controller := Controller{
		0,
		[]*processes.PCB{},
		mmu,
		freelist,
	}
	return &controller
}

func (controller *Controller) SetPageTabletoMMU(pcb *processes.PCB) {
	controller.mmu.SetPageTable(&pcb.PageTable)
	fmt.Printf("Switched to Process %d\n", pcb.Pid)
}

// TODO: Self explanatory
func (controller *Controller) MakeTestProcessFromFile() {
	logger.Log.Println("INFO: MakeTestProcessFromFile()")
	pcb, err := controller.MakeProcess()
	if err != nil {
		fmt.Println(err)
		return
	}
	pcb.Name = "Manual"
	fmt.Println(pcb)
}

func (controller *Controller) MakeTestProcessBasic() *processes.PCB {
	logger.Log.Println("INFO: MakeTestProcessBasic()")
	pcb, err := controller.MakeProcess()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pcb.Name = "Increment"
	instructionAdd := cpu.Instruction{0, cpu.ADD, 1}
	instructionStore := cpu.Instruction{0, cpu.STORE, 65536}
	instructionJump := cpu.Instruction{0, cpu.JUMP, 0}

	instructionAddBytes := instructionAdd.ToInt()
	instructionStoreBytes := instructionStore.ToInt()
	instructionJumpBytes := instructionJump.ToInt()

	controller.StoreInstruction(instructionAddBytes, pcb.Pid)
	controller.StoreInstruction(instructionStoreBytes, pcb.Pid)
	controller.StoreInstruction(instructionJumpBytes, pcb.Pid)
	fmt.Println(pcb)
	return pcb
}

func (controller *Controller) MakeTestProcessBasic2() *processes.PCB {
	logger.Log.Println("INFO: MakeTestProcessBasic()")
	pcb, err := controller.MakeProcess()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pcb.Name = "Increment"
	instructionAdd := cpu.Instruction{0, cpu.ADD, 10}
	instructionStore := cpu.Instruction{0, cpu.STORE, 65536}
	instructionJump := cpu.Instruction{0, cpu.JUMP, 0}

	instructionAddBytes := instructionAdd.ToInt()
	instructionStoreBytes := instructionStore.ToInt()
	instructionJumpBytes := instructionJump.ToInt()

	controller.StoreInstruction(instructionAddBytes, pcb.Pid)
	controller.StoreInstruction(instructionStoreBytes, pcb.Pid)
	controller.StoreInstruction(instructionJumpBytes, pcb.Pid)
	fmt.Println(pcb)
	return pcb
}
