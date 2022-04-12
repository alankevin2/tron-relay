package config

import (
	"fmt"
)

type ChainID uint16

const (
	Mainnet ChainID = 1
	Shasta  ChainID = 2
)

type ChainInfo struct {
	Name         string
	ID           ChainID
	ProviderFile string
	Decimal      int8
}

func RetrieveChainInfo(id ChainID) (ChainInfo, error) {
	var info ChainInfo
	switch id {
	case Mainnet:
		info = ChainInfo{"mainnet", Mainnet, "provider-mainnet.yml", 6}
	case Shasta:
		info = ChainInfo{"shasta", Shasta, "provider-testnet-shasta.yml", 6}
	default:
		return info, fmt.Errorf("no support yet for chain id : %d", id)
	}

	return info, nil
}
