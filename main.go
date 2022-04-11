package main

import (
	"fmt"

	"gitlab.inlive7.com/crypto/tron-relay/pkg/api"
)

func main() {
	fmt.Println(api.GetBalance(1, "TJZeQH5YT7wUVvFiKB8ofoQCvX9RAhyLyQ"))
	// pk, pub, addr := api.CreateAccount()
	// fmt.Println(pk, pub, addr)
	// fmt.Println(api.GetBalance(1, addr))
}
