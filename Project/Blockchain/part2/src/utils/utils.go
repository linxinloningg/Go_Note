package utils

import (
	"bytes"
	"encoding/binary"
	"log"
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
