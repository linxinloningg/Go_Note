# 深入底层：Go语言从零构建区块链（五）：密钥对、钱包与地址

## 前言

此前我们已经理解了交易信息与UTXO模型，本章开始我们将进入区块链系统的另一大板块即密码学相关的部分。这一部分内容较多，将分为两个章节进行讲解。本章主要涉及对非对称密钥的简单介绍（特别是椭圆曲线加密算法），讲解钱包是什么以及区块链中的地址究竟是什么。

## 非对称密钥与椭圆曲线加密算法

非对称密钥算法是区块链系统中的核心技术之一。使用非对称密钥算法，我们能够获得一密钥对，即公钥与私钥，公钥是可以公开给所有人的并作为用户的唯一标识符，而私钥则只能由用户自己保管。使用公钥加密的信息需要使用私钥才能够解密，而使用私钥签名的信息可以通过公钥进行验证。常见的非对称密钥算法包括RSA,Elgamal,ECC等。由于本教程以比特币系统为蓝本，故后续使用的非对称密钥算法均为ECC（椭圆曲线加密算法）。

介于椭圆曲线加密算法的原理网络上已经有很多博文与视频进行讲解，这里本人就不再班门弄斧了，读者需要自己搜索资料学习一下椭圆曲线加密算法的原理，这里推荐一篇讲解的比较透彻的博文： https://www.jianshu.com/p/e41bc1eb1d81。

为方便后续讲解，这里总结一下椭圆曲线加密算法的特点：

- 私钥是一个倍数，公钥是一个点（x,y）。
- 知道基点G（Base Point）和倍数x（私钥）可以在椭圆曲线上计算xG（公钥），而知道G与xG则几乎不可能推测x。
- 使用私钥对信息签名得到的是两个大数，一个是随机数r，一个是计算得到的s。

## 钱包与地址

区块链中的钱包与现实生活中的钱包有很大不同，这里我说说本人自己的理解。可以认为区块链中的一个钱包就对应了一个区块链用户，用户的密钥对都由钱包保存。现实中的钱包保存的资产就是真实的现金，在需要做交易时直接划拨现金即可；而在区块链中的钱包保存的资产不是现金（真实的货币），而是密钥对所指向的区块链中的相应UTXO总额，划拨资产要先通过密钥对证明这些UTXO的所属权。换句话说，区块链中的钱包根本不存储资产（UTXO自始至终都保存于区块链中），它只是使用密钥对帮助管理用户的个人资产。

在前面的章节中我们经常提到一个概念，那就是地址。区块链中的地址作用只有一个，那就是唯一指向一个用户，使得UTXO可以通过指向该地址来流向该用户。我们现阶段的send命令是用的昵称来作为地址，一个昵称指向一个特定的用户，但是这种地址表示方式终究很简陋，首先昵称太容易重复，其次没有一种手段来证明昵称的所属权（如何证明你的昵称是你的）。结合非对称密钥，我们知道公钥是可以公开给任何人的，同时也能够作为身份标识，用户通过掌握私钥也能够证明公钥的所属权，那为什么不使用公钥来作为地址了。

在最初的比特币中，的确是使用公钥来作为地址的，也即构造的交易信息是从公钥指向公钥。而在后续的版本更迭中，逐渐不再直接使用公钥作为地址，而是使用公钥进行一些列哈希操作得到的值作为地址，这是因为公钥哈希值能够在指向原用户个体的同时提升匿名性。

公钥哈希就是将公钥连续做了两次哈希操作得到（一次sha256一次ripemd160）。在公钥哈希的基础上还生成了钱包地址，钱包地址其实就是公钥哈希增加一个版本号位与四个检查位生成，最后转为比特币专门的Base58编码输出。

总结一下，公钥可以推得公钥哈希，但是公钥哈希不能推回公钥，同时公钥哈希和钱包地址可以互相转化。

## 实现准备

前文我们已经粗略地讲解了钱包与地址地作用，可能读者还是有很多不明白的问题，无妨，我们可以通过后面的代码辅助加深理解。

我们首先设置一下新的全局常量。打开constcoe.go，并增加如下常量。

```go
package constcode

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
)

```

在tmp文件夹下新创建wallets与ref_list两个文件夹。

在项目中创建wallet包，并创建两个go文件，即wallet.go与walletmanager.go。

接下来我们开始实现区块链系统的钱包功能。

## 钱包实现

我们首先打开wallet.go，引入以下包。

```go
package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"io/ioutil"
	"part5/src/constcode"
	"part5/src/utils"

	"errors"
)
```

