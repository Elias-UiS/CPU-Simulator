package scheduler

import (
	"CPU-Simulator/simulator/pkg/bindings"
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/processes"
)

type FifoScheduler struct {
	readyQueue []*processes.PCB
	currentPCB *processes.PCB
}

func NewScheduler() *FifoScheduler {
	return &FifoScheduler{
		readyQueue: make([]*processes.PCB, 0),
	}
}

func (s *FifoScheduler) AddProcess(pcb *processes.PCB) {
	logger.Log.Printf("DEBUGGING: Addprocess: pcb: %d", pcb)
	logger.Log.Printf("DEBUGGING: ReadyeQueueLEngth:  %d", len(s.readyQueue))
	s.readyQueue = append(s.readyQueue, pcb)
	logger.Log.Printf("DEBUGGING: ReadyeQueueLEngth:  %d", len(s.readyQueue))
	SyncReadyQueue(s.readyQueue)
}

func (s *FifoScheduler) GetNextProcess() *processes.PCB {
	if len(s.readyQueue) == 0 {
		return nil
	}
	s.currentPCB = s.readyQueue[0]
	s.readyQueue = s.readyQueue[1:] // Remove first process (FIFO scheduling)
	SyncReadyQueue(s.readyQueue)
	return s.currentPCB
}

func (s *FifoScheduler) GetRunningProcess() *processes.PCB {
	return s.currentPCB
}

func (s *FifoScheduler) GetReadyQueue() []*processes.PCB {
	return s.readyQueue
}

// Syncs the bound list with cpu.ReadyQueue
func SyncReadyQueue(readyQueue []*processes.PCB) {
	var tempList []interface{}
	for _, pcb := range readyQueue {
		tempList = append(tempList, pcb) // Convert []*PCB to []interface{}
	}
	bindings.ReadyQueueBinding.Set(tempList) // Update binding list
}
