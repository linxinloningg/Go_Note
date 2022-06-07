# 深入底层：Go语言从零构建区块链（六）附1：梅克尔树和SPV

## 前言
在第六章我们完成了该教程的上半部分，在进入下一部分前，作为过渡章节，我期望再为我们的goblockchain添加一些常见的功能，就比如本文即将介绍的区块中的Merkle Tree（梅克尔树，后面简称MT）和区块链系统的SPV（快速交易验证）功能。对于接触了区块链一段时间的读者应该都有了解过MT，尽管MT的数据结构非常容易理解，但是我认为MT对于区块链的意义以及其涉及的SPV功能是认知整个区块链技术的关键步骤之一,通过对MT与SPV的深入理解与思考，我们也许能够窥探区块链在未来的应用场景以及更好地判断区块链技术当前面临的技术瓶颈。同时我相信大多数读者对于MT在区块链中的具体实现以及如何通过MT树实现SPV的过程是存在理解偏差的，本文将通过代码的形式讲述MT与SPV的实现，希望读者能够重点关注一些唯独在具体实现才会遇见的问题，这些隐藏在代码里的问题的解决可以帮助读者宏观把控区块链系统的设计理念，这也是本教程的初衷。

## 区块链中的Merkle Tree
MT是一种树形结构，其节点保存的值都是哈希值，所以也可以叫做哈希树。要理解MT，我们首先应理解哈希值有什么用。哈希值的一般特点可以总结如下。
* 哈希值是对一段信息的总结性描述，当原信息变动哈希值也会改变。可以用于验证信息的完整性。
* 知道哈希值几乎不可能逆推导原信息。
* 不同长度的信息使用哈希算法得到的哈希值长度相同。

MT由哈希值构成，MT可以是二叉树也可以多叉树，其任意父节点存储的哈希值为两子节点哈希值拼接后做哈希得到。要理解MT就是要理解MT中的叶子节点，MT中的任意叶子节点保存的哈希值改变则MT根节点的哈希值也会改变。

想象这样一种场景，即一个客户端向分布式文件存储服务器请求下载一个大型文件，客户端在服务器引导下从多个节点下载了文件片段，那么客户端怎样才能正确的拼凑出这个大型文件并检验其完整性了？在这种情况下，将每个文件片段按拼凑顺序排列计算哈希值作为MT的叶节点构造MT，那么服务器就能够通过给客户端发送该大型文件对应的MT来辅助客户端对该大型文件的片段排序，同时若大型文件的任一片段出错，MT的根节点存储的哈希值都会改变，客户端可以比对自己构造的MT根节点哈希值与服务器发来的MT根节点哈希值判断下载的大型文件的完整性，如果是非完整的也能能够通过比对MT树快速定位到是那一片段的文件出错。

结合上述场景，我们可以理解区块链中的MT的作用。将大型文件看作是一个区块的区块体，文件片段看作交易信息，则区块体是由很多交易信息按顺序排列组成。如下图所示，以交易信息的哈希值为叶子节点构造MT，那么通过比对区块体的MT树根（区块头中有该哈希值的副本）就能够判断整个区块体的完整性。

## 快速交易验证SPV
如果没有SPV功能，区块链中MT的作用也就到上文为止了，而且你可能或多或少地认为MT树的作用甚微，因为如果要检查区块体的完整性，可以直接使用区块头部中保存的哈希值直接验证整个区块的完整性。事实上，要理解区块链中的MT就要深入理解SPV是什么，有什么用。

