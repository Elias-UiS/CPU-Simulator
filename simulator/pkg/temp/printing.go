package temp

import (
	"fmt"
	"time"
)

func Run(stopChan chan bool) {
	for i := 0; i < 100; i++ {
		select {
		case <-stopChan: // If stop signal is received, exit the goroutine
			fmt.Println("Simulation stopped.")
			return
		default:
			// Simulate work (e.g., a long computation)
			fmt.Println("Simulation running:", i)
			time.Sleep(time.Second) // Simulate a time delay
		}
	}
	fmt.Println("Simulation finished.")
}
