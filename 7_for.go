package main

import "fmt"

/*
Go 语言的 For 循环有 3 种形式，只有其中的一种使用分号。

和 C 语言的 for 一样：
for init; condition; post { }

和 C 的 while 一样：
for condition { }

和 C 的 for(;;) 一样：
for { }
*/

/*
for语句执行过程如下：

    1、先对表达式 1 赋初值；

    2、判别赋值表达式 init 是否满足给定条件，若其值为真，满足循环条件，则执行循环体内语句，然后执行 post，进入第二次循环，再判别 condition；否则判断 condition 的值为假，不满足条件，就终止for循环，执行循环体外语句。

*/
func main() {
	/*
		sum := 0
		for i := 0; i <= 10; i++ {
			sum += i
			fmt.Println(sum)
		}
	*/

	/*
		无限循环
	*/
	/*
		package main

		import "fmt"

		func main() {
		        sum := 0
		        for {
		            sum++ // 无限循环下去
		        }
		        fmt.Println(sum) // 无法输出
		}
	*/
	/*
		For-each range 循环

		这种格式的循环可以对字符串、数组、切片等进行迭代输出元素。
	*/

	/*strings := []string{"hello", "world"}
	for i, s := range strings {
		fmt.Println(i, s)
	}*/

	strings := []string{"hello", "world"}
	for i := range strings {
		fmt.Println(strings[i])
	}

	/*
		numbers := [6]int{1, 2, 3, 5}
		for i, x := range numbers {
			fmt.Printf("第 %d 位 x 的值 = %d\n", i, x)
		}

	*/

	/*
		以下为 Go 语言嵌套循环的格式：
		for [condition |  ( init; condition; increment ) | Range]
		{
		   for [condition |  ( init; condition; increment ) | Range]
		   {
		      statement(s);
		   }
		   statement(s);
		}
	*/

	/* 定义局部变量 */
	var i, j int

	for i = 2; i < 100; i++ {
		for j = 2; j <= (i / j); j++ {
			if i%j == 0 {
				break // 如果发现因子，则不是素数
			}
		}
		if j > (i / j) {
			fmt.Printf("%d  是素数\n", i)
		}
	}
}
