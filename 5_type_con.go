package main

/*
变量的类型转换
- Go中不存在隐式转换，所有类型转换必须显式声明
- 转换只能发生在两种相互兼容的类型之间
- 类型转换的格式：
<ValueA> [:]= <TypeOfValueA>(<ValueB>)
*/

func main() {
	var _1 int
	_2 := int16(_1)
	print(_2)
}
