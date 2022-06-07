# 深入底层：Go语言从零构建区块链（二）：PoW工作证明机制

## 共识机制

我们常说区块链是一个分布式系统，系统中每个节点都有机会储存数据信息构造一个区块然后追加到区块链尾部。这里就存在一个问题，那就是当区块链系统中有多个节点都想将自己的区块追加到区块链是我们该怎么办？我们将这些等待添加的区块统称为候选区块，显然我们不能对候选区块全盘照收，否则区块链就不再是一条链而是不同分叉成区块树。那么我们如何确定一种方法来从候选区块中选择一个加入到区块链中了？这里就需要用到区块链的共识机制，后文将以比特币使用的最经典PoW共识机制进行讲解。

共识机制说的通俗明白一点就是要在相对公平的条件下让想要添加区块进区块链的节点内卷，通过竞争选择出一个大家公认的节点添加它的区块进入区块链。整个共识机制被分为两部分，首先是竞争，然后是共识。中本聪在比特币中设计了如下的一个Game来实现竞争：每个节点去寻找一个随机值（也就是nonce），将这个随机值作为候选区块的头部信息属性之一，要求候选区块对自身信息（注意这里是包含了nonce的）进行哈希后表示为数值要小于一个难度目标值（也就是Target），最先寻找到nonce的节点即为卷王，可以将自己的候选区块发布并添加到区块链尾部。这个Game设计的非常巧妙，首先每个节点要寻找到的nonce只对自己候选区块有效，防止了其它节点同学抄答案；其次，nonce的寻找是完全随机的没有技巧，寻找到nonce的时间与目标难度值与节点本身计算性能有关，但不妨碍性能较差的节点也有机会获胜；最后寻找nonce可能耗费大量时间与资源，但是验证卷王是否真的找到了nonce却非常却能够很快完成并几乎不需要耗费资源，这个寻找到的nonce可以说就是卷王真的是卷王的证据。现在我们就来一步一步实现这个Game。

## 添加Nonce

如前文所说，我们要先增加一些区块的头部信息。

```go
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
```

Nonce就是节点寻找到的作为POW的验证。Target就是我们前文说到的目标难度值，将它保存到区块中便于其他节点快速验证Nonce是否正确。

这样一来之前创建的几个函数会报错，我们先暂时不理会。

## PoW实现

新建一个constcode文件夹，在constcode文件下新建constcode.go用于存放一下全局变量

```go
package constcode

const (
	Difficulty = 12
)
```

在blockchain文件夹下创建proofofwork.go，我们来实现之前说到的功能。首先引入以下包。

```go
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"part2/src/constcode"
	"part2/src/utils"
)
```

我们现在来构建一个可以返回目标难度值的函数。我们这里使用的之前设定的一个常量Difficulty来构造目标难度值，但是在实际的区块链中目标难度值会根据网络情况定时进行调整，且能够保证各节点在同一时间在同一难度下进行竞争，故这里的GetTarget可以理解为预留API，期待一下之后的分布式网络实现。

```go
//返回区块目标难度值
func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcode.Difficulty))
	return target.Bytes()
}
```

Lsh函数就是向左移位，移的越多目标难度值越大，哈希取值落在的空间就更多就越容易找到符合条件的nonce。

每次我们输入一个nonce对应的区块的哈希值都会改变，如下。

```go
//获取区块信息总和
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		utils.ToHexInt(int64(nonce)),
		b.Target,
		b.Data,
	},
		[]byte{},
	)
	return data
}
```

现在对于任意一个区块，我们都能去寻找一个合适的nonce了。

举个例子 ：

必须计算出连续17个`0`开头的哈希值，矿工先确定PrevHash，MerkleHash，Timestamp，bits，
然后，不断变化`nonce`来计算哈希，直到找出连续17个`0`开头的哈希值。
我们可以大致推算一下，17个十六进制的`0`相当于计算了1617次，大约需要计算2.9万亿亿次

```go
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
```

可以看到，神秘的nonce不过是从0开始取的整数而已，随着不断尝试，每次失败nonce就加1直到由当前nonce得到的区块哈希转化为数值小于目标难度值为止。

我们再来实现一个快速验证POW函数。

```go
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
```

好了，PoW我们已经实现了。回到block.go，调整以下函数。

```go
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
```

## 调试带PoW的区块链系统

现在打开main.go，我们可以编写程序启动我们的区块链系统了。

```go
//main.go

package main

import (
	"fmt"
	"part2/src/blockchain"
)

func main() {
	chain := blockchain.CreateBlockChain()
	chain.AddBlock("第一个区块")

	for _, block := range chain.Blocks {

		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
		fmt.Println("")

	}
}

```



