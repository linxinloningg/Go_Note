package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"part7/src/constcode"
	"part7/src/utils"
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
	txIn := TxInput{[]byte{}, -1, []byte{}, nil}
	txOut := TxOutput{constcode.InitCoin, toaddress}
	tx := Transaction{[]byte("这是创始区块交易！"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

//检查是否为创始交易
func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}


/*
用户想要使用指向特定地址的资产时就需要通过签名来证明自己对这些地址的拥有权。
想像这样的一个场景，用户A拥有a这一钱包，用户B拥有b这一钱包。
在区块链中有3个UTXO流向了a对应的公钥哈希地址，总值为5币。
现在用户A想要转账5币给用户B，需要生成交易信息，于是A便构建三个Input来引用前述的3个UTXO，
同时将用户B提供的b钱包地址计算得到的公钥哈希作为Output的地址。A为了证明对3个UTXO的使用权，
需要使用私钥对整个交易过程进行签名，并将签名信息作为交易信息的一部分向整个区块链扩散。
为了验证这样的一个交易信息的有效性，需要同时获得交易信息代表的交易过程，
三个UTXO指向的公钥哈希地址（也即是a所对应的哈希公钥），A的公钥，A提供的签名信息。
 */

//描述一个交易信息的交易过程
func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutIdx, nil, nil})
	}

	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

//PlainHash函数用以辅助对交易信息进行签名
func (tx *Transaction) PlainHash(inidx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKey
	return txCopy.TxHash()
}

//交易信息的签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey) // 这是因为我们要单独对输入进行签名！
		signature := utils.Sign(plainhash, privKey)
		tx.Inputs[idx].Sig = signature
	}
}

//交易信息的验证
func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		if !utils.Verify(plainhash, input.PubKey, input.Sig) {
			return false
		}
	}
	return true
}