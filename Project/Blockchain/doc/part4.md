# 深入底层：Go语言从零构建区块链（四）：区块链的存储、读取与管理

## 前言

在前面的章节中我们了解了交易信息与UTXO模型，这样就掌握了区块链系统的基本数据结构。你可能已经发现了我们在前几章对区块链系统进行调试时每次都需要重新创建区块链，区块链并没有得到保存，这与实际的区块链系统不符。同时伴随着我们区块链的组件越来越多，我们需要一个统一的功能管理模块来操作区块链，而不是手动地去调用一个又一个的函数。介于此，本章我们将会实现区块链以及交易信息的存储、读取功能，并设计一个命令行模块来管理区块链系统。

## Badger键值对数据库

对于区块链的储存，一种自然的想法就是把一个一个的区块序列化，然后使用每个区块的哈希值作为文件名称在磁盘上进行保存。但是考虑到我们可能会频繁的对已存储的区块进行查询，本教程将使用数据库的方式存储并管理区块。我在参考了一些国外go语言实现区块链系统相关的资料后，决定使用Badger数据库。

Badger是dgraph.io开发的一款基于 Log Structured Merge (LSM) Tree 的 key-value 本地数据库，其官网在这里https://github.com/dgraph-io/badger。使用Badger有几个好处，首先它是由go语言编写地，然后它功能单一专注于键值对的存储，最后是安装简单不占地。

## 存储地址的设定

我们在项目目录下创建一个tmp文件夹，专门用来存储区块链系统中需要保存的文件。在tmp文件夹下再创建blocks文件夹用以存储区块链中的区块。

打开constcoe.go设定一些全局变量

```go
package constcode

const (
	Difficulty = 12
	InitCoin   = 1000
	TransactionPoolFile = "./tmp/transaction_pool.data"
	BCPath              = "./tmp/blocks"
	BCFile              = "./tmp/blocks/MANIFEST"
)
```

中TransactionPoolFile相当于一个缓冲池，存放一个节点收集到的交易信息，在后文会详细讲到。BCPath与BCFile都是指向的我们即将创建的区块链数据库的相关地址。

在utils.go中增加一个检查文件地址下文件是否存在的函数。

```go
package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

//检查文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}
```

## 区块改动

虽然我们要实现的是区块链的存储，但是对于区块我们并没有太大的改动，当然qvq，还是有一点咯。打开block.go引入如下的库。

```go
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"part4/src/transaction"
	"part4/src/utils"
	"time"
)
```

首先是创始区块的构建函数，我们希望他现在可以更加个性化一点，让创始交易信息指向我们提供的地址。

```go
//创世区块创建
func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte{}, []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}
```

添加两个新的函数。

```go
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
```

这两个函数一个是用于序列化区块生成字节串的，一个是反序列化区块的。Badger的键值对只支持字节串存储形式，所以我们需要这两个函数。

## 区块链结构体重构

接下来是对区块链结构体的重构。我们打开blockchain.go，引入下面的库。

```go
package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"part4/src/constcode"
	"part4/src/transaction"
	"part4/src/utils"
	"runtime"
)
```

此前我们的区块链结构体BlockChain的属性之一为区块组成的切片，而现在我们的区块存储在数据库中，所以结构体BlockChain的属性之一应该指向存储区块的数据库。重构我们的区块链结构体如下。

```go
//链结构体
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}
```

LastHash属性是指当前区块链最后一个区块的哈希值，它并不是必须的，但可以避免我们在后面的函数编写中每次都去数据库中查找LastHash。

现在可以看到我们的项目又开始飘红了，没有关系，我们将一步一步更改函数。此前我们有一个CreateBlockChain函数，该函数可以创建一个区块链并返回该区块链的指针。现在我们既然要实现区块链的创建、存储与读取功能，就需要摒弃这个函数。我们注释掉该函数，然后创建两个新的函数，分别为InitBlockChain与ContinueBlockChain。InitBlockChain可以初始化我们的区块链并创建一个数据库保存，而ContinueBlockChain可以读取已有的数据库并加载区块链。

InitBlockChain函数:

```go
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
```

