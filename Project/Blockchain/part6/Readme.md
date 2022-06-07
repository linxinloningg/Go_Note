# 深入底层：Go语言从零构建区块链（六）：签名与验证

## 前言
在第五章我们理解学习了密钥对，在密钥对的基础上理解了钱包的概念，并且阐述了各种地址（钱包地址）的产生方式与能否可逆变换。事实上，我们在上一章只使用了公钥就实现了地址的产生，并使得区块链中的资产可以流入所产生的特定的地址。为使用流向这些地址的资产，用户需要证明对这些特定地址的拥有权。为了实现这一目的，需要配合地使用密钥对的私钥与公钥对交易信息进行签名与验证。

## 公钥与公钥哈希
在上一章中，交易信息的Input与Output地址我们都暂时的使用了钱包地址来进行表征。对于接触区块链不深的读者而言，可能会自然的想到以密钥对的公钥来表征Input与Output的地址，这样配合密钥对的签名与验证能够实现用户对UTXO流向的公钥地址的拥有权证明，这种方式称为P2PK（pay to public key）。事实上，主流的区块链系统选择的是用公钥表征Input的地址，用公钥哈希表征Output的地址，这种方式称为P2PKH（pay to public key hash）。本教程后续也是使用P2PKH的方式，读者可以思考一下使用P2PKH的好处，这里我也说说我的理解。使用公钥哈希作为Output的地址可以进一步提升区块链系统中交易的匿名性，用户的公钥只会在需要转移其资产时才会披露，同时公钥哈希占据更小的存储位。

就此，我们需要在项目上进行一些小改动。打开inoutput.go，进行以下修改。

```go
package transaction

import (
	"bytes"
	"part6/src/utils"
)

/*
TxOutput将记录HashPubkey（公钥哈希）作为地址，TxInput将记录PubKey（公钥）作为地址
*/

//转
type TxOutput struct {
	Value      int    //转出的资产值
	HashPubKey []byte //资产的接收者的地址哈希（ToAddress）
}

//收
type TxInput struct {
	TxID   []byte //指明支持本次交易的前置交易信息
	OutIdx int    //具体指明是前置交易信息中的第几个Output
	PubKey []byte //资产转出者的地址（FromAddress）
	Sig    []byte //签名认证
}

//验证FromAddress是否正确
func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.PubKey, address)
}

//验证ToAddress是否正确
func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.HashPubKey, utils.PublicKeyHash(address))
}
```

可以看到，TxOutput将记录HashPubkey（公钥哈希）作为地址，TxInput将记录PubKey（公钥）作为地址。同时我们发现，TxInput结构体多了一个Sig属性，该属性即是将要讲到的签名信息。

由于TxInput与TxOutput结构体的改变，项目中很多已有的函数开始报错，我们先只做以下修改。

打开transaction.go，修改BaseTx函数。

```go
//创区块交易
func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}, nil}
	txOut := TxOutput{constcode.InitCoin, toaddress}
	tx := Transaction{[]byte("这是创始区块交易！"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}
```

## ECDSA签名算法
非对称密钥技术的一个重要应用就是对信息进行签名,不同的种类的密钥技术签名算法不同，在ECC则是ECDSA。ECDSA的理论基础与原理需要读者自行学习，在第五章中也些有介绍，本教程在此只以结合代码的方式来讲解ECDSA的实现方式。

在util.go中，更新引入下列库包。
```go
//util.go
package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"goblockchain/constcoe"
	"log"
	"math/big"
	"os"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)
```

构造签名函数。

```go
//签名函数
/*
privKey是私钥，msg则是待签名的信息。
注意到r是一个随机的大数，s则是基于r、私钥与待签名信息经过一系列计算得到的另一个大数，由(r,s)组成的元组就是签名信息
*/
func Sign(msg []byte, privKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, msg)
	Handle(err)
	signature := append(r.Bytes(), s.Bytes()...) //(r,s)组成的元组
	return signature
}
```

其中privKey是私钥，msg则是待签名的信息。注意到r是一个随机的大数，s则是基于r、私钥与待签名信息经过一系列计算得到的另一个大数，由(r,s)组成的元组就是签名信息。

我们可以通过公钥，待签名的信息，签名信息，来验证签名的真实性与权威性，构造验证函数如下。

