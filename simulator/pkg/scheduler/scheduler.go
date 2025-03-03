package scheduler

import (
	"CPU-Simulator/simulator/pkg/processes"
)

// Scheduler defines the interface that all scheduling algorithms must implement.
type SchedulerInterface interface {
	AddProcess(p *processes.PCB)       // Add process to the queue
	GetNextProcess() *processes.PCB    // Decides the next process to run
	GetRunningProcess() *processes.PCB // Return current running process.
}
