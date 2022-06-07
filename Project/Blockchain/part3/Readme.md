# 深入底层：Go语言从零构建区块链（三）：交易信息与UTXO模型

## 从一次转账说起

区块中的数据往往以交易信息（Transaction）的形式存储。交易信息顾名思义，最初指的就是bitcoin中的各个用户的转账信息。这里提醒一下，随着区块链的发展，在非金融领域，人们还是习惯于将区块中储存的一条一条的有用数据称为交易信息。

既然在比特币中交易信息就是转账信息，我们不妨思考一下如何将“A把五块钱转给B”这个转账信息表示出来。也许你的表示如下：Sender：A；Reciever：B；Amount：5。很好，这很符合我们的直觉。在日常生活中要确认上述转账信息是否有效，只需要通过银行这个可信第三方机构就可以实现，因为银行记录了A与B各自的资产信息。现在回到区块链中，区块链作为一个去中心的分布式系统，其目的就是去掉可信第三方，此时我们如何确认这样一个转账信息有效了？聪明的你可能想到了和中本聪一样的办法，那就是在交易信息中向前回溯，找到以A作Reciever的前置交易信息，加和它们的Amout是否大于5，如果大于5本次转账就是有效的。

到这里，有的小朋友就有话要说了，如果我再进行“一次A把五块钱转给B”的转账，这次转账肯定也会被认为是有效的，一直重复，都将是有效的。很好，这个问题的关键所在就是在进行了交易回溯后，那些支持本次交易的前置交易信息没有被标记，导致这些前置交易信息被无限次的用于支持其它交易的回溯。我们需要做的就是在本次交易信息中标记出那些用于支持本次交易的前置交易信息。

可以看到，我们在确认转账信息进行回溯时，我们其实根本不关心前置交易信息的Sender是谁，我们只关心它们的Reciever和Amount，这就是就是比特币中强调的UTXO（Unspent Transaction Outputs）模型的基本思路，现在让我们看看比特币中的UTXO模型究竟是如何构建来实现上述功能的。（关于区块链为何使用UTXO我举了上述这个例子来讲解，可能一些同学还是不能理解，无妨，我们直接阅读并理解后文的代码来直观感受UTXO的精妙之处，这个东西有时候就是有点无法言传，这才需要看代码，talk is cheap, show me the code!）

## 交易信息

我们在项目下新建文件夹叫transaction，然后在其中建立两个文件，分别为transaction.go和inoutput.go

在比特币中，交易信息被切分为两部分，分别为input与output。在inoutput中我们构造结构体，使得TxOutput代表交易信息的Output，TxInput代表交易信息的Input。

```go
package transaction

import "bytes"

type TxOutput struct {
	Value     int    //转出的资产值
	ToAddress []byte //资产的接收者的地址
}

type TxInput struct {
	TxID        []byte //指明支持本次交易的前置交易信息
	OutIdx      int    //具体指明是前置交易信息中的第几个Output
	FromAddress []byte //资产转出者的地址
}
```

TxOutput包含Value与ToAddress两个属性，前者是转出的资产值，后者是资产的接收者的地址。TxInput包含的TxID用于指明支持本次交易的前置交易信息，OutIdx是具体指明是前置交易信息中的第几个Output，FromAddress就是资产转出者的地址。

然后我们转到transaction.go，构造结构体。

```go
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
```

这里解释一下，Transaction由自身的ID值（其实就是哈希值），一组TxInput与一组TxOutput构成。TxInput用于标记支持我们本次转账的前置的交易信息的TxOutput，而TxOutput记录我们本次转账的amount和Reciever。你可能对于为什么要记录一组TxOutput有疑惑，这是因为在寻找到了足够多的未使用的TxOutput（后面全部简称UTXO）后，其资产总量可能大于我们本次交易的转账总量，我们可以将找零计入本次的TxOutput中，设置其流入方向就是本次交易的Sender（一定要好好理解！），这样就实现了找零。

这里我们首先实现计算每个transaction的哈希值的功能，如下。

```go
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
```

这里用到了gob，其功能主要是序列化结构体，与json有些像但是更方便。TxHash返回交易信息的哈希值，SetID设置每个交易信息的ID值，也就是哈希值。

## 创建交易信息

之前的章节我们实现了创始区块的构建，现在我们已经了解到区块存储的数据信息其实就是交易信息，那么在创建创始区块的时候我们是不是也应该创建创始的交易信息？答案是肯定的，世界上本没有比特币，随着创始区块的产生才有了币。代码如下。

```go
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
```

代码中使用了constcoe.InitCoin这个常量它代表了区块链在创建时的总的比特币数目，我们需要转到constcoe.go进行设置。

```go
package constcoe

const (
	Difficulty = 12
	InitCoin   = 1000// This line is new
)

```

