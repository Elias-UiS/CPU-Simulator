package os

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/memory"
	"CPU-Simulator/simulator/pkg/metrics"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/settings"
	"CPU-Simulator/simulator/pkg/translator"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Main struct
type Controller struct {
	nextFreeID       int
	ProcessTable     *processes.ProcessTable
	mmu              *memory.MMU
	freelist         *memory.FreeList
	InstructionList  *[]cpu.Instruction
	MemoryController *memory.MemoryController
}

func (controller *Controller) MakeProcess(instructionCount int) (*processes.PCB, error) {
	logger.Log.Println("INFO: MakeProcess()")
	instructionPages := (instructionCount*2 + settings.PageSize - 1) / settings.PageSize // gives needed pages for instructions
	stackPages := 2
	heapPages := 2
	pageNum := instructionPages + stackPages + heapPages + 2 // +2 for guard pages
	metrics := metrics.Metrics{
		CreationTime:        time.Now(),
		WaitingStartTime:    time.Now(),
		SimulationStartTime: time.Now(),
	}
	pcb := &processes.PCB{
		Pid:          controller.nextFreeID,
		Name:         "Default",
		State:        processes.Ready,
		ProcessState: processes.ProcessState{},
		PageTable: &memory.PageTable{
			Entries: make(map[int]*memory.PTE),
		},
		NextFreeCodeAddress: 0,
		PageAmount:          pageNum,
		Metrics:             metrics,
	}
	for i := range pageNum {
		pte := &memory.PTE{
			FrameNumber: -1,
		}
		pcb.PageTable.Entries[i] = pte
	}
	// allocates frame(s) for code
	list, err := controller.freelist.AllocateFrame(instructionPages)
	if err != nil {
		logger.Log.Println("INFO: MakeProcess() 4")
		fmt.Println("Error: ProcessController | addPages\n", err)
		return nil, fmt.Errorf("makeProcess failed\n", err)
	}
	for i := range len(list) {
		pcb.PageTable.Entries[i].FrameNumber = list[i]
		pcb.PageTable.Entries[i].Valid = true
	}
	// allocates frame for heap
	list, err = controller.freelist.AllocateFrame(1)
	if err != nil {
		logger.Log.Println("INFO: MakeProcess() 4")
		fmt.Println("Error: ProcessController | addPages\n", err)
		return nil, fmt.Errorf("makeProcess failed\n", err)
	}
	for i := range len(list) {
		pcb.PageTable.Entries[instructionPages+1].FrameNumber = list[i]
		pcb.PageTable.Entries[instructionPages+1].Valid = true
	}

	// allocates frame for stack
	list, err = controller.freelist.AllocateFrame(1)
	if err != nil {
		logger.Log.Println("INFO: MakeProcess() 4")
		fmt.Println("Error: ProcessController | addPages\n", err)
		return nil, fmt.Errorf("makeProcess failed\n", err)
	}
	for i := range len(list) {
		pcb.PageTable.Entries[pageNum-1].FrameNumber = list[i]
		pcb.PageTable.Entries[pageNum-1].Valid = true
	}
	address := translator.TranslateVPNandOffsetToAddress(pageNum-1, 8)
	pcb.ProcessState.SP = int(address) // sets the stack pointer to the one less than the bottom of stack
	controller.ProcessTable.AddProcessToTable(pcb)
	controller.nextFreeID += 1

	codeEnd := translator.TranslateVPNandOffsetToAddress(instructionPages-1, settings.PageSize-1)
	heapStart := translator.TranslateVPNandOffsetToAddress(instructionPages+1, 0)
	heapEnd := translator.TranslateVPNandOffsetToAddress(instructionPages+heapPages, settings.PageSize-1)
	stackStart := translator.TranslateVPNandOffsetToAddress(instructionPages+heapPages+stackPages, 0)
	stackEnd := translator.TranslateVPNandOffsetToAddress(instructionPages+heapPages+stackPages+1, settings.PageSize-1)

	pcb.Limits.CodeStart = 0
	pcb.Limits.CodeEnd = codeEnd

	pcb.Limits.HeapStart = heapStart
	pcb.Limits.HeapEnd = heapEnd

	pcb.Limits.StackStart = stackStart
	pcb.Limits.StackEnd = stackEnd
	return pcb, nil
}

