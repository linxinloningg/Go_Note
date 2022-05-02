### go语言数据类型及定义

* **Go基本类型：**

```go
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
```

* **字符串**

```go
var str = "hello world"
print(str)
```

* **数组**

```go
/*
声明数组
Go 语言数组声明需要指定元素类型及元素个数，语法格式如下：

var variable_name [SIZE] variable_type
*/
//初始化数组
	//以下演示了数组初始化：

	var _1 = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	//我们也可以通过字面量在声明数组的同时快速初始化数组：
	print(_1)

	_2 := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	//如果数组长度不确定，可以使用 ... 代替数组的长度，编译器会根据元素个数自行推断数组的长度：
	print(_2)

	var _3 = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	print(_3)
	//或
	_4 := [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	print(_4)
	//如果设置了数组的长度，我们还可以通过指定下标来初始化元素：

	//  将索引为 1 和 3 的元素初始化
	_5 := [5]float32{1: 2.0, 3: 7.0}
	print(_5)
```

* **切片**

```go
/*
   定义切片
   你可以声明一个未指定大小的数组来定义切片：

   var identifier []type
   切片不需要说明长度。

   或使用 make() 函数来创建切片:

   var slice []type = make([]type, len)

   也可以简写为

   slice := make([]type, len)
   也可以指定容量，其中 capacity 为可选参数。

   make([]T, length, capacity)
   这里 len 是数组的长度并且也是切片的初始长度。
*/

var _1 []int = make([]int, 0)
print(_1)
_2 := make([]int, 0)
print(_2)
```

* **结构体** 

```go
/*
定义结构体
结构体定义需要使用 type 和 struct 语句。struct 语句定义一个新的数据类型，结构体中有一个或多个成员。type 语句设定了结构体的名称。结构体的格式如下：

type struct_variable_type struct {
   member definition
   member definition
   ...
   member definition
}
*/

/*
一旦定义了结构体类型，它就能用于变量的声明，语法格式如下：

variable_name := structure_variable_type {value1, value2...valuen}
或
variable_name := structure_variable_type { key1: value1, key2: value2..., keyn: valuen}
*/
type Books struct {
	title   string
	author  string
	subject string
	book_id int
}
book := Books{"Go 语言", "www.runoob.com", "Go 语言教程", 6495407}
	print(&book)	
```

* **指针**

```go
/*
指针声明格式如下：

var var_name *var-type
*/
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
```

* **集合**

```go
//定义 Map
//可以使用内建函数 make 也可以使用 map 关键字来定义 Map:

/* 声明变量，默认 map 是 nil */
//var map_variable map[key_data_type]value_data_type

/* 使用 make 函数 */
//map_variable := make(map[key_data_type]value_data_type)

/*
   //创建集合
   var countryCapitalMap map[string]string

   countryCapitalMap = make(map[string]string)*/

countryCapitalMap := make(map[string]string)
/* map插入key - value对,各个国家对应的首都 */
countryCapitalMap["France"] = "巴黎"
countryCapitalMap["Italy"] = "罗马"
countryCapitalMap["Japan"] = "东京"
countryCapitalMap["India "] = "新德里"
```

* **常量定义**

```go
常量是一个简单值的标识符，在程序运行时，不会被修改的量。

常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型。

常量的定义格式：

const identifier [type] = value
你可以省略类型说明符 [type]，因为编译器可以根据变量的值来推断其类型。

显式类型定义： const b string = "abc"
隐式类型定义： const b = "abc"
多个相同类型的声明可以简写为：

const c_name1, c_name2 = value1, value2
```

