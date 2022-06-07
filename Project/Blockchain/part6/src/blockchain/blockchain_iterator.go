package blockchain

import (
	"github.com/dgraph-io/badger"
	"part6/src/utils"
)

//基于区块的迭代器
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//创建迭代器的初始化函数
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}
	return &iterator
}

//迭代函数:让每次迭代返回一个block，然后迭代器指向前一个区块的哈希值
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

//判断迭代器是否终止
/*
通过比较迭代器的CurrentHash与数据库存储的OgPrevHash是否相等就能够判断迭代器是否已经迭代到创始区块
*/
func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})

		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	return ogprevhash
}
