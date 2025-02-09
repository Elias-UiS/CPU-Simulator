package memory

import (
	"fmt"
)

type MMU struct {
	TLB       int    // doesnt store int, only temp.
	pageTable *[]PTE // Pages
	memory    *Memory
}

func (mmu *MMU) TranslateAddress(virtualAddr uint32) (int, error) {
	vpn := uint16(virtualAddr >> 16)
	offset := uint16(virtualAddr & 0xFFFF)

	if offset < 0 || offset >= PageSize {
		return -1, fmt.Errorf("offset: address out of bounds")
	}
	if int(vpn) >= len(*mmu.pageTable) {
		return -1, fmt.Errorf("VPN: address out of bounds")
	}

	pte := (*mmu.pageTable)[vpn]
	frame := pte.FrameNumber
	physicalAddr := (uint32(frame) << 16) | uint32(offset)

	// Return the physical address (offset from base)
	return int(physicalAddr), nil
}

func (mmu *MMU) Read(physicalAddr uint32) (int, error) {
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= PageSize {
		return -1, fmt.Errorf("offset: address out of bounds")
	}
	if int(pfn) >= NumFrames {
		return -1, fmt.Errorf("PFN: address out of bounds")
	}

	data := mmu.memory.Frames[pfn][offset]

	// Return the physical address (offset from base)
	return int(data), nil
}

func (mmu *MMU) Write(physicalAddr uint32, value byte) error {
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= PageSize {
		return fmt.Errorf("offset: address out of bounds")
	}
	if int(pfn) >= NumFrames {
		return fmt.Errorf("PFN: address out of bounds")
	}

	mmu.memory.Frames[pfn][offset] = value

	// Return the physical address (offset from base)
	return nil
}

func (mmu *MMU) StoreInstruction(pc int, opType int, opcode int, value byte) error {
	physicalAddr, err := mmu.TranslateAddress(uint32(pc))
	if physicalAddr == -1 {
		fmt.Println(err)
	}
	pfn := uint16(physicalAddr >> 16)
	offset := uint16(physicalAddr & 0xFFFF)

	if offset < 0 || offset >= PageSize {
		return fmt.Errorf("offset: address out of bounds")
	}
	if int(pfn) >= NumFrames {
		return fmt.Errorf("PFN: address out of bounds")
	}

	instructionByte := (uint8(opType) << 7) | (uint8(opcode) & 0x7F)

	mmu.memory.Frames[pfn][offset] = instructionByte
	mmu.memory.Frames[pfn][offset+1] = uint8(value)

	// Return the physical address (offset from base)
	return nil
}
