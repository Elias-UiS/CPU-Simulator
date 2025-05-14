package memory

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/settings"
	"fmt"
)

type MemoryController struct {
	memory *Memory
}

func (memoryController *MemoryController) Read(physicalAddr uint32) (int, error) {
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= settings.PageSize {
		err := fmt.Errorf("ERROR: memoryController_Read() | offset: address out of bounds")
		logger.Log.Println(err)
		logger.Log.Println("Offset: %d", offset)
		logger.Log.Println("PFN: %d", pfn)
		return -1, err
	}

	if int(pfn) >= settings.NumFrames {
		err := fmt.Errorf("ERROR: memoryController_Read() | pfn: address out of bounds")
		logger.Log.Println(err)
		return -1, err
	}

	data := memoryController.memory.Frames[pfn][offset]

	// Return the physical address (offset from base)
	return int(data), nil
}

func (memoryController *MemoryController) Write(physicalAddr uint32, value uint32) error {
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= settings.PageSize {
		err := fmt.Errorf("ERROR: mmu_Write() | offset: address out of bounds")
		logger.Log.Println("Offset: %d", offset)
		logger.Log.Println(err)
		return err
	}
	if int(pfn) >= settings.NumFrames {
		err := fmt.Errorf("ERROR: mmu_Write() | pfn: address out of bounds")
		logger.Log.Println(err)
		return err
	}

	memoryController.memory.Frames[pfn][offset] = value

	// Return the physical address (offset from base)
	return nil
}

func NewMemoryController(mem *Memory) *MemoryController {
	memoryController := &MemoryController{
		memory: mem,
	}
	return memoryController
}
