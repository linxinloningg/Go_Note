package main

// 通过在函数体外部使用 var 关键字来进行全局变量的声明与赋值
/*
单个变量的声明与赋值
- 变量的声明格式：var <变量名称> <变量类型>
- 变量的赋值格式：<变量名称> = <表达式>
- 声明的同时赋值：var <变量名称> [变量类型] = <表达式>
*/

func main() {
	// Go基本类型：

	// 布尔型：bool - 长度：1字节
	//- 取值范围：true, false
	var bool_value bool = true
	println(bool_value)

	// 整型：int/uint - 根据运行平台可能为32或64位
	var int_value int = 1
	println(int_value)

	//8 位整型：int8/uint8
	//- 长度：1字节
	//- 取值范围：-128~127/0~255
	var int8_value int8 = 1
	println(int8_value)

	// 字节型：byte（uint8别名）
	var byte_value byte
	println(byte_value)

	// 16位整型：int16/uint16
	//- 长度：2字节
	//- 取值范围：-32768~32767/0~65535
	var int16_value int16 = 1
	println(int16_value)

	//  32位整型：int32（rune）/uint32
	//- 长度：4字节
	//- 取值范围：-2^32/2~2^32/2-1/0~2^32-1
	var int32_value int32 = 1
	println(int32_value)

	// 64位整型：int64/uint64
	//- 长度：8字节
	//- 取值范围：-2^64/2~2^64/2-1/0~2^64-1
	var int64_value int64 = 1
	println(int64_value)

	// 浮点型：float32/float64
	//- 长度：4/8字节
	//- 小数位：精确到7/15小数位
	var float32_value float32 = 1
	println(float32_value)
	var float64_value float64 = 1
	println(float64_value)

	// 复数：complex64/complex128
	//- 长度：8/16字节
	var complex64_value complex64 = 1
	println(complex64_value)
	var complex128_value complex128 = 1
	println(complex128_value)

	// 足够保存指针的 32 位或 64 位整数型：uintptr
	var uintptr_value uintptr = 1
	println(uintptr_value)

	/*
		其它值类型：
		- array、struct、string
		- 引用类型：
		- slice、map、chan
		- 接口类型：inteface
		-
	*/

}
