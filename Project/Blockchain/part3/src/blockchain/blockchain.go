package blockchain

import (
	"encoding/hex"
	"fmt"
	"part3/src/transaction"
	"part3/src/utils"
)

//链结构体
type BlockChain struct {
	Blocks []*Block
}

//添加区块
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

/*
在真实区块链中，一个节点会维护一个候选区块，候选区块会维持一个交易信息池（Transaction Pool），
然后在挖矿时将交易池中的交易信息打包进行挖矿（PoW过程）。
*/
func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}

//创建区块链
func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}

//根据目标地址寻找可用交易信息
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	/*
		unSpentTxs就是我们要返回包含指定地址的可用交易信息的切片。
		spentTxs用于记录遍历区块链时那些已经被使用的交易信息的Output，
		key值为交易信息的ID值（需要转成string），value值为Output在该交易信息中的序号
	*/
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int) // 不能使用类型 []byte 作为键值

	//从最后一个区块开始向前遍历区块链，然后遍历每一个区块中的交易信息
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]

		for _, TX := range block.Transactions {
			ID := hex.EncodeToString(TX.ID)

			//检查当前交易信息是否为Base Transaction（主要是它没有input），
			//如果不是就检查当前交易信息的input中是否包含目标地址，有的话就将指向的Output信息加入到spentTxs中
			if !TX.IsBase() {
				for _, txInput := range TX.Inputs {
					//是否包含目标地址
					if txInput.FromAddressRight(address) {
						TxID := hex.EncodeToString(txInput.TxID)
						spentTxs[TxID] = append(spentTxs[TxID], txInput.OutIdx)
					}
				}
			}

		IterOutputs:
			//遍历交易信息的Output，如果该Output在spentTxs中就跳过，说明该Output已被消费
			for i, txOutput := range TX.Outputs {
				if spentTxs[ID] != nil {
					for _, outIdx := range spentTxs[ID] {
						if outIdx == i {
							continue IterOutputs
						}
					}
				}

				//否则确认ToAddress正确与否，正确就是我们要找的可用交易信息
				if txOutput.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *TX)
				}
			}

		}
	}
	return unSpentTxs
}

//找到一个地址的所有UTXO以及该地址对应的资产总和
func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, TX := range unspentTxs {
		ID := hex.EncodeToString(TX.ID)
		for outIdx, out := range TX.Outputs {
			if out.ToAddressRight(address) {
				accumulated += out.Value
				unspentOuts[ID] = outIdx
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}
	return accumulated, unspentOuts
}

//找到资产总量大于本次交易转账额的一部分UTXO
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, TX := range unspentTxs {
		ID := hex.EncodeToString(TX.ID)
		for i, txOutput := range TX.Outputs {
			if txOutput.ToAddressRight(address) && accumulated < amount {
				accumulated += txOutput.Value
				unspentOuts[ID] = i
				if accumulated >= amount {
					break Work
				}
				continue Work // 一笔交易只能有一个输出引用地址
			}
		}
	}
	return accumulated, unspentOuts
}

//创建交易
//可以用一个输入对于多个输出
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	accumulated, unspentOuts := bc.FindSpendableOutputs(from, amount)

	//没有足够数量的余额
	if accumulated < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}

	//转
	for ID, i := range unspentOuts {
		txID, err := hex.DecodeString(ID)
		utils.Handle(err)
		input := transaction.TxInput{txID, i, from}
		inputs = append(inputs, input)
	}

	//收
	output := transaction.TxOutput{amount, to}
	outputs = append(outputs, output)

	//找零
	if accumulated > amount {
		output := transaction.TxOutput{accumulated - amount, from}
		outputs = append(outputs, output)
	}

	//一个输入对应多个输出
	tx := transaction.Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx, true
}
