package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
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

func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}
