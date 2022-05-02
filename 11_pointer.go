package main

/*
指针声明格式如下：

var var_name *var-type
*/
func main() {
	var INT int = 10
	var INTPOINTER *int = &INT
	println(INTPOINTER)

	/*
		Go 空指针
		当一个指针被定义后没有分配到任何变量时，它的值为 nil。

		nil 指针也称为空指针。

		nil在概念上和其它语言的null、None、nil、NULL一样，都指代零值或空值。
	*/
	var NULLPOINTER *int
	println(NULLPOINTER)

	/*
		空指针判断：
	*/

	/* ptr 不是空指针 */
	if NULLPOINTER != nil {
		print("NOT NULL")
	}
	/* ptr 是空指针 */
	if NULLPOINTER == nil {
		print("NULL")
	}
}
