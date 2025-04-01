package dashboard

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/systemState"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// setupCpuTab creates a table-like UI for the CPU registers
func setupSystemStateTab(pubSub *systemState.PubSub[systemState.State]) fyne.CanvasObject {

	// Define register labels and corresponding bindings
	labels := []string{"Name:", "PC:", "AC:", "MAR:", "OpType:", "Opcode:", "Operand:", "IsInstruction:", "Data:", "Instruction Count:", "SP:"}
	values := []binding.DataItem{
		bindings.NameBinding,
		binding.IntToString(bindings.SystemStatePcBinding),
		binding.IntToString(bindings.SystemStateAcBinding),
		binding.IntToString(bindings.SystemStateMarBinding),
		binding.IntToString(bindings.SystemStateInstructionOpTypeBinding),
		binding.IntToString(bindings.SystemStateInstructionOpCodeBinding),
		binding.IntToString(bindings.SystemStateInstructionOperandBinding),
		binding.BoolToString(bindings.SystemStateMdrIsInstructionBinding),
		binding.IntToString(bindings.SystemStateMdrDataBinding),
		binding.IntToString(bindings.SystemStateInstructionCount),
		binding.IntToString(bindings.SystemStateSpBinding),
	}

	// Create the table
	table := widget.NewTable(
		func() (int, int) { return len(labels), 2 }, // Rows x Columns
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Create an empty label for each cell
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			if id.Col == 0 {
				// Set register labels in the first column (bold)
				obj.(*widget.Label).SetText(labels[id.Row])
				obj.(*widget.Label).TextStyle.Bold = true
			} else {
				// Bind register values in the second column
				obj.(*widget.Label).Bind(values[id.Row].(binding.String))
			}
		},
	)

	// Set table column widths
	table.SetColumnWidth(0, 150) // Label column width
	table.SetColumnWidth(1, 150) // Value column width
	go UpdateState(pubSub)
	// Return the table inside a container
	return container.NewBorder(
		widget.NewLabelWithStyle("CPU Registers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil, // No bottom, left, or right content
		table,
	)
}

func UpdateState(pubSub *systemState.PubSub[systemState.State]) {
	channel := pubSub.Subscribe()
	var state systemState.State
	loop := 1
	for {
		logger.Log.Println("INFO: StateTab loop", loop)
		state = <-channel
		logger.Log.Println("StateTab AC from channel: ", state.Registers.AC)
		bindings.MdrIsInstructionBinding.Set(state.Registers.MDR.IsInstruction)
		bindings.MdrInstructionOpTypeBinding.Set(state.Registers.MDR.Instruction.OpType)
		bindings.MdrInstructionOpCodeBinding.Set(state.Registers.MDR.Instruction.Opcode)
		bindings.MdrInstructionOperandBinding.Set(state.Registers.MDR.Instruction.Operand)
		bindings.MdrDataBinding.Set(state.Registers.MDR.Data)

		bindings.InstructionOpTypeBinding.Set(state.Registers.IR.OpType)
		bindings.InstructionOpCodeBinding.Set(state.Registers.IR.Opcode)
		bindings.InstructionOperandBinding.Set(state.Registers.IR.Operand)

		bindings.MarBinding.Set(state.Registers.MAR)
		bindings.AcBinding.Set(state.Registers.AC)
		bindings.PcBinding.Set(state.Registers.PC)
		loop += 1
	}
}
