package tinyeth

import (
	"github.com/ethereum/go-ethereum/common"
)

// SimpleHex returns a normalized version of a hexadecimal string.
// without prefix ("0x").
func SimpleHex(addrHex string) string {
	if addrHex[1] == 'x' || addrHex[1] == 'X' {
		return addrHex[2:]
	}
	return addrHex
}

// BuildAddress builds a common.Address from a hexadecimal string.
func BuildAddress(addrHex string) common.Address {
	return common.HexToAddress(addrHex)
}

func ParseAddress(addr string) common.Address {
	return common.HexToAddress(addr)
}
