package translator

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/settings"
)

// vpn and frame is the same thing in this context
func TranslateVPNandOffsetToAddress(vpn int, offset int) uint32 {
	address := (uint32(vpn) << 16) | uint32(offset)
	logger.Log.Printf("VPN: %d\n Offset: %d\n Address: %d", vpn, offset, address)
	return address
}

func TranslateAddressToVPNandOffset(address int) (int, int) {
	vpn := uint16(address >> 16)
	offset := uint16(address & 0xFFFF)
	logger.Log.Printf("Address: %d\n VPN: %d\n Offset: %d", address, vpn, offset)
	return int(vpn), int(offset)
}

func FindNextInstructionAddress(memType int, pc int) int {
	if memType == 0 {
		vpn, offset := TranslateAddressToVPNandOffset(pc)
		if offset >= settings.PageSize {
			address := TranslateVPNandOffsetToAddress(vpn+1, 0)
			return int(address)
		} else {
			return pc
		}

	} else {
		return pc
	}
}

// TODO: Implement this function
// func TranslateInstructionIntToInstruction(intructionValue int) (int, string, int) {

// 	return int(vpn), int(offset)
// }
