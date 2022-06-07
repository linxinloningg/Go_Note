package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"part6/src/constcode"
	"part6/src/utils"
)

//返回区块目标难度值
func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcode.Difficulty))
	return target.Bytes()
}

//获取区块信息总和
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		utils.ToHexInt(int64(nonce)),
		b.Target,
		b.SerializeTransaction(),
	},
		[]byte{},
	)
	return data
}

/*
必须计算出连续17个`0`开头的哈希值，矿工先确定PrevHash，MerkleHash，Timestamp，bits，
然后，不断变化`nonce`来计算哈希，直到找出连续17个`0`开头的哈希值。
我们可以大致推算一下，17个十六进制的`0`相当于计算了1617次，大约需要计算2.9万亿亿次
*/
//寻找一个合适正确的nonce
func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	var nonce int64
	nonce = 0
	intTarget.SetBytes(b.Target)

	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

//快速验证POW
func (b *Block) ValidatePoW() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTarget) == -1 {
		return true
	}
	return false
}