SPV英文为Simplified Payment Verification，即快速交易验证。下面例举SPV的使用场景以及实现流程。我们假设小C响应国家号召推了个小车在夜市摆摊卖炸串，小C的手机性能孱弱而且没有多少剩余流量，但是小C也想使用区块链，那么他就将整条区块链的区块的头部信息下载下来维护，同时也不会去参与挖矿与共识。我们将小C这样的手机称为轻节点，与之相对的要参与挖矿并共识的节点称为全节点，轻节点可以挂载到多个全节点下面请求全节点帮助其维护区块头部信息。拿比特币为例，一个区块的大小为1MB左右，而区块头信息只有80KB，相较而言轻节点的通信消耗和计算消耗都大大减少，就算把整个区块链的头部信息存储下来都不过30G（事实上不需要全部存储，全部存储的话安全性最大），相较而言现在的手机设备能够勉强负担得起这个存储大小。当有人来买炸串时，小C指了指招牌上的二维码，客户使用手机扫描这个二维码后发现是这是一笔交易，交易的中输出方和金额已经填好，客户使用手机钱包将该交易补全后直接将该交易信息发给小C的手机同时在区块链网络中进行发布，过了几秒钟（实际的区块链网络要10分钟）告诉小C该笔交易已经完成，并将该笔交易信息直接发给了小C。但是小C多了个心眼，他想知道这笔交易是否真的完成了，现在的情况下，小C的手机上包括：已经更新的只包含区块头的区块链，客户发来的交易信息。小C接下来进行SPV，他向任一全节点发起对该交易信息的SPV请求，全节点向小C返回一区块ID以及一个MT验证路径。小C首先在自己只包含区块头的区块链上查找区块ID找到一区块头部信息，该区块头部信息中包含一MT树根哈希值，同时小C再根据从客户那收到的交易信息按照MT验证路径重新计算得到一MT树根哈希值，比对两MT树根哈希值就能够判断该笔交易信息是否真的已经完成。以上图进行讲解，如果该笔交易信息为交易3，则全节点返回的MT验证路径就应该是Hash4->Hash12->Hash5~8。在这个SPV过程中需要注意：
* 整个MT验证路径不大于1KB，通信消耗小。
* 小C不需要担心全节点伪造了一个验证路径，因为知道交易信息3的哈希值Hash3与MT根的哈希值Hash1\~8是不可能逆推导出Hash4，Hash12，Hash5\~8的。

以上就是SPV的一个经典流程。

不知道读者有没有思考过，当我们不使用MT，则全节点也可以发送一个包含交易3的区块给小C实现对该笔交易的确认。小C在收到一完整区块后首先确认其是否在自己维护的只有区块头的区块链上，然后确认交易3是否在该区块中，最后通过计算整个区块的哈希值并与区块头中的哈希值比对来确认交易3已经完成上链共识。相较于启用MT后传输的不到1KB的MT验证路径，1MB的区块在后5G时代似乎也不是不能接收。常规的解释会认为这就是一个量变引起质变的问题，当小C确认交易信息这一操作变得频繁时，SPV可以节省巨大的通信消耗。毫无疑问这样的解释是合理的，但是我认为SPV的实现让我们有了遐想区块链在未来应用的空间。

严谨的SPV讲解就到上文为止了，以下观点为我个人的思考，还请读者自行斟酌。

SPV的实现也许可以为我们提供一种交易线下验证的途径，在上述小C的例子中，假如小C不仅维护一只保留区块头的区块链，还把与他相关的交易信息用MT验证路径的形式存储，另外一个小K也是如此。假设在线下的情况下，即小C和小K面交，手机皆不联网的情况下，小K想要购买小C的炸串，于是小K根据小C所给的钱包地址创建了一笔交易信息，将该笔交易以及能够验证这笔交易信息中输入的MT验证路径一并以蓝牙（NFC）空投给小C，那么小C就能够通过自身保存的只保留区块头的区块链以及小K所给的MT验证路径判断该笔交易信息是否有效，如有效小K即可离开，小C在摆摊结束回家后一次性将形如小K一样的交易信息一并在区块链网络发布即可。如此的一个交易逻辑是否已经和我们线下使用RMB交易已经很相似了，小K提供MT验证路径相当于在给小K确认RMB的面值，然后创建一笔交易信息相当于将RMB交给了小C。这个交易流程存在着两个关键技术难点，首先是如何保证小K提供了MT验证路径的交易信息没有被使用过，其次是如何保证小C回到家的这段时间小K没有再用这些交易信息去做其它买卖。一种简单的保证方法，那就是小C和小K都使用一种权威可信的APP管理钱包和做交易，这个APP的出品方一定要是大型权威机构（没错，他可以是国家和ZF，这是不是就像数字RMB的双离线功能了，当然数字RMB的原理细节是不可知的，但我们可以猜测）。

