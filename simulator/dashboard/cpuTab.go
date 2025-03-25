package dashboard

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// setupCpuTab creates a table-like UI for the CPU registers
func setupCpuTab(os *os.OS) fyne.CanvasObject {
	cpu := os.GetCpu()

	// Ensure CPU is available before proceeding
	for cpu == nil {
		cpu = os.GetCpu()
	}

	// Define register labels and corresponding bindings
	labels := []string{"Name:", "PC:", "AC:", "MAR:", "OpType:", "Opcode:", "Operand:", "IsInstruction:", "Data:", "Instruction Count:", "SP:"}
	values := []binding.DataItem{
		bindings.NameBinding,
		binding.IntToString(bindings.PcBinding),
		binding.IntToString(bindings.AcBinding),
		binding.IntToString(bindings.MarBinding),
		binding.IntToString(bindings.InstructionOpTypeBinding),
		binding.IntToString(bindings.InstructionOpCodeBinding),
		binding.IntToString(bindings.InstructionOperandBinding),
		binding.BoolToString(bindings.MdrIsInstructionBinding),
		binding.IntToString(bindings.MdrDataBinding),
		binding.IntToString(bindings.InstructionCount),
		binding.IntToString(bindings.SpBinding),
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

	// Return the table inside a container
	return container.NewBorder(
		widget.NewLabelWithStyle("CPU Registers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil, // No bottom, left, or right content
		table,
	)
}
