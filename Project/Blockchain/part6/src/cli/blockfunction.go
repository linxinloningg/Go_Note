package cli

import (
	"bytes"
	"fmt"
	"part6/src/blockchain"
	"part6/src/utils"
	"part6/src/wallet"
	"strconv"
)

//createblockchain
/*
调用在blockchain.go中编写的InitBlockChain函数即可实现
*/
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

//mine
func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("完成挖矿")
}
