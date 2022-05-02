package main

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

func main() {

	book := Books{"Go 语言", "www.runoob.com", "Go 语言教程", 6495407}
	print(&book)

	/*
		var Books book
		book.title = "Go 语言"
		book.author = "www.runoob.com"
		book.subject = "Go 语言教程"
		book.book_id = 6495407
	*/
}

func print(book *Books) {
	println("title:", book.title)
	println("author:", book.author)
	println("subject", book.subject)
	println("book_id", book.book_id)
}
