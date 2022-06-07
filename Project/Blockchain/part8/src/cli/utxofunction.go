package cli

import (
	"fmt"
	"part8/src/blockchain"
	"part8/src/wallet"
)

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
