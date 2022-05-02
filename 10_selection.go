package main

import "fmt"

/*
定义切片
你可以声明一个未指定大小的数组来定义切片：

var identifier []type
切片不需要说明长度。

或使用 make() 函数来创建切片:

var slice1 []type = make([]type, len)

也可以简写为

slice1 := make([]type, len)
也可以指定容量，其中 capacity 为可选参数。

make([]T, length, capacity)
这里 len 是数组的长度并且也是切片的初始长度。
*/
func main() {
	// 创建栈
	stack := make([]int, 0)
	// push压入
	//stack = append(stack, 10)

	/*
		// pop弹出
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// 检查栈空
		len(stack) == 0
	*/

	/*for i := 0; i < 10; i++ {
		stack = append(stack, i)
	}
	for i := range stack {
		print(stack[i])
	}*/

	/*
			len() 和 cap() 函数
			切片是可索引的，并且可以由 len() 方法获取长度。
		stack
			切片提供了计算容量的方法 cap() 可以测量切片最长可以达到多少。
	*/
	fmt.Printf("len=%d cap=%d slice=%v\n", len(stack), cap(stack), stack)

	/*
		append() 和 copy() 函数
		如果想增加切片的容量，我们必须创建一个新的更大的切片并把原分片的内容都拷贝过来。

		下面的代码描述了从拷贝切片的 copy 方法和向切片追加新元素的 append 方法。
	*/

	/* 允许追加空切片 */
	stack = append(stack, 0)
	fmt.Printf("len=%d cap=%d slice=%v\n", len(stack), cap(stack), stack)

	/* 向切片添加一个元素 */
	stack = append(stack, 1)
	fmt.Printf("len=%d cap=%d slice=%v\n", len(stack), cap(stack), stack)

	/* 同时添加多个元素 */
	stack = append(stack, 2, 3, 4)
	fmt.Printf("len=%d cap=%d slice=%v\n", len(stack), cap(stack), stack)

	///* 创建切片 stack_copy 是之前切片的两倍容量*/
	//stack_copy := make([]int, len(stack), (cap(stack))*2)
	//
	///* 拷贝 stack 的内容到 stack_copy */
	//copy(stack_copy, stack)
	//fmt.Printf("len=%d cap=%d slice=%v\n", len(stack), cap(stack), stack)
	//
	stack_copy := make([]int, len(stack)) /* 拷贝 stack 的内容到 stack_copy */
	copy(stack_copy, stack)
	fmt.Printf("len=%d cap=%d slice=%v\n", len(stack_copy), cap(stack_copy), stack_copy)

}