可以看到，InitBlockChain函数会先检查是否有存储区块链的数据库存在，如果存在将会给出警告并退出（注意这里我们使用的是runtime.Goexit()进行退出,相对来说更加安全,当然你也可以使用一般的退出方式）。opts就是启动Badger的配置，我们这里全部使用默认配置，并将地址指定为constcoe.BCPath即可。opts.Logger =nil可以使数据库的操作信息不输出到标准输出中，当然你也可以不进行该设置，方便在调用数据库时进行debug。badger.Open(opts)就是按照我们的配置启动一个数据库（如果没有现成的数据库就会初始化一个），将会返回该数据库的指针。db.Update()是Badger中对数据库进行更新操作的函数，介于这种函数作为参数的写法可能刚接触go语言的同学不是很明白，我这里把db.Update()的源码放出来。

```go
// Update executes a function, creating and managing a read-write transaction
// for the user. Error returned by the function is relayed by the Update method.
// Update cannot be used with managed transactions.
func (db *DB) Update(fn func(txn *Txn) error) error {
	if db.opt.managedTxns {
		panic("Update can only be used with managedDB=false.")
	}
	txn := db.NewTransaction(true)
	defer txn.Discard()

	if err := fn(txn); err != nil {
		return err
	}

	return txn.Commit()
}

```

可以看到，在db.Update()中，以函数作为参数最主要的目的就是可以先构造一个事务（这里的Txn全称也就是Transaction，不同于我们区块链中的交易信息，在这里指数据库中的事务）给内部函数调用，然后再在外部实现事务的提交与结束。

可以看到在db.Update()中，我们利用提供的事务txn，添加了创始区块，其具体方式是将哈希值作为Key值，区块序列化后的数据作为Value进行存储，这也是Badger惟一支持的存储方式。我们将当前区块链最后一个区块的哈希值存储在"lh"这个Key值中，同时我们把创始区块的PrevHash存在"ogprevhash"这个Key值中。InitBlockChain函数最后返回以我们打开的数据库构造的BlockChain。这样我们就能够完成区块链的创建（包括数据库的创建）。

我们现在期望能够通过已有的数据库读取并加载我们的区块链，这就需要构建ContinueBlockChain函数。

```go
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
```

与InitBlockChain类似，我们首先需要设定opts，然后打开存储区块链的数据库。注意我们这里不再使用db.Update()函数对数据库进行访问，而是使用db.View()函数来调取视图，读取当前区块链的最后一个区块的哈希值。

现在我们如果要将一个区块加入到区块链中，就需要通过数据库来完成。修改AddBlock函数。

```go
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
```

可以看到，我们的AddBlock函数会先检查区块链中的LastHash与即将加入的区块的PrevHash是否一致，如果一致才会将其加入到区块链（数据库）中，并更新数据库中的"lh"。

## 区块链的遍历

尽管拥有ContinueBlockChain函数后我们实现了对区块链的加载，但是我们发现对于区块链的遍历不像之前那么方便了，FindUnspentTransactions函数一时不知如何修改。这里我们创建一个基于区块的迭代器来实现区块链的遍历。在blockchain_iterator.go中创建结构体。

引入下列包

```go
package blockchain

import (
	"github.com/dgraph-io/badger"
	"part4/src/utils"
)
```

```go
//基于区块的迭代器
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}
```

创建迭代器的初始化函数。

```go
//创建迭代器的初始化函数
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}
	return &iterator
}
```

创建迭代器的迭代函数，让每次迭代返回一个block，然后迭代器指向前一个区块的哈希值。

```go
//迭代函数:让每次迭代返回一个block，然后迭代器指向前一个区块的哈希值
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}
```

创建一个辅助函数来帮助判断迭代器是否终止

```go
//判断迭代器是否终止
/*
通过比较迭代器的CurrentHash与数据库存储的OgPrevHash是否相等就能够判断迭代器是否已经迭代到创始区块
*/
func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})

		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	return ogprevhash
}
```

通过比较迭代器的CurrentHash与数据库存储的OgPrevHash是否相等就能够判断迭代器是否已经迭代到创始区块。

现在我们可以在blockchain.go中修改FindUnspentTransactions函数了。

```go
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

```

可以看到我们通过迭代器完成了对区块链中区块由后到前的遍历。