现在让我们的思路再拓宽一些，回归到SPV功能的本质。SPV功能的本质就是能够根据MT验证路径确认一交易信息存在于区块链上，在上述的几个例子中，区块链中的交易信息保存的是实实在在的交易，而在如今广阔的区块链系统中，交易信息可以保存任何有意义的信息，不再局限于金融方面的应用。继续上述小C与小K的例子，小K说他要用一本电子书《Just for Fun》的资源与小C交易，即小C给小K炸串，小K给小C生成一条交易信息，该交易信息写明了批准小C从一权威服务器上下载小K上传的电子书《Just for Fun》，同时小K提供了一些附带MT验证路径的交易信息，这些交易信息包括的内容是电子书被上传至服务器、电子书长期有效、电子书不可撤回。与前述例子不同的是，此时我们将不关心小K会再用这些具有MT验证路径的交易信息去进行其它买卖，因为电子书本来就具有可复制性，以这种思路出发，小C与小K间的交易不但可以线下进行，而且根本不用在意区块链的共识，小C可以不急不缓的在需要的时候在区块链网络上发布小K所给的交易信息即可获得电子书《Just for Fun》。由这个例子，读者还可以继续自己思考，将下载一次电子书这一操作看作是任何一种服务操作，窥探一下区块链在未来可能的应用场景。

## MT的代码实现
在完成了对MT的原理理解后，我们便可以为我们的goblockchain项目添加MT了。创建两个文件夹，分别为merkletree和test。

在merkletree文件夹中创建文件merkletree.go，写入我们将要用到的包。

```go
package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"part7/src/transaction"
	"part7/src/utils"
)
```
首先创建结构体MerkleNode来表征MT中的各节点，然后创建一个结构体MerkleTree用于表示MT，我们知道一般树形结构可以使用根节点来存储，所以MerkleTree中只保存了一个根节点。

```go
//结构体MerkleTree
type MerkleTree struct {
	RootNode *MerkleNode //根节点
}

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte // Hash
}
```

在创建了MT节点的结构体后，我们需要编写一个构造函数。我们知道每个MT枝节点值是由其左右子节点的哈希值拼凑再做哈希得到，那么就需要输入左右子节点，而对于MT的叶子节点不存在只有节点，而是直接保存交易信息的哈希值，所以需要输入交易信息的哈希值。以下为一个通用的MT节点构造函数。

```go
//先判断当前创建的是否为叶子节点，然后再根据节点类型不同创建MT节点
//MT节点构造
func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	tempNode := MerkleNode{}

	//叶子节点
	if left == nil && right == nil { //The leaf
		tempNode.Data = data
	} else {
		catenateHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(catenateHash)
		tempNode.Data = hash[:]
	}

	tempNode.LeftNode = left
	tempNode.RightNode = right

	return &tempNode
}
```
可以看到该函数首先判断当前创建的是否为叶子节点，然后再根据节点类型不同创建MT节点。

有了MT节点构造函数后，我们需编写MT构造函数，该函数的最终返回类型是一MerkleTree（其实就是MT的根节点）。一般的MT原理介绍一般会举八个叶子节点的例子，在这种情况下每一层都必有偶数个的节点可供拼接计算出更上层的节点，但实际上MT在构造时每一层的节点数目都可能为奇数。解决方法是，如果叶子节点为奇数个，我们复制最后一个叶子节点，这样就有偶数个叶子节点了。当发现某一层枝节点为奇数个时，我们将最后一个枝节点不做处理直接加入到上层枝节点最前面，然后剩余的本层枝节点两两合并生成上层枝节点。该MT构造函数如下。