可以看到，创始交易信息因为是凭空产生的，其Input指向一个为空的交易信息中的序号为-1的Output，充分体现了它的特殊之处。创始交易信息的Output不为空，它将指向一个地址（暂且先这样叫吧），该地址可以指向中本聪也可以指向创建区块链的你自己（有没有体会到一种当上帝的感觉qvq）。IsBase函数是用于检验一个交易信息是否为创始交易信息的。***注意代码中可以使用SetID函数设置创世交易信息的ID，但我并没有这样做，而是以"This is the Base Transaction!“为其ID值，这也是当上帝后的特权。

在修改block.go之前我们先转到inoutput.go写几个后文要用的函数。

```go
//验证FromAddress是否正确
func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

//验证ToAddress是否正确
func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
```

这两个函数还是比较好理解的，就是验证ToAddress和FromAddress是否正确。

## 代码重构

转到block.go，引入下面的包。

```go
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"part3/src/transaction"
	"part3/src/utils"
	"time"
)
```

修改Block结构体。

```go
//区块的结构体
type Block struct {
	Timestamp    int64                      //时间戳
	Hash         []byte                     //本身的哈希值
	PrevHash     []byte                     //指向上一个区块的哈希
	Target       []byte                     //目标难度值
	Nonce        int64                      //POW
	Transactions []*transaction.Transaction //交易事务
}
```

创建一个函数来协助处理区块中交易信息的序列化。

```go
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
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.SerializeTransaction()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}
```

修改CreateBlock与GenesisBlock两个函数

```go
//区块创建
func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

//创世区块创建
func GenesisBlock() *Block {

	//初始的所有比特币转给了神秘的人
	tx := transaction.BaseTx([]byte("创始人"))
	return CreateBlock([]byte{}, []*transaction.Transaction{tx})
}
```

创建创始区块的同时一并创建了创始交易信息，并且把初始的所有比特币转给了神秘的创始人

转到proofofwork.go，修改GetBase4Nonce函数

```go
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"part3/src/constcode"
	"part3/src/utils"
)

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
```

在blockchain.go引入下列库

```go
package blockchain

import (
	"encoding/hex"
	"fmt"
	"part3/src/transaction"
	"part3/src/utils"
)
```

修改AddBlock与CreateBlockChain两个函数

```go
//添加区块
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

//创建区块链
func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}
```

现在我们将会来一起构建CreateTransaction这样一个函数。该函数用于创建一个交易信息，其输入为资产转出者地址，资产接收者地址，以及转出资产总值。如前文所说，我们在创建一个交易信息，首先应该找到足够的可用的前置交易信息的Output来支撑本次交易，然后在构造交易信息时在TxInput中体现出来。

## 创建交易信息

我们将一步一步实现寻找前置可用交易信息Output这一功能。首先实现的函数是根据目标地址寻找可用交易信息的函数。

```go
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
```

这是一个返回区块链中一个地址的可用交易信息的函数。这里用到了go语言的label特性，需要你去先了解一下其使用方法。我们拆解这个函数来讲解。

```go
//从最后一个区块开始向前遍历区块链，然后遍历每一个区块中的交易信息
var unSpentTxs []transaction.Transaction
spentTxs := make(map[string][]int) // can't use type []byte as key value

//遍历交易信息的Output，如果该Output在spentTxs中就跳过，说明该Output已被消费。
for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}

//遍历交易信息的Output，如果该Output在spentTxs中就跳过，说明该Output已被消费。
if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
//否则确认ToAddress正确与否，正确就是我们要找的可用交易信息。
if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}

```

检查当前交易信息是否为Base Transaction（主要是它没有input），如果不是就检查当前交易信息的input中是否包含目标地址，有的话就将指向的Output信息加入到spentTxs中。

利用FindUnspentTransactions函数，我们可以找到一个地址的所有UTXO以及该地址对应的资产总和。

```go
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
```

当然，我们在实际应用中不需要每次都要找到所有UTXO，我们只需找到资产总量大于本次交易转账额的一部分UTXO就行。代码如下。

```go
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
```

构造CreateTransaction函数

```go
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
```

在真实区块链中，一个节点会维护一个候选区块，候选区块会维持一个交易信息池（Transaction Pool），然后在挖矿时将交易池中的交易信息打包进行挖矿（PoW过程）。

我们现在不希望再使用AddBlock直接添加区块进入到区块链中，而是预留一个函数模拟一下从交易信息池中获取交易信息打包并挖矿这个过程。

```go
/*
在真实区块链中，一个节点会维护一个候选区块，候选区块会维持一个交易信息池（Transaction Pool），
然后在挖矿时将交易池中的交易信息打包进行挖矿（PoW过程）。
*/
func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
```

