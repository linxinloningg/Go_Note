# 深入底层：Go语言从零构建区块链（六）附2：建立本地UTXO集加速可用交易信息查找

## 前言
本章内容较少，只是简单补充介绍一下区块链钱包客户端在实际运作时如何做到较快地维护与钱包地址有关的交易信息的。本章将会是本教程上半部分的最后一章，之后该教程将开启下半部分逐步建立一个P2P网络（已经建文件夹了QAQ），实现区块链系统的分布式部署。

## 本地UTXO集的作用
区块链作为一种链式结构，随着区块数量的不断增长，我们对于区块链中的某一账户的UTXO查找与验证所需要的时间与计算消耗都会线性增长。因此在实际应用时区块链系统的客户端一般会建立一个本地UTXO集，来加快对可用交易信息的查找与验证。UTXO集也就是区块链中没有被使用的Output，本教程在第三章有讲到。对于一个参与挖矿过程的区块链全节点而言，建立本地UTXO数据库可以实现对新区快是否合法的快速验证，因为在验证区块中的交易信息是否正确时不需要再遍历区块链验证是否存在双花；对于一个只是维护区块链但不参与挖矿的非全节点而言（如钱包客户端），维护一个与钱包地址相关地本地UTXO集可以更加便捷地管理钱包余额，以及生成交易信息（或者说这就是钱包客户端主要作用）。

## 区块结构体重构
在实现本地UTXO集前，我们先对区块结构体添加一个Height属性，也就是区块在区块链中的高度。该属性可以便于系统中的各设备节点更新维护区块链，也可以辅助我们检查本地UTXO集是否已经过时了。

打开blockchain文件夹下的block.go，修改结构体，以及相应的实例化函数。
```go
//区块的结构体
type Block struct {
	Timestamp    int64                      //时间戳
	Hash         []byte                     //本身的哈希值
	PrevHash     []byte                     //指向上一个区块的哈希
	Height       int64                      //区块在区块链中的高度
	Target       []byte                     //目标难度值
	Nonce        int64                      //POW
	Transactions []*transaction.Transaction //交易事务
	MTree        *merkletree.MerkleTree     //MT
}

//区块创建
func CreateBlock(prevhash []byte, height int64, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, height, []byte{}, 0, txs, merkletree.CrateMerkleTree(txs)}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

//创世区块创建
func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte{}, 0, []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}
```

打开blockchain.go，修改并添加一些与Height有关的函数。
```go
//获取当前区块链最后一个区块
func (chain *BlockChain) GetCurrentBlock() *Block {
	var block *Block
	err := chain.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get(chain.LastHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
	return block
}

//获取区块高度
func (bc *BlockChain) BackHeight() int64 {
	return bc.GetCurrentBlock().Height
}

```

打开mine.go，修改相应的函数。
```go
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()

	//验证交易信息池信息是否合法
	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}

	//挖矿完成，创建一个新区块
	candidateBlock := CreateBlock(bc.LastHash, bc.BackHeight()+1, transactionPool.PubTx) //PoW 已经在这里完成.

	//验证新区块是否合法，合法即添加到链上
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("区块拥有无效 nonce.")
		return
	}
}
```

Height属性我们定义为一int64类型，由于后文涉及字节串类型与int64类型的互转换，打开utils文件夹下的util.go增加如下函数。
```go
//ToInt64将节串类型转换为int64类型
func ToInt64(num []byte) int64 {
	var num64 int64
	buff := bytes.NewBuffer(num)
	err := binary.Read(buff, binary.BigEndian, &num64)
	Handle(err)
	return num64
}
```
ToInt64可以将字节串转换为int64类型，而util.go中存在另一此前就有的函数ToHexInt可以将int64编码为字节串。


## UTXO集构建
新建一个utxoset文件夹，然后在底下创建utxoset.go，在其中定义我们的UTXO结构体,以及相应后续会使用的一些函数。要唯一指明一个UTXO，应该包括该output的所有信息（value与pubkeyHash），该output在所在交易信息的ID值以及位于该交易信息中output的序号。

首先导入下面的包,并进行一些常量的声明。

```go
package utxoset

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"part8/src/transaction"
	"part8/src/utils"
	"runtime"
)

var (
	info         = "INFO:"
	infoname     = info + "NAME"
	infoheight   = info + "HIGT"
	utxokey      = "UTXO:"
	utxokeyorder = ":ORDER:"
)
```



```go
//UTXO集结构体
type UTXOSet struct {
	Name   []byte     //用于辨别UTXO的名称
	DB     *badger.DB //UTXO 数据库
	Height int64
}
```
回到bolck.go更新相关函数

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