创建椭圆曲线密钥对的生成函数。

```go
//椭圆曲线生成密钥对
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}
```

通过NewKeyPair函数，我们能够得到一对密钥对。可以看到，公钥是一个点，需要将横纵坐标拼接起来保存。

前文说过，钱包的主要作用就是保存一密钥对。我们构造wallet结构体。

```go
//钱包结构体
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}
```

接下来创建钱包生成函数。

```go
//钱包生成
func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}
```

接下来我们打开util.go，增加一些用于构造公钥哈希和钱包地址的函数。

引入下面的库包。

```go
//util.go
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"goblockchain/constcoe"
	"log"
	"os"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

```

增加将公钥转为公钥哈希的函数。

```go
//将公钥转为公钥哈希
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}
```

增加检查位生成函数。

```go
//检查位生成
func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:constcode.ChecksumLength]
}
```

增加Base256转Base58函数，以及其反函数。

```go
//Base256转Base58
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

//Base58转Base256
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}
```

现在就可以编写公钥哈希生成钱包地址的函数了。

```go
//公钥哈希生成钱包地址
func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{constcode.NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}
```

然后再增加一个钱包地址转公钥哈希的函数。

```go
//钱包地址转公钥哈希
func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constcode.ChecksumLength]
	return pubKeyHash
}
```

回到wallet.go 文件编写

钱包地址生成函数

```go
func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

```

我们希望在创建了一个钱包后可以保存这个钱包。

```go
//保存钱包
func (w *Wallet) Save() {
	filename := constcode.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}
```

可以看到我们的钱包将保存在.wlt文件中，其文件名使用的是钱包地址。这里需要注意使用gob时要先注册elliptic.P256()声明elliptic.Curve接口，否则会报错。

同时我们期望能够加载已经保存的钱包（也即密钥对）。

```go
//加载钱包
func LoadWallet(address string) *Wallet {
	filename := constcode.Wallets + address + ".wlt"
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("没有这个地址的钱包"))
	}
	var w Wallet
	gob.Register(elliptic.P256())
	fileContent, err := ioutil.ReadFile(filename)
	utils.Handle(err)
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}
```

到此wallet.go中的代码也就完成了。虽然一个钱包对应一个用户，但一个用户可以拥有多个钱包并将它们保存在一个机器上，我们需要建立一个钱包管理模块来管理一台机器上保存的所有钱包。walletmanager.go将实现一简易的管理模块。

打开walletmanager.go，引入以下包。

```go
package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"part5/src/constcode"
	"part5/src/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)
```

我们构建一字典型RefList来记录机器上记录的钱包。其中key值为钱包地址，value为钱包的别名。

```go
//RefList来记录机器上记录的钱包。其中key值为钱包地址，value为钱包的别名
type RefList map[string]string
```

RefList应该能够被保存，故实现Save函数。

```go
//保存RefList
func (r *RefList) Save() {
	filename := constcode.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}
```

RefList应实现一更新函数，用于扫描机器上保存的所有钱包文件（特别是检查是否存在从其他机器上拷贝的钱包）。

```go
//RefList应实现一更新函数，用于扫描机器上保存的所有钱包文件（特别是检查是否存在从其他机器上拷贝的钱包）
func (r *RefList) Update() {
	err := filepath.Walk(constcode.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			_, ok := (*r)[fileName[:len(fileName)-4]]
			if !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}
```

同时我们期望能够加载已保存的RefList。

```go
//加载已保存的RefList
func LoadRefList() *RefList {
	filename := constcode.WalletsRefList + "ref_list.data"
	var reflist RefList
	if utils.FileExists(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&reflist)
		utils.Handle(err)
	} else {
		/*
			如果没有可以可以加载的RefList文件，LoadRefList会自动创建一个新的RefList并调用Update函数扫描本机的所有.wlt文件
		*/
		reflist = make(RefList)
		reflist.Update()
	}
	return &reflist
}
```

可以看到如果没有可以可以加载的RefList文件，LoadRefList会自动创建一个新的RefList并调用Update函数扫描本机的所有.wlt文件。

为了方便我们区块链系统的演示，我们希望通过用别名的方式指向本地钱包（***注意：实际的区块链系统中钱包是没有别名的，这里完全是方便我们后续的演示。）

构建别名绑定函数。

```go
/*
用别名的方式指向本地钱包（***注意：实际的区块链系统中钱包是没有别名的，这里完全是方便我们后续的演示。）
 */
func (r *RefList) BindRef(address, refname string) {
	(*r)[address] = refname
}
```

