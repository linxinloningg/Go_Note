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