接下来打开blockchain文件夹下的blockchain.go增加如下函数。BackUTXOs可以根据所给的钱包地址返回区块链中所有该地址的UTXO，用以辅助构建UTXO集。

```go
//根据所给的钱包地址返回区块链中所有该地址的UTXO
func (bc *BlockChain) BackUTXOs(address []byte) []transaction.UTXO {
	var UTXOs []transaction.UTXO
	unspentTxs := bc.FindUnspentTransactions(address)

Work:
	for _, tx := range unspentTxs {
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				UTXOs = append(UTXOs, transaction.UTXO{tx.ID, outIdx, out})
				continue Work // one transaction can only have one output referred to adderss
			}
		}
	}

	return UTXOs
}
```
如果是阅读了本教程前面章节的读者应该会发现BackUTXOs与之前写的FindUTXOs函数作用极其极其相似。事实上FindUTXOs在此前的作用只是通过遍历区块链查询某一钱包地址的余额，这对于之后通过BackUTXOs函数建立了本地UTXO集后便没有被使用的意义了（因为可以通过遍历UTXO集快速获得钱包余额），也就是说FindUTXOs函数被弃用了。

接下来我们正式建立UTXO集。UTXO集的建立其实就是重新建立一个数据库，该数据库可以及时更新其存储的UTXO。与本教程前面章节的选择一样，这里通过badger数据库进行实现。

可以看到，一个UTXO集应该包括一个数据库用于维护其中的UTXO，一个Name用于辨别这是什么的UTXO集，以及一个Height用于说明该UTXO集维护的UTXO的有效性，便于判断UTXO集是否需要更新。

接下来就需要实现UTXO集的创建与加载函数。如下所示。

```go
//utxoset.go
func GetUtxoSetFile(dir string) string {
	fileAddress := dir + "/" + "MANIFEST"
	return fileAddress
}

func ToUtxoKey(txID []byte, order int) []byte {
	utxoKey := bytes.Join([][]byte{[]byte(utxokey), txID, []byte(utxokeyorder), utils.ToHexInt(int64(order))}, []byte{})
	return utxoKey
}
```

```go
//UTXO集的创建
func CreateUTXOSet(name []byte, dir string, utxos []transaction.UTXO, height int64) *UTXOSet {
	if utils.FileExists(GetUtxoSetFile(dir)) {
		fmt.Println("UTXOSet has already existed, now rebuild it.")
		err := os.RemoveAll(dir)
		utils.Handle(err)
	}

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	utxoSet := UTXOSet{name, db, height}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(infoname), name)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(infoheight), utils.ToHexInt(height))
		if err != nil {
			return err
		}
		for _, utxo := range utxos {
			utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
			err = txn.Set(utxoKey, utxo.Serialize())
			return err
		}
		return nil
	})
	utils.Handle(err)
	return &utxoSet

}

//UTXO集的加载
func LoadUTXOSet(dir string) *UTXOSet {
	if !utils.FileExists(GetUtxoSetFile(dir)) {
		fmt.Println("No UTXOSet found, please create one first")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)
	var name []byte
	var height int64
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(infoname))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			name = val
			return nil
		})
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte(infoheight))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			height = utils.ToInt64(val)
			return nil
		})

		return err
	})
	utils.Handle(err)

	utxoSet := UTXOSet{name, db, height}
	return &utxoSet
}
```



如此便实现了UTXO集的建立。接下来还需要实现一些辅助函数便于对UTXO集的调用。

```go
//添加新的UTXO进数据库
func (us *UTXOSet) AddUTXO(txID []byte, outIdx int, output transaction.TxOutput) {
	utxo := transaction.UTXO{txID, outIdx, output}
	//us.AddUtxo(&utxo)
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
		err := txn.Set(utxoKey, utxo.Serialize())
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

//在数据库删除UTXO
func (us *UTXOSet) DelUTXO(txID []byte, order int) {
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(txID, order)
		err := txn.Delete(utxoKey)
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

func (us *UTXOSet) UpdateHeight(height int64) {
	us.Height = height
	err := us.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(infoheight), utils.ToHexInt(height))
		return err
	})
	utils.Handle(err)
}

func IsInfo(inkey []byte) bool {
	if bytes.HasPrefix(inkey, []byte(info)) {
		return true
	} else {
		return false
	}
}
```

其中AddUTXO函数可以添加新的UTXO进数据库，而DelUTXO可以删除UTXO。UpdateHeight函数可以更新UTXO集的Height信息，IsInfo函数主要是判断数据库中所存储的某一键值对信息是否是Height或Name一类的描述信息还是一般的UTXO信息。

## 钱包客户端
到上文为止UTXO集的基本框架已经搭建完毕，但是我们还没有使用UTXO集去做一些实际应用。现在以钱包客户端的UTXO集应用为例进行后续实例化的讲解。

