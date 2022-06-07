package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"part6/src/merkletree"
	"part6/src/transaction"
	"part6/src/utils"
	"time"
)

//区块的结构体
type Block struct {
	Timestamp    int64                      //时间戳
	Hash         []byte                     //本身的哈希值
	PrevHash     []byte                     //指向上一个区块的哈希
	Target       []byte                     //目标难度值
	Nonce        int64                      //POW
	Transactions []*transaction.Transaction //交易事务
	MTree        *merkletree.MerkleTree     //MT
}

//协助处理区块中交易信息的序列化
func (b *Block) SerializeTransaction() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

//哈希构造函数
func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.SerializeTransaction(), b.MTree.RootNode.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

//区块创建
func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, txs, merkletree.CrateMerkleTree(txs)}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

//创世区块创建
func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte{}, []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}

//序列化区块生成字节串
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

//反序列化区块
func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
