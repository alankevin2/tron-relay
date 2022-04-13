package main

import (
	"gitlab.inlive7.com/crypto/tron-relay/pkg/api"

	"fmt"
)

func main() {
	// fmt.Println(api.GetBalance(2, "TJZeQH5YT7wUVvFiKB8ofoQCvX9RAhyLyQ"))
	pk, pub, addr := api.CreateNewAccount()
	fmt.Println(pk, pub, addr)
}
