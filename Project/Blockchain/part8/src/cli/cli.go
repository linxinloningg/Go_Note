package cli

import (
	"fmt"
)

type CommandLine struct{}

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
