package tinyeth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"strings"
)

var (
	zero = big.NewInt(0)
)

// Client represents an Ethereum client with some helper methods for apification
type Client struct {
	PeerUrl string
	// client is the implementation of ethereum client
	client    *ethclient.Client
	rpcClient *rpc.Client
}

// Connect establish a new connection with provided ethereum peer URL
func (ctl *Client) Connect(nodeUrl string) error {
	rpcClient, err := rpc.DialContext(context.Background(), nodeUrl)
	if err != nil {
		return err
	}
	ctl.client = ethclient.NewClient(rpcClient)
	ctl.rpcClient = rpcClient
	return nil
}

// LatestBlock returns the latest block of the Blockchain
func (ctl *Client) LatestBlock() (string, error) {
	if ctl.client == nil {
		_ = ctl.Init()
	}
	// will call eth_getBlockByNumber
	header, err := ctl.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return "unknown", err
	}
	return header.Number.String(), nil
}

// PendingCallContract makes a smart contract method call
// Executes a new message call immediately without creating a transaction on the block chain.
// result - the return value of executed contract.
func (ctl *Client) pendingCallContract(ctx context.Context, data *ethereum.CallMsg) ([]byte, error) {
	if ctl.client == nil {
		if err := ctl.Init(); err != nil {
			return nil, err
		}
	}
	// makes a JSON rpc call with eth_call method
	return ctl.client.PendingCallContract(ctx, *data)
}

func (ctl *Client) suggestGasPrice(ctx context.Context) (*big.Int, error) {
	if ctl.client == nil {
		if err := ctl.Init(); err != nil {
			return nil, err
		}
	}
	return ctl.client.SuggestGasTipCap(ctx)
}

// Init initializes EthereumController internal values and parameters
func (ctl *Client) Init() error {
	return ctl.Connect(ctl.PeerUrl)
}

// QueryBlockchain makes a blockchain query to read data from it
func (ctl *Client) QueryBlockchain(ctx context.Context, from *Account, to common.Address, abiData *abi.ABI, methodName string, params ...any) ([]interface{}, error) {
	if ctl.client == nil {
		if err := ctl.Init(); err != nil {
			return nil, err
		}
	}
	calldata, method, err := ctl.prepareParams(ctx, abiData, *from.Address, to, methodName, true, params...)
	if err != nil {
		return nil, err
	}
	result, err := ctl.client.CallContract(ctx, calldata, nil)
	if err != nil {
		return nil, err
	}
	// now we need to decode contract output using the abi
	return method.Outputs.Unpack(result)
}

// SendTransaction makes a blockchain transaction to send new data to it
func (ctl *Client) SendTransaction(ctx context.Context, acc *Account, to common.Address, abiData *abi.ABI, methodName string, params ...any) (string, error) {
	if ctl.client == nil {
		if err := ctl.Init(); err != nil {
			return "", err
		}
	}
	// 1 setup transaction content payload abi-encoded
	from := *acc.Address
	txData, _, err := ctl.prepareParams(ctx, abiData, from, to, methodName, false, params...)
	if err != nil {
		return "", err
	}
	// 3 get the sender address to compute its nonce
	nonce, err := ctl.client.PendingNonceAt(ctx, from)
	if err != nil {
		return "", err
	}
	// 4 fetch connected network id
	chainID, err := ctl.client.ChainID(ctx)
	if err != nil {
		return "", err
	}
	// 5 create the tx signer
	signer := types.NewLondonSigner(chainID)
	// 6 sign the transaction
	txTemp := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: txData.GasPrice,
		Gas:      txData.Gas,
		To:       txData.To,
		Value:    txData.Value,
		Data:     txData.Data,
	})
	// Add V, R and S values as signature
	signedTx, err := types.SignTx(txTemp, signer, acc.PrivateKey)
	if err != nil {
		return "", err
	}
	// 7 verify the sender
	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return "", err
	}
	// safety check here
	if sender != from {
		log.Fatalf("signer mismatch: expected %s, got %s", from.Hex(), sender.Hex())
	}
	// another safety check
	v, r, s := signedTx.RawSignatureValues()
	vraw := v.Bytes()
	if crypto.ValidateSignatureValues(vraw[0], r, s, false) {
		return "", errors.New("signature invalid")
	}
	// 7 send the signed transaction
	// If the transaction was a contract creation use the TransactionReceipt method to get the
	// contract address after the transaction has been mined.
	err = ctl.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().String(), err
}

// ParseAbi attempts to parse input JSON abi to object
func (ctl *Client) ParseAbi(abiData string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(abiData))
}

// PrepareParams prepare transaction params according to given ABI specification
func (ctl *Client) prepareParams(ctx context.Context, abi *abi.ABI, from, to common.Address, methodName string, isQuery bool, params ...any) (ethereum.CallMsg, *abi.Method, error) {
	/*
		Parameters
		Object - The transaction call object
		from: DATA, 20 Bytes - (optional) The address the transaction is sent from.
		to: DATA, 20 Bytes - The address the transaction is directed to.
		gas: QUANTITY - (optional) Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
		gasPrice: QUANTITY - (optional) Integer of the gasPrice used for each paid gas
		value: QUANTITY - (optional) Integer of the value sent with this transaction
		data: DATA - (optional) Hash of the method signature and encoded parameters. For details see Ethereum Contract ABI
		QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending", see the default block parameter
	*/
	msg := ethereum.CallMsg{
		From: from,
		To:   &to,
		// The gas limit for a standard ETH transfer is 21000 units.
		Gas:   0,
		Value: zero, // no value
	}
	// on queries we set gas to zero for unlimited execution
	if !isQuery {
		// now estimate call gas
		gasEstimated, err := ctl.client.EstimateGas(ctx, msg)
		if err != nil {
			switch err.Error() {
			case "base fee exceeds gas limit":
				return msg, nil, err
			case "VM Exception while processing transaction: revert":
				// return msg, nil, err
				msg.Gas = 210000
			default:
				return msg, nil, err
			}
		} else {
			msg.Gas = gasEstimated
			msg.GasPrice = big.NewInt(30000000000) // in wei (30 gwei),
		}
		// suggest a better price based on last block data
		gasPrice, err := ctl.client.SuggestGasPrice(ctx)
		if err != nil {
			return msg, nil, err
		}
		msg.GasPrice = gasPrice
		msg.GasFeeCap = gasPrice
		// gas tip cap
		gasTipCap, err := ctl.client.SuggestGasTipCap(ctx)
		if err != nil {
			switch err.Error() {
			case "Method eth_maxPriorityFeePerGas not supported.":
				// do not propagate the error and set the same value as gas
				// do not set a gas tip cap if not supported
				msg.GasTipCap = nil
			default:
				return msg, nil, err
			}
		} else {
			// if method is supported use the result as gas tip cap
			msg.GasTipCap = gasTipCap
		}
	}

	// pick user selected method name
	targetMethod, found := abi.Methods[methodName]
	if !found {
		return msg, nil, errors.New("method not found in ABI specification")
	}
	id := targetMethod.ID
	fmt.Println("target method id:", hex.EncodeToString(id))
	if params != nil && len(params) > 0 {
		inputs, err := targetMethod.Inputs.Pack(params...)
		if err != nil {
			return msg, nil, err
		}
		// append transaction method id to packed data (4bytes stuff)
		inputs = append(id[:], inputs[:]...)
		msg.Data = inputs
	} else {
		msg.Data = id[:]
	}
	return msg, &targetMethod, nil
}
