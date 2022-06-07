package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"part6/src/transaction"
	"part6/src/utils"
)


//验证交易池中的交易信息有效性
func (bc *BlockChain) VerifyTransactions(txs []*transaction.Transaction) bool {
	if len(txs) == 0 {
		return true
	}
	//TODO: The following method to verify the transactions is to query the blockchain for
	//unspent outputs and is definitely right. However, I believe in the near future, an unspent
	//outputs database can be maintained to accelerate the verification.
	spentOutputs := make(map[string]int)
	for _, tx := range txs {
		pubKey := tx.Inputs[0].PubKey
		unspentOutputs := bc.FindUnspentTransactions(pubKey)
		inputAmount := 0
		OutputAmount := 0

		for _, input := range tx.Inputs {
			if outidx, ok := spentOutputs[hex.EncodeToString(input.TxID)]; ok && outidx == input.OutIdx {
				return false
			}
			ok, amount := isInputRight(unspentOutputs, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TxID)] = input.OutIdx
		}

		for _, output := range tx.Outputs {
			OutputAmount += output.Value
		}
		if inputAmount != OutputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}
	}
	return true
}

func isInputRight(txs []transaction.Transaction, in transaction.TxInput) (bool, int) {
	for _, tx := range txs {
		if bytes.Equal(tx.ID, in.TxID) {
			return true, tx.Outputs[in.OutIdx].Value
		}
	}
	return false, 0
}

/*
在真实区块链中，一个节点会维护一个候选区块，候选区块会维持一个交易信息池（Transaction Pool），
然后在挖矿时将交易池中的交易信息打包进行挖矿（PoW过程）。
*/
/*
Mine过程在构造候选区块时应该先检查交易信息池中的所有交易信息的有效性，
这包括验证是否引用了已花费的Output，是否重复引用了同一UTXO，Input与Output资产总额是否对应，交易信息的签名信息
 */
//挖矿
/*func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("属于交易验证")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}

	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx) //PoW has been done here.
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("区块拥有无效的 nonce.")
		return
	}
}*/
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}

	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx) //PoW has been done here.
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}

