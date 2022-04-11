package api

import (
	"gitlab.inlive7.com/crypto/tron-relay/internal/relay"
)

type TronRelay interface {
	GetBalance(uint8, string) uint64
}

func GetBalance(chainID uint8, address string) (balance int64, err error) {
	return relay.Shared(chainID).GetBalance(address)
}

func CreateAccount() (privateKey string, publicKey string, publicAddress string) {
	// walletClient api required online, we use chainID 1 for the mainnet
	// account should be activated seperately on mainnet and testnet
	pk, pb, addr, _ := relay.Shared(1).CreateAccount()
	return pk, pb, addr
}

func InitRelay(chainIds []uint8) {
	for i := range chainIds {
		// first time call Shared inits the instance
		relay.Shared(chainIds[i])
	}
}
