package dashboard

// // setupMemoryTest initializes the UI for memory visualization.
// // func setupProcessorTab(os *os.OS) fyne.CanvasObject {
// // 	// Running

// // 	// IN queue (queue list)
// // 	// Click to show information
// // 	// Also shows virtual memory os said processor
// // }

// import (
// 	"CPU-Simulator/simulator/pkg/memory"
// 	"CPU-Simulator/simulator/pkg/processes"
// 	"strconv"

// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/widget"
// )

// // setupMemoryTest initializes the UI for memory visualization.
// func setupProcessTab(pcb *map[uint32]*processes.PCB) fyne.CanvasObject {
// 	// Create a binding for the selected frame values
// 	frameValuesBinding := binding.NewStringList()
// 	// Right list: displays the bound frame values
// 	rightList := widget.NewListWithData(frameValuesBinding,
// 		func(item binding.DataItem, obj fyne.CanvasObject) {
// 			// We need to access the HBox container to get its objects
// 			hbox := obj.(*fyne.Container)
// 			label := hbox.Objects[0].(*widget.Label) // Access the Label

// 			// Bind the text of the label to the corresponding item in the frameValuesBinding
// 			label.Bind(item.(binding.String))

// 			// Set the OnTapped function for the button to copy the label text
// 		})

// 	// Left list: displays the frames available
// 	pcbListBinding := binding.NewStringList() // Binding for the left list
// 	updatepcbListBinding(pcbListBinding, pcb) // Set initial data for frame list

// 	// Create the frame list with bindings
// 	pcbList := widget.NewListWithData(pcbListBinding,
// 		func() fyne.CanvasObject {
// 			return widget.NewLabel("")
// 		},
// 		func(item binding.DataItem, obj fyne.CanvasObject) {
// 			obj.(*widget.Label).Bind(item.(binding.String))
// 		},
// 	)

// 	// Track which frame is currently selected
// 	selectedPcb := binding.NewInt() // Binding to store selected frame index

// 	// When a frame is selected, bind it to the right list
// 	pcbList.OnSelected = func(index int) {
// 		selectedPcb.Set(index) // Store selected frame index
// 		// Automatically update the right list with selected frame values
// 		updatePcbValuesBinding(index, frameValuesBinding, pcb)
// 	}

// 	// Top container: holds the button
// 	topContainer := container.NewVBox()

// 	// Main split: left for frame list, right for value list
// 	split := container.NewHSplit(pcbList, rightList)
// 	split.SetOffset(0.3)

// 	// Overall layout: buttons at top and then split container
// 	content := container.NewBorder(topContainer, nil, nil, nil, split)

// 	// Start periodic updates for the selected frame every 100ms
// 	go updateFrameValuesPeriodically(frameValuesBinding, selectedPcb, pcb)

// 	return content
// }

// // updateFrameListBinding updates the left list UI when the pcb entries change
// func updatePcbListBinding(pcbListBinding binding.StringList, pcb *map[uint32]*processes.PCB) {
// 	// Create a new list to hold the frame labels
// 	frameLabels := make([]string, len(*pcb))
// 	counter := 0
// 	for _, value := range *pcb {
// 		// Set the frame list labels based on memory frames
// 		frameLabels[counter] = value.Name + " " + strconv.Itoa(value.Pid)
// 		counter += 1
// 	}
// 	// Bind the frame labels data to the frame list directly
// 	pcbListBinding.Set(frameLabels)
// }

// // Convert []uint32 to []string
// func getStringList(data []uint32) []string {
// 	strData := make([]string, len(data))
// 	for i, v := range data {
// 		strData[i] = strconv.Itoa(int(v))
// 	}
// 	return strData
// }

// // updateFrameValuesBinding updates the right list UI when Frames is changed.
// func updatepcbValuesBinding(index int, frameValuesBinding binding.StringList, mem *memory.Memory) {
// 	// This will automatically update the values when the selected frame changes
// 	if index >= 0 && index < len(c.Frames) {
// 		frameValuesBinding.Set(getStringList(mem.Frames[index]))
// 	}
// }

// // updateFrameValuesPeriodically updates the right list UI every 100ms.
// func updatePcbValuesPeriodically(pcbValuesBinding binding.StringList, selectedPcb binding.Int, pcb *map[uint32]*processes.PCB) {
// 	for {
// 		// Sleep for 100ms before updating again
// 		time.Sleep(100 * time.Millisecond)

// 		// Get the currently selected frame index
// 		index, _ := selectedPcb.Get()
// 		if index >= 0 && index < len(*pcb) {
// 			// Automatically update the right list with selected frame values
// 			pcbValuesBinding.Set(getString(mem.Frames[index]))
// 		}
// 	}
// }
