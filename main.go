package main

import (
	"fmt"

	"gitlab.inlive7.com/crypto/tron-relay/pkg/api"
)

func main() {
	// fmt.Println(api.GetBalance(2, "TJZeQH5YT7wUVvFiKB8ofoQCvX9RAhyLyQ"))
	// pk, pub, addr := api.CreateNewAccount()
	// fmt.Println(pk, pub, addr)

	// txid, err := api.TransferTokenUsingPrivateKey(2, "47c8c6a372e106d0095ff223c89418095de5f36f501f7b75b345bfe8f2cca9fe", &types.TransactionRaw{
	// 	Value:          new(big.Int).SetInt64(3333333),
	// 	To:             "TH6QdQ3jkBj76JQvV42SkqTHsTx1DnQuEK",
	// 	PreferredLimit: 1000000000,
	// 	TokenSymbol:    "USDT",
	// })
	// fmt.Println(txid, err)

	t, err := api.QueryTransaction(2, "44a2a52c90f4fab6bb8db87200d45c9fc1dcbb6cc2a0fb06b16e2004ba2315e9")
	fmt.Println(t, err)
}
