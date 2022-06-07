package blockchain

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"part5/src/constcode"
	"part5/src/transaction"
	"part5/src/utils"
)

//交易信息池的结构体
type TransactionPool struct {
	PubTx []*transaction.Transaction //PubTx用于储存节点收集到的交易信息
}

//创建交易信息池
func CreateTransactionPool() *TransactionPool {
	transactionPool := TransactionPool{}
	err := transactionPool.LoadFile()
	utils.Handle(err)
	return &transactionPool
}

//清空交易信息池
func RemoveTransactionPoolFile() error {
	err := os.Remove(constcode.TransactionPoolFile)
	return err
}

//添加新交易信息
func (tp *TransactionPool) AddTransaction(tx *transaction.Transaction) {
	tp.PubTx = append(tp.PubTx, tx)
}

/*
每次都将交易信息池保存到constcode.TransactionPoolFile这个地址中。
0644是指八进制的644（110，100，100），指明了不同用户对文件读写执行的权限
*/
//交易信息池保存
func (tp *TransactionPool) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(tp)
	utils.Handle(err)
	err = ioutil.WriteFile(constcode.TransactionPoolFile, content.Bytes(), 0644)
	utils.Handle(err)
}

//交易信息池加载
func (tp *TransactionPool) LoadFile() error {
	if !utils.FileExists(constcode.TransactionPoolFile) {
		return nil
	}

	var transactionPool TransactionPool

	fileContent, err := ioutil.ReadFile(constcode.TransactionPoolFile)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&transactionPool)

	if err != nil {
		return err
	}

	tp.PubTx = transactionPool.PubTx
	return nil
}


