package main

import (
	"fmt"
	"part3/src/blockchain"
	"part3/src/transaction"
)

func main() {
	//交易信息池
	txPool := make([]*transaction.Transaction, 0)

	chain := blockchain.CreateBlockChain()

	tempTx, status := chain.CreateTransaction([]byte("创始人"), []byte("linxinloningg"), 100)
	if status {
		txPool = append(txPool, tempTx)
	}

	//挖矿打包新区块
	chain.Mine(txPool)

	property, _ := chain.FindUTXOs([]byte("linxinloningg"))
	fmt.Println("Balance of linxinloningg: ", property)

}