然后构建通过别名调取钱包地址的函数。

```go
//构建通过别名调取钱包地址
func (r *RefList) FindRef(refname string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refname {
			temp = key
			break
		}
	}
	if temp == "" {
		err := errors.New("the refname is not found")
		return temp, err
	}
	return temp, nil
}
```

到此，我们的钱包就已完全实现。

## 命令行程序变更

在构建了钱包模块后，我们期望不再使用昵称来指代用户进行交易，而是使用钱包地址（在实际的区块链中应该是使用公钥哈希值，但是介于二者可以互相转化，本章先使用钱包地址，在下一章讲解签名与验证时再进行代码重构）。为此，我们的命令行程序应该添加一些新功能。

在cli文件夹下新建walletfunction.go

引入以下包。

```go
package cli

import (
	"fmt"
	"part5/src/utils"
	"part5/src/wallet"
)
```

在cli.go更新我们的printUsage函数。

```go
func (cli *CommandLine) printUsage() {
	fmt.Println("欢迎来到微型区块链系统，用法如下：")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("您只需要首先创建一个区块链并声明所有者。")
	fmt.Println("然后你就可以进行交易了。")
	fmt.Println("进行交易以扩展区块链。")
	fmt.Println("另外，收集交易后不要忘记挖矿功能，以打包交易信息。")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("createblockchain -address ADDRESS                   ----> 读取一个地址，然后以该地址创建创始交易信息与创始区块完成区块链的初始化")
	fmt.Println("balance -address ADDRESS                            ----> 读取一个地址，然后在区块链中找到该地址的UTXO并统计出其余额")
	fmt.Println("blockchaininfo                                      ----> 打印区块链中的所有区块")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> 产生一个交易信息并将该交易信息存储到交易信息池中")
	fmt.Println("mine                                                ----> 模拟挖矿过程，将交易信息池中的交易打包成区块加入区块链中")

	fmt.Println("createwallet -refname REFNAME                       ----> 创建并保存钱包。 refname（别名） 是可选的。")
	fmt.Println("walletinfo -refname NAME -address Address           ----> 打印钱包信息。至少需要引用别名和地址之一。")
	fmt.Println("walletsupdate                                       ----> 注册并更新所有钱包（尤其是当您添加了现有的 .wlt 文件时）。")
	fmt.Println("walletslist                                         ----> 列出所有找到的钱包（确保您已先运行 walletsupdate）。")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
}
```

可以看到，我们增加了createwallet命令用于创建新的钱包，增加了walletinfo命令来打印指定钱包的基本信息，增加walletsupdate命令来扫描更新本机上存放的钱包文件，增加walletslist来打印本机上存放的所有钱包的基本信息， 增加sendbyrefname来通过钱包别名实现交易创建的功能（这个功能在实际区块链系统中不需要，这里只是为了演示方便）。同时createblockchain，balance功能也将支持用别名指定钱包。

回到walletfunction.go 构建createwallet功能实现函数。

```go
//createwallet
//用于创建新的钱包
func (cli *CommandLine) createwallet(refname string) {
	newWallet := wallet.NewWallet()
	newWallet.Save()
	refList := wallet.LoadRefList()
	refList.BindRef(string(newWallet.Address()), refname)
	refList.Save()
	fmt.Println("成功创建钱包。")
}
```

构建walletinfo功能实现函数。

```go
//walletinfo
//打印指定钱包的基本信息(别名或者地址)
func (cli *CommandLine) walletInfoRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.walletInfo(address)
}
func (cli *CommandLine) walletInfo(address string) {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()
	fmt.Printf("钱包地址:%x\n", wlt.Address())
	fmt.Printf("公钥:%x\n", wlt.PublicKey)
	fmt.Printf("别名:%s\n", (*refList)[address])
}
```

构建walletsUpdate功能实现函数

```go
//walletsupdate
//扫描更新本机上存放的钱包文件
func (cli *CommandLine) walletsUpdate() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("成功更新钱包。")
}
```

构建walletslist功能实现函数。

```go
//walletslist
//walletslist来打印本机上存放的所有钱包的基本信息
func (cli *CommandLine) walletslist() {
	refList := wallet.LoadRefList()
	for address, _ := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Wallet address:%s\n", address)
		fmt.Printf("Public Key:%x\n", wlt.PublicKey)
		fmt.Printf("Reference Name:%s\n", (*refList)[address])
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}
}
```

回到blockfunction.go

