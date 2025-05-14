package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/os"
	"CPU-Simulator/simulator/pkg/scheduler"
	"CPU-Simulator/simulator/pkg/settings"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func setupMainTab(os *os.OS) fyne.CanvasObject {

	schedulers := scheduler.ListSchedulers()
	schedulerTitle := widget.NewLabel("Scheduler: ")
	schedulerSelect := widget.NewSelect(schedulers, nil)
	schedulerSelect.PlaceHolder = "FiFo"
	schedulerSelect.Selected = schedulers[0] // Default to the first scheduler

	timeBetweenInstructionsTitle := widget.NewLabel("Time between instructions (ms): ")
	timeBetweenInstructions := widget.NewEntry()
	timeBetweenInstructions.PlaceHolder = "1000"
	timeBetweenInstructions.SetText("1000") // Default value

	saveSettings := widget.NewButton("Save", func() {
		schedulerName := schedulerSelect.Selected
		scheduler, err := scheduler.CreateScheduler(schedulerName)
		if err != nil {
			widget.NewLabel("Error: " + err.Error())
			return
		}
		logger.Log.Printf("%v", os.Scheduler)
		logger.Log.Printf("%v", scheduler)
		delay, err := strconv.Atoi(timeBetweenInstructions.Text)
		if err != nil {
			logger.Log.Printf("DEBUGGING: setupMainTab")
			widget.NewLabel("Error: Invalid time input")
			return
		}
		settings.CpuFetchDecodeExecuteDelay = (time.Duration(delay)) / 3
		timeBetweenInstructions.PlaceHolder = timeBetweenInstructions.Text
		os.Scheduler = scheduler
		schedulerSelect.PlaceHolder = schedulerName
	})

	selectSchedulerContainer := container.NewVBox(
		timeBetweenInstructionsTitle,
		timeBetweenInstructions,
		schedulerTitle,
		schedulerSelect,
		saveSettings,
	)

	return selectSchedulerContainer
}
