package main

import (
	"CPU-Simulator/simulator/dashboard"
	"CPU-Simulator/simulator/pkg/logger"
	"fmt"
)

func main() {
	fmt.Println("INFO: main() fmt")
	defer fmt.Println("INFO: main() fmt defer")
	logger.Init()
	logger.Log.Println("INFO: main()")
	dashboard.Dashboard()
}
