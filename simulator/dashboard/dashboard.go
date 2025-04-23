package dashboard

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/os"
	"CPU-Simulator/simulator/pkg/systemState"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DashboardStruct struct {
	SystemState systemState.State
	IsRunnning  bool
	clock       int // In seconds
}

func (dash *DashboardStruct) startClock(clock *widget.Label) {
	dash.clock = 0
	for {
		time.Sleep(1000 * time.Millisecond)
		if !dash.IsRunnning {
			fmt.Println("Clock paused.")
		} else {
			dash.clock += 1
			duration := time.Duration(dash.clock) * time.Second
			timeStr := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
			clock.SetText(timeStr)
		}
	}
}

func Dashboard(dash *DashboardStruct) {
	// Fyne app
	logger.Log.Println("INFO: dashboard_Dashboard() - Starting Dashboard.")

	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")
	myWindow.Resize(fyne.NewSize(800, 600))

	os := os.NewOS()
	go os.TestNumer()
	//var cpuInstance = os.CPU[0]

	clock := widget.NewLabel("00:00:00")
	status := widget.NewLabel("Not Running")

	var scheduler *container.TabItem

	startSimulationButton := widget.NewButton("Start", func() {
		if dash.IsRunnning {
			return
		}
		go dash.startClock(clock)
		os.StartSimulation()
		status.SetText("Running")
		dash.IsRunnning = true
		scheduler.Content = CreateSchedulerTab(os.Scheduler)

		go dash.SystemState.UpdateState(os)
	})

	stopSimulationButton := widget.NewButton("Stop", func() {
		status.SetText("Not Running")
		os.StopSimulation()
	})

	resumetSimulationButton := widget.NewButton("Resume", func() {
		if dash.IsRunnning {
			return
		}
		os.StepMode = false
		go os.ResumeSimulation()
		status.SetText("Running")
		dash.IsRunnning = true
	})

	pauseSimulationButton := widget.NewButton("Pause", func() {
		if !dash.IsRunnning {
			return
		}
		go os.PauseSimulation()
		status.SetText("Paused")
		dash.IsRunnning = false
	})

	nextStepSimulationButton := widget.NewButton("Next Step", func() {
		if dash.IsRunnning {
			return
		}
		os.StepMode = true
		go os.ResumeSimulation()
		status.SetText("Step Mode")
	})

	nextProcessSimulationButton := widget.NewButton("Next Process", func() {
		fmt.Println("test")
		if os.MidContextSwitch {
			return
		}
		os.ContextSwitch(os.GetCpu())
	})

	topBarContainer := container.NewHBox(
		startSimulationButton,
		stopSimulationButton,
		resumetSimulationButton,
		pauseSimulationButton,
		nextStepSimulationButton,
		nextProcessSimulationButton,
		status,
		clock,
	)

	main := container.NewTabItem("Main (Settings)", setupMainTab(os))
	cpu := container.NewTabItem("CPU", setupCpuTab(os))
	memory := container.NewTabItem("MEMORY", setupMemoryTab(os.Memory))
	processes := container.NewTabItem("Processes", CreatePCBUI(os, os.ProcessTable))
	placeholderContent := widget.NewLabel("Start simulation first")
	scheduler = container.NewTabItem("Scheduler", placeholderContent)
	calculator := container.NewTabItem("Calculator", setupCalculatorTab())
	processCreation := container.NewTabItem("Process Creator", ProcessCreationTab(os))
	systemState := container.NewTabItem("System State", setupSystemStateTab(dash.SystemState.PubSub))

	tabs := container.NewAppTabs(
		main,
		cpu,
		memory,
		processes, // Todo: Implement processes tab
		// Den skal inneholde en liste over alle prosessene som kjører, og informasjon om hver enkelt prosess.
		// Som Pages, og memoriet til den.
		scheduler,
		calculator,
		processCreation,
		systemState,
	)

	tabsContainer := container.NewStack(tabs)

	tabs.SetTabLocation(container.TabLocationLeading)

	mainContainer := container.NewVBox(
		topBarContainer,
		tabsContainer,
	)

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()

}
