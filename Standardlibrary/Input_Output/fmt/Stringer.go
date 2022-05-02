package main

import (
	"bytes"
	"fmt"
	"strconv"
)

/*
实现了Stringer接口的类型（即有String方法），定义了该类型值的原始显示。
当采用任何接受字符的verb（%v %s %q %x %X）动作格式化一个操作数时，或者被不使用格式字符串如Print函数打印操作数时，
会调用String方法来生成输出的文本。
*/
type Person struct {
	Name string
	Age  int
	Sex  int
}

//为Person增加String方法
func (this *Person) String() string {
	buffer := bytes.NewBufferString("This is ")
	buffer.WriteString(this.Name + ", ")
	if this.Sex == 0 {
		buffer.WriteString("He ")
	} else {
		buffer.WriteString("She ")
	}

	buffer.WriteString("is ")
	buffer.WriteString(strconv.Itoa(this.Age))
	buffer.WriteString(" years old.")
	return buffer.String()
}
func main() {

	// people := Person{"polaris", 28, 0}

	test := &Person{"polaris", 28, 0}
	fmt.Println(test)
	/*
		&{polaris 28 0}
	*/

}
