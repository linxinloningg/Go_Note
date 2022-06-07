package transaction

import (
	"bytes"
	"encoding/gob"
	"part8/src/utils"
)


type UTXO struct {
	TxID   []byte
	OutIdx int
	TxOutput
}

func (u *UTXO) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(u)
	utils.Handle(err)
	return res.Bytes()
}

func DeserializeUTXO(data []byte) *UTXO {
	var utxo UTXO
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&utxo)
	utils.Handle(err)
	return &utxo
}
