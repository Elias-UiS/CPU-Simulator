package settings

var UpdateTimer int = 2

const (
	PageSize   = 8                    // size of each page (elements)
	NumFrames  = 16                   // number of frames(pages) in physical memory
	MemorySize = PageSize * NumFrames // Total physical memory size
	WordSize   = 32                   // Size of the element at an address (bits)
)
