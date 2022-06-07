package main

import (
	"fmt"
	"part3/src/blockchain"
)

func main() {

	chain := blockchain.CreateBlockChain()
	property, _ := chain.FindUTXOs([]byte("创始人"))
	fmt.Println("创始人的余额: ", property)

}
