### 接口(interfaces)

一般而言，接口定义了一组方法的集合，接口不能被实例化，一个类型可以实现多个接口。

举一个简单的例子，定义一个接口 `Person`和对应的方法 `getName()` 和 `getAge()`：

```
type Person interface {
	getName() string
}

type Student struct {
	name string
	age  int
}

func (stu *Student) getName() string {
	return stu.name
}

type Worker struct {
	name   string
	gender string
}

func (w *Worker) getName() string {
	return w.name
}

func main() {
	var p Person = &Student{
		name: "Tom",
		age:  18,
	}

	fmt.Println(p.getName()) // Tom
}
```

- Go 语言中，并不需要显式地声明实现了哪一个接口，只需要直接实现该接口对应的方法即可。
- 实例化 `Student`后，强制类型转换为接口类型 Person。

在上面的例子中，我们在 main 函数中尝试将 Student 实例类型转换为 Person，如果 Student 没有完全实现 Person 的方法，比如我们将 `(*Student).getName()` 删掉，编译时会出现如下报错信息。

```
*Student does not implement Person (missing getName method)
```

但是删除 `(*Worker).getName()` 程序并不会报错，因为我们并没有在 main 函数中使用。这种情况下我们如何确保某个类型实现了某个接口的所有方法呢？一般可以使用下面的方法进行检测，如果实现不完整，编译期将会报错。

```
var _ Person = (*Student)(nil)
var _ Person = (*Worker)(nil)
```

- 将空值 nil 转换为 *Student 类型，再转换为 Person 接口，如果转换失败，说明 Student 并没有实现 Person 接口的所有方法。
- Worker 同上。

实例可以强制类型转换为接口，接口也可以强制类型转换为实例。

```
func main() {
	var p Person = &Student{
		name: "Tom",
		age:  18,
	}

	stu := p.(*Student) // 接口转为实例
	fmt.Println(stu.getAge())
}
```

### 6.3 空接口

如果定义了一个没有任何方法的空接口，那么这个接口可以表示任意类型。例如

```
func main() {
	m := make(map[string]interface{})
	m["name"] = "Tom"
	m["age"] = 18
	m["scores"] = [3]int{98, 99, 85}
	fmt.Println(m) // map[age:18 name:Tom scores:[98 99 85]]
}
```