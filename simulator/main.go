package main

import (
	"CPU-Simulator/simulator/dashboard"
	"CPU-Simulator/simulator/pkg/logger"
)

func main() {
	logger.Init()
	logger.Log.Println("INFO: main()")
	dashboardStruct := dashboard.DashboardStruct{}
	dashboard.Dashboard(&dashboardStruct)
}