首先打开constcoe文件夹下的constcoe.go文件加入一个新的全局变量。

```go
//constcoe.go
const (
	Difficulty          = 12
	InitCoin            = 1000
	TransactionPoolFile = "./tmp/transaction_pool.data"
	BCPath              = "./tmp/blocks"
	BCFile              = "./tmp/blocks/MANIFEST"
	ChecksumLength      = 4
	NetworkVersion      = byte(0x00)
	Wallets             = "./tmp/wallets/"
	WalletsRefList      = "./tmp/ref_list/"
	UTXOSet             = "./tmp/utxoset/" //This is new
)
```

然后打开wallet文件夹，创建一个属于wallet包的utxoset.go文件。引入以下包。

```go
package wallet

import (
	"bytes"
	"fmt"
	"part8/src/blockchain"
	"part8/src/constcode"
	"part8/src/transaction"
	"part8/src/utils"
	"part8/src/utxoset"

	"github.com/dgraph-io/badger"
)
```

通过调用utxoset包，重构wallet自己的UTXO集创建与加载函数。
```go
func (wt *Wallet) GetUtxoSetDir() string {
	strAddress := string(wt.Address())
	dirAddress := constcode.UTXOSet + strAddress
	return dirAddress
}

//UTXO集创建
func (wt *Wallet) CreateUTXOSet(chain *blockchain.BlockChain) *utxoset.UTXOSet {
	UTXOs := chain.BackUTXOs(wt.PublicKey)
	utxoSet := utxoset.CreateUTXOSet(wt.Address(), wt.GetUtxoSetDir(), UTXOs, chain.BackHeight())
	return utxoSet
}

//UTXO集加载
func (wt *Wallet) LoadUTXOSet() *utxoset.UTXOSet {
	utxoSet := utxoset.LoadUTXOSet(wt.GetUtxoSetDir())
	return utxoSet
}
```

通过建立一个特定钱包地址的UTXO集，实际上就是实现了一个钱包客户端的基本功能。现在，我们可以通过遍历UTXO集快速得到钱包余额而不需要再去遍历整个区块链了。以下为余额查询的函数实现。

```go
/获取余额
func (wt *Wallet) GetBalance() int {
	amount := 0
	us := wt.LoadUTXOSet()
	defer us.DB.Close()

	err := us.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if utxoset.IsInfo(item.Key()) {
				continue
			}
			err := item.Value(func(v []byte) error {
				tmpUTXO := transaction.DeserializeUTXO(v)
				amount += tmpUTXO.Value
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	utils.Handle(err)
	return amount
}
```
随着区块链的增长，我们的本地UTXO集也是需要更新的，因此这里构造一个简单的UTXO集更新函数。

```go
/UTXO集更新
/*
通过输入一个比当前UTXO集Height高一个单位的区块来实现UTXO集的更新
 */
func (w *Wallet) ScanBlock(block *blockchain.Block) {
	utxoSet := w.LoadUTXOSet()
	defer utxoSet.DB.Close()

	if block.Height > (utxoSet.Height + 1) {
		fmt.Println("UTXO Set is out of date!")
		return
	}

	for _, tx := range block.Transactions {
		for _, in := range tx.Inputs {
			if bytes.Equal(in.PubKey, w.PublicKey) {
				utxoSet.DelUTXO(in.TxID, in.OutIdx)
			}
		}

		for outIdx, out := range tx.Outputs {
			if bytes.Equal(out.HashPubKey, utils.PublicKeyHash(w.PublicKey)) {
				utxoSet.AddUTXO(tx.ID, outIdx, out)
			}
		}
	}
	utxoSet.UpdateHeight(block.Height)
}
```

可以看到，该函数是通过输入一个比当前UTXO集Height高一个单位的区块来实现UTXO集的更新的。

接下来，我们将修改cli命令行程序的balance实现方法，让它通过钱包的UTXO集实现。

在cli文件夹下的utxofunction.go文件。增加一个iniUtxoSet函数。

```go
//为所有已知的钱包创建各自的UTXO集
func (cli *CommandLine) iniUtxoSet() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	refList := wallet.LoadRefList()
	for addr, _ := range *refList {
		wlt := wallet.LoadWallet(addr)
		utxoSet := wlt.CreateUTXOSet(chain)
		utxoSet.DB.Close()
	}
	fmt.Println("Succeed in initializing UTXO sets.")
}

```
该函数可以为所有已知的钱包创建各自的UTXO集。

然后我们注释掉之前的balance函数，并重新构造。
```go
func (cli *CommandLine) balance(address string) {
	wlt := wallet.LoadWallet(address)
	balance := wlt.GetBalance()
	fmt.Printf("地址:%s, 余额:%d \n", address, balance)
}
```

