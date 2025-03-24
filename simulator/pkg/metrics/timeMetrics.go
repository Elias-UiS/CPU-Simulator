package metrics

import "time"

type Metrics struct {
	ArrivalTime         time.Time // When it first arrives in the queue
	CpuStartTime        time.Time // When it enters the cpu. Gets overwritten every time it enters the cpu
	CreationTime        time.Time // When it gets created
	CompletionTime      time.Time // When it finishes (or is terminated)
	SimulationStartTime time.Time // Time when simulation is resumed or started. Gets overwritten every time the simulation is resumed
	WaitingStartTime    time.Time // When it starts waiting in the ready queue. Gets overwritten every time it enters the queue

	BurstTime      time.Duration // Cpu time required to finish the process (Approx.)(At creation)
	CpuTime        time.Duration // Total time spent in the cpu
	SimulationTime time.Duration // Total time spent in the simulator
	WaitingTime    time.Duration // Total time spent waiting for the CPU
	ResponseTime   time.Duration // Time spent waiting in ready queue before first run (StartTime - ArrivalTime)
	TurnaroundTime time.Duration // Time between process creation and termination : (Completion - Arrival)

	InstructionAmount int // Total number of instructions executed
}
