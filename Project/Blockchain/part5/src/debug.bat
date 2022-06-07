rd /s /q tmp
md tmp\blocks
md tmp\wallets
md tmp\ref_list
go run main.go createwallet
go run main.go walletslist
go run main.go createwallet -refname LeoCao
go run main.go walletinfo -refname LeoCao
go run main.go createwallet -refname Krad
go run main.go createwallet -refname Exia
go run main.go createwallet
go run main.go walletslist
go run main.go createblockchain -refname LeoCao
go run main.go blockchaininfo
go run main.go balance -refname LeoCao
go run main.go sendbyrefname -from LeoCao -to Krad -amount 100
go run main.go balance -refname Krad
go run main.go mine
go run main.go blockchaininfo
go run main.go balance -refname LeoCao
go run main.go balance -refname Krad
go run main.go sendbyrefname -from LeoCao -to Exia -amount 100
go run main.go sendbyrefname -from Krad -to Exia -amount 30
go run main.go mine
go run main.go blockchaininfo
go run main.go balance -refname LeoCao
go run main.go balance -refname Krad
go run main.go balance -refname Exia
go run main.go sendbyrefname -from Exia -to LeoCao -amount 90
go run main.go sendbyrefname -from Exia -to Krad -amount 90
go run main.go mine
go run main.go blockchaininfo
go run main.go balance -refname LeoCao
go run main.go balance -refname Krad
go run main.go balance -refname Exia