func (controller *Controller) AllocateFrameToPage(pcb *processes.PCB, vpn int, pageNum int) error {
	list, err := controller.freelist.AllocateFrame(pageNum)
	if err != nil {
		fmt.Println("Error: ProcessController | addPages\n", err)
		return fmt.Errorf("ERROR: Not enough free frames.")
	}
	for _, value := range list {
		pcb.PageTable.Entries[vpn].FrameNumber = value
		pcb.PageTable.Entries[vpn].Valid = true
	}
	return nil
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
	for i := range len(controller.ProcessTable.ProcessMap) {
		if controller.ProcessTable.ProcessMap[i].Pid == id {
			return controller.ProcessTable.ProcessMap[i], nil
		}
	}
	return nil, fmt.Errorf("could not find PCB, with Pid: %v", id)
}

func (controller *Controller) StoreInstruction(instruction uint64, id int) {
	pcb, err := controller.FindPCB(id)
	if err != nil {
		fmt.Println(err)
	}
	vpn, offset := translator.TranslateAddressToVPNandOffset(int(pcb.NextFreeCodeAddress))
	frame := pcb.PageTable.Entries[vpn].FrameNumber
	logger.Log.Println("Frame Number: %d", frame)
	opcodeType := uint32(instruction >> 32)
	operand := uint32(instruction & 0xFFFFFFFF)
	physicalAddr := (uint32(frame) << 16) | uint32(offset) // First 16 bits: frame, Last 16 bits: offset

	logger.Log.Printf("Info: opcodeType: %d", opcodeType)
	logger.Log.Printf("Info: StoreInstruction(). Physical address: %d", physicalAddr)
	controller.MemoryController.Write(uint32(physicalAddr), opcodeType)
	if err != nil {
		logger.Log.Println(err)
	}
	offset += 1

	physicalAddr = (uint32(frame) << 16) | uint32(offset) // First 16 bits: frame, Last 16 bits: offset

	if err != nil {
		logger.Log.Println("Error: StoreInstruction failed, TranslateAddress()")
	}
	controller.MemoryController.Write(uint32(physicalAddr), operand)
	if err != nil {
		logger.Log.Println(err)
	}
	logger.Log.Print("Done with adding instruction")
	address := translator.FindNextInstructionAddress(settings.MemType, int(pcb.NextFreeCodeAddress+2))
	pcb.NextFreeCodeAddress = uint32(address)
	logger.Log.Print("Done with setting NextFreeCodeAddress")
}

func createController(mmu *memory.MMU, freelist *memory.FreeList, processTableStruct *processes.ProcessTable, memoryController *memory.MemoryController) *Controller {
	instructionList := &[]cpu.Instruction{}
	controller := Controller{
		nextFreeID:       0,
		ProcessTable:     processTableStruct,
		mmu:              mmu,
		freelist:         freelist,
		InstructionList:  instructionList,
		MemoryController: memoryController,
	}
	return &controller
}

func (controller *Controller) SetPageTabletoMMU(pcb *processes.PCB) {
	controller.mmu.SetPageTable(pcb.PageTable)
	fmt.Printf("Switched to Process %d\n", pcb.Pid)
}

