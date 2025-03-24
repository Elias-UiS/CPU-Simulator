package dashboard

// import (
// 	"CPU-Simulator/simulator/pkg/scheduler"
// 	"fmt"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/widget"
// )

// type ReadyQueueItem struct {
// 	Pid  int
// 	Name string
// }

// func (item *ReadyQueueItem) ID() string {
// 	return fmt.Sprintf("%d", item.Pid)
// }

// func (item *ReadyQueueItem) String() string {
// 	return fmt.Sprintf("%d | %s", item.Pid, item.Name)
// }

// // Function to create the Scheduler UI
// func CreateSchedulerTab(scheduler scheduler.SchedulerInterface) fyne.CanvasObject {
// 	// Create the title label (Bold and Left-aligned)
// 	titleLabel := widget.NewLabelWithStyle("Ready Queue", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

// 	// Initialize the ready queue data
// 	readyQueueData := binding.NewStringList()
// 	updateReadyQueueContent(readyQueueData, scheduler)

// 	// Create the list widget
// 	list := widget.NewListWithData(readyQueueData, func() fyne.CanvasObject {
// 		return widget.NewLabel("")
// 	}, func(item binding.DataItem, obj fyne.CanvasObject) {
// 		obj.(*widget.Label).Bind(item.(binding.String))
// 	})

// 	// Create the container
// 	contentContainer := container.NewVBox(titleLabel, list)

// 	// Auto-refresh mechanism
// 	go autoRefreshReadyQueue(contentContainer, scheduler, readyQueueData)

// 	return contentContainer
// }

// // Function to update the content based on the current state of the ready queue
// func updateReadyQueueContent(readyQueueData binding.StringList, scheduler scheduler.SchedulerInterface) {
// 	readyQueueData.Set([]string{})
// 	for _, process := range scheduler.GetReadyQueue() {
// 		readyQueueData.Append(fmt.Sprintf("%d | %s", process.Pid, process.Name))
// 	}
// }

// // Auto-refresh function
// func autoRefreshReadyQueue(contentContainer *fyne.Container, scheduler scheduler.SchedulerInterface, readyQueueData binding.StringList) {
// 	for {
// 		time.Sleep(1 * time.Second) // Refresh every second
// 		updateReadyQueueContent(readyQueueData, scheduler)
// 		contentContainer.Refresh()
// 	}
// }
