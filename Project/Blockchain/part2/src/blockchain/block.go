package blockchain

import (
	"bytes"
	"crypto/sha256"
	"part2/src/utils"
	"time"
)

//区块的结构体
type Block struct {
	Timestamp int64  //时间戳
	Hash      []byte //本身的哈希值
	PrevHash  []byte //指向上一个区块的哈希
	Target    []byte //目标难度值
	Nonce     int64  //POW
	Data      []byte //区块中的数据
}

//哈希构造函数
func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

//区块创建
func CreateBlock(prevhash, data []byte) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, data}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

//创世区块创建
func GenesisBlock() *Block {
	genesisWords := "创世区块"
	return CreateBlock([]byte{}, []byte(genesisWords))
}
