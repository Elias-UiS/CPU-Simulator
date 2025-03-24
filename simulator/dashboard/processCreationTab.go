package dashboard

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/os"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var dynamicOpCodes []string

// ProcessCreationTab creates a UI for adding instructions and displaying the instruction list
func ProcessCreationTab(os *os.OS) fyne.CanvasObject {

	for _, opcode := range cpu.OpcodeNames {
		dynamicOpCodes = append(dynamicOpCodes, opcode)
	}

	// Create input fields
	opTypeOptions := []string{"0", "1"}
	opTypeSelect := widget.NewSelect(opTypeOptions, nil)
	opTypeSelect.PlaceHolder = "OpType"

	opCodeSelect := widget.NewSelect(dynamicOpCodes, nil)
	opCodeSelect.PlaceHolder = "OpCode"

	operandEntry := widget.NewEntry()
	operandEntry.SetPlaceHolder("Operand")

	processNameEntry := widget.NewEntry()
	processNameEntry.SetPlaceHolder("Process Name")

	// List binding for instruction list
	instructionList := binding.NewStringList()

	updateInstructionList := func() {
		var items []string
		for _, inst := range *os.ProcessController.InstructionList {
			items = append(items,
				fmt.Sprintf("Type: %d, Code: %d, Operand: %d", inst.OpType, inst.Opcode, inst.Operand))
		}
		instructionList.Set(items)
	}

	// Button to add instruction
	addButton := widget.NewButton("Add", func() {
		opType, err1 := strconv.Atoi(opTypeSelect.Selected)
		opCodeName := opCodeSelect.Selected
		opCode := cpu.OpcodeValues[opCodeName] // Get the integer value for selected OpCode
		operand, err3 := strconv.Atoi(operandEntry.Text)

		if err1 == nil && err3 == nil {
			os.ProcessController.AddInstructionToList(opType, opCode, operand)
			updateInstructionList()
			opTypeSelect.SetSelected("")
			opCodeSelect.SetSelected("")
			operandEntry.SetText("")
		}
	})

	// Button to create process
	createProcessButton := widget.NewButton("Create Process", func() {
		processName := processNameEntry.Text
		if processName != "" && len(*os.ProcessController.InstructionList) > 0 {
			pcb := os.ProcessController.CreateProcessFromInstructionList(processName)
			os.AddProcessToSchedulerQueue(pcb)
			processNameEntry.SetText("")
			updateInstructionList()
		}
	})

	deleteInstructionListButton := widget.NewButton("Delete List", func() {
		os.ProcessController.DeleteInstructionList()
		updateInstructionList() // Optionally update instruction list or any UI state
	})

	// List widget to show instructions
	instructionListView := widget.NewListWithData(
		instructionList,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			str, _ := i.(binding.String).Get()
			o.(*widget.Label).SetText(str)
		},
	)

	instructionListContainer := container.NewVScroll(instructionListView)
	instructionListContainer.SetMinSize(fyne.NewSize(400, 300))

	topContainer := container.NewVBox(
		container.NewHBox(opTypeSelect, opCodeSelect),
		operandEntry,
	)
	gap := widget.NewLabel("")

	return container.NewVBox(
		topContainer,
		addButton,
		gap,
		processNameEntry,
		createProcessButton,
		gap,
		deleteInstructionListButton,
		instructionListContainer,
	)
}