```go
//MT构造
func CrateMerkleTree(txs []*transaction.Transaction) *MerkleTree {
	txslen := len(txs)
	if txslen%2 != 0 {
		txs = append(txs, txs[txslen-1])
	}

	var nodePool []*MerkleNode

	for _, tx := range txs {
		nodePool = append(nodePool, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodePool) > 1 {
		var tempNodePool []*MerkleNode
		poolLen := len(nodePool)
		if poolLen%2 != 0 { //Notice here, we place the remained node at the head of the upper layer
			tempNodePool = append(tempNodePool, nodePool[poolLen-1])
		}
		for i := 0; i < poolLen/2; i++ {
			tempNodePool = append(tempNodePool, CreateMerkleNode(nodePool[2*i], nodePool[2*i+1], nil))
		}
		nodePool = tempNodePool
	}

	merkleTree := MerkleTree{nodePool[0]}

	return &merkleTree
}
```
有的读者可能有疑惑，为什么当发现某一层枝节点为奇数个时，我们将最后一个枝节点不做处理直接加入到上层枝节点最前面而不是放在最后面。对于这两种方式的区别读者可以用9个叶子节点的情况来比较下异同。举个更加极端的例子，当有257个叶子节点时，采用后者构造MT时，第257的叶子节点的验证路径长度只有2，而其他节点均为9，而采用前者则能很好地避免这种极端情况的出现。

接下来我们需要构造一个路径搜索函数，用于返回一个叶子节点的MT验证路径。验证路径由两部分组成，一是指定用于拼凑的哈希值，二是该哈希值是拼在前面还是后面。整个路径搜索函数就是一个深度优先搜索算法，如下所示。

```go
//路径搜索
/*
输入以此为寻找的目标哈希值，已经走过的路径（route与hashroute，其中route为方向，0为左，1为右，hashroute保存的是哈希值）
*/
func (mn *MerkleNode) Find(data []byte, route []int, hashroute [][]byte) (bool, []int, [][]byte) {
	findFlag := false

	if bytes.Equal(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashroute
	} else {
		if mn.LeftNode != nil {
			route_t := append(route, 0)
			hashroute_t := append(hashroute, mn.RightNode.Data)
			findFlag, route_t, hashroute_t = mn.LeftNode.Find(data, route_t, hashroute_t)
			if findFlag {
				return findFlag, route_t, hashroute_t
			} else {
				if mn.RightNode != nil {
					route_t = append(route, 1)
					hashroute_t = append(hashroute, mn.LeftNode.Data)
					findFlag, route_t, hashroute_t = mn.RightNode.Find(data, route_t, hashroute_t)
					if findFlag {
						return findFlag, route_t, hashroute_t
					} else {
						return findFlag, route, hashroute
					}

				}
			}
		} else {
			return findFlag, route, hashroute
		}
	}
	return findFlag, route, hashroute
}
```
Find函数的输入以此为寻找的目标哈希值，已经走过的路径（route与hashroute，其中route为方向，0为左，1为右，hashroute保存的是哈希值）。Find函数是一个递归函数，我们再创建一个 BackValidationRoute函数来包装一下Find，这样代码的可读性更高。

```go
//输入为一交易信息的ID（也即交易信息的哈希值），返回的是验证路径与一个是否找到该交易信息的信号
func (mt *MerkleTree) BackValidationRoute(txid []byte) ([]int, [][]byte, bool) {
	ok, route, hashroute := mt.RootNode.Find(txid, []int{}, [][]byte{})
	return route, hashroute, ok
}
```
可以看到，BackValidationRoute的输入为一交易信息的ID（也即交易信息的哈希值），返回的是验证路径与一个是否找到该交易信息的信号。

