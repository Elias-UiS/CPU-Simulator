package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/processes"
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Function to create the PCB UI
func CreatePCBUI(processTable *processes.ProcessTable) fyne.CanvasObject {
	if processTable.ProcessMap == nil {
		logger.Log.Println("processList is empty")
		return widget.NewLabel("Error: Process list is nil")
	}

	// Left side: List of processes
	processListWidget := widget.NewList(
		func() int {
			return len(processTable.ProcessMap) // Number of processes
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Process ID | Process Name") // Updated header
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			keys := getProcessKeys(*processTable)
			selectedProcess := (processTable.ProcessMap)[keys[i]]
			// Displaying Process ID and Process Name
			obj.(*widget.Label).SetText(fmt.Sprintf("%d | %s", selectedProcess.Pid, selectedProcess.Name))
		},
	)

	// Right side: Process Details
	infoLabel := widget.NewLabel("Select a process to see details")

	// Handle selection
	processListWidget.OnSelected = func(id widget.ListItemID) {
		keys := getProcessKeys(*processTable)
		selectedProcess := (processTable.ProcessMap)[keys[id]]
		infoLabel.SetText(formatPCBDetails(selectedProcess))
	}

	// Auto-refresh mechanism
	go autoRefreshProcessList(processListWidget, *processTable)

	// Layout: Split into two sections
	split := container.NewHSplit(
		processListWidget,            // Left: List of processes
		container.NewVBox(infoLabel), // Right: PCB details
	)
	split.SetOffset(0.3) // Adjust ratio

	return split
}

// Helper: Get process keys
func getProcessKeys(processTable processes.ProcessTable) []int {
	keys := make([]int, 0, len(processTable.ProcessMap))
	for k := range processTable.ProcessMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] }) // Ensure consistent order
	return keys
}

// Helper: Format PCB details
func formatPCBDetails(pcb *processes.PCB) string {
	return fmt.Sprintf("PID: %d\nName: %s\nState: %s\nPriority: %d\nPC: %d\nAC: %d\nNextFreeCodeAddress: %d\nPageAmount: %d\n",
		pcb.Pid, pcb.Name, pcb.State.String(), pcb.Priority, pcb.ProcessState.PC, pcb.ProcessState.AC, pcb.NextFreeCodeAddress, pcb.PageAmount,
	)
}

// Auto-refresh function
func autoRefreshProcessList(list *widget.List, processTable processes.ProcessTable) {
	// Only refresh if there is any change in the process list
	var id int
	for k := range processTable.ProcessMap {
		id = k
		break
	}
	logger.Log.Println("Info: autoRefreshProcessList() | ID: %d", id)
	lastLen := len(processTable.ProcessMap)
	for {
		time.Sleep(1 * time.Second) // Refresh every second
		for _, pcb := range processTable.ProcessMap {
			logger.Log.Printf("Info: autoRefreshProcessList() %d", pcb.Pid)
		}
		// Only refresh if the length of the process list has changed
		if len(processTable.ProcessMap) != lastLen {
			list.Refresh()
			lastLen = len(processTable.ProcessMap)
		}
	}
}
