package api

import (
	"math/big"

	"gitlab.inlive7.com/crypto/tron-relay/config"
	"gitlab.inlive7.com/crypto/tron-relay/internal/relay"
	"gitlab.inlive7.com/crypto/tron-relay/pkg/types"
)

type TronRelay interface {
	// GetBalance(uint8, string) uint64
	GetBalanceForToken(chainID uint16, address string, symbol string) (balance *big.Int, decimal uint8, err error)
	CreateNewAccount() (privateKey string, publicKey string, publicAddress string)
	QueryTransaction(chainID uint16, txn string) (*types.TransactionState, bool, error)
	// TransferValueUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error)
	TransferTokenUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error)
	GetFeeLimit(chainID uint16) (uint64, error)
	GetTokenAddress(symbol string) (address string)
	InitRelay(chainIds []config.ChainID)
}

// func GetBalance(chainID uint16, address string) (balance int64, err error) {
// 	return relay.Shared(config.ChainID(chainID)).GetBalance(address)
// }

func GetBalanceForToken(chainID uint16, address string, symbol string) (balance *big.Int, decimal uint8, err error) {
	return relay.Shared(config.ChainID(chainID)).GetBalanceForToken(address, symbol)
}

func CreateNewAccount() (privateKey string, publicKey string, publicAddress string) {
	// walletClient api required online, we use chainID 1 for the mainnet
	// account should be activated seperately on mainnet and testnet
	pk, pb, addr, _ := relay.Shared(1).CreateNewAccount()
	return pk, pb, addr
}

func QueryTransaction(chainID uint16, txn string) (t *types.TransactionState, isPending bool, err error) {
	return relay.Shared(config.ChainID(chainID)).QueryTransaction(txn)
}

// func TransferValueUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error) {
// 	return relay.Shared(config.ChainID(chainID)).TransferValueUsingPrivateKey(privateKey, data)
// }

func TransferTokenUsingPrivateKey(chainID uint16, privateKey string, data *types.TransactionRaw) (hash string, err error) {
	return relay.Shared(config.ChainID(chainID)).TransferTokenUsingPrivateKey(privateKey, data)
}

func GetFeeLimit(chainID uint16) (uint64, error) {
	return relay.Shared(config.ChainID(chainID)).GetFeeLimit()
}

func GetTokenAddress(chainID uint16, symbol string) (address string) {
	return relay.Shared(config.ChainID(chainID)).GetTokenAddress(symbol)
}

func InitRelay(chainIds []config.ChainID) {
	for i := range chainIds {
		// first time call Shared inits the instance
		relay.Shared(chainIds[i])
	}
}