最后我们实现一下SPV函数，这里的SPV函数主要功能就是按照MT验证路径验证交易信息是否有效，如果成功则返回True，否则返回False。

```go
func SimplePaymentValidation(txid, mtroothash []byte, route []int, hashroute [][]byte) bool {
	routeLen := len(route)
	var tempHash []byte
	tempHash = txid

	for i := routeLen - 1; i >= 0; i-- {
		if route[i] == 0 {
			catenateHash := append(tempHash, hashroute[i]...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else if route[i] == 1 {
			catenateHash := append(hashroute[i], tempHash...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else {
			utils.Handle(errors.New("error in validation route"))
		}
	}
	return bytes.Equal(tempHash, mtroothash)
}
```
至此，MT树与SPV的结构体与函数就完成了。

## 区块结构体重构
MT根节点应该存储在区块的头部信息中，所以我们需要重构我们的区块数据结构来引入MT与SPV。打开blockchain中的block.go文件，添加merkletree包。

```go
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"part6/src/transaction"
	"part6/src/utils"
	"time"
)
```
为Block结构体添加一个属性MTree。

```go
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
```
修改SetHash和CreateBlock函数以适配。

```go
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
```
至此我们的代码改动就已完成。

## SPV测试
由于我们现在的命令行程序功能有限，不能很好的模仿SPV的应用场景，故在本章使用go自带的go test来测试merkletree模块与SPV功能。在test文件夹下创建SPV_test.go文件，粘贴如下测试代码，读者如有兴趣可以读一读本人写的测试代码，增加更多的测试用例，在本教程以后的章节中使用go test对模块进行调试的情况会越来越多。

