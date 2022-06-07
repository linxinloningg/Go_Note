package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"part8/src/constcode"
	"part8/src/transaction"
	"part8/src/utils"
	"runtime"
)

//链结构体
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

//添加区块
func (bc *BlockChain) AddBlock(newBlock *Block) {
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)

		return err
	})
	utils.Handle(err)
	if !bytes.Equal(newBlock.PrevHash, lastHash) {
		fmt.Println("此区块已过期")
		runtime.Goexit()
	}

	err = bc.Database.Update(func(transaction *badger.Txn) error {
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

//初始化区块链并创建一个数据库保存
func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte

	if utils.FileExists(constcode.BCFile) {
		fmt.Println("区块链已经存在")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcode.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("创世区块创建")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash) //store the hash of the block in blockchain
		utils.Handle(err)
		err = txn.Set([]byte("ogprevhash"), genesis.PrevHash) //store the prevhash of genesis(original) block
		utils.Handle(err)
		lastHash = genesis.Hash
		return err
	})
	utils.Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

//通过已有的数据库读取并加载区块链
func ContinueBlockChain() *BlockChain {
	if utils.FileExists(constcode.BCFile) == false {
		fmt.Println("没有找到区块链，请先创建一个")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constcode.BCPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

//根据目标地址寻找可用交易信息
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int) // can't use type []byte as key value

	iter := bc.Iterator()

all:
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}

				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}
		if bytes.Equal(block.PrevHash, bc.BackOgPrevHash()) {
			break all
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
func (bc *BlockChain) CreateTransaction(from_PubKey, to_HashPubKey []byte, amount int, privkey ecdsa.PrivateKey) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	accumulated, validOutputs := bc.FindSpendableOutputs(from_PubKey, amount)

	//没有足够数量的余额
	if accumulated < amount {
		fmt.Println("没有足够数量的余额!")
		return &transaction.Transaction{}, false
	}

	//转
	for ID, i := range validOutputs {
		txID, err := hex.DecodeString(ID)
		utils.Handle(err)
		input := transaction.TxInput{txID, i, from_PubKey, nil}
		inputs = append(inputs, input)
	}

	//收
	outputs = append(outputs, transaction.TxOutput{amount, to_HashPubKey})

	//找零
	if accumulated > amount {
		outputs = append(outputs, transaction.TxOutput{accumulated - amount, utils.PublicKeyHash(from_PubKey)})
	}

	//一个输入对应多个输出
	tx := transaction.Transaction{nil, inputs, outputs}

	tx.SetID()

	tx.Sign(privkey)
	return &tx, true
}

//根据所给的钱包地址返回区块链中所有该地址的UTXO
func (bc *BlockChain) BackUTXOs(address []byte) []transaction.UTXO {
	var UTXOs []transaction.UTXO
	unspentTxs := bc.FindUnspentTransactions(address)

Work:
	for _, tx := range unspentTxs {
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				UTXOs = append(UTXOs, transaction.UTXO{tx.ID, outIdx, out})
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}

	return UTXOs
}

//获取当前区块链最后一个区块
func (chain *BlockChain) GetCurrentBlock() *Block {
	var block *Block
	err := chain.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get(chain.LastHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
	return block
}

//获取区块高度
func (bc *BlockChain) BackHeight() int64 {
	return bc.GetCurrentBlock().Height
}
