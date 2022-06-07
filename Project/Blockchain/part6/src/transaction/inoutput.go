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
