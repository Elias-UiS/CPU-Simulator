package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Function to create the Scheduler UI
func CreateSchedulerTab(scheduler scheduler.SchedulerInterface) fyne.CanvasObject {
	// Create the title label (Bold and Left-aligned)
	titleLabel := widget.NewLabelWithStyle("Ready Queue", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	// Create the container that will hold the ready queue list or empty message
	emptyLabel := widget.NewLabel("Queue is empty")
	readyQueueContainer := container.NewBorder(titleLabel, emptyLabel, nil, nil, nil)

	runningTitleLabel := widget.NewLabel("Running Process")
	runningProcessContainer := container.NewVBox(
		runningTitleLabel,                     // Running Process title
		widget.NewLabel("No running process"), // Process details placeholder
	)

	contentContainer := container.NewGridWithRows(2, runningProcessContainer, readyQueueContainer)

	// Initial call to load the ready queue (empty or not)
	updateReadyQueueContent(contentContainer, scheduler)

	// Auto-refresh mechanism
	go autoRefreshReadyQueue(contentContainer, scheduler)

	return contentContainer
}

// Function to update the content based on the current state of the ready queue
func updateReadyQueueContent(contentContainer *fyne.Container, scheduler scheduler.SchedulerInterface) {
	readyQueue := scheduler.GetReadyQueue()

	// If the ready queue is empty
	if len(readyQueue) == 0 || readyQueue == nil {
		logger.Log.Println("Queue is empty")
		if container, ok := contentContainer.Objects[1].(*fyne.Container); ok {
			container.Objects[1].(*widget.Label).SetText("Queue is empty")
		}
		contentContainer.Objects = []fyne.CanvasObject{
			contentContainer.Objects[0],
			contentContainer.Objects[1], // Add "Queue is empty" label

		}
		contentContainer.Refresh()
		return
	}

	// If there are processes in the ready queue
	readyQueueWidget := widget.NewList(
		func() int {
			return len(readyQueue) // Number of processes
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Process ID | Process Name")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			selectedProcess := readyQueue[i] // Maintain the insertion order
			obj.(*widget.Label).SetText(fmt.Sprintf("%d | %s", selectedProcess.Pid, selectedProcess.Name))
		},
	)

	// Right side: Process Details
	infoLabel := widget.NewLabel("Select a process to see details")

	// Handles selection (Onclick)
	readyQueueWidget.OnSelected = func(id widget.ListItemID) {
		selectedProcess := readyQueue[id] // Use the process in the current order
		infoLabel.SetText(formatPCBDetailsScheduler(selectedProcess))
	}

	split := container.NewHSplit(
		readyQueueWidget,             // Left: List of processes
		container.NewVBox(infoLabel), // Right: PCB details
	)

	runningProcess := scheduler.GetRunningProcess()
	var runningProcessDetails *widget.Label
	if runningProcess == nil {
		runningProcessDetails = widget.NewLabel("No details")
	} else {
		runningProcessDetails = widget.NewLabel(formatPCBDetailsScheduler(runningProcess))
	}

	// Create the Running Process container
	runningProcessContainer := container.NewVBox(
		widget.NewLabelWithStyle("Running Process", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		runningProcessDetails,
	)

	// Update content container with the list and process details
	contentContainer.Objects = []fyne.CanvasObject{
		runningProcessContainer,
		split,
	}
	contentContainer.Refresh()
}

// Helper: Format PCB details
func formatPCBDetailsScheduler(pcb *processes.PCB) string {
	return fmt.Sprintf("PID: %d\nName: %s\nState: %s\nPriority: %d\nPC: %d\nAC: %d\nNextFreeCodeAddress: %d\nPageAmount: %d\n",
		pcb.Pid, pcb.Name, pcb.State.String(), pcb.Priority, pcb.ProcessState.PC, pcb.ProcessState.AC, pcb.NextFreeCodeAddress, pcb.PageAmount,
	)
}

// Auto-refresh function
func autoRefreshReadyQueue(contentContainer *fyne.Container, scheduler scheduler.SchedulerInterface) {
	for {
		time.Sleep(300 * time.Millisecond)                   // Refresh every second
		updateReadyQueueContent(contentContainer, scheduler) // Update the content dynamically
	}
}
