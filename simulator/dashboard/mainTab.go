package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/os"
	"CPU-Simulator/simulator/pkg/scheduler"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func setupMainTab(os *os.OS) fyne.CanvasObject {

	schedulers := scheduler.ListSchedulers()
	for i, value := range schedulers {
		logger.Log.Printf("DEBUGGING: scheduler: %d: %s", i, value)
	}
	schedulerSelect := widget.NewSelect(schedulers, nil)
	schedulerSelect.PlaceHolder = "FiFo"

	createProcessButtonFromFile := widget.NewButton("Save", func() {
		schedulerName := schedulerSelect.Selected
		scheduler, err := scheduler.CreateScheduler(schedulerName)
		if err != nil {
			widget.NewLabel("Error: " + err.Error())
			return
		}
		logger.Log.Printf("DEBUGGING: scheduler: %s", schedulerName)
		logger.Log.Printf("%v", os.Scheduler)
		logger.Log.Printf("%v", scheduler)
		os.Scheduler = scheduler
		schedulerSelect.PlaceHolder = schedulerName
	})

	selectSchedulerContainer := container.NewVBox(
		schedulerSelect,
		createProcessButtonFromFile,
	)

	return selectSchedulerContainer
}
