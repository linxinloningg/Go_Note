package utxoset

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"part8/src/transaction"
	"part8/src/utils"
	"runtime"
)

var (
	info         = "INFO:"
	infoname     = info + "NAME"
	infoheight   = info + "HIGT"
	utxokey      = "UTXO:"
	utxokeyorder = ":ORDER:"
)

//UTXO集结构体
type UTXOSet struct {
	Name   []byte     //用于辨别UTXO的名称
	DB     *badger.DB //UTXO 数据库
	Height int64
}

func GetUtxoSetFile(dir string) string {
	fileAddress := dir + "/" + "MANIFEST"
	return fileAddress
}

func ToUtxoKey(txID []byte, order int) []byte {
	utxoKey := bytes.Join([][]byte{[]byte(utxokey), txID, []byte(utxokeyorder), utils.ToHexInt(int64(order))}, []byte{})
	return utxoKey
}

//UTXO集的创建
func CreateUTXOSet(name []byte, dir string, utxos []transaction.UTXO, height int64) *UTXOSet {
	if utils.FileExists(GetUtxoSetFile(dir)) {
		fmt.Println("UTXOSet has already existed, now rebuild it.")
		err := os.RemoveAll(dir)
		utils.Handle(err)
	}

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	utxoSet := UTXOSet{name, db, height}

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(infoname), name)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(infoheight), utils.ToHexInt(height))
		if err != nil {
			return err
		}
		for _, utxo := range utxos {
			utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
			err = txn.Set(utxoKey, utxo.Serialize())
			return err
		}
		return nil
	})
	utils.Handle(err)
	return &utxoSet

}

//UTXO集的加载
func LoadUTXOSet(dir string) *UTXOSet {
	if !utils.FileExists(GetUtxoSetFile(dir)) {
		fmt.Println("No UTXOSet found, please create one first")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)
	var name []byte
	var height int64
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(infoname))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			name = val
			return nil
		})
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte(infoheight))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			height = utils.ToInt64(val)
			return nil
		})

		return err
	})
	utils.Handle(err)

	utxoSet := UTXOSet{name, db, height}
	return &utxoSet
}

/*func (us *UTXOSet) AddUtxo(utxo *transaction.UTXO) {
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
		err := txn.Set(utxoKey, utxo.Serialize())
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}*/

//添加新的UTXO进数据库
func (us *UTXOSet) AddUTXO(txID []byte, outIdx int, output transaction.TxOutput) {
	utxo := transaction.UTXO{txID, outIdx, output}
	//us.AddUtxo(&utxo)
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(utxo.TxID, utxo.OutIdx)
		err := txn.Set(utxoKey, utxo.Serialize())
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

//在数据库删除UTXO
func (us *UTXOSet) DelUTXO(txID []byte, order int) {
	err := us.DB.Update(func(txn *badger.Txn) error {
		utxoKey := ToUtxoKey(txID, order)
		err := txn.Delete(utxoKey)
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

func (us *UTXOSet) UpdateHeight(height int64) {
	us.Height = height
	err := us.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(infoheight), utils.ToHexInt(height))
		return err
	})
	utils.Handle(err)
}

func IsInfo(inkey []byte) bool {
	if bytes.HasPrefix(inkey, []byte(info)) {
		return true
	} else {
		return false
	}
}
