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
	t, isPending, err := api.QueryTransaction(65534, "e722bbf8949f009ce60e61ca0d391c728a2b6695671733af27020f0ada0d8687")
	fmt.Println(t, isPending, err)

	balance, decimal, err := api.GetBalanceForToken(65534, "TK7q7c6RRSjTvuzmVmZNgq18nQrmx1UZtc", "USDT")
	fmt.Println(balance, decimal, err)

	fmt.Println(api.GetFeeLimit(65534))
}
