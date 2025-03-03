package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/processes"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Function to create the PCB UI
func CreatePCBUI(processList *map[uint32]*processes.PCB) fyne.CanvasObject {
	if processList == nil {
		return widget.NewLabel("Error: Process list is nil")
	}

	// Left side: List of processes
	processListWidget := widget.NewList(
		func() int {
			return len(*processList) // Number of processes
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Process ID | Process Name") // Updated header
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			keys := getProcessKeys(processList)
			selectedProcess := (*processList)[keys[i]]
			// Displaying Process ID and Process Name
			obj.(*widget.Label).SetText(fmt.Sprintf("%d | %s", selectedProcess.Pid, selectedProcess.Name))
		},
	)

	// Right side: Process Details
	infoLabel := widget.NewLabel("Select a process to see details")

	// Handle selection
	processListWidget.OnSelected = func(id widget.ListItemID) {
		keys := getProcessKeys(processList)
		selectedProcess := (*processList)[keys[id]]
		infoLabel.SetText(formatPCBDetails(selectedProcess))
	}

	// Auto-refresh mechanism
	go autoRefreshProcessList(processListWidget, processList)

	// Layout: Split into two sections
	split := container.NewHSplit(
		processListWidget,            // Left: List of processes
		container.NewVBox(infoLabel), // Right: PCB details
	)
	split.SetOffset(0.3) // Adjust ratio

	return split
}

// Helper: Get process keys
func getProcessKeys(processList *map[uint32]*processes.PCB) []uint32 {
	keys := make([]uint32, 0, len(*processList))
	for k := range *processList {
		keys = append(keys, k)
	}
	return keys
}

// Helper: Format PCB details
func formatPCBDetails(pcb *processes.PCB) string {
	return fmt.Sprintf("PID: %d\nName: %s\nState: %d\nPriority: %d\nPC: %d\nAC: %d",
		pcb.Pid, pcb.Name, pcb.State, pcb.Priority, pcb.ProcessState.PC, pcb.ProcessState.AC,
	)
}

// Auto-refresh function
func autoRefreshProcessList(list *widget.List, processList *map[uint32]*processes.PCB) {
	// Only refresh if there is any change in the process list
	var id uint32
	for k := range *processList {
		id = k
		break
	}
	logger.Log.Println("IDDDDD")
	logger.Log.Println(id)
	lastLen := len(*processList)
	for {
		time.Sleep(1 * time.Second) // Refresh every second

		// Only refresh if the length of the process list has changed
		if len(*processList) != lastLen {
			list.Refresh()
			lastLen = len(*processList)
		}
	}
}
