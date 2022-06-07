# 深入底层：Go语言从零构建区块链（一）：Hello, Blockchain

## 区块与区块链

区块链以区块（block）的形式储存数据信息，一个区块记录了一段时间内系统或网络中产生的重要数据信息，区块通过引用上一个区块的hash值来连接上一个区块这样区块就按时间顺序排列形成了一条链。每个区块应该包含头部（head）信息用于总结性的描述这个区块，然后在区块的数据存放区（body）中存放要保存的重要数据。首先我们需要初始化main.go，并导入一些基本的包。

```go
package main

import (
	"fmt"
	"part1/src/blockchain"
)

func main {
	
}
```

新建一个文件夹blockchain，然后在底下新建一个block.go

开始定义区块的结构体。

```go
//区块的结构体
type Block struct {
	Timestamp int64  //时间戳
	Hash      []byte //本身的哈希值
	PrevHash  []byte //指向上一个区块的哈希
	Data      []byte //区块中的数据
}
```

我们定义的区块中有时间戳，本身的哈希值，指向上一个区块的哈希这三个属性构成头部信息，而区块中的数据以Data属性表示。

再在底下新建blockchain.go

在获得了区块后，我们可以定义区块链。

```go
//链结构体
type BlockChain struct {
	Blocks []*Block
}
```

可以看到我们这里的区块链就是区块的一个集合。好了，现在你已经掌握了区块与区块链了，现在就可以去搭建自己的区块链系统了。

## 哈希

现在来给我们的区块增加点细节，来看看它们是怎么连接起来的。对于一个区块而言，可以通过哈希算法概括其所包含的所有信息，哈希值就相当于区块的ID值，同时也可以用来检查区块所包含信息的完整性。

在这之前先创建一个utils文件夹用于存放一些工具函数，在util文件下创建utils.go

引入必要的包

```go
package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)
```

然后编写第一个需要用到的工具函数

```go
//ToHexInt将int64转换为字节串类型
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	Handle(err)
	return buff.Bytes()
}
```

还有一个用于处理错误的函数

```go
//处理错误
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
```

回到blockchain.go,开始哈希函数的构造，如下：

在头部引入`part1/src/utils`

```go
//哈希构造函数
func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}
```

information变量是将区块的各项属性串联之后的字节串。这里提醒一下bytes.Join可以将多个字节串连接，第二个参数是将字节串连接时的分隔符，这里设置为[]byte{}即为空，ToHexInt将int64转换为字节串类型。然后我们对information做哈希就可以得到区块的哈希值了。

## 区块创建与创始区块

既然我们可以获得区块的哈希值了，我们就能够创建区块了。

```go
//区块创建
func CreateBlock(prevhash, data []byte) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, data}
	block.SetHash()
	return &block
}
```

可以看到在创建一个区块时一定要引用前一个区块的哈希值，这里会有一个问题，那就是区块链中的第一个区块怎么创建？其实，在区块链中有一个创世区块，随着区块链的创建而添加，它指向的上一个区块的哈希值为空。

```go
//创世区块创建
func GenesisBlock() *Block {
	genesisWords := "创世区块"
	return CreateBlock([]byte{}, []byte(genesisWords))
}
```

可以看到我们在创始区块中存放了 *Hello, blockchain!* 这段信息。现在我们来构建函数，使得区块链可以根据其它信息创建区块进行储存。

```go
//添加区块
func (bc *BlockChain) AddBlock(data string) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
}
```

最后我们构建一个区块链初始化函数，使其返回一个包含创始区块的区块链。

```go
//创建区块链
func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}

```

## 运行区块链系统

现在我们已经拥有了所有创建区块链需要的函数了，来看看我们的区块链是怎么运作的。

```go
package main

import (
	"fmt"
	"part1/src/blockchain"
)

func main() {
	bc := blockchain.CreateBlockChain()
	bc.AddBlock("第一个区块")

	for _, block := range bc.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)

	}

}

```

你需要注意的是创始区块没有Previous Hash，同时后面的每一个区块都保留了前一个区块的哈希值。

## 总结

在本章中，我们构建了一个最简单的区块链模型。本章需要重点理解区块与区块链的关系，区块的哈希值的意义，以及创世区块的构建。在下一章中，我们将讲解PoW(Proof of Work)共识机制，并增加一些区块结构体的头部信息。