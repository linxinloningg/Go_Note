package cli

import (
	"bytes"
	"fmt"
	"part4/src/blockchain"
	"strconv"
)

type CommandLine struct{}

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

//createblockchain
/*
调用在blockchain.go中编写的InitBlockChain函数即可实现
 */
func (cli *CommandLine) createblockchain(address string) {
	newChain := blockchain.InitBlockChain([]byte(address))
	/*
	注意在使用完数据库后，需要使用newChain.Database.Close()函数关闭数据库
	 */
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

	balance, _ := chain.FindUTXOs([]byte(address))
	fmt.Printf("地址:%s, 余额:%d \n", address, balance)
}

//blockchaininfo
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

//send
/*
调用CreateTransaction函数，并将创建的交易信息保存到交易信息池中
 */
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

//mine
func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("完成挖矿")
}