```go
//验证函数
/*
公钥，待签名的信息，签名信息，来验证签名的真实性与权威性
 */
func Verify(msg []byte, pubkey []byte, signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	return ecdsa.Verify(&rawPubKey, msg, &r, &s)
}
```

## 交易信息的签名与验证
比特币中的签名与验证是通过锁定脚本与解锁脚本实现的，读者可以自行搜索资料进行学习，本教程将二者的功能抽离出来进行实现，但不会出现严格的脚本概念。

用户想要使用指向特定地址的资产时就需要通过签名来证明自己对这些地址的拥有权。想像这样的一个场景，用户A拥有a这一钱包，用户B拥有b这一钱包。在区块链中有3个UTXO流向了a对应的公钥哈希地址，总值为5币。现在用户A想要转账5币给用户B，需要生成交易信息，于是A便构建三个Input来引用前述的3个UTXO，同时将用户B提供的b钱包地址计算得到的公钥哈希作为Output的地址。A为了证明对3个UTXO的使用权，需要使用私钥对整个交易过程进行签名，并将签名信息作为交易信息的一部分向整个区块链扩散。为了验证这样的一个交易信息的有效性，需要同时获得交易信息代表的交易过程，三个UTXO指向的公钥哈希地址（也即是a所对应的哈希公钥），A的公钥，A提供的签名信息。

打开transaction.go，我们构造PlainCopy函数来描述一个交易信息的交易过程。
```go
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
```