## 交易信息池

我们知道区块链中每个区块可以存放复数个交易信息，而不是每生成一个交易信息就创建一个区块加入区块链。对于一个区块链节点而言，它会将自己生成的交易信息与从其它节点收集的交易信息存储于一个交易信息池中（Transaction Pool），当存储的交易信息数量达到一阈值或者等待时间超过一阈值就会将交易池中的交易信息打包为候选区块参与PoW共识机制，这一过程称为挖矿（mine），是节点争取将自己的候选区块加入到区块链中的过程。

上一章中我们已经预留了Mine函数来表征上述过程，现在我们可以更清晰地实现这个过程。

我们在blockchain文件夹下创建transactionpool.go文件，并且引入如下的库。

```go
package blockchain

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"part4/src/constcode"
	"part4/src/transaction"
	"part4/src/utils"
)
```

我们创建一个交易信息池的结构体。

```go
//交易信息池的结构体
type TransactionPool struct {
	PubTx []*transaction.Transaction //PubTx用于储存节点收集到的交易信息
}
```

PubTx用于储存节点收集到的交易信息。

```go
//添加新交易信息
func (tp *TransactionPool) AddTransaction(tx *transaction.Transaction) {
	tp.PubTx = append(tp.PubTx, tx)
}
```

以上函数实现新交易信息的添加。

我们希望存储我们新收集到的交易信息，就需要能够保存我们的交易信息池这一结构体。

```go
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
```

可以看到我们每次都将交易信息池保存到constcoe.TransactionPoolFile这个地址中。0644是指八进制的644（110，100，100），指明了不同用户对文件读写执行的权限。

在能够保存交易信息池后，我们也需要能够加载我们的交易信息池。

```go
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
```

这样使用CreateTransactionPool函数，我们可以创建或者加载一个交易信息池了。

```go
//创建交易信息池
func CreateTransactionPool() *TransactionPool {
	transactionPool := TransactionPool{}
	err := transactionPool.LoadFile()
	utils.Handle(err)
	return &transactionPool
}
```

最后考虑当节点在mine后需要清空交易信息池，我们还需要以下函数。

```go
//清空交易信息池
func RemoveTransactionPoolFile() error {
	err := os.Remove(constcode.TransactionPoolFile)
	return err
}
```

通过删除constcoe.TransactionPoolFile，即可实现交易信息池的清空。

## 挖矿Mine

前文我们已经说到将交易信息池打包为候选区块参与PoW共识机制的过程称为挖矿，现阶段我们无法模拟与其它网络节点挖矿时的竞争过程，所以我们假设现在的单个节点每次挖矿都能胜出并将自己的候选区块加入到区块链中。

我们将原blockchain.go下的Mine函数删除，在blockchain文件夹下创建mine.go文件，并创建RunMine函数，如下。

```go
package blockchain

import (
	"fmt"
	"part4/src/utils"
)

/*
在真实区块链中，一个节点会维护一个候选区块，候选区块会维持一个交易信息池（Transaction Pool），
然后在挖矿时将交易池中的交易信息打包进行挖矿（PoW过程）。
*/
//挖矿
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	//在不久的将来，我们必须先在这里验证交易。
	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx) //打包交易信息，挖矿寻找nonce.
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("区块有无效的 nonce.")
		return
	}
}
```

可以看到我们现阶段还是省略了挖矿前对候选区块中交易信息的验证过程，这留待以后的章节进行实现。

## 命令行Command Line

为了使我们的区块链系统运作的更像一个系统而非一个脚本，我们期望实现一个命令行（Command Line）模块，通过该模块来管理我们的区块链系统。

在我们的项目下创建文件夹cli，然后在cli下创建cli.go文件。

为我们的cli.go导入以下库。

```go
package cli

import (
	"bytes"
	"fmt"
	"part4/src/blockchain"
	"strconv"
)
```

我们创建一个空结构体来表征我们的命令行。

```go
type CommandLine struct{}
```

当我们打开一个命令行程序不知道做什么时，命令行程序首先应打印所有的命令及其用法。我们创建printUsage函数。

