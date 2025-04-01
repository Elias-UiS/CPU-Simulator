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

	// /////////////////////////////////
	// // Dropdown to select memory model
	// memoryType := widget.NewSelect([]string{"Paging", "Segmentation"}, func(selected string) {
	// 	fmt.Println("Selected memory type:", selected)
	// })
	// memoryType.SetSelected("Paging") // Default selection

	// // Button to create memory based on selection
	// var physicalMemory *memory.Memory
	// setupMemoryButton := widget.NewButton("Setup Memory", func() {
	// 	if memoryType.Selected == "Paging" {
	// 		physicalMemory = memory.NewMemory()
	// 	} else {
	// 		physicalMemory = memory.NewMemory()
	// 	}
	// 	fmt.Println("Memory initialized:", memoryType.Selected)
	// })
	// fmt.Printf("TESTING HERE4")
	// var cpuInstance *cpu.CPU
	// setupCPUButton := widget.NewButton("Setup CPU", func() {
	// 	if physicalMemory == nil {
	// 		fmt.Println("Error: Setup memory first!")
	// 		return
	// 	}
	// 	cpuInstance = cpu.NewCPU()
	// 	fmt.Println("CPU initialized with memory:", memoryType.Selected)
	// })

	// configContainer := container.NewVBox(
	// 	widget.NewLabel("Select Memory Type"),
	// 	memoryType,
	// 	setupMemoryButton,
	// 	setupCPUButton,
	// )

	//////////////////////////////////

	os := os.NewOS()
	go os.TestNumer()
	//var cpuInstance = os.CPU[0]

	clock := widget.NewLabel("00:00:00")
	status := widget.NewLabel("Not Running")

	startSimulationButton := widget.NewButton("Start", func() {
		if dash.IsRunnning {
			return
		}
		go dash.startClock(clock)
		os.StartSimulation()
		status.SetText("Running")
		dash.IsRunnning = true
		go dash.SystemState.UpdateState(os)
	})

	stopSimulationButton := widget.NewButton("Stop", func() {

		//stopClockChan <- true
		status.SetText("Stopped")
	})

	resumetSimulationButton := widget.NewButton("Resume", func() {
		if dash.IsRunnning {
			return
		}
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

	cpu := container.NewTabItem("CPU", setupCpuTab(os))
	memory := container.NewTabItem("MEMORY", setupMemoryTab(os.Memory))
	processes := container.NewTabItem("Processes", CreatePCBUI(os, os.ProcessTable))
	scheduler := container.NewTabItem("Scheduler", CreateSchedulerTab(os.Scheduler))
	calculator := container.NewTabItem("Calculator", setupCalculatorTab())
	processCreation := container.NewTabItem("Process Creator", ProcessCreationTab(os))
	systemState := container.NewTabItem("System State", setupSystemStateTab(dash.SystemState.PubSub))

	tabs := container.NewAppTabs(
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

// // Registers
// var registerHeader *widget.Label = widget.NewLabel("Registers\n")
// registerHeader.TextStyle.Monospace = true
// registerHeader.TextStyle.Bold = true

// var registerDisplay string = cpu.GetRegisters()
// var registerDisplayWidget *widget.Label = widget.NewLabel(registerDisplay)
// registerDisplayWidget.TextStyle.Monospace = true
// registerDisplayWidget.TextStyle.Bold = true
// var registerContainer = container.NewStack(
// 	container.NewVBox(
// 		registerHeader,
// 		registerDisplayWidget,
// 	))

// accumulatorValue := binding.BindInt(&cpu.Registers.AC)
// s, _ := boundString.Get()
// log.Printf("Bound = '%s'", s)
