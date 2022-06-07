package main

import (
	"fmt"
	"part3/src/blockchain"
	"part3/src/transaction"
)

func main() {
	txPool := make([]*transaction.Transaction, 0)
	var tempTx *transaction.Transaction
	var ok bool
	var property int
	chain := blockchain.CreateBlockChain()
	property, _ = chain.FindUTXOs([]byte("创始人"))
	fmt.Println("Balance of 创始人: ", property)

	tempTx, ok = chain.CreateTransaction([]byte("创始人"), []byte("第一个人"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	property, _ = chain.FindUTXOs([]byte("创始人"))
	fmt.Println("Balance of 创始人: ", property)

	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("第二个人"), 200) // this transaction is invalid
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("第二个人"), 50)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("创始人"), []byte("第二个人"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)
	property, _ = chain.FindUTXOs([]byte("创始人"))
	fmt.Println("Balance of 创始人: ", property)
	property, _ = chain.FindUTXOs([]byte("第一个人"))
	fmt.Println("Balance of 第一个人: ", property)
	property, _ = chain.FindUTXOs([]byte("第二个人"))
	fmt.Println("Balance of 第二个人: ", property)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
		fmt.Println("")
	}

	//I want to show the bug at this version.

	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("第二个人"), 30)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("创始人"), 30)
	if ok {
		txPool = append(txPool, tempTx)
	}

	chain.Mine(txPool)
	txPool = make([]*transaction.Transaction, 0)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}

	property, _ = chain.FindUTXOs([]byte("创始人"))
	fmt.Println("Balance of 创始人: ", property)
	property, _ = chain.FindUTXOs([]byte("第一个人"))
	fmt.Println("Balance of 第一个人: ", property)
	property, _ = chain.FindUTXOs([]byte("第二个人"))
	fmt.Println("Balance of 第二个人: ", property)
}
