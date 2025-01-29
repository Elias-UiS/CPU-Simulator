package dashboard

import (
	"CPU-Simulator/simulator/pkg/cpu"
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
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")
	myWindow.Resize(fyne.NewSize(800, 600))

	cpuInstance := cpu.NewCPU()

	registerHeader := widget.NewLabel("Registers")
	registerHeader.TextStyle.Bold = true

	registerDisplayWidget := widget.NewLabel("")
	registerDisplayWidget.TextStyle.Monospace = true

	updateRegisterDisplay(cpuInstance, registerDisplayWidget)

	clock := widget.NewLabel("00:00:00")
	status := widget.NewLabel("Stopped")

	stopClockChan := make(chan bool)

	startSimulationButton := widget.NewButton("Start", func() {
		go cpu.Run(cpuInstance)
		go startClock(clock, stopClockChan, cpuInstance, registerDisplayWidget)
		status.SetText("Running")
	})

	stopSimulationButton := widget.NewButton("Stop", func() {
		stopClockChan <- true
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

	// Organize the register UI
	registerContainer := container.NewVBox(
		registerHeader,
		registerDisplayWidget,
	)

	cpu := container.NewTabItem("CPU", registerContainer)
	memory := container.NewTabItem("MEMORY", widget.NewLabel("memory"))
	disk := container.NewTabItem("DISK", widget.NewLabel("disk"))
	test := container.NewTabItem("TEST", buttonsContainer)

	tabs := container.NewAppTabs(
		cpu,
		memory,
		disk,
		test,
	)

	tabs.SetTabLocation(container.TabLocationLeading)

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

	mainContainer := container.NewVBox(
		topBarContainer,
		tabs,
	)

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}