接着构造一PlainHash函数用以辅助对交易信息进行签名。
```go
//PlainHash函数用以辅助对交易信息进行签名
func (tx *Transaction) PlainHash(inidx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKey
	return txCopy.TxHash()
}
```
现在可以构造交易信息的签名函数了。
```go
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
```
这里可能有的读者就会问了，为什么不直接对整个交易信息进行一次签名，而是要对每个input做一次签名。这个是一个有意义的问题，同时也极少被讨论，可以参考下面这个知乎问题。[对于区块中的一个交易，比特币为什么不直接对这个交易整体进行签名，而是对该交易中的每个交易输入分别签名？ - 知乎](https://www.zhihu.com/question/315268017)

这里也说说我自己的观点，首先是每个input指向的UTXO的锁定脚本版本号与方式可能不同，需要对每一个input单独进行签名，其次是从宏观的角度讲，对整个交易信息签名本身是没有意义的，签名的意义在于使用每个input引用的UTXO ，该过程的参与者应该是当前交易信息的创建者与包含前述UTXO的交易信息的创建者。

接下来，我们可以构造交易信息的验证函数。

```go
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
```

在实现了交易信息的签名与验证后，我们需要重构我们的交易生成函数。打开blockchain.go，重构CreateTransaction函数。

```go
//创建交易
//可以用一个输入对于多个输出
func (bc *BlockChain) CreateTransaction(from_PubKey, to_HashPubKey []byte, amount int, privkey ecdsa.PrivateKey) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	accumulated, validOutputs := bc.FindSpendableOutputs(from_PubKey, amount)

	//没有足够数量的余额
	if accumulated < amount {
		fmt.Println("没有足够数量的余额!")
		return &transaction.Transaction{}, false
	}

	//转
	for ID, i := range validOutputs {
		txID, err := hex.DecodeString(ID)
		utils.Handle(err)
		input := transaction.TxInput{txID, i, from_PubKey, nil}
		inputs = append(inputs, input)
	}

	//收
	outputs = append(outputs, transaction.TxOutput{amount, to_HashPubKey})

	//找零
	if accumulated > amount {
		outputs = append(outputs, transaction.TxOutput{accumulated - amount, utils.PublicKeyHash(from_PubKey)})
	}

	//一个输入对应多个输出
	tx := transaction.Transaction{nil, inputs, outputs}

	tx.SetID()

	tx.Sign(privkey)
	return &tx, true
}
```

## 重构Mine
在第三章我们说到此前的区块链系统存在双花现象，这是我们在Mine的时候没有对交易信息进行验证就构造候选区块造成的。一个完整的Mine过程在构造候选区块时应该先检查交易信息池中的所有交易信息的有效性，这包括验证是否引用了已花费的Output，是否重复引用了同一UTXO，Input与Output资产总额是否对应，交易信息的签名信息。

打开mine.go，更新下列库包。

```go
package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"part6/src/transaction"
	"part6/src/utils"
)
```

构建验证交易池中的交易信息有效性的函数。

```go
//验证交易池中的交易信息有效性
func (bc *BlockChain) VerifyTransactions(txs []*transaction.Transaction) bool {
	if len(txs) == 0 {
		return true
	}
	//TODO: The following method to verify the transactions is to query the blockchain for
	//unspent outputs and is definitely right. However, I believe in the near future, an unspent
	//outputs database can be maintained to accelerate the verification.
	spentOutputs := make(map[string]int)
	for _, tx := range txs {
		pubKey := tx.Inputs[0].PubKey
		unspentOutputs := bc.FindUnspentTransactions(pubKey)
		inputAmount := 0
		OutputAmount := 0

		for _, input := range tx.Inputs {
			if outidx, ok := spentOutputs[hex.EncodeToString(input.TxID)]; ok && outidx == input.OutIdx {
				return false
			}
			ok, amount := isInputRight(unspentOutputs, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TxID)] = input.OutIdx
		}

		for _, output := range tx.Outputs {
			OutputAmount += output.Value
		}
		if inputAmount != OutputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}
	}
	return true
}

```
正如注释所说，这里采用的验证是否引用了已花费的Output的方法是非常简单的，每次验证需要询问并查找整个区块链，在未来的章节中我们可以建立一个本地UTXO数据库来进行优化。

现在更新我们的Mine函数。

```go
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}

	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx) //PoW has been done here.
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
```
可以看到我们就是在构造候选区块前调用了我们构造的VerifyTransactions函数。

## CLI更新与调试
由于我们不再以钱包地址作为交易信息中input与output的地址，故需要更新命令行模块中的各函数以实现原有功能。打开cli.go，将以下函数更新。

```go
func (cli *CommandLine) createblockchain(address string) {
	newChain := blockchain.InitBlockChain(utils.Address2PubHash([]byte(address)))

	//注意在使用完数据库后，需要使用newChain.Database.Close()函数关闭数据库
	newChain.Database.Close()
	fmt.Println("完成创建区块链，所有者是：", address)
}

//balance
/*
先使用ContinueBlockChain函数接入数据库，然后使用FindUTXO函数统计余额并打印
*/
func (cli *CommandLine) balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	wlt := wallet.LoadWallet(address)

	balance, _ := chain.FindUTXOs(wlt.PublicKey)
	fmt.Printf("地址:%s, 余额:%d \n", address, balance)
}

//send
//调用CreateTransaction函数，并将创建的交易信息保存到交易信息池中
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	fromWallet := wallet.LoadWallet(from)
	tx, ok := chain.CreateTransaction(fromWallet.PublicKey, utils.Address2PubHash([]byte(to)), amount, fromWallet.PrivateKey)
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

Congratulations！一个完整的区块链系统就此诞生（当然是没有考虑分布式的情况下qvq）。

这里我们在test.bat中增加几条命令来测试双花问题（其实在实现了分布式网络部署后仍然存在双花风险，但这留待我们日后再说）是否已经解决。

```bash
rd /s /q tmp
md tmp\blocks
md tmp\wallets
md tmp\ref_list
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
测试结果如下。

    D:\learngo\goblockchain>test.bat
    
    D:\learngo\goblockchain>rd /s /q tmp
    
    D:\learngo\goblockchain>md tmp\blocks
    
    D:\learngo\goblockchain>md tmp\wallets
    
    D:\learngo\goblockchain>md tmp\ref_list
    
    D:\learngo\goblockchain>main.exe createwallet
    Succeed in creating wallet.
    
    D:\learngo\goblockchain>main.exe walletslist
    --------------------------------------------------------------------------------------------------------------
    Wallet address:1C1SvPjPmSuBMZziNvvTDjVZLFgNX4Z5AS
    Public Key:bbeadcb894f618e5987db38299b82e88d6191efbaa70b3d57d2ebcf87dde3be75b93ce6d16c838e3408a72e326948be0694674ddd29a7e47df8927d145a1771a
    Reference Name:
    --------------------------------------------------------------------------------------------------------------


    D:\learngo\goblockchain>main.exe createwallet -refname LeoCao
    Succeed in creating wallet.
    
    D:\learngo\goblockchain>main.exe walletinfo -refname LeoCao
    Wallet address:314c73676173577571314a56783958756d793976727357625431574741525931434b
    Public Key:4aa42bb071009d3b15c67305e860dc9a96eb24ef227f8459f70aeea8ccbbfa9da30da207047096de05362d2d0ce25f1c1232c5b0941dc37bf3c44d2eafe4735b
    Reference Name:LeoCao
    
    D:\learngo\goblockchain>main.exe createwallet -refname Krad
    Succeed in creating wallet.
    
    D:\learngo\goblockchain>main.exe createwallet -refname Exia
    Succeed in creating wallet.
    
    D:\learngo\goblockchain>main.exe createwallet
    Succeed in creating wallet.
    
    D:\learngo\goblockchain>main.exe walletslist
    --------------------------------------------------------------------------------------------------------------
    Wallet address:1C1SvPjPmSuBMZziNvvTDjVZLFgNX4Z5AS
    Public Key:bbeadcb894f618e5987db38299b82e88d6191efbaa70b3d57d2ebcf87dde3be75b93ce6d16c838e3408a72e326948be0694674ddd29a7e47df8927d145a1771a
    Reference Name:
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Wallet address:1JTEZEvMPyBqqajgn2j7qrLc74j2sRgqPH
    Public Key:01b460490915b31f223e957d09ea124dbcb85a4e5a8421e6592496363fd9abcd955caace1411064469b27552922e8fe39f9467d0e615f805a6ff02bbb3ce7f85
    Reference Name:Exia
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Wallet address:12ZCV3cXjm3LrpJQRUKXbk1NZJo7usKdEr
    Public Key:98f603302df5cc17c40147b0de41394c485bb140697b768a76ccb93aae90465cbb40b66801764e9de6c181c43a90e3528d46ed8d0d6ca2cbcaca5bacf4711d15
    Reference Name:
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Wallet address:1LsgasWuq1JVx9Xumy9vrsWbT1WGARY1CK
    Public Key:4aa42bb071009d3b15c67305e860dc9a96eb24ef227f8459f70aeea8ccbbfa9da30da207047096de05362d2d0ce25f1c1232c5b0941dc37bf3c44d2eafe4735b
    Reference Name:LeoCao
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Wallet address:1LNagMDsH2DGRALq8hghSAGTDqN5tgf7FC
    Public Key:10766f210599b21bcb23151afcfb7335ecdeef1335d1bffa44853da5660794c8ca9b0a916b79f2f0b7f8665a2c11bdb97a2d3487176ae4d30c7d7f523d123d25
    Reference Name:Krad
    --------------------------------------------------------------------------------------------------------------


    D:\learngo\goblockchain>main.exe createblockchain -refname LeoCao
    Genesis Created
    Finished creating blockchain, and the owner is:  1LsgasWuq1JVx9Xumy9vrsWbT1WGARY1CK
    
    D:\learngo\goblockchain>main.exe blockchaininfo
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:4c656f2043616f20697320617765736f6d6521
    Transactions:[0xc0000541e0]
    hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Pow: true
    --------------------------------------------------------------------------------------------------------------


    D:\learngo\goblockchain>main.exe balance -refname LeoCao
    Address:1LsgasWuq1JVx9Xumy9vrsWbT1WGARY1CK, Balance:1000
    
    D:\learngo\goblockchain>main.exe sendbyrefname -from LeoCao -to Krad -amount 100
    Success!
    
    D:\learngo\goblockchain>main.exe balance -refname Krad
    Address:1LNagMDsH2DGRALq8hghSAGTDqN5tgf7FC, Balance:0
    
    D:\learngo\goblockchain>main.exe mine
    Finish Mining
    
    D:\learngo\goblockchain>main.exe blockchaininfo
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Transactions:[0xc0000f0190]
    hash:0c3261f50269e03d8612c50cc62631194707569e7665fa77658ddcf012f6668a
    Pow: true
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:4c656f2043616f20697320617765736f6d6521
    Transactions:[0xc0000f0230]
    hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Pow: true
    --------------------------------------------------------------------------------------------------------------


    D:\learngo\goblockchain>main.exe balance -refname LeoCao
    Address:1LsgasWuq1JVx9Xumy9vrsWbT1WGARY1CK, Balance:900
    
    D:\learngo\goblockchain>main.exe balance -refname Krad
    Address:1LNagMDsH2DGRALq8hghSAGTDqN5tgf7FC, Balance:100
    
    D:\learngo\goblockchain>main.exe sendbyrefname -from LeoCao -to Exia -amount 100
    Success!
    
    D:\learngo\goblockchain>main.exe sendbyrefname -from Krad -to Exia -amount 30
    Success!
    
    D:\learngo\goblockchain>main.exe mine
    Finish Mining
    
    D:\learngo\goblockchain>main.exe blockchaininfo
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335743
    Previous hash:0c3261f50269e03d8612c50cc62631194707569e7665fa77658ddcf012f6668a
    Transactions:[0xc000184140 0xc0001841e0]
    hash:25dd758d306eaa82672ad3656fec14520bac03b315f132d4034d7adfa10f78b8
    Pow: true
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Transactions:[0xc000184280]
    hash:0c3261f50269e03d8612c50cc62631194707569e7665fa77658ddcf012f6668a
    Pow: true
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:4c656f2043616f20697320617765736f6d6521
    Transactions:[0xc000184320]
    hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Pow: true
    --------------------------------------------------------------------------------------------------------------


    D:\learngo\goblockchain>main.exe balance -refname LeoCao
    Address:1LsgasWuq1JVx9Xumy9vrsWbT1WGARY1CK, Balance:800
    
    D:\learngo\goblockchain>main.exe balance -refname Krad
    Address:1LNagMDsH2DGRALq8hghSAGTDqN5tgf7FC, Balance:70
    
    D:\learngo\goblockchain>main.exe balance -refname Exia
    Address:1JTEZEvMPyBqqajgn2j7qrLc74j2sRgqPH, Balance:130
    
    D:\learngo\goblockchain>main.exe sendbyrefname -from Exia -to LeoCao -amount 90
    Success!
    
    D:\learngo\goblockchain>main.exe sendbyrefname -from Exia -to Krad -amount 90
    Success!
    
    D:\learngo\goblockchain>main.exe mine
    2021/11/19 23:29:03 falls in transactions verification
    Finish Mining
    
    D:\learngo\goblockchain>main.exe blockchaininfo
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335743
    Previous hash:0c3261f50269e03d8612c50cc62631194707569e7665fa77658ddcf012f6668a
    Transactions:[0xc0000f0190 0xc0000f0230]
    hash:25dd758d306eaa82672ad3656fec14520bac03b315f132d4034d7adfa10f78b8
    Pow: true
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Transactions:[0xc0000f02d0]
    hash:0c3261f50269e03d8612c50cc62631194707569e7665fa77658ddcf012f6668a
    Pow: true
    --------------------------------------------------------------------------------------------------------------
    
    --------------------------------------------------------------------------------------------------------------
    Timestamp:1637335742
    Previous hash:4c656f2043616f20697320617765736f6d6521
    Transactions:[0xc0000f0370]
    hash:a7bbf71e7144084a57a0995c013c86c1b4e9acc66e583e9f5d0df21b519fbc61
    Pow: true
    --------------------------------------------------------------------------------------------------------------


    D:\learngo\goblockchain>main.exe balance -refname LeoCao
    Address:1LsgasWuq1JVx9Xumy9vrsWbT1WGARY1CK, Balance:800
    
    D:\learngo\goblockchain>main.exe balance -refname Krad
    Address:1LNagMDsH2DGRALq8hghSAGTDqN5tgf7FC, Balance:70
    
    D:\learngo\goblockchain>main.exe balance -refname Exia
    Address:1JTEZEvMPyBqqajgn2j7qrLc74j2sRgqPH, Balance:130

从Exia转账90给LeoCao，然后再转账90给Krad其实都是引用的相同的UTXO，也就构成了双花，而在测试中注意到记录了“2021/11/19 23:29:03 falls in transactions verification”这一错误，说明该双花被检测出来，程序功能正常。

## 总结
本章着重讲解了如何使用私钥签名并配合公钥证明并验证交易信息的有效性，就此本教程建立的区块链系统已能够能够在本地完整地运行了，这也是本教程第一部分的全部内容，此后会介绍一下Merkle Tree和SPV的实现，以及建立本地UTXO数据库优化本地交易信息与区块的查找效率，在本教程的第二部分将全面进入分布式网络相关功能的实现
