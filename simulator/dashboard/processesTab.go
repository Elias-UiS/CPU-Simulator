package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/os"
	"CPU-Simulator/simulator/pkg/processes"
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var selectedProcess *processes.PCB
var virtualPagesListView *widget.List
var virtualPagesListViewContainer *container.Scroll
var virtualPagesEntriesView *widget.List
var listOfMemoryFrames map[int][]uint32
var pageInfoLabel *widget.Label

// Function to create the PCB UI
func CreatePCBUI(os *os.OS, processTable *processes.ProcessTable) fyne.CanvasObject {
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
			selectedProcess = (processTable.ProcessMap)[keys[i]]
			// Displaying Process ID and Process Name
			obj.(*widget.Label).SetText(fmt.Sprintf("%d | %s", selectedProcess.Pid, selectedProcess.Name))
		},
	)

	// Right side: Process Details
	infoLabel := widget.NewLabel("Select a process to see details")
	pageInfoLabel = widget.NewLabel("Select a page to see details")

	virtualPagesListView = widget.NewList(
		func() int {
			return 0 // Show one entry when no process is selected
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("No process selected")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText("No process selected")
		},
	)
	virtualPagesListViewContainer = container.NewVScroll(virtualPagesListView)
	virtualPagesListViewContainer.SetMinSize(fyne.NewSize(400, 300))
	// Handle selection
	processListWidget.OnSelected = func(id widget.ListItemID) {
		logger.Log.Println("Selected process ID:", id)
		keys := getProcessKeys(*processTable)
		selectedProcess = (processTable.ProcessMap)[keys[id]]
		infoLabel.SetText(formatPCBDetails(selectedProcess))
		updateVirtualPagesListView(os)

	}

	// Auto-refresh mechanism
	go autoRefreshProcessList(processListWidget, *processTable)

	// Layout: Split into two sections
	split2 := container.NewHSplit(
		container.NewVBox(infoLabel, virtualPagesListViewContainer),
		pageInfoLabel,
	)
	split := container.NewHSplit(
		processListWidget, // Left: List of processes
		split2,
	)
	//split.SetOffset(0.3) // Adjust ratio

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
	return "PID: " + fmt.Sprintf("%d", pcb.Pid) +
		"\nName: " + pcb.Name +
		"\nState: " + pcb.State.String() +
		"\nPriority: " + fmt.Sprintf("%d", pcb.Priority) +
		"\nPC: " + fmt.Sprintf("%d", pcb.ProcessState.PC) +
		"\nAC: " + fmt.Sprintf("%d", pcb.ProcessState.AC) +
		"\nNextFreeCodeAddress: " + fmt.Sprintf("%d", pcb.NextFreeCodeAddress) +
		"\nPageAmount: " + fmt.Sprintf("%d", pcb.PageAmount) +
		"\n" +
		"\nCPU Time: " + fmt.Sprintf("%.2f", pcb.Metrics.CpuTime.Seconds()) +
		"\nResponse Time: " + fmt.Sprintf("%.2f", pcb.Metrics.ResponseTime.Seconds()) +
		"\nTurnaround Time: " + fmt.Sprintf("%.2f", pcb.Metrics.TurnaroundTime.Seconds()) +
		"\nWaiting Time: " + fmt.Sprintf("%.2f", pcb.Metrics.WaitingTime.Seconds())
}

// Auto-refresh function
func autoRefreshProcessList(list *widget.List, processTable processes.ProcessTable) {
	// Only refresh if there is any change in the process list
	lastLen := len(processTable.ProcessMap)
	for {
		time.Sleep(1 * time.Second) // Refresh every second
		// Only refresh if the length of the process list has changed
		if len(processTable.ProcessMap) != lastLen {
			list.Refresh()
			lastLen = len(processTable.ProcessMap)
		}
	}
}

func updateVirtualPagesListView(os *os.OS) {
	if selectedProcess == nil {
		return
	}

	// Clear previous memory frames
	listOfMemoryFrames = make(map[int][]uint32)

	// Create a list of virtual page numbers (VPN)
	vpnList := []uint32{}
	for i := 0; i < selectedProcess.PageAmount; i++ {
		entry := selectedProcess.PageTable.Entries[i]
		vpnList = append(vpnList, uint32(entry.FrameNumber))
		if entry.Valid == false {
			continue // Skip invalid entries
		}
		listOfMemoryFrames[i] = os.Memory.Frames[entry.FrameNumber] // Get the memory frame for the entry
		// Append the frame number to the list

	}

	logger.Log.Println("listOfMemoryFrames:", len(listOfMemoryFrames))

	// Initialize virtualPagesEntriesView (List of entries within selected frame)
	virtualPagesEntriesView = widget.NewList(
		func() int {
			if len(listOfMemoryFrames) > 0 {
				// Return the number of entries in the selected memory frame
				return len(vpnList) // Assuming the first frame is selected
			}
			return 0 // No entries if no memory frames
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("No entries available")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			if len(listOfMemoryFrames) > 0 {
				selectedFrame := listOfMemoryFrames[0]
				// Update the label with the entry data
				obj.(*widget.Label).SetText(fmt.Sprintf("Entry: %d", selectedFrame[i])) // Display each entry
			}
		},
	)

	// Update virtualPagesListView with vpnList data
	virtualPagesListView.Length = func() int {
		return len(vpnList) // Dynamic size based on vpnList
	}

	virtualPagesListView.UpdateItem = func(i widget.ListItemID, obj fyne.CanvasObject) {
		// Update the list item with the correct details
		entry := selectedProcess.PageTable.Entries[i]
		obj.(*widget.Label).SetText(fmt.Sprintf("VPN: %d | Frame: %d | Valid: %t", i, entry.FrameNumber, entry.Valid))
	}

	// Handle selection in the virtualPagesListView
	virtualPagesListView.OnSelected = func(id widget.ListItemID) {
		logger.Log.Println("Selected VPN:", id)
		// When an item is selected, get the memory frame details
		if selectedProcess.PageTable.Entries[id].Valid == false {
			pageInfoLabel.SetText(fmt.Sprintf("Selected Page is not valid \n Either not allocated or is a guard page."))
			return

		}
		selectedFrame := listOfMemoryFrames[id]
		// Show details of the selected memory frame
		logger.Log.Println("Selected Memory Frame:", selectedFrame)
		updateVirtualPagesEntriesView(os, selectedFrame)
		// Here you can create a label or a dialog to display the details of the selected frame
		pageInfoLabel.SetText(fmt.Sprintf("Selected Memory Frame: %v", selectedFrame))
	}

	// Refresh both list views to show updated data
	virtualPagesListView.Refresh()
	virtualPagesEntriesView.Refresh() // Refresh entries view
}

func updateVirtualPagesEntriesView(os *os.OS, frame []uint32) {
	// Initialize the virtualPagesEntriesView (List of entries within the selected memory frame)
	virtualPagesEntriesView = widget.NewList(
		func() int {
			// Return the number of entries in the selected memory frame
			return len(frame)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("No entries available") // Default label if no entries
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			// Display each entry in the selected frame
			obj.(*widget.Label).SetText(fmt.Sprintf("Entry: %d", frame[i])) // Display each entry as a new line
		},
	)

	// Refresh the virtualPagesEntriesView to show updated data
	virtualPagesEntriesView.Refresh()
}