接下来更新一下usage说明以及将新的iniutxoset注入为命令（这里大家就照着此前的注入方式注入一下就行了，这里就不再详述，当然也可以参考前文中所给的代码地址进行比照）。
```go
func (cli *CommandLine) printUsage() {
	fmt.Println("欢迎来到微型区块链系统，用法如下：")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("您只需要首先创建一个区块链并声明所有者。")
	fmt.Println("然后你就可以进行交易了。")
	fmt.Println("进行交易以扩展区块链。")
	fmt.Println("另外，收集交易后不要忘记挖矿功能，以打包交易信息。")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("createblockchain -refname NAME -address ADDRESS                   ----> 使用您输入的所有者（地址或引用名称）创建一个区块链。")
	fmt.Println("balance -refname NAME -address ADDRESS                            ----> 使用您输入的地址（或引用名称）返回钱包的余额。")
	fmt.Println("blockchaininfo                                      ----> 打印区块链中的所有区块")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> 产生一个交易信息并将该交易信息存储到交易信息池中")
	fmt.Println("sendbyrefname -from NAME1 -to NAME2 -amount AMOUNT  ----> 进行交易并使用 refname 将其放入候选块中。")
	fmt.Println("mine                                                ----> 模拟挖矿过程，将交易信息池中的交易打包成区块加入区块链中")

	fmt.Println("createwallet -refname REFNAME                       ----> 创建并保存钱包。 refname（别名） 是可选的。")
	fmt.Println("walletinfo -refname NAME -address Address           ----> 打印钱包信息。至少需要引用别名和地址之一。")
	fmt.Println("walletsupdate                                       ----> 注册并更新所有钱包（尤其是当您添加了现有的 .wlt 文件时）。")
	fmt.Println("walletslist                                         ----> 列出所有找到的钱包（确保您已先运行 walletsupdate）。")

	fmt.Println("initutxoset                                         ----> 初始化所有已知钱包的 UTXO 集。")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
}
```
这时我们还需要更改一下mine函数，让我们的每个wallet都在有新的区块加入时可以更新自己的本地UTXO集。

```go
func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("完成挖矿")

	//每个wallet都在有新的区块加入时可以更新自己的本地UTXO集
	newblock := chain.GetCurrentBlock()
	refList := wallet.LoadRefList()
	for k, _ := range *refList {
		wlt := wallet.LoadWallet(k)
		wlt.ScanBlock(newblock)
	}
	fmt.Println("完成更新 UTXO 集")
}
```

到此我们就实现了钱包客户端的UTXO集建立与应用了。

## 测试
修改我们的bat文件，主要是增加一个initutxoset命令。

```
rd /s /q tmp
md tmp\blocks
md tmp\wallets
md tmp\ref_list
md tmp\utxoset
main.exe createwallet 
main.exe walletslist
main.exe createwallet -refname LeoCao
main.exe walletinfo -refname LeoCao
main.exe createwallet -refname Krad
main.exe createwallet -refname Exia
main.exe createwallet 
main.exe walletslist
main.exe createblockchain -refname LeoCao
main.exe blockchaininfo
main.exe initutxoset
main.exe balance -refname LeoCao
main.exe sendbyrefname -from LeoCao -to Krad -amount 100
main.exe balance -refname Krad
main.exe mine
main.exe blockchaininfo
main.exe balance -refname LeoCao
main.exe balance -refname Krad
main.exe sendbyrefname -from LeoCao -to Exia -amount 100
main.exe sendbyrefname -from Krad -to Exia -amount 30
main.exe mine
main.exe blockchaininfo
main.exe balance -refname LeoCao
main.exe balance -refname Krad
main.exe balance -refname Exia
main.exe sendbyrefname -from Exia -to LeoCao -amount 90
main.exe sendbyrefname -from Exia -to Krad -amount 90
main.exe mine
main.exe blockchaininfo
main.exe balance -refname LeoCao
main.exe balance -refname Krad
main.exe balance -refname Exia
```

运行test.bat，如果运行结果符合预期，则各钱包的本地UTXO集成功建立并能够正常工作。

## 总结
本章只是以钱包客户端的本地UTXO集应用举例建立并说明了UTXO集，对于一个区块链的全节点而言，也应该建立UTXO集用于对可用UTXO的快速查找与验证。有兴趣的读者可以拓展本章建立的utxo包实现全节点的UTXO集应用，通过UTXO集实现mine.go中的VerifyTransactions函数。

由于在本教程的第二部分将涉及网络编程，整个cli程序将会重构，钱包客户端与区块链主程序将分离，因此对于全节点的UTXO集建立就留待届时。







