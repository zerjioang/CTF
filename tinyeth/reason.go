package tinyeth

import (
	"context"
)

// Struct used for outlining errors that ocurred when mining a transaction.
// It makes the wrapper more verbose when something went wrong by giving
// as output the revert reason.
type ReceiptError struct {
	Status       bool   // The status of the transaction. If success true, false otherwise.
	RevertReason string // The revert reason given as a string.
}

// Struct used for outlining errors that ocurred when querying a contract.
// It makes the wrapper more verbose when something went wrong by giving
// as output the revert reason.
type QueryError struct {
	Code         int    `json:"code"`           // Error code.
	Message      string `json:"message"`        // Error message (usually 'Execution reverted').
	RevertReason string `json:"data,omitempty"` // Revert reason given as a string.
}

// Struct used for outlining errors that ocurred when querying a contract.
// It makes the wrapper more verbose when something went wrong by giving
// as output the revert reason.
type EstimateGasError struct {
	Code         int    `json:"code"`           // Error code.
	Message      string `json:"message"`        // Error message (usually 'Execution reverted').
	RevertReason string `json:"data,omitempty"` // Revert reason given as a string.
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ctl *Client) TransactionReceipt(ctx context.Context, txHash string) (interface{}, error) {
	dst := map[string]interface{}{}
	err := ctl.rpcClient.CallContext(ctx, &dst, "eth_getTransactionReceipt", txHash)
	return dst, err
}

// GetStorageAt returns slot information data
func (ctl *Client) GetStorageAt(ctx context.Context, address string, slot string) (interface{}, error) {
	var dst string
	err := ctl.rpcClient.CallContext(ctx, &dst, "eth_getStorageAt", address, slot)
	return dst, err
}
