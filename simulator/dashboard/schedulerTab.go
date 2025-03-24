package dashboard

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// Global variables
var (
	selectedPCB    *processes.PCB
	details        *widget.Label
	runningDetails *widget.Label
)

// Function to create the Scheduler UI
func CreateSchedulerTab(scheduler scheduler.SchedulerInterface) fyne.CanvasObject {
	// Create the Ready Queue container
	readyQueueContainer := ReadyQueueContainer(scheduler)

	// Create the Running Process container
	runningProcessContainer := RunningProcessContainer(scheduler)

	// Right side: Process Details
	details = widget.NewLabel("Select a process to view details.")
	runningDetails = widget.NewLabel("No process currently running.")

	// Create the content container (tab)
	contentContainer := container.NewGridWithRows(2, runningProcessContainer, readyQueueContainer)

	// Initial call to load the ready queue (empty or not)
	updateReadyQueueContent(contentContainer, scheduler)

	// Auto-refresh mechanism
	go autoRefreshReadyQueue(contentContainer, scheduler)

	return contentContainer
}

func RunningProcessContainer(scheduler scheduler.SchedulerInterface) *fyne.Container {
	runningTitleLabel := widget.NewLabel("Running Process")
	runningProcessContainer := container.NewVBox(
		runningTitleLabel,                     // Running Process title
		widget.NewLabel("No running process"), // Process details placeholder
	)
	return runningProcessContainer
}

func ReadyQueueContainer(scheduler scheduler.SchedulerInterface) *fyne.Container {
	// Create the title label (Bold and Left-aligned)
	titleLabel := widget.NewLabelWithStyle("Ready Queue", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	// Create the container that will hold the ready queue list or empty message
	emptyLabel := widget.NewLabel("Queue is empty")
	readyQueueContainer := container.NewBorder(titleLabel, emptyLabel, nil, nil, nil)

	return readyQueueContainer
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
	// Ready Queue List (Bound)
	readyQueueWidget := widget.NewListWithData(
		bindings.ReadyQueueBinding,
		func() fyne.CanvasObject {
			return widget.NewLabel("Process ID | Process Name")
		},
		func(di binding.DataItem, obj fyne.CanvasObject) {
			item, _ := di.(binding.Untyped).Get()
			if pcb, ok := item.(*processes.PCB); ok {
				obj.(*widget.Label).SetText(fmt.Sprintf("%d | %s", pcb.Pid, pcb.Name))
			}
		},
	)

	// Handles selection (Onclick)
	readyQueueWidget.OnSelected = func(id widget.ListItemID) {
		selectedPCB = readyQueue[id] // Use the process in the current order
		updateDetails()
	}

	split := container.NewHSplit(
		readyQueueWidget,           // Left: List of processes
		container.NewVBox(details), // Right: PCB details
	)

	runningProcess := scheduler.GetRunningProcess()
	if runningProcess == nil {
		runningDetails = widget.NewLabel("No details")
	} else {
		updateRunningDetails(runningProcess)
	}

	// Create the Running Process container
	runningProcessContainer := container.NewVBox(
		widget.NewLabelWithStyle("Running Process", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		runningDetails,
	)

	// Update content container with the list and process details
	contentContainer.Objects = []fyne.CanvasObject{
		runningProcessContainer,
		split,
	}
	contentContainer.Refresh()
}

// Auto-refresh function
func autoRefreshReadyQueue(contentContainer *fyne.Container, scheduler scheduler.SchedulerInterface) {
	for {
		time.Sleep(1000 * time.Millisecond)                  // Refresh every second
		updateReadyQueueContent(contentContainer, scheduler) // Update the content dynamically
	}
}

func updateDetails() {
	if selectedPCB == nil {
		details.SetText("No process selected.")
		return
	}

	details.SetText(
		"Process: " + selectedPCB.Name +
			"\nPID: " + fmt.Sprintf("%d", selectedPCB.Pid) +
			"\nState: " + fmt.Sprintf("%d", selectedPCB.State) +
			"\nPriority: " + fmt.Sprintf("%d", selectedPCB.Priority) +
			"\n" +
			"\nCreation Time: " + selectedPCB.Metrics.CreationTime.Format("00:00:00") +
			"\nArrival Time: " + selectedPCB.Metrics.ArrivalTime.Format("00:00:00") +
			"\n" +
			"\nBurst Time: " + fmt.Sprintf("%.2f", selectedPCB.Metrics.BurstTime.Seconds()) + "s" +
			"\nCPU Time: " + fmt.Sprintf("%.2f", selectedPCB.Metrics.CpuTime.Seconds()) + "s" +
			"\nResponse Time: " + fmt.Sprintf("%.2f", selectedPCB.Metrics.ResponseTime.Seconds()) + "s" +
			"\nTurnaround Time: " + fmt.Sprintf("%.2f", selectedPCB.Metrics.TurnaroundTime.Seconds()) + "s" +
			"\nSimulation Time: " + fmt.Sprintf("%.2f", selectedPCB.Metrics.SimulationTime.Seconds()) + "s" +
			"\nWaiting Time: " + fmt.Sprintf("%.2f", selectedPCB.Metrics.WaitingTime.Seconds()) + "s" +
			"\n" +
			"\nInstuction Count: " + fmt.Sprintf("%d", selectedPCB.Metrics.InstructionAmount),
	)
}

func updateRunningDetails(pcb *processes.PCB) {
	runningDetails.SetText(
		"Process: " + pcb.Name +
			"\nPID: " + fmt.Sprintf("%d", pcb.Pid) +
			"\nState: " + fmt.Sprintf("%d", pcb.State) +
			"\nPriority: " + fmt.Sprintf("%d", pcb.Priority) +
			"\n" +
			"\nCreation Time: " + pcb.Metrics.CreationTime.Format("00:00:00") +
			"\nArrival Time: " + pcb.Metrics.ArrivalTime.Format("00:00:00") +
			"\n" +
			"\nBurst Time: " + fmt.Sprintf("%.2f", pcb.Metrics.BurstTime.Seconds()) + "s" +
			"\nCPU Time: " + fmt.Sprintf("%.2f", pcb.Metrics.CpuTime.Seconds()) + "s" +
			"\nResponse Time: " + fmt.Sprintf("%.2f", pcb.Metrics.ResponseTime.Seconds()) + "s" +
			"\nTurnaround Time: " + fmt.Sprintf("%.2f", pcb.Metrics.TurnaroundTime.Seconds()) + "s" +
			"\nSimulation Time: " + fmt.Sprintf("%.2f", pcb.Metrics.SimulationTime.Seconds()) + "s" +
			"\nWaiting Time: " + fmt.Sprintf("%.2f", pcb.Metrics.WaitingTime.Seconds()) + "s" +
			"\n" +
			"\nInstuction Count: " + fmt.Sprintf("%d", pcb.Metrics.InstructionAmount),
	)
}