```go
//SPV_test.go
package test

import (
	"crypto/sha256"
	"fmt"
	"part7/src/blockchain"
	"part7/src/merkletree"
	"part7/src/transaction"
	"strconv"
	"strings"
	"testing"
)

func GenerateTransaction(outCash int, inAccount string, toAccount string, prevTxID string, outIdx int) *transaction.Transaction {
	prevTxIDHash := sha256.Sum256([]byte(prevTxID))
	inAccountHash := sha256.Sum256([]byte(inAccount))
	toAccountHash := sha256.Sum256([]byte(toAccount))
	txIn := transaction.TxInput{prevTxIDHash[:], outIdx, inAccountHash[:], nil}
	txOut := transaction.TxOutput{outCash, toAccountHash[:]}
	tx := transaction.Transaction{[]byte("This is the Base Transaction!"),
		[]transaction.TxInput{txIn}, []transaction.TxOutput{txOut}} //Whether set ID is not nessary
	tx.SetID()                                                      //Here the ID is reset to the hash of the whole transaction. Signature is skipped
	return &tx
}

var transactionTests = []struct {
	outCash   int
	inAccount string
	toAccount string
	prevTxID  string
	outIdx    int
}{
	{
		outCash:   10,
		inAccount: "LLL",
		toAccount: "CCC",
		prevTxID:  "prev1",
		outIdx:    1,
	},
	{
		outCash:   20,
		inAccount: "EEE",
		toAccount: "OOO",
		prevTxID:  "prev2",
		outIdx:    1,
	},
	{
		outCash:   30,
		inAccount: "OOO",
		toAccount: "EEE",
		prevTxID:  "prev3",
		outIdx:    0,
	},
	{
		outCash:   100,
		inAccount: "CCC",
		toAccount: "LLL",
		prevTxID:  "prev4",
		outIdx:    1,
	},
	{
		outCash:   50,
		inAccount: "AAA",
		toAccount: "OOO",
		prevTxID:  "prev5",
		outIdx:    1,
	},
	{
		outCash:   110,
		inAccount: "OOO",
		toAccount: "AAA",
		prevTxID:  "prev6",
		outIdx:    0,
	},
	{
		outCash:   200,
		inAccount: "LLL",
		toAccount: "CCC",
		prevTxID:  "prev7",
		outIdx:    1,
	},
	{
		outCash:   500,
		inAccount: "EEE",
		toAccount: "OOO",
		prevTxID:  "prev8",
		outIdx:    1,
	},
}

func GenerateBlock(txs []*transaction.Transaction, prevBlock string) *blockchain.Block {
	prevBlockHash := sha256.Sum256([]byte(prevBlock))
	testblock := blockchain.CreateBlock(prevBlockHash[:], txs)
	return testblock
}

var spvTests = []struct {
	txContained []int
	prevBlock   string
	findTX      []int
	wants       []bool
}{
	{
		txContained: []int{0},
		prevBlock:   "prev1",
		findTX:      []int{0, 1},
		wants:       []bool{true, false},
	},
	{
		txContained: []int{0, 1, 2, 3, 4, 5, 6, 7},
		prevBlock:   "prev2",
		findTX:      []int{3, 7, 5},
		wants:       []bool{true, true, true},
	},
	{
		txContained: []int{0, 1, 2, 3},
		prevBlock:   "prev3",
		findTX:      []int{0, 1, 5},
		wants:       []bool{true, true, false},
	},
	{
		txContained: []int{0, 3, 5, 6, 7},
		prevBlock:   "prev4",
		findTX:      []int{0, 1, 6, 7},
		wants:       []bool{true, false, true, true},
	},
	{
		txContained: []int{0, 1, 2, 4, 5, 6, 7},
		prevBlock:   "prev5",
		findTX:      []int{0, 1, 3},
		wants:       []bool{true, true, false},
	},
}

func TestSPV(t *testing.T) {
	primeTXs := []*transaction.Transaction{}
	for _, tx := range transactionTests {
		tx := GenerateTransaction(tx.outCash, tx.inAccount, tx.toAccount, tx.prevTxID, tx.outIdx)
		primeTXs = append(primeTXs, tx)
	}

	fmt.Println("TestSPV 开始...")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	for idx, test := range spvTests {
		fmt.Println("当前测试 No: ", idx)
		fmt.Println("默克尔树就像:")
		mtGraphPaint(test.txContained)
		txs := []*transaction.Transaction{}
		for _, txidx := range test.txContained {
			txs = append(txs, primeTXs[txidx])
		}
		testBlock := GenerateBlock(txs, test.prevBlock)
		fmt.Println("------------------------------------------------------------------")
		for num, findidx := range test.findTX {
			fmt.Println("查找交易:", findidx)
			fmt.Printf("交易ID: %x\n", primeTXs[findidx].ID)
			route, hashroute, ok := testBlock.MTree.BackValidationRoute(primeTXs[findidx].ID)
			if ok {
				fmt.Println("已找到验证路线: ", route)
				fmt.Println("路线就像:")
				routePaint(route)
			} else {
				fmt.Println("没有找到引用的交易")
			}
			spvRes := merkletree.SimplePaymentValidation(primeTXs[findidx].ID, testBlock.MTree.RootNode.Data, route, hashroute)
			fmt.Println("SPV 结果: ", spvRes, ", 想要的结果: ", test.wants[num])
			if spvRes != test.wants[num] {
				t.Errorf("测试 %d 发现 %d: SPV 不正确", idx, findidx)
			}
			fmt.Println("------------------------------------------------------------------")
		}
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	}
}

func mtGraphPaint(txContained []int) {
	mtLayer := [][]string{}
	bottomLayer := []string{}
	for i := 0; i < len(txContained); i++ {
		bottomLayer = append(bottomLayer, strconv.Itoa(txContained[i]))
	}
	if len(bottomLayer)%2 == 1 {
		bottomLayer = append(bottomLayer, bottomLayer[len(bottomLayer)-1])
	}
	mtLayer = append(mtLayer, bottomLayer)

	for len(mtLayer[len(mtLayer)-1]) != 1 {
		tempLayer := []string{}
		if len(mtLayer[len(mtLayer)-1])%2 == 1 {
			tempLayer = append(tempLayer, mtLayer[len(mtLayer)-1][len(mtLayer[len(mtLayer)-1])-1])
			mtLayer[len(mtLayer)-1][len(mtLayer[len(mtLayer)-1])-1] = "->"
		}
		for i := 0; i < len(mtLayer[len(mtLayer)-1])/2; i++ {
			tempLayer = append(tempLayer, mtLayer[len(mtLayer)-1][2*i]+mtLayer[len(mtLayer)-1][2*i+1])
		}

		mtLayer = append(mtLayer, tempLayer)
	}

	layers := len(mtLayer)
	fmt.Println(strings.Repeat(" ", layers-1), mtLayer[layers-1][0])
	foreSpace := 0
	for i := layers - 2; i >= 0; i-- {
		var str1, str2 string
		str1 += strings.Repeat(" ", foreSpace)
		str2 += strings.Repeat(" ", foreSpace)

		for j := 0; j < len(mtLayer[i]); j++ {
			str1 += strings.Repeat(" ", i+1)
			if j%2 == 0 {
				if mtLayer[i][j] == "->" {
					foreSpace += (i+1)*2 + 1
					str1 = strings.Repeat(" ", foreSpace) + str1
					str2 = strings.Repeat(" ", foreSpace) + str2
				} else {
					str1 += "/"
				}

			} else {
				str1 += "\\"
			}
			str1 += strings.Repeat(" ", len(mtLayer[i][j])-1)
			str2 += strings.Repeat(" ", i+1)
			str2 += mtLayer[i][j]
		}
		fmt.Println(str1)
		fmt.Println(str2)
	}

}

func routePaint(route []int) {
	probe := len(route)
	fmt.Println(strings.Repeat(" ", probe) + "*")
	for i := 0; i < len(route); i++ {
		var str1, str2 string
		str1 += strings.Repeat(" ", probe)
		if route[i] == 0 {
			str1 += "/"
			probe -= 1
		} else {
			str1 += "\\"
			probe += 1
		}
		str2 += strings.Repeat(" ", probe) + "*"
		fmt.Println(str1)
		fmt.Println(str2)
	}
}
```
在终端中（VS code使用Ctrl+`快捷命令打开）从项目主目录进入test目录，然后输入go test -v，go会自动运行test文件夹下所有的测试代码

运行测试后，如果所有测试用例都通过，终端最后会输出“PASS: TestSPV”的语句，否则会出现错误提示。这里以测试用例3来强调一下MT的构造。测试用例3对应的测试结果输出如下：

```bash
/*
TestSPV 开始...
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
当前测试 No:  0
默克尔树就像:
  00
 / \
 0 0
------------------------------------------------------------------
查找交易: 0
交易ID: 8bb1fcfcdcaa57352e751a2209dfb868849b3fce82411cbdd5003e25d87f1f5c
已找到验证路线:  [0]
路线就像:
 *
 /
*
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 1
交易ID: 3b08c1962c5d875e7b6138e6fbcae11e1e046d48a7abb7b3fc8c04b82b77984c
没有找到引用的交易
SPV 结果:  false , 想要的结果:  false
------------------------------------------------------------------
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
当前测试 No:  1
默克尔树就像:
    01234567
   /      \
   0123   4567
  /   \   /   \
  01  23  45  67
 / \ / \ / \ / \
 0 1 2 3 4 5 6 7
------------------------------------------------------------------
查找交易: 3
交易ID: df3b67915d3db06422960fc673d2b606daa2b97ef022d608eaa06e59716c691b
已找到验证路线:  [0 1 1]
路线就像:
   *
   /
  *
  \
   *
   \
    *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 7
交易ID: 0e00f37b86a9895a1ca96e1f8c755923872c45a525b797ef7f14a056a4ce5849
已找到验证路线:  [1 1 1]
路线就像:
   *
   \
    *
    \
     *
     \
      *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 5
交易ID: d8e1e04504270b6e98420f481833fcb646cbe43bcba6f6de930701dbe2641c2a
已找到验证路线:  [1 0 1]
路线就像:
   *
   \
    *
    /
   *
   \
    *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
当前测试 No:  2
默克尔树就像:
   0123
  /   \
  01  23
 / \ / \
 0 1 2 3
------------------------------------------------------------------
查找交易: 0
交易ID: 8bb1fcfcdcaa57352e751a2209dfb868849b3fce82411cbdd5003e25d87f1f5c
已找到验证路线:  [0 0]
路线就像:
  *
  /
 *
 /
*
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 1
交易ID: 3b08c1962c5d875e7b6138e6fbcae11e1e046d48a7abb7b3fc8c04b82b77984c
已找到验证路线:  [0 1]
路线就像:
  *
  /
 *
 \
  *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 5
交易ID: d8e1e04504270b6e98420f481833fcb646cbe43bcba6f6de930701dbe2641c2a
没有找到引用的交易
SPV 结果:  false , 想要的结果:  false
------------------------------------------------------------------
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
当前测试 No:  3
默克尔树就像:
    770356
   /    \
   77   0356
       /   \
       03  56  ->
      / \ / \ / \
      0 3 5 6 7 7
------------------------------------------------------------------
查找交易: 0
交易ID: 8bb1fcfcdcaa57352e751a2209dfb868849b3fce82411cbdd5003e25d87f1f5c
已找到验证路线:  [1 0 0]
路线就像:
   *
   \
    *
    /
   *
   /
  *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 1
交易ID: 3b08c1962c5d875e7b6138e6fbcae11e1e046d48a7abb7b3fc8c04b82b77984c
没有找到引用的交易
SPV 结果:  false , 想要的结果:  false
------------------------------------------------------------------
查找交易: 6
交易ID: b96a485521c99693b3ccbe58a8f6d41c68e955992dcadb20f81d4d9d85a5a034
已找到验证路线:  [1 1 1]
路线就像:
   *
   \
    *
    \
     *
     \
      *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 7
交易ID: 0e00f37b86a9895a1ca96e1f8c755923872c45a525b797ef7f14a056a4ce5849
已找到验证路线:  [0 0]
路线就像:
  *
  /
 *
 /
*
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
当前测试 No:  4
默克尔树就像:
    01245677
   /      \
   0124   5677
  /   \   /   \
  01  24  56  77
 / \ / \ / \ / \
 0 1 2 4 5 6 7 7
------------------------------------------------------------------
查找交易: 0
交易ID: 8bb1fcfcdcaa57352e751a2209dfb868849b3fce82411cbdd5003e25d87f1f5c
已找到验证路线:  [0 0 0]
路线就像:
   *
   /
  *
  /
 *
 /
*
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 1
交易ID: 3b08c1962c5d875e7b6138e6fbcae11e1e046d48a7abb7b3fc8c04b82b77984c
已找到验证路线:  [0 0 1]
路线就像:
   *
   /
  *
  /
 *
 \
  *
SPV 结果:  true , 想要的结果:  true
------------------------------------------------------------------
查找交易: 3
交易ID: df3b67915d3db06422960fc673d2b606daa2b97ef022d608eaa06e59716c691b
没有找到引用的交易
SPV 结果:  false , 想要的结果:  false
------------------------------------------------------------------

*/

```
可以看到交易7的验证路径长度只有2，而且在第三层最后一个的77被移动到了第二层的第一个枝节点，与本文前述的MT构造方法相符。

## 总结
本文详细讲解了Merkle Tree与SPV在区块链中的作用与实现细节，尽管本文只作为本教程的一个补充章节，但是在帮助读者理解区块链及区块链系统方面的意义是重大的，也是本教程最大的初衷。本文之后还会出一个小的附加章节来讲解如何建立本地UTXO数据库优化本地交易信息与区块的查找效率。

至于本教程的第二部分何时开启还得从长计议qvq。

