package memory

import "fmt"

// kan endre til å lagre bits?
type PTE struct { // Page table entry
	//Valid       bool   // Page is valid in memory
	FrameNumber uint32 // Physical frame number
	//Read        bool   // Read permission
	//Write       bool   // Write permission
	//Execute     bool   // Execute permission
	//Dirty       bool   // Modified flag
	//Referenced  bool   // Recently accessed flag
	//SwapLoc     uint32 // Location on disk if swapped
}

type PCB struct {
	Pid            int   // id of the process
	PageTable      []PTE // index is the same as the virtual page number
	State          string
	ProgramCounter int
	StackPointer   int
}

type Controller struct {
	nextFreeID   int
	ProcessTable []*PCB
}

func (controller *Controller) MakeProcess() *PCB {
	pageNum := 2
	pcb := &PCB{
		Pid:            controller.nextFreeID,
		State:          "New",
		ProgramCounter: 0,
		StackPointer:   128,
		PageTable:      []PTE{}, // TODO change later
	}
	controller.addPages(pcb, pageNum)
	controller.ProcessTable = append(controller.ProcessTable, pcb)

	return pcb
}

func (controller *Controller) addPages(pcb *PCB, pageNum int) {
	list, err := FreelistObject.allocateFrame(pageNum)
	if err != nil {
		fmt.Println("Error: ProcessController | addPages\n", err)
		return
	}
	for i := range len(list) {
		pte := PTE{
			FrameNumber: uint32(list[i]),
		}
		index := len(pcb.PageTable)
		pcb.PageTable[index] = pte
	}
}

func (controller *Controller) deallocateFrameForProcess(pcb *PCB) {
	list := []int{}
	for i := range len(pcb.PageTable) {
		num := pcb.PageTable[i].FrameNumber
		list = append(list, int(num))
	}
	FreelistObject.deallocateFrame(list)
}

func createController() Controller {
	controller := Controller{0, []*PCB{}}
	return controller
}
