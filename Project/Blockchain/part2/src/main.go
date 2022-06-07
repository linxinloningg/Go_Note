//main.go

package main

import (
	"fmt"
	"part2/src/blockchain"
)

func main() {
	chain := blockchain.CreateBlockChain()
	chain.AddBlock("第一个区块")

	for _, block := range chain.Blocks {

		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
		fmt.Println("")

	}
}
