package scheduler

import (
	"CPU-Simulator/simulator/pkg/processes"
	"fmt"
)

// Scheduler defines the interface that all scheduling algorithms must implement.
type SchedulerInterface interface {
	AddProcess(*processes.PCB)         // Add process to the queue
	GetNextProcess() *processes.PCB    // Decides the next process to run
	GetRunningProcess() *processes.PCB // Return current running process.
	GetReadyQueue() []*processes.PCB
}

// Global map to store registered scheduler implementations
var Schedulers = make(map[string]func() SchedulerInterface)

// RegisterScheduler adds a new scheduler to the registry
func RegisterScheduler(name string, newScheduler func() SchedulerInterface) {
	Schedulers[name] = newScheduler
}

// ListSchedulers returns all available scheduler names
func ListSchedulers() []string {
	keys := make([]string, 0, len(Schedulers))
	for k := range Schedulers {
		keys = append(keys, k)
	}
	return keys
}

// CreateScheduler creates a scheduler by name
func CreateScheduler(name string) (SchedulerInterface, error) {
	if constructor, found := Schedulers[name]; found {
		return constructor(), nil
	}
	return nil, fmt.Errorf("scheduler not found: %s", name)
}
