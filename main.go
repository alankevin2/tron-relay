package main

import (
	"fmt"

	"gitlab.inlive7.com/crypto/tron-relay/pkg/api"
)

func main() {
	fmt.Println(api.GetBalance(2, "TJZeQH5YT7wUVvFiKB8ofoQCvX9RAhyLyQ"))
	pk, pub, addr := api.CreateNewAccount()
	fmt.Println(pk, pub, addr)
	fmt.Println(api.GetBalance(1, "TJZeQH5YT7wUVvFiKB8ofoQCvX9RAhyLyQ"))

	// fmt.Println(api.GetBalance(1, addr))
}
