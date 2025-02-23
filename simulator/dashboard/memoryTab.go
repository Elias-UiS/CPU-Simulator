package dashboard

import (
	"CPU-Simulator/simulator/pkg/memory"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// setupMemoryTest initializes the UI for memory visualization.
func setupMemoryTab(mem *memory.Memory) fyne.CanvasObject {
	Clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	// Create a binding for the selected frame values
	frameValuesBinding := binding.NewStringList()

	// Right list: displays the bound frame values
	rightList := widget.NewListWithData(frameValuesBinding,
		func() fyne.CanvasObject {
			// Each item will have a Label and a Button in an HBox
			return container.NewHBox(widget.NewLabel(""), widget.NewButton("Copy", nil))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			// We need to access the HBox container to get its objects
			hbox := obj.(*fyne.Container)
			label := hbox.Objects[0].(*widget.Label)   // Access the Label
			button := hbox.Objects[1].(*widget.Button) // Access the Button

			// Bind the text of the label to the corresponding item in the frameValuesBinding
			label.Bind(item.(binding.String))

			// Set the OnTapped function for the button to copy the label text
			button.OnTapped = func() {
				Clipboard.SetContent(label.Text) // Copy the label text to the clipboard
			}
		})

	// Left list: displays the frames available
	frameListBinding := binding.NewStringList()   // Binding for the left list
	updateFrameListBinding(frameListBinding, mem) // Set initial data for frame list

	// Create the frame list with bindings
	frameList := widget.NewListWithData(frameListBinding,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(item.(binding.String))
		},
	)

	// Track which frame is currently selected
	selectedFrame := binding.NewInt() // Binding to store selected frame index

	// When a frame is selected, bind it to the right list
	frameList.OnSelected = func(index int) {
		selectedFrame.Set(index) // Store selected frame index
		// Automatically update the right list with selected frame values
		updateFrameValuesBinding(index, frameValuesBinding, mem)
	}

	// Top container: holds the button
	topContainer := container.NewVBox()

	// Main split: left for frame list, right for value list
	split := container.NewHSplit(frameList, rightList)
	split.SetOffset(0.3)

	// Overall layout: buttons at top and then split container
	content := container.NewBorder(topContainer, nil, nil, nil, split)

	// Start periodic updates for the selected frame every 100ms
	go updateFrameValuesPeriodically(frameValuesBinding, selectedFrame, mem)

	return content
}

// updateFrameListBinding updates the left list UI when the memory frames change
func updateFrameListBinding(frameListBinding binding.StringList, mem *memory.Memory) {
	// Create a new list to hold the frame labels
	frameLabels := make([]string, len(mem.Frames))
	for i := range mem.Frames {
		// Set the frame list labels based on memory frames
		frameLabels[i] = "Frame " + strconv.Itoa(i)
	}
	// Bind the frame labels data to the frame list directly
	frameListBinding.Set(frameLabels)
}

// Convert []uint32 to []string
func getStringList(data []uint32) []string {
	strData := make([]string, len(data))
	for i, v := range data {
		strData[i] = strconv.Itoa(int(v))
	}
	return strData
}

// updateFrameValuesBinding updates the right list UI when Frames is changed.
func updateFrameValuesBinding(index int, frameValuesBinding binding.StringList, mem *memory.Memory) {
	// This will automatically update the values when the selected frame changes
	if index >= 0 && index < len(mem.Frames) {
		frameValuesBinding.Set(getStringList(mem.Frames[index]))
	}
}

// updateFrameValuesPeriodically updates the right list UI every 100ms.
func updateFrameValuesPeriodically(frameValuesBinding binding.StringList, selectedFrame binding.Int, mem *memory.Memory) {
	for {
		// Sleep for 100ms before updating again
		time.Sleep(100 * time.Millisecond)

		// Get the currently selected frame index
		index, _ := selectedFrame.Get()
		if index >= 0 && index < len(mem.Frames) {
			// Automatically update the right list with selected frame values
			frameValuesBinding.Set(getStringList(mem.Frames[index]))
		}
	}
}
