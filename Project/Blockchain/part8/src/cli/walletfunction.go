package cli

import (
	"fmt"
	"part8/src/utils"
	"part8/src/wallet"
)

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

//walletsupdate
//扫描更新本机上存放的钱包文件
func (cli *CommandLine) walletsUpdate() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("成功更新钱包。")
}

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

//通过别名（创始人的）创建区块链
func (cli *CommandLine) createBlockChainRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.createblockchain(address)
}

//用别名调用balance命令查询余额
func (cli *CommandLine) balanceRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.balance(address)
}

