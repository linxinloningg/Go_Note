package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"io/ioutil"
	"part5/src/constcode"
	"part5/src/utils"

	"errors"
)

//钱包结构体
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

//钱包生成
func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

//钱包地址生成函数
func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

//保存钱包
func (w *Wallet) Save() {
	filename := constcode.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

//加载钱包
func LoadWallet(address string) *Wallet {
	filename := constcode.Wallets + address + ".wlt"
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("没有这个地址的钱包"))
	}
	var w Wallet
	gob.Register(elliptic.P256())
	fileContent, err := ioutil.ReadFile(filename)
	utils.Handle(err)
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}

//椭圆曲线生成密钥对
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}