## 调试

* 创建创始区块并查看创始人余额

  ```go
  package main
  
  import (
  	"fmt"
  	"part3/src/blockchain"
  )
  
  func main() {
  
  	chain := blockchain.CreateBlockChain()
  	property, _ := chain.FindUTXOs([]byte("创始人"))
  	fmt.Println("创始人的余额: ", property)
  
  }
  
  ```

* 创建交易

  ```go
  package main
  
  import (
  	"fmt"
  	"part3/src/blockchain"
  	"part3/src/transaction"
  )
  
  func main() {
  	//交易信息池
  	txPool := make([]*transaction.Transaction, 0)
  
  	chain := blockchain.CreateBlockChain()
  
  	tempTx, status := chain.CreateTransaction([]byte("创始人"), []byte("linxinloningg"), 100)
  	if status {
  		txPool = append(txPool, tempTx)
  	}
  
  	//挖矿打包新区块
  	chain.Mine(txPool)
  
  	property, _ := chain.FindUTXOs([]byte("linxinloningg"))
  	fmt.Println("Balance of linxinloningg: ", property)
  
  }
  
  ```

* 出现“双花”bug

  ```go
  package main
  
  import (
  	"fmt"
  	"part3/src/blockchain"
  	"part3/src/transaction"
  )
  
  func main() {
  	txPool := make([]*transaction.Transaction, 0)
  	var tempTx *transaction.Transaction
  	var ok bool
  	var property int
  	chain := blockchain.CreateBlockChain()
  	property, _ = chain.FindUTXOs([]byte("创始人"))
  	fmt.Println("Balance of 创始人: ", property)
  
  	tempTx, ok = chain.CreateTransaction([]byte("创始人"), []byte("第一个人"), 100)
  	if ok {
  		txPool = append(txPool, tempTx)
  	}
  	chain.Mine(txPool)
  	txPool = make([]*transaction.Transaction, 0)
  	property, _ = chain.FindUTXOs([]byte("创始人"))
  	fmt.Println("Balance of 创始人: ", property)
  
  	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("第二个人"), 200) // this transaction is invalid
  	if ok {
  		txPool = append(txPool, tempTx)
  	}
  
  	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("第二个人"), 50)
  	if ok {
  		txPool = append(txPool, tempTx)
  	}
  
  	tempTx, ok = chain.CreateTransaction([]byte("创始人"), []byte("第二个人"), 100)
  	if ok {
  		txPool = append(txPool, tempTx)
  	}
  	chain.Mine(txPool)
  	txPool = make([]*transaction.Transaction, 0)
  	property, _ = chain.FindUTXOs([]byte("创始人"))
  	fmt.Println("Balance of 创始人: ", property)
  	property, _ = chain.FindUTXOs([]byte("第一个人"))
  	fmt.Println("Balance of 第一个人: ", property)
  	property, _ = chain.FindUTXOs([]byte("第二个人"))
  	fmt.Println("Balance of 第二个人: ", property)
  
  	for _, block := range chain.Blocks {
  		fmt.Printf("Timestamp: %d\n", block.Timestamp)
  		fmt.Printf("hash: %x\n", block.Hash)
  		fmt.Printf("Previous hash: %x\n", block.PrevHash)
  		fmt.Printf("nonce: %d\n", block.Nonce)
  		fmt.Println("Proof of Work validation:", block.ValidatePoW())
  		fmt.Println("")
  	}
  
  	//I want to show the bug at this version.
  
  	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("第二个人"), 30)
  	if ok {
  		txPool = append(txPool, tempTx)
  	}
  
  	tempTx, ok = chain.CreateTransaction([]byte("第一个人"), []byte("创始人"), 30)
  	if ok {
  		txPool = append(txPool, tempTx)
  	}
  
  	chain.Mine(txPool)
  	txPool = make([]*transaction.Transaction, 0)
  
  	for _, block := range chain.Blocks {
  		fmt.Printf("Timestamp: %d\n", block.Timestamp)
  		fmt.Printf("hash: %x\n", block.Hash)
  		fmt.Printf("Previous hash: %x\n", block.PrevHash)
  		fmt.Printf("nonce: %d\n", block.Nonce)
  		fmt.Println("Proof of Work validation:", block.ValidatePoW())
  	}
  
  	property, _ = chain.FindUTXOs([]byte("创始人"))
  	fmt.Println("Balance of 创始人: ", property)
  	property, _ = chain.FindUTXOs([]byte("第一个人"))
  	fmt.Println("Balance of 第一个人: ", property)
  	property, _ = chain.FindUTXOs([]byte("第二个人"))
  	fmt.Println("Balance of 第二个人: ", property)
  }
  
  ```

  



