package translator

import (
	"CPU-Simulator/simulator/pkg/logger"
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

// TODO: Implement this function
// func TranslateInstructionIntToInstruction(intructionValue int) (int, string, int) {

// 	return int(vpn), int(offset)
// }
