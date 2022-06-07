package main

import (
	"fmt"
	"part1/src/blockchain"
)

func main() {
	bc := blockchain.CreateBlockChain()
	bc.AddBlock("第一个区块")

	for _, block := range bc.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)

	}

}
