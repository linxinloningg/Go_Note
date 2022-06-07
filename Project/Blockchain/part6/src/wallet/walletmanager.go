package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"part6/src/constcode"
	"part6/src/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//RefList来记录机器上记录的钱包。其中key值为钱包地址，value为钱包的别名
type RefList map[string]string

//保存RefList
func (r *RefList) Save() {
	filename := constcode.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

//加载已保存的RefList
func LoadRefList() *RefList {
	filename := constcode.WalletsRefList + "ref_list.data"
	var reflist RefList
	if utils.FileExists(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&reflist)
		utils.Handle(err)
	} else {
		/*
			如果没有可以可以加载的RefList文件，LoadRefList会自动创建一个新的RefList并调用Update函数扫描本机的所有.wlt文件
		*/
		reflist = make(RefList)
		reflist.Update()
	}
	return &reflist
}

//RefList应实现一更新函数，用于扫描机器上保存的所有钱包文件（特别是检查是否存在从其他机器上拷贝的钱包）
func (r *RefList) Update() {
	err := filepath.Walk(constcode.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			_, ok := (*r)[fileName[:len(fileName)-4]]
			if !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}


/*
用别名的方式指向本地钱包（***注意：实际的区块链系统中钱包是没有别名的，这里完全是方便我们后续的演示。）
 */
func (r *RefList) BindRef(address, refname string) {
	(*r)[address] = refname
}

//构建通过别名调取钱包地址
func (r *RefList) FindRef(refname string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refname {
			temp = key
			break
		}
	}
	if temp == "" {
		err := errors.New("the refname is not found")
		return temp, err
	}
	return temp, nil
}