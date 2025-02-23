package os

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"fmt"
)

// Main struct
type Controller struct {
	nextFreeID   int
	ProcessTable []*PCB
	mmu          *memory.MMU
}

type PCB struct {
	Pid                 int           // id of the process
	Name                string        // Name of process, only for showing in list.
	PageTable           []*memory.PTE // index is the same as the virtual page number
	State               string
	ProgramCounter      int
	StackPointer        int
	NextFreeCodeAddress uint32 // next address for the storing instructions
}

func (controller *Controller) MakeProcess() (*PCB, error) {
	logger.Log.Println("INFO: MakeProcess()")
	pageNum := 3
	pcb := &PCB{
		Pid:                 controller.nextFreeID,
		Name:                "Default",
		State:               "New",
		ProgramCounter:      0,
		StackPointer:        128,
		PageTable:           []*memory.PTE{}, // TODO change later
		NextFreeCodeAddress: 0,
	}
	logger.Log.Println("INFO: MakeProcess() 2")
	list, err := memory.FreelistObject.AllocateFrame(pageNum)
	logger.Log.Println("INFO: MakeProcess() 3")
	if err != nil {
		logger.Log.Println("INFO: MakeProcess() 4")
		fmt.Println("Error: ProcessController | addPages\n", err)
		return nil, fmt.Errorf("makeProcess failed\n", err)
	}
	logger.Log.Println("INFO: MakeProcess() 5")
	for i := range len(list) - 1 {
		logger.Log.Println("INFO: MakeProcess() loop")
		segmentType := memory.PageType(i)
		logger.Log.Println("INFO: MakeProcess() loop 2")
		pte := &memory.PTE{
			FrameNumber: uint32(list[i]),
			Type:        segmentType,
		}
		logger.Log.Println("INFO: MakeProcess() loop 3")
		pcb.PageTable = append(pcb.PageTable, pte)
		logger.Log.Println("INFO: MakeProcess() loop 5")
	}
	logger.Log.Println("INFO: MakeProcess() 6")
	controller.ProcessTable = append(controller.ProcessTable, pcb)
	logger.Log.Println("INFO: MakeProcess() 7")

	return pcb, nil
}

func (controller *Controller) AddPages(pcb *PCB, pageNum int, segmentType int) {
	list, err := memory.FreelistObject.AllocateFrame(pageNum)
	if err != nil {
		fmt.Println("Error: ProcessController | addPages\n", err)
		return
	}
	for i := range len(list) {
		pte := &memory.PTE{
			FrameNumber: uint32(i),
			Type:        memory.PageType(segmentType),
		}
		index := len(pcb.PageTable)
		pcb.PageTable[index] = pte
	}
}

func (controller *Controller) DeallocateFrameForProcess(pcb *PCB) {
	list := []int{}
	for i := range len(pcb.PageTable) {
		num := pcb.PageTable[i].FrameNumber
		list = append(list, int(num))
	}
	memory.FreelistObject.DeallocateFrame(list)
}

func (controller *Controller) FindPCB(id int) (*PCB, error) {
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

func createController(mmu *memory.MMU) Controller {
	controller := Controller{0, []*PCB{}, mmu}
	return controller
}

func (controller *Controller) ScheduleProcess(pcb *PCB) {
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

func (controller *Controller) MakeTestProcessBasic() *PCB {
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

func MakeController(mmu *memory.MMU) *Controller {
	controller := Controller{0, []*PCB{}, mmu}
	return &controller
}
