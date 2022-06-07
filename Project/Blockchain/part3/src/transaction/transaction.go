package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"part3/src/constcode"
	"part3/src/utils"
)

type Transaction struct {
	ID      []byte     //自身的ID值（其实就是哈希值）
	Inputs  []TxInput  //用于标记支持我们本次转账的前置的交易信息的TxOutput
	Outputs []TxOutput //TxOutput记录我们本次转账的amount和Reciever
}

//TxHash返回交易信息的哈希值
func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

//SetID设置每个交易信息的ID值，也就是哈希值
func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

//创区块交易
func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}}
	txOut := TxOutput{constcode.InitCoin, toaddress}
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

//检查是否为创始交易
func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

