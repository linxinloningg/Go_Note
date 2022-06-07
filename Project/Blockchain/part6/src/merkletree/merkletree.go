package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"part6/src/transaction"
	"part6/src/utils"
)

//结构体MerkleTree
type MerkleTree struct {
	RootNode *MerkleNode //根节点
}

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte // Hash
}

//先判断当前创建的是否为叶子节点，然后再根据节点类型不同创建MT节点
//MT节点构造
func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	tempNode := MerkleNode{}

	//叶子节点
	if left == nil && right == nil { //The leaf
		tempNode.Data = data
	} else {
		catenateHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(catenateHash)
		tempNode.Data = hash[:]
	}

	tempNode.LeftNode = left
	tempNode.RightNode = right

	return &tempNode
}

//MT构造
func CrateMerkleTree(txs []*transaction.Transaction) *MerkleTree {
	txslen := len(txs)
	if txslen%2 != 0 {
		txs = append(txs, txs[txslen-1])
	}

	var nodePool []*MerkleNode

	for _, tx := range txs {
		nodePool = append(nodePool, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodePool) > 1 {
		var tempNodePool []*MerkleNode
		poolLen := len(nodePool)
		if poolLen%2 != 0 { //Notice here, we place the remained node at the head of the upper layer
			tempNodePool = append(tempNodePool, nodePool[poolLen-1])
		}
		for i := 0; i < poolLen/2; i++ {
			tempNodePool = append(tempNodePool, CreateMerkleNode(nodePool[2*i], nodePool[2*i+1], nil))
		}
		nodePool = tempNodePool
	}

	merkleTree := MerkleTree{nodePool[0]}

	return &merkleTree
}

//路径搜索
/*
输入以此为寻找的目标哈希值，已经走过的路径（route与hashroute，其中route为方向，0为左，1为右，hashroute保存的是哈希值）
*/
func (mn *MerkleNode) Find(data []byte, route []int, hashroute [][]byte) (bool, []int, [][]byte) {
	findFlag := false

	if bytes.Equal(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashroute
	} else {
		if mn.LeftNode != nil {
			route_t := append(route, 0)
			hashroute_t := append(hashroute, mn.RightNode.Data)
			findFlag, route_t, hashroute_t = mn.LeftNode.Find(data, route_t, hashroute_t)
			if findFlag {
				return findFlag, route_t, hashroute_t
			} else {
				if mn.RightNode != nil {
					route_t = append(route, 1)
					hashroute_t = append(hashroute, mn.LeftNode.Data)
					findFlag, route_t, hashroute_t = mn.RightNode.Find(data, route_t, hashroute_t)
					if findFlag {
						return findFlag, route_t, hashroute_t
					} else {
						return findFlag, route, hashroute
					}

				}
			}
		} else {
			return findFlag, route, hashroute
		}
	}
	return findFlag, route, hashroute
}

//输入为一交易信息的ID（也即交易信息的哈希值），返回的是验证路径与一个是否找到该交易信息的信号
func (mt *MerkleTree) BackValidationRoute(txid []byte) ([]int, [][]byte, bool) {
	ok, route, hashroute := mt.RootNode.Find(txid, []int{}, [][]byte{})
	return route, hashroute, ok
}

func SimplePaymentValidation(txid, mtroothash []byte, route []int, hashroute [][]byte) bool {
	routeLen := len(route)
	var tempHash []byte
	tempHash = txid

	for i := routeLen - 1; i >= 0; i-- {
		if route[i] == 0 {
			catenateHash := append(tempHash, hashroute[i]...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else if route[i] == 1 {
			catenateHash := append(hashroute[i], tempHash...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else {
			utils.Handle(errors.New("error in validation route"))
		}
	}
	return bytes.Equal(tempHash, mtroothash)
}
