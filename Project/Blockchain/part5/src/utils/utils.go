package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/ripemd160"
	"log"
	"os"
	"part5/src/constcode"

	"github.com/mr-tron/base58"
)

//ToHexInt将int64转换为字节串类型
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	Handle(err)
	return buff.Bytes()
}

//处理错误
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

//判断文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

//将公钥转为公钥哈希
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

//检查位生成
func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:constcode.ChecksumLength]
}

//Base256转Base58
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

//Base58转Base256
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

//公钥哈希生成钱包地址
func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{constcode.NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

//钱包地址转公钥哈希
func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constcode.ChecksumLength]
	return pubKeyHash
}