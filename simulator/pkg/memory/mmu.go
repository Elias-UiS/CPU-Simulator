package memory

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/settings"
	"fmt"
)

type MMU struct {
	TLB                  int        // doesnt store int, only temp. cache
	PageTable            *PageTable // Pages for the cpu
	memory               *Memory
	PageTableForCreation *PageTable // Used accessing memory not by cpu.
}

func (mmu *MMU) TranslateAddress(virtualAddr uint32) (int, error) {
	for key, value := range mmu.PageTable.Entries {
		logger.Log.Printf("MMU - PageTableEntrie %d -> %d", key, value.FrameNumber)
	}
	vpn := uint16(virtualAddr >> 16)
	offset := uint16(virtualAddr & 0xFFFF)
	logger.Log.Printf("VPN: %d\n", vpn)
	logger.Log.Printf("Offset: %d\n", offset)

	if offset < 0 || offset >= settings.PageSize {
		logger.Log.Printf("INFO: TranslateAddress() - Offset: %d\n", offset)
		err := fmt.Errorf("ERROR: mmu_TranslateAddress() | offset: address out of bounds")
		logger.Log.Println(err)
		return -1, err
	}
	if int(vpn) >= len(mmu.PageTable.Entries) {

		logger.Log.Printf("INFO: TranslateAddress() - VPN: %d\n", int(vpn))
		logger.Log.Printf("INFO: TranslateAddress() - PageTableSize: %d\n", len(mmu.PageTable.Entries))
		err := fmt.Errorf("ERROR: mmu_TranslateAddress() | pfn: address out of bounds")
		logger.Log.Println(err)
		return -1, err
	}

	frame := mmu.PageTable.Entries[vpn].FrameNumber
	physicalAddr := (uint32(frame) << 16) | uint32(offset)

	// Return the physical address (offset from base)
	return int(physicalAddr), nil
}

func (mmu *MMU) Read(physicalAddr uint32) (int, error) {
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= settings.PageSize {
		err := fmt.Errorf("ERROR: mmu_Read() | offset: address out of bounds")
		logger.Log.Println(err)
		return -1, err
	}

	if int(pfn) >= settings.NumFrames {
		err := fmt.Errorf("ERROR: mmu_Read() | pfn: address out of bounds")
		logger.Log.Println(err)
		return -1, err
	}

	data := mmu.memory.Frames[pfn][offset]

	// Return the physical address (offset from base)
	return int(data), nil
}

func (mmu *MMU) Write(physicalAddr uint32, value uint32) error {
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= settings.PageSize {
		err := fmt.Errorf("ERROR: mmu_Write() | offset: address out of bounds")
		logger.Log.Println(err)
		return err
	}
	if int(pfn) >= settings.NumFrames {
		err := fmt.Errorf("ERROR: mmu_Write() | pfn: address out of bounds")
		logger.Log.Println(err)
		return err
	}

	mmu.memory.Frames[pfn][offset] = value

	// Return the physical address (offset from base)
	return nil
}

func NewMMU(mem *Memory) *MMU {
	mmu := &MMU{
		TLB: settings.NumFrames,
		PageTable: &PageTable{
			Entries: make(map[uint16]*PTE),
		},
		memory: mem,
		PageTableForCreation: &PageTable{
			Entries: make(map[uint16]*PTE),
		},
	}

	return mmu
}

func (mmu *MMU) SetPageTable(pageTable *PageTable) {
	mmu.PageTable = pageTable
}

// func (mmu *MMU) StoreInstruction(pc int, opType int, opcode int, value byte) error {
// 	physicalAddr, err := mmu.TranslateAddress(uint32(pc))
// 	if physicalAddr == -1 {
// 		fmt.Println(err)
// 	}
// 	pfn := uint16(physicalAddr >> 16)
// 	offset := uint16(physicalAddr & 0xFFFF)

// 	if offset < 0 || offset >= PageSize {
// 		return fmt.Errorf("offset: address out of bounds")
// 	}
// 	if int(pfn) >= NumFrames {
// 		return fmt.Errorf("PFN: address out of bounds")
// 	}

// 	instructionByte := (uint32(opType) << 31) | (uint32(opcode) & 0x31F)

// 	mmu.memory.Frames[pfn][offset] = instructionByte
// 	mmu.memory.Frames[pfn][offset+1] = uint32(value)

// 	// Return the physical address (offset from base)
// 	return nil
// }
