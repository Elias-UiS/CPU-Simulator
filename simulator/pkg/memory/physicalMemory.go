package memory

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/settings"
)

type Memory struct {
	Frames [][]uint32 // Represents physical memory
}

var physicalMemory *Memory

func NewMemory() *Memory {
	logger.Log.Println("INFO: physicalMemory NewMemory()")
	frame := make([][]uint32, settings.NumFrames)
	for i := range frame {
		frame[i] = make([]uint32, settings.PageSize)
	}

	physicalMemory = &Memory{
		Frames: frame,
	}

	return physicalMemory

}

func GetMemory() *Memory {
	logger.Log.Println("INFO: GetMemory called.")

	if physicalMemory == nil {
		logger.Log.Println("INFO: physicalMemory is nil. Initializing...")
		return NewMemory() // Initializes physicalMemory
	}

	logger.Log.Println("INFO: physicalMemory already initialized.")
	return physicalMemory
}
