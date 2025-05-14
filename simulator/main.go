package main

import (
	"CPU-Simulator/simulator/dashboard"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/systemLog"
	"CPU-Simulator/simulator/pkg/systemState"
)

func main() {
	logger.Init()
	logger.Log.Println("INFO: main()")
	systemState := systemState.CreateState()
	systemStateLog := systemLog.NewSystemStateLog(systemState.PubSub)
	go systemStateLog.LogSystemState()
	dashboardStruct := dashboard.DashboardStruct{
		SystemState:    systemState,
		SystemStateLog: systemStateLog,
	}
	dashboard.Dashboard(&dashboardStruct)
}
