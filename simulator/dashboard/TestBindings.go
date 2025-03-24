package dashboard

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// TestBinding creates a UI that automatically updates when os.Test changes
func TestBinding(os *os.OS) fyne.CanvasObject {

	label := widget.NewLabelWithData(binding.IntToString(bindings.SharedValue))
	return label
}
