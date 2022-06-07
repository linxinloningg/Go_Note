package blockchain

import (
	"fmt"
	"part5/src/utils"
)

/*
在真实区块链中，一个节点会维护一个候选区块，候选区块会维持一个交易信息池（Transaction Pool），
然后在挖矿时将交易池中的交易信息打包进行挖矿（PoW过程）。
*/
//挖矿
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	//在不久的将来，我们必须先在这里验证交易。
	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx) //打包交易信息，挖矿寻找nonce.
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("区块有无效的 nonce.")
		return
	}
}
