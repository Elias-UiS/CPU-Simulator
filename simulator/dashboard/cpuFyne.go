package dashboard

import (
	"CPU-Simulator/simulator/pkg/os"
	"time"

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

	// Create bindings for register values
	boundPC := binding.NewInt()
	boundAC := binding.NewInt()
	boundMAR := binding.NewInt()
	boundOpType := binding.NewInt()
	boundOpcode := binding.NewInt()
	boundOperand := binding.NewInt()
	boundIsInstruction := binding.NewBool()
	boundData := binding.NewInt()

	// Define register labels and corresponding bindings
	labels := []string{"PC:", "AC:", "MAR:", "OpType:", "Opcode:", "Operand:", "IsInstruction:", "Data:"}
	values := []binding.DataItem{
		binding.IntToString(boundPC),
		binding.IntToString(boundAC),
		binding.IntToString(boundMAR),
		binding.IntToString(boundOpType),
		binding.IntToString(boundOpcode),
		binding.IntToString(boundOperand),
		binding.BoolToString(boundIsInstruction),
		binding.IntToString(boundData),
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
	table.SetColumnWidth(0, 100) // Label column width
	table.SetColumnWidth(1, 150) // Value column width

	// Start a goroutine to update the values periodically
	go func() {
		for {
			// Update bindings with the latest register values
			boundPC.Set(cpu.Registers.PC)
			boundAC.Set(cpu.Registers.AC)
			boundMAR.Set(cpu.Registers.MAR)
			boundOpType.Set(cpu.Registers.IR.OpType)
			boundOpcode.Set(cpu.Registers.IR.Opcode)
			boundOperand.Set(cpu.Registers.IR.Operand)
			boundIsInstruction.Set(cpu.Registers.MDR.IsInstruction)
			boundData.Set(cpu.Registers.MDR.Data)

			time.Sleep(10 * time.Millisecond) // Adjust refresh rate as needed
		}
	}()

	// Return the table inside a container
	return container.NewBorder(
		widget.NewLabelWithStyle("CPU Registers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil, // No bottom, left, or right content
		table,
	)
}
