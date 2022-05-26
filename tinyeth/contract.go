package tinyeth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Contract represents the target contract to interact
type Contract struct {
	Address common.Address
	Abi     *abi.ABI
}