```go
func (cli *CommandLine) printUsage() {
	fmt.Println("欢迎来到微型区块链系统，用法如下：")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("您只需要首先创建一个区块链并声明所有者。")
	fmt.Println("然后你就可以进行交易了。")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("createblockchain -address ADDRESS                   ----> 读取一个地址，然后以该地址创建创始交易信息与创始区块完成区块链的初始化")
	fmt.Println("balance -address ADDRESS                            ----> 读取一个地址，然后在区块链中找到该地址的UTXO并统计出其余额")
	fmt.Println("blockchaininfo                                      ----> 打印区块链中的所有区块")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> 产生一个交易信息并将该交易信息存储到交易信息池中")
	fmt.Println("mine                                                ----> 模拟挖矿过程，将交易信息池中的交易打包成区块加入区块链中")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
}
```

由上可以看到，我们现在的命令行将管理实现五种功能。createblockchain命令会读取一个地址，然后以该地址创建创始交易信息与创始区块完成区块链的初始化。balance命令读取一个地址，然后在区块链中找到该地址的UTXO并统计出其余额。blockchaininfo将会打印区块链中的所有区块。send命令可以产生一个交易信息并将该交易信息存储到交易信息池中。mine命令模拟挖矿过程，将交易信息池中的交易打包成区块加入区块链中。

接下来，我们来一步一步的实现这五个功能。

首先是createblockchain，调用在blockchain.go中编写的InitBlockChain函数即可实现。

```go
func (cli *CommandLine) createblockchain(address string) {
	newChain := blockchain.InitBlockChain([]byte(address))
	/*
	注意在使用完数据库后，需要使用newChain.Database.Close()函数关闭数据库
	 */
	newChain.Database.Close()
	fmt.Println("完成创建区块链，所有者是：", address)
}
```

注意我们在使用完数据库后，需要使用newChain.Database.Close()函数关闭数据库。

```go
func (cli *CommandLine) balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	balance, _ := chain.FindUTXOs([]byte(address))
	fmt.Printf("地址:%s, 余额:%d \n", address, balance)
}
```

在实现balance功能时，我们先使用ContinueBlockChain函数接入数据库，然后使用FindUTXO函数统计余额并打印。注意我们这里使用了go语言的defer关键字，其后的代码将会在函数运行结束前最后执行，也就是我们最后将关闭数据库。

```go
func (cli *CommandLine) blockchaininfo() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	iterator := chain.Iterator()
	ogprevhash := chain.BackOgPrevHash()
	for {
		block := iterator.Next()
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("Transactions:%v\n", block.Transactions)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		if bytes.Equal(block.PrevHash, ogprevhash) {
			break
		}
	}
}
```

blockchaininfo命令需要使用我们之前设计的迭代器遍历区块链。

```go
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	tx, ok := chain.CreateTransaction([]byte(from), []byte(to), amount)
	if !ok {
		fmt.Println("创建交易失败")
		return
	}
	tp := blockchain.CreateTransactionPool()
	tp.AddTransaction(tx)
	tp.SaveFile()
	fmt.Println("成功!")
}
```

send命令将会调用CreateTransaction函数，并将创建的交易信息保存到交易信息池中。

```go
func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("完成挖矿")
}
```

对于mine命令，我们只需要使用此前编写的RunMine函数即可实现。

现在，我们只需要使用go语言自带的flag库将各个命令注册就行了。

```go
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	balanceAddress := balanceCmd.String("address", "", "Who need to get balance amount")
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "blockchaininfo":
		err := getBlockChainInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		utils.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createblockchain(*createBlockChainOwner)
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			balanceCmd.Usage()
			runtime.Goexit()
		}
		cli.balance(*balanceAddress)
	}

	if sendCmd.Parsed() {
		if *sendFromAddress == "" || *sendToAddress == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	}

	if getBlockChainInfoCmd.Parsed() {
		cli.blockchaininfo()
	}

	if mineCmd.Parsed() {
		cli.mine()
	}
}
```

这段代码将各项命令进行注册，不了解的同学可以去看一下go语言的flag库的文档，这里不再赘述。

## 系统调试

现在我们的区块链系统通过cli模块来进行管理，main函数就只需要提供一个入口就行了。

```go
package main

import (
	"os"
	"part4/src/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}

```

