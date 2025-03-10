package dashboard

import (
	"CPU-Simulator/simulator/pkg/os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// TestBinding creates a UI that automatically updates when os.Test changes
func TestBinding(os *os.OS) fyne.CanvasObject {
	// Step 1: Create a binding for the os.Test value
	boundTest := binding.NewInt()
	boundTest.Set(os.Test)
	// Step 2: Bind the value to the label's data
	label := widget.NewLabelWithData(binding.IntToString(boundTest))
	go func() {
		for {
			// Simulate change in os.Test for demonstration purposes
			// Replace this with actual logic to monitor changes in os.Test
			boundTest.Set(os.Test) // Updates the label when os.Test changes
		}
	}()

	return label
}
