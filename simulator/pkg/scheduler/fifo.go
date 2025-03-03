package scheduler

import (
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
	s.readyQueue = append(s.readyQueue, pcb)
}

func (s *FifoScheduler) GetNextProcess() *processes.PCB {
	if len(s.readyQueue) == 0 {
		return nil
	}
	s.currentPCB = s.readyQueue[0]
	s.readyQueue = s.readyQueue[1:] // Remove first process (FIFO scheduling)
	return s.currentPCB
}

func (s *FifoScheduler) GetRunningProcess() *processes.PCB {
	return s.currentPCB
}
