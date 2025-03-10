package processes

type ProcessTable struct {
	ProcessMap map[int]*PCB
}

func CreateProcessTable() *ProcessTable {
	return &ProcessTable{
		ProcessMap: make(map[int]*PCB),
	}
}

func (pt *ProcessTable) AddProcessToTable(pcb *PCB) {
	pt.ProcessMap[pcb.Pid] = pcb
}

func (pt *ProcessTable) RemoveProcessFromTable(pcb *PCB) {
	delete(pt.ProcessMap, pcb.Pid)
}