增加sendbyrefname功能实现函数。

```go
//sendbyrefname
//通过钱包别名实现交易创建的功能（这个功能在实际区块链系统中不需要，这里只是为了演示方便）
func (cli *CommandLine) sendRefName(fromRefname, toRefname string, amount int) {
	refList := wallet.LoadRefList()
	fromAddress, err := refList.FindRef(fromRefname)
	utils.Handle(err)
	toAddress, err := refList.FindRef(toRefname)
	utils.Handle(err)
	cli.send(fromAddress, toAddress, amount)
}
```

可以看见sendbyrefname也是通过调用send函数实现的。

增加支持用别名调用createblockchain命令的函数。

```go
//通过别名（创始人的）创建区块链
func (cli *CommandLine) createBlockChainRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.createblockchain(address)
}
```

增加支持用别名调用balance命令的函数。

```go
//用别名调用balance命令查询余额
func (cli *CommandLine) balanceRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.balance(address)
}
```

回到register.go重新实现Run函数。

```go
func (cli *CommandLine) Run() {
	cli.validateArgs()

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)    // new command
	walletInfoCmd := flag.NewFlagSet("walletinfo", flag.ExitOnError)        // new command
	walletsUpdateCmd := flag.NewFlagSet("walletsupdate ", flag.ExitOnError) // new command
	walletsListCmd := flag.NewFlagSet("walletslist", flag.ExitOnError)      // new command
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendByRefNameCmd := flag.NewFlagSet("sendbyrefname", flag.ExitOnError) // this is never needed in real blockchain system.
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	createWalletRefName := createWalletCmd.String("refname", "", "The refname of the wallet, and this is optimal") // this line is new
	walletInfoRefName := walletInfoCmd.String("refname", "", "The refname of the wallet")                          // this line is new
	walletInfoAddress := walletInfoCmd.String("address", "", "The address of the wallet")                          // this line is new
	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	createBlockChainByRefNameOwner := createBlockChainCmd.String("refname", "", "The name refer to the owner of blockchain") // this line is new
	balanceAddress := balanceCmd.String("address", "", "Who needs to get balance amount")
	balanceRefName := balanceCmd.String("refname", "", "Who needs to get balance amount") // this line is new
	sendByRefNameFrom := sendByRefNameCmd.String("from", "", "Source refname")            // this line is new
	sendByRefNameTo := sendByRefNameCmd.String("to", "", "Destination refname")           // this line is new
	sendByRefNameAmount := sendByRefNameCmd.Int("amount", 0, "Amount to send")            // this line is new
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createwallet": // this case is new
		err := createWalletCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletinfo": // this case is new
		err := walletInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletsupdate": // this case is new
		err := walletsUpdateCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "walletslist": // this case is new
		err := walletsListCmd.Parse(os.Args[2:])
		utils.Handle(err)

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

	case "sendbyrefname": // this case is new
		err := sendByRefNameCmd.Parse(os.Args[2:])
		utils.Handle(err)

	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		utils.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createWalletCmd.Parsed() {
		cli.createwallet(*createWalletRefName)
	}

	if walletInfoCmd.Parsed() {
		if *walletInfoAddress == "" {
			if *walletInfoRefName == "" {
				walletInfoCmd.Usage()
				runtime.Goexit()
			} else {
				cli.walletInfoRefName(*walletInfoRefName)
			}
		} else {
			cli.walletInfo(*walletInfoAddress)
		}
	}

	if walletsUpdateCmd.Parsed() {
		cli.walletsUpdate()
	}

	if walletsListCmd.Parsed() {
		cli.walletslist()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			if *createBlockChainByRefNameOwner == "" {
				createBlockChainCmd.Usage()
				runtime.Goexit()
			} else {
				cli.createBlockChainRefName(*createBlockChainByRefNameOwner)
			}
		} else {
			cli.createblockchain(*createBlockChainOwner)
		}
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			if *balanceRefName == "" {
				balanceCmd.Usage()
				runtime.Goexit()
			} else {
				cli.balanceRefName(*balanceRefName)
			}
		} else {
			cli.balance(*balanceAddress)
		}
	}

	if sendByRefNameCmd.Parsed() {
		if *sendByRefNameFrom == "" || *sendByRefNameTo == "" || *sendByRefNameAmount <= 0 {
			sendByRefNameCmd.Usage()
			runtime.Goexit()
		}
		cli.sendRefName(*sendByRefNameFrom, *sendByRefNameTo, *sendByRefNameAmount)
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

