package memory

const (
	PageSize   = 128                  // size of each page
	NumFrames  = 8                    // number of frames(pages) in physical memory
	MemorySize = PageSize * NumFrames // Total physical memory size
)

type Memory struct {
	Frames [][]byte // Represents physical memory
}

func NewMemory() *Memory {
	frame := make([][]byte, NumFrames)
	for i := range frame {
		frame[i] = make([]byte, PageSize)
	}

	return &Memory{
		Frames: frame,
	}
}
