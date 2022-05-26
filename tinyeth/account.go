package tinyeth

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents an Ethereum account
type Account struct {
	Address    *common.Address
	PrivateKey *ecdsa.PrivateKey
}

// CreateAccount creates a new Ethereum account (address and private key)
func CreateAccount() (*Account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	addr := common.HexToAddress(address)
	return &Account{&addr, privateKey}, nil
}

// LoadAccount loads the account from private key
func LoadAccount(key string) (*Account, error) {
	key = SimpleHex(key)
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &Account{&fromAddress, privateKey}, nil
}

func (a *Account) Sign() {

}

// Export returns the address and the private key of the account hex encoded
func (a *Account) Export() (string, string) {
	privateKeyBytes := crypto.FromECDSA(a.PrivateKey)
	keyBytes := hexutil.Encode(privateKeyBytes)
	addrHex := a.Address.Hex()
	return addrHex, keyBytes
}
