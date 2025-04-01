package memory

type PageType int

// Use const to define the possible types: Stack, Code, Heap
const (
	Code  PageType = 0
	Heap  PageType = 1
	Stack PageType = 2
)

// kan endre til å lagre bits?
type PTE struct { // Page table entry
	Valid       bool // Page is mapped or not
	FrameNumber int  // Physical frame number
	Type        PageType
	//Dirty       bool   // Modified flag
	//Referenced  bool   // Recently accessed flag
	//SwapLoc     uint32 // Location on disk if swapped
}

type PageTable struct {
	Entries map[int]*PTE // Maps virtual page numbers to PTEs
}
