package main

import (
	"CPU-Simulator/pkg/cpu"
	"CPU-Simulator/pkg/temp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")
	myWindow.Resize(fyne.NewSize(800, 600))

	stopChan := make(chan bool)
	isRunning := false
	// Create buttons for starting and stopping the simulation
	startSimulationButton := widget.NewButton("Start Simulation", func() {
		go cpu.Run()
	})

	startTestSimulationButton := widget.NewButton("Start Simulation", func() {
		if isRunning == false {
			go temp.Run(stopChan)
			isRunning = true
		}
	})
	stopTestSimulationButton := widget.NewButton("Stop Simulation", func() {
		if isRunning {
			stopChan <- true
			isRunning = false
		}
	})

	// Create more buttons for CPU settings
	resetCPUButton := widget.NewButton("Reset CPU", func() {
		// Simulate reset button click
	})
	viewStatsButton := widget.NewButton("View CPU Stats", func() {
		// Simulate stats view button click
	})

	// Organize buttons into a VBox (vertical layout)
	buttonsContainer := container.NewVBox(
		startSimulationButton,
		resetCPUButton,
		viewStatsButton,
		startTestSimulationButton,
		stopTestSimulationButton,
	)

	cpu := container.NewTabItem("CPU", buttonsContainer)
	memory := container.NewTabItem("MEMORY", widget.NewLabel("memory"))
	disk := container.NewTabItem("DISK", widget.NewLabel("disk"))

	tabs := container.NewAppTabs(
		cpu,
		memory,
		disk,
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
