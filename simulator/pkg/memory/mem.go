package memory

const (
	PageSize   = 256                  // Each page is 256 bytes
	NumPages   = 16                   // 16 pages in virtual memory
	NumFrames  = 8                    // 8 frames in physical memory
	MemorySize = PageSize * NumFrames // Total physical memory
)

// Page Table Entry
type PageTableEntry struct {
	FrameNumber int  // Maps to physical frame
	Valid       bool // True if the page is loaded in memory
}

// Memory with Paging
type Memory struct {
	PhysicalMemory [MemorySize]int          // Physical memory (array of frames)
	PageTable      [NumPages]PageTableEntry // Page table
	FrameUsage     [NumFrames]bool          // Tracks used frames
}
