package schedulerFiles

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/processes"
	"CPU-Simulator/simulator/pkg/scheduler"
	"CPU-Simulator/simulator/pkg/settings"
)

type RoundRobinScheduler struct {
	readyQueue []*processes.PCB
	currentPCB *processes.PCB
	quantum    int // Quantum for round-robin scheduling. In cpu cycles, not time.
}

func NewRoundRobinScheduler() *RoundRobinScheduler {
	settings.InstructionLimitPerRun = 2 // Set the instruction limit per run to 1 for round-robin scheduling
	return &RoundRobinScheduler{
		readyQueue: make([]*processes.PCB, 0),
		quantum:    1, // Default quantum value
	}
}

func (s *RoundRobinScheduler) AddProcess(pcb *processes.PCB) {
	logger.Log.Printf("DEBUGGING: RoundRobin: Addprocess: pcb: %d", pcb)
	s.readyQueue = append(s.readyQueue, pcb)
	syncReadyQueue2(s.readyQueue)
}

func (s *RoundRobinScheduler) GetNextProcess() *processes.PCB {
	if len(s.readyQueue) == 0 {
		return nil
	}
	if s.currentPCB != nil {
		if s.currentPCB.State != processes.Terminated {
			s.AddProcess(s.currentPCB) // Re-add the current process to the end of the queue
		}
	}
	s.currentPCB = s.readyQueue[0]
	s.readyQueue = s.readyQueue[1:] // Remove first process

	syncReadyQueue2(s.readyQueue)
	return s.currentPCB
}

func (s *RoundRobinScheduler) GetRunningProcess() *processes.PCB {
	return s.currentPCB
}

func (s *RoundRobinScheduler) GetReadyQueue() []*processes.PCB {
	return s.readyQueue
}

// Syncs the bound list with cpu.ReadyQueue
func syncReadyQueue2(readyQueue []*processes.PCB) {
	var tempList []interface{}
	for _, pcb := range readyQueue {
		tempList = append(tempList, pcb) // Convert []*PCB to []interface{}
	}
	bindings.ReadyQueueBinding.Set(tempList) // Update binding list
}

func init() {
	scheduler.RegisterScheduler("RoundRobin", func() scheduler.SchedulerInterface { return NewRoundRobinScheduler() })
}
