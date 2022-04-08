package main

import (
	"fmt"

	"gitlab.inlive7.com/crypto/tron-relay/pkg/api"
)

func main() {
	fmt.Println(api.GetBalance(1, "TNQ33j2jiST9AkJ1P5ZNvqa53mkVPmv7cN"))
	fmt.Println(api.GetBalance(2, "TNQ33j2jiST9AkJ1P5ZNvqa53mkVPmv7cN"))
	fmt.Println(api.GetBalance(2, "TNQ33j2jiST9AkJ1P5ZNvqa53mkVPmv7cN"))
	fmt.Println(api.GetBalance(2, "TNQ33j2jiST9AkJ1P5ZNvqa53mkVPmv7cN"))
}
