package cli

import (
	"flag"
	"os"
	"part7/src/utils"
	"runtime"
)

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}


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