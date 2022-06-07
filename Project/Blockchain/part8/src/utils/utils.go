package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
	"os"
	"part8/src/constcode"

	"github.com/mr-tron/base58"
)

//ToHexInt将int64转换为字节串类型
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	Handle(err)
	return buff.Bytes()
}

//ToInt64将节串类型转换为int64类型
func ToInt64(num []byte) int64 {
	var num64 int64
	buff := bytes.NewBuffer(num)
	err := binary.Read(buff, binary.BigEndian, &num64)
	Handle(err)
	return num64
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

//签名函数
/*
privKey是私钥，msg则是待签名的信息。
注意到r是一个随机的大数，s则是基于r、私钥与待签名信息经过一系列计算得到的另一个大数，由(r,s)组成的元组就是签名信息
*/
func Sign(msg []byte, privKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, msg)
	Handle(err)
	signature := append(r.Bytes(), s.Bytes()...) //(r,s)组成的元组
	return signature
}

//验证函数
/*
公钥，待签名的信息，签名信息，来验证签名的真实性与权威性
 */
func Verify(msg []byte, pubkey []byte, signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	return ecdsa.Verify(&rawPubKey, msg, &r, &s)
}
