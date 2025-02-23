package os

type Scheduler struct {
	readyQueue []*PCB
	currentPCB *PCB
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		readyQueue: make([]*PCB, 0),
	}
}

func (s *Scheduler) AddProcess(pcb *PCB) {
	s.readyQueue = append(s.readyQueue, pcb)
}

func (s *Scheduler) GetNextProcess() *PCB {
	if len(s.readyQueue) == 0 {
		return nil
	}
	s.currentPCB = s.readyQueue[0]
	s.readyQueue = s.readyQueue[1:] // Remove first process (FIFO scheduling)
	return s.currentPCB
}
