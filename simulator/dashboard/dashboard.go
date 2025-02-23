package dashboard

import (
	"CPU-Simulator/simulator/pkg/cpu"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/os"
	"CPU-Simulator/simulator/pkg/temp"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func load() {
	fmt.Println("Test")
}

func run() {
	fmt.Println("Test")
}

func next() {
	fmt.Println("Test")
}

func pause() {
	fmt.Println("Test")
}

func exit() {
	fmt.Println("Test")
}

func updateTime(clock *widget.Label, start time.Time) {
	elapsed := time.Since(start)
	elapsedStr := fmt.Sprintf("%02d:%02d:%02d",
		int(elapsed.Hours())%24,
		int(elapsed.Minutes())%60,
		int(elapsed.Seconds())%60)
	clock.SetText(elapsedStr)
}

func startClock(clock *widget.Label, stopClockChan chan bool, cpuInstance *cpu.CPU, registerDisplayWidget *widget.Label) {
	start := time.Now()
	for range time.Tick(time.Second) {
		select {
		case <-stopClockChan:
			fmt.Println("Clock Stopped.")
			return
		default:
			updateTime(clock, start)
			updateRegisterDisplay(cpuInstance, registerDisplayWidget)
		}
	}
}

func updateRegisterDisplay(cpuInstance *cpu.CPU, registerDisplayWidget *widget.Label) {
	registers := cpuInstance.Registers
	displayText := fmt.Sprintf(
		"AC: %d\nPC: %d\nMAR: %d\nMDR: %d\nIR_Opcode: %d\nIR_Operand: %d",
		registers.AC, registers.PC, registers.MAR, registers.MDR.Data, registers.IR.Opcode, registers.IR.Operand,
	)
	registerDisplayWidget.SetText(displayText)
}

func Dashboard() {
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
	//var cpuInstance = os.CPU[0]

	clock := widget.NewLabel("00:00:00")
	status := widget.NewLabel("Stopped")

	//stopClockChan := make(chan bool)

	startSimulationButton := widget.NewButton("Start", func() {
		//go startClock(clock, stopClockChan, cpuInstance, registerDisplayWidget)
		os.StartSimulation()
		status.SetText("Running")
	})

	stopSimulationButton := widget.NewButton("Stop", func() {
		//stopClockChan <- true
		status.SetText("Stopped")
	})

	pauseSimulationButton := widget.NewButton("Pause", func() {
		status.SetText("Paused")
	})

	nextStepSimulationButton := widget.NewButton("Next Step", func() {
		fmt.Println("test")
	})

	topBarContainer := container.NewHBox(
		startSimulationButton,
		stopSimulationButton,
		pauseSimulationButton,
		nextStepSimulationButton,
		status,
		clock,
	)

	stopChan := make(chan bool)
	isRunning := false

	// Create buttons for starting and stopping the simulation
	startTestSimulationButton := widget.NewButton("Start Test Simulation", func() {
		if isRunning == false {
			go temp.Run(stopChan)
			isRunning = true
		}
	})
	stopTestSimulationButton := widget.NewButton("Stop Test Simulation", func() {
		if isRunning {
			stopChan <- true
			isRunning = false
		}
	})

	resetCPUButton := widget.NewButton("Reset CPU", func() {
		// Todo?
	})
	viewStatsButton := widget.NewButton("View CPU Stats", func() {
		// Todo?
	})

	// Organize buttons into a VBox (vertical layout)
	buttonsContainer := container.NewVBox(
		startTestSimulationButton,
		stopTestSimulationButton,
		resetCPUButton,
		viewStatsButton,
	)
	cpu := container.NewTabItem("CPU", setupCpuTab(os))
	memory := container.NewTabItem("MEMORY", setupMemoryTab(os.Memory))
	disk := container.NewTabItem("DISK", widget.NewLabel("disk"))
	processes := container.NewTabItem("Processes", widget.NewLabel("Process list"))
	test := container.NewTabItem("TEST", buttonsContainer)
	calculator := container.NewTabItem("Calculator", setupCalculatorTab())

	tabs := container.NewAppTabs(
		cpu,
		memory,
		processes, // Todo: Implement processes tab
		// Den skal inneholde en liste over alle prosessene som kjører, og informasjon om hver enkelt prosess.
		// Som Pages, og memoriet til den.
		disk,
		test,
		calculator,
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
