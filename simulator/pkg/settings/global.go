package settings

import "time"

const (
	PageSize   = 8                    // size of each page (elements)
	NumFrames  = 64                   // number of frames(pages) in physical memory
	MemorySize = PageSize * NumFrames // Total physical memory size
	WordSize   = 32                   // Size of the element at an address (bits)
	MemType    = 0                    // 0 = paging,
)

var CpuFetchDecodeExecuteDelay time.Duration = 100 // Delay between the cycle steps, in milliseconds.
var InstructionLimitPerRun int = 100               // How many instructions a process gets before a context switch.
