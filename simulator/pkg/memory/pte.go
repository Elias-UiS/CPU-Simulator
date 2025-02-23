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
	//Valid       bool   // Page is valid in memory
	FrameNumber uint32 // Physical frame number
	Type        PageType
	//Read        bool   // Read permission
	//Write       bool   // Write permission
	//Execute     bool   // Execute permission
	//Dirty       bool   // Modified flag
	//Referenced  bool   // Recently accessed flag
	//SwapLoc     uint32 // Location on disk if swapped
}