func (controller *Controller) MakeTestProcessBasic() *processes.PCB {
	logger.Log.Println("INFO: MakeTestProcessBasic()")
	pcb, err := controller.MakeProcess(3)
	logger.Log.Printf("DEBUG: PageTable Entry count after creation of process: %d", len(pcb.PageTable.Entries))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pcb.Name = "Increment"
	instructionAdd := cpu.Instruction{0, cpu.ADD, 1}
	instructionStore := cpu.Instruction{0, cpu.STORE, 131072}
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
	pcb, err := controller.MakeProcess(3)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pcb.Name = "Increment 2"
	instructionAdd := cpu.Instruction{0, cpu.ADD, 10}
	instructionStore := cpu.Instruction{0, cpu.STORE, 131072}
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

func (controller *Controller) MakeTestProcessBasic3() *processes.PCB {
	logger.Log.Println("INFO: MakeTestProcessBasic3()")
	pcb, err := controller.MakeProcess(2)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pcb.Name = "Increment without jump"
	instructionAdd := cpu.Instruction{0, cpu.ADD, 10}
	instructionStore := cpu.Instruction{0, cpu.STORE, 131072}

	instructionAddBytes := instructionAdd.ToInt()
	instructionStoreBytes := instructionStore.ToInt()

	controller.StoreInstruction(instructionAddBytes, pcb.Pid)
	controller.StoreInstruction(instructionStoreBytes, pcb.Pid)
	fmt.Println(pcb)
	return pcb
}

func (controller *Controller) MakeTestProcessStackBasic() *processes.PCB {
	logger.Log.Println("INFO: MakeTestProcessStackBasic()")
	pcb, err := controller.MakeProcess(5)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pcb.Name = "Stack Test"
	instructionPush := cpu.Instruction{0, cpu.PUSH, 10}
	instructionPush2 := cpu.Instruction{0, cpu.PUSH, 25}
	instructionPush3 := cpu.Instruction{0, cpu.PUSH, 35}
	instructionPop := cpu.Instruction{0, cpu.POP, 0}
	instructionJump := cpu.Instruction{0, cpu.JUMP, 0}

	instructionPushytes := instructionPush.ToInt()
	instructionPush2Bytes := instructionPush2.ToInt()
	instructionPush3Bytes := instructionPush3.ToInt()
	instructionPopBytes := instructionPop.ToInt()
	instructionJumpBytes := instructionJump.ToInt()

	controller.StoreInstruction(instructionPushytes, pcb.Pid)
	controller.StoreInstruction(instructionPush2Bytes, pcb.Pid)
	controller.StoreInstruction(instructionPush3Bytes, pcb.Pid)
	controller.StoreInstruction(instructionPopBytes, pcb.Pid)
	controller.StoreInstruction(instructionJumpBytes, pcb.Pid)
	fmt.Println(pcb)
	return pcb
}

func (controller *Controller) AddInstructionToList(opType int, opCode int, operand int) {
	instruction := cpu.Instruction{opType, opCode, operand}
	*controller.InstructionList = append(*controller.InstructionList, instruction)
}

func (controller *Controller) CreateProcessFromInstructionList(name string, isFromFile bool) *processes.PCB {
	instructionCount := len(*controller.InstructionList)
	process, err := controller.MakeProcess(instructionCount)
	if err != nil {
		logger.Log.Printf("DEBUG: CreateProcessFromInstructionList()")
		return nil
	}

	process.Name = name

	for _, value := range *controller.InstructionList {
		instructionInBytes := value.ToInt()
		controller.StoreInstruction(instructionInBytes, process.Pid)
	}
	if !isFromFile {
		controller.WriteInstructionListToFile(name)
	}
	*controller.InstructionList = []cpu.Instruction{} // Clear the instruction list after creating the process

	return process
}

func (controller *Controller) DeleteInstructionList() {

	*controller.InstructionList = []cpu.Instruction{}
	return
}

func (controller *Controller) WriteInstructionListToFile(name string) error {

	timestamp := time.Now().Format("2006.01.02_1504")
	filepath := fmt.Sprintf("simulator/pkg/processes/processFiles/%s_%s.json", name, timestamp)

	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("ERROR: Could not create file: %v", err)
		return err
	}
	defer file.Close()

	// Write header
	_, err = fmt.Fprintln(file, "OpType\tOpCode\tOperand")
	if err != nil {
		return err
	}

	// Iterate over instructions and write them to the file
	for _, instruction := range *controller.InstructionList {

		_, err := fmt.Fprintf(file, "%d\t%d\t%d\n",
			instruction.OpType,
			instruction.Opcode,
			instruction.Operand)
		if err != nil {
			logger.Log.Printf("ERROR: Could not write instruction to file: %v", err)
			return err
		}
	}

	logger.Log.Println("Instruction list written to file:", filepath)
	return nil
}

func (controller *Controller) CreateProcessFromFile(filename string) *processes.PCB {
	logger.Log.Println("INFO: CreateProcessFromFile()", filename)
	filepath := fmt.Sprintf("simulator/pkg/processes/processFiles/%s", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("ERROR: Could not read file: %v", err)
	}

	// Split the file contents by newlines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	// Iterate through each line, skipping the header
	for i, line := range lines {
		if i == 0 {
			continue // Skip header row
		}

		// Split line into parts by the tab character
		parts := strings.Split(line, "\t")

		opType, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("ERROR: Invalid OPType format in line: %s", line)
			continue
		}

		opCode, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("ERROR: Invalid OPCode format in line: %s", line)
			continue
		}

		operand, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("ERROR: Invalid Operand format in line: %s", line)
			continue
		}
		controller.AddInstructionToList(opType, opCode, operand)
	}
	name := strings.Split(filename, "_")[0]
	pcb := controller.CreateProcessFromInstructionList(name, true)
	return pcb
}
