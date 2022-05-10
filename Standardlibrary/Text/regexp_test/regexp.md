工具网址：https://c.runoob.com/front-end/854/

regexp包实现了正则表达式搜索。

正则表达式采用RE2语法（除了\c、\C），和Perl、Python等语言的正则基本一致。

参见http://code.google.com/p/re2/wiki/Syntax。

#### Syntax

本包采用的正则表达式语法，默认采用perl标志。某些语法可以通过切换解析时的标志来关闭。

单字符：

```
        .              任意字符（标志s==true时还包括换行符）
        [xyz]          字符族
        [^xyz]         反向字符族
        \d             Perl预定义字符族
        \D             反向Perl预定义字符族
        [:alpha:]      ASCII字符族
        [:^alpha:]     反向ASCII字符族
        \pN            Unicode字符族（单字符名），参见unicode包
        \PN            反向Unicode字符族（单字符名）
        \p{Greek}      Unicode字符族（完整字符名）
        \P{Greek}      反向Unicode字符族（完整字符名）
```

结合：

```
        xy             匹配x后接着匹配y
        x|y            匹配x或y（优先匹配x）
```

重复：

```
        x*             重复>=0次匹配x，越多越好（优先重复匹配x）
        x+             重复>=1次匹配x，越多越好（优先重复匹配x）
        x?             0或1次匹配x，优先1次
        x{n,m}         n到m次匹配x，越多越好（优先重复匹配x）
        x{n,}          重复>=n次匹配x，越多越好（优先重复匹配x）
        x{n}           重复n次匹配x
        x*?            重复>=0次匹配x，越少越好（优先跳出重复）
        x+?            重复>=1次匹配x，越少越好（优先跳出重复）
        x??            0或1次匹配x，优先0次
        x{n,m}?        n到m次匹配x，越少越好（优先跳出重复）
        x{n,}?         重复>=n次匹配x，越少越好（优先跳出重复）
        x{n}?          重复n次匹配x
```

实现的限制：计数格式x{n}等（不包括x*等格式）中n最大值1000。负数或者显式出现的过大的值会导致解析错误，返回ErrInvalidRepeatSize。

分组：

```
        (re)           编号的捕获分组
        (?P<name>re)   命名并编号的捕获分组
        (?:re)         不捕获的分组
        (?flags)       设置当前所在分组的标志，不捕获也不匹配
        (?flags:re)    设置re段的标志，不捕获的分组
```

标志的语法为xyz（设置）、-xyz（清楚）、xy-z（设置xy，清楚z），标志如下：

```
        I              大小写敏感（默认关闭）
        m              ^和$在匹配文本开始和结尾之外，还可以匹配行首和行尾（默认开启）
        s              让.可以匹配\n（默认关闭）
        U              非贪婪的：交换x*和x*?、x+和x+?……的含义（默认关闭）
```

边界匹配：

```
        ^              匹配文本开始，标志m为真时，还匹配行首
        $              匹配文本结尾，标志m为真时，还匹配行尾
        \A             匹配文本开始
        \b             单词边界（一边字符属于\w，另一边为文首、文尾、行首、行尾或属于\W）
        \B             非单词边界
        \z             匹配文本结尾
```

转义序列：

```
        \a             响铃符（\007）
        \f             换纸符（\014）
        \t             水平制表符（\011）
        \n             换行符（\012）
        \r             回车符（\015）
        \v             垂直制表符（\013）
        \123           八进制表示的字符码（最多三个数字）
        \x7F           十六进制表示的字符码（必须两个数字）
        \x{10FFFF}     十六进制表示的字符码
        \*             字面值'*'
        \Q...\E        反斜线后面的字符的字面值
```

字符族（预定义字符族之外，方括号内部）的语法：

```
        x              单个字符
        A-Z            字符范围（方括号内部才可以用）
        \d             Perl字符族
        [:foo:]        ASCII字符族
        \pF            单字符名的Unicode字符族
        \p{Foo}        完整字符名的Unicode字符族
```

预定义字符族作为字符族的元素：

```
        [\d]           == \d
        [^\d]          == \D
        [\D]           == \D
        [^\D]          == \d
        [[:name:]]     == [:name:]
        [^[:name:]]    == [:^name:]
        [\p{Name}]     == \p{Name}
        [^\p{Name}]    == \P{Name}
```

Perl字符族：

```
        \d             == [0-9]
        \D             == [^0-9]
        \s             == [\t\n\f\r ]
        \S             == [^\t\n\f\r ]
        \w             == [0-9A-Za-z_]
        \W             == [^0-9A-Za-z_]
```

ASCII字符族：

```
        [:alnum:]      == [0-9A-Za-z]
        [:alpha:]      == [A-Za-z]
        [:ascii:]      == [\x00-\x7F]
        [:blank:]      == [\t ]
        [:cntrl:]      == [\x00-\x1F\x7F]
        [:digit:]      == [0-9]
        [:graph:]      == [!-~] == [A-Za-z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~]
        [:lower:]      == [a-z]
        [:print:]      == [ -~] == [ [:graph:]]
        [:punct:]      == [!-/:-@[-`{-~]
        [:space:]      == [\t\n\v\f\r ]
        [:upper:]      == [A-Z]
        [:word:]       == [0-9A-Za-z_]
        [:xdigit:]     == [0-9A-Fa-f]
```

本包的正则表达式保证搜索复杂度为O(n)，其中n为输入的长度。这一点很多其他开源实现是无法保证的。参见：

```
http://swtch.com/~rsc/regexp/regexp1.html
```

### func [QuoteMeta](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#581) 

```
func QuoteMeta(s string) string
```

QuoteMeta返回将s中所有正则表达式元字符都进行转义后字符串。该字符串可以用在正则表达式中匹配字面值s。例如，QuoteMeta(`[foo]`)会返回`\[foo\]`。

### func [Match](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#433) 

```
func Match(pattern string, b []byte) (matched bool, err error)
```

Match检查b中是否存在匹配pattern的子序列。更复杂的用法请使用Compile函数和Regexp对象。

### func [MatchString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#422) 

```
func MatchString(pattern string, s string) (matched bool, err error)
```

MatchString类似Match，但匹配对象是字符串。

Example

### func [MatchReader](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#411) 

```
func MatchReader(pattern string, r io.RuneReader) (matched bool, err error)
```

MatchReader类似Match，但匹配对象是io.RuneReader。

### type [Regexp](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#82) 

```
type Regexp struct {
    // 内含隐藏或非导出字段
}
```

Regexp代表一个编译好的正则表达式。Regexp可以被多线程安全地同时使用。

#### func [Compile](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#117) 

```
func Compile(expr string) (*Regexp, error)
```

Compile解析并返回一个正则表达式。如果成功返回，该Regexp就可用于匹配文本。

在匹配文本时，该正则表达式会尽可能早的开始匹配，并且在匹配过程中选择回溯搜索到的第一个匹配结果。这种模式被称为“leftmost-first”，Perl、Python和其他实现都采用了这种模式，但本包的实现没有回溯的损耗。对POSIX的“leftmost-longest”模式，参见CompilePOSIX。

#### func [CompilePOSIX](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#140) 

```
func CompilePOSIX(expr string) (*Regexp, error)
```

类似Compile但会将语法约束到POSIX ERE（egrep）语法，并将匹配模式设置为leftmost-longest。

在匹配文本时，该正则表达式会尽可能早的开始匹配，并且在匹配过程中选择搜索到的最长的匹配结果。这种模式被称为“leftmost-longest”，POSIX采用了这种模式（早期正则的DFA自动机模式）。

然而，可能会有多个“leftmost-longest”匹配，每个都有不同的组匹配状态，本包在这里和POSIX不同。在所有可能的“leftmost-longest”匹配里，本包选择回溯搜索时第一个找到的，而POSIX会选择候选结果中第一个组匹配最长的（可能有多个），然后再从中选出第二个组匹配最长的，依次类推。POSIX规则计算困难，甚至没有良好定义。

参见http://swtch.com/~rsc/regexp/regexp2.html#posix获取细节。

#### func [MustCompile](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#218) 

```
func MustCompile(str string) *Regexp
```

MustCompile类似Compile但会在解析失败时panic，主要用于全局正则表达式变量的安全初始化。

#### func [MustCompilePOSIX](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#229) 

```
func MustCompilePOSIX(str string) *Regexp
```

MustCompilePOSIX类似CompilePOSIX但会在解析失败时panic，主要用于全局正则表达式变量的安全初始化。

#### func (*Regexp) [String](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#103) 

```
func (re *Regexp) String() string
```

String返回用于编译成正则表达式的字符串。

#### func (*Regexp) [LiteralPrefix](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#388) 

```
func (re *Regexp) LiteralPrefix() (prefix string, complete bool)
```

LiteralPrefix返回一个字符串字面值prefix，任何匹配本正则表达式的字符串都会以prefix起始。 如果该字符串字面值包含整个正则表达式，返回值complete会设为真。

#### func (*Regexp) [NumSubexp](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#245) 

```
func (re *Regexp) NumSubexp() int
```

NumSubexp返回该正则表达式中捕获分组的数量。

#### func (*Regexp) [SubexpNames](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#254) 

```
func (re *Regexp) SubexpNames() []string
```

SubexpNames返回该正则表达式中捕获分组的名字。第一个分组的名字是names[1]，因此，如果m是一个组匹配切片，m[i]的名字是SubexpNames()[i]。因为整个正则表达式是无法被命名的，names[0]必然是空字符串。该切片不应被修改。

Example

#### func (*Regexp) [Longest](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#148) 

```
func (re *Regexp) Longest()
```

Longest让正则表达式在之后的搜索中都采用"leftmost-longest"模式。在匹配文本时，该正则表达式会尽可能早的开始匹配，并且在匹配过程中选择搜索到的最长的匹配结果。

#### func (*Regexp) [Match](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#404) 

```
func (re *Regexp) Match(b []byte) bool
```

Match检查b中是否存在匹配pattern的子序列。

#### func (*Regexp) [MatchString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#399) 

```
func (re *Regexp) MatchString(s string) bool
```

MatchString类似Match，但匹配对象是字符串。

#### func (*Regexp) [MatchReader](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#394) 

```
func (re *Regexp) MatchReader(r io.RuneReader) bool
```

MatchReader类似Match，但匹配对象是io.RuneReader。

#### func (*Regexp) [Find](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#663) 

```
func (re *Regexp) Find(b []byte) []byte
```

Find返回保管正则表达式re在b中的最左侧的一个匹配结果的[]byte切片。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#688) 

```
func (re *Regexp) FindString(s string) string
```

Find返回保管正则表达式re在b中的最左侧的一个匹配结果的字符串。如果没有匹配到，会返回""；但如果正则表达式成功匹配了一个空字符串，也会返回""。如果需要区分这种情况，请使用FindStringIndex 或FindStringSubmatch。

Example

#### func (*Regexp) [FindIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#675) 

```
func (re *Regexp) FindIndex(b []byte) (loc []int)
```

Find返回保管正则表达式re在b中的最左侧的一个匹配结果的起止位置的切片（显然len(loc)==2）。匹配结果可以通过起止位置对b做切片操作得到：b[loc[0]:loc[1]]。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindStringIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#700) 

```
func (re *Regexp) FindStringIndex(s string) (loc []int)
```

Find返回保管正则表达式re在b中的最左侧的一个匹配结果的起止位置的切片（显然len(loc)==2）。匹配结果可以通过起止位置对b做切片操作得到：b[loc[0]:loc[1]]。如果没有匹配到，会返回nil。

Example

#### func (*Regexp) [FindReaderIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#713) 

```
func (re *Regexp) FindReaderIndex(r io.RuneReader) (loc []int)
```

Find返回保管正则表达式re在b中的最左侧的一个匹配结果的起止位置的切片（显然len(loc)==2）。匹配结果可以在输入流r的字节偏移量loc[0]到loc[1]-1（包括二者）位置找到。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindSubmatch](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#726) 

```
func (re *Regexp) FindSubmatch(b []byte) [][]byte
```

Find返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的[][]byte切片。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindStringSubmatch](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#882) 

```
func (re *Regexp) FindStringSubmatch(s string) []string
```

Find返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的[]string切片。如果没有匹配到，会返回nil。

Example

#### func (*Regexp) [FindSubmatchIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#873) 

```
func (re *Regexp) FindSubmatchIndex(b []byte) []int
```

Find返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的起止位置的切片。匹配结果和分组匹配结果可以通过起止位置对b做切片操作得到：b[loc[2*n]:loc[2*n+1]]。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindStringSubmatchIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#901) 

```
func (re *Regexp) FindStringSubmatchIndex(s string) []int
```

Find返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的起止位置的切片。匹配结果和分组匹配结果可以通过起止位置对b做切片操作得到：b[loc[2*n]:loc[2*n+1]]。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindReaderSubmatchIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#910) 

```
func (re *Regexp) FindReaderSubmatchIndex(r io.RuneReader) []int
```

Find返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的起止位置的切片。匹配结果和分组匹配结果可以在输入流r的字节偏移量loc[0]到loc[1]-1（包括二者）位置找到。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindAll](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#920) 

```
func (re *Regexp) FindAll(b []byte, n int) [][]byte
```

Find返回保管正则表达式re在b中的所有不重叠的匹配结果的[][]byte切片。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindAllString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#956) 

```
func (re *Regexp) FindAllString(s string, n int) []string
```

Find返回保管正则表达式re在b中的所有不重叠的匹配结果的[]string切片。如果没有匹配到，会返回nil。

Example

#### func (*Regexp) [FindAllIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#938) 

```
func (re *Regexp) FindAllIndex(b []byte, n int) [][]int
```

Find返回保管正则表达式re在b中的所有不重叠的匹配结果的起止位置的切片。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindAllStringIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#974) 

```
func (re *Regexp) FindAllStringIndex(s string, n int) [][]int
```

Find返回保管正则表达式re在b中的所有不重叠的匹配结果的起止位置的切片。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindAllSubmatch](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#992) 

```
func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte
```

Find返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的[][][]byte切片。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindAllStringSubmatch](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#1034) 

```
func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string
```

Find返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的[][]string切片。如果没有匹配到，会返回nil。

Example

#### func (*Regexp) [FindAllSubmatchIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#1016) 

```
func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int
```

Find返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的起止位置的切片（第一层表示第几个匹配结果，完整匹配和分组匹配的起止位置对在第二层）。如果没有匹配到，会返回nil。

#### func (*Regexp) [FindAllStringSubmatchIndex](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#1059) 

```
func (re *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int
```

Find返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的起止位置的切片（第一层表示第几个匹配结果，完整匹配和分组匹配的起止位置对在第二层）。如果没有匹配到，会返回nil。

Example

#### func (*Regexp) [Split](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#1088) 

```
func (re *Regexp) Split(s string, n int) []string
```

Split将re在s中匹配到的结果作为分隔符将s分割成多个字符串，并返回这些正则匹配结果之间的字符串的切片。

返回的切片不会包含正则匹配的结果，只包含匹配结果之间的片段。当正则表达式re中不含正则元字符时，本方法等价于strings.SplitN。

举例：

```
s := regexp.MustCompile("a*").Split("abaabaccadaaae", 5)
// s: ["", "b", "b", "c", "cadaaae"]
```

参数n绝对返回的子字符串的数量：

```
n > 0 : 返回最多n个子字符串，最后一个子字符串是剩余未进行分割的部分。
n == 0: 返回nil (zero substrings)
n < 0 : 返回所有子字符串
```

#### func (*Regexp) [Expand](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#757) 

```
func (re *Regexp) Expand(dst []byte, template []byte, src []byte, match []int) []byte
```

Expand返回新生成的将template添加到dst后面的切片。在添加时，Expand会将template中的变量替换为从src匹配的结果。match应该是被FindSubmatchIndex返回的匹配结果起止位置索引。（通常就是匹配src，除非你要将匹配得到的位置用于另一个[]byte）

在template参数里，一个变量表示为格式如：$name或${name}的字符串，其中name是长度>0的字母、数字和下划线的序列。一个单纯的数字字符名如$1会作为捕获分组的数字索引；其他的名字对应(?P<name>...)语法产生的命名捕获分组的名字。超出范围的数字索引、索引对应的分组未匹配到文本、正则表达式中未出现的分组名，都会被替换为空切片。

$name格式的变量名，name会尽可能取最长序列：$1x等价于${1x}而非${1}x，$10等价于${10}而非${1}0。因此$name适用在后跟空格/换行等字符的情况，${name}适用所有情况。

如果要在输出中插入一个字面值'$'，在template里可以使用$$。

#### func (*Regexp) [ExpandString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#764) 

```
func (re *Regexp) ExpandString(dst []byte, template string, src string, match []int) []byte
```

ExpandString类似Expand，但template和src参数为字符串。它将替换结果添加到切片并返回切片，以便让调用代码控制内存申请。

#### func (*Regexp) [ReplaceAllLiteral](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#556) 

```
func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte
```

ReplaceAllLiteral返回src的一个拷贝，将src中所有re的匹配结果都替换为repl。repl参数被直接使用，不会使用Expand进行扩展。

#### func (*Regexp) [ReplaceAllLiteralString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#458) 

```
func (re *Regexp) ReplaceAllLiteralString(src, repl string) string
```

ReplaceAllLiteralString返回src的一个拷贝，将src中所有re的匹配结果都替换为repl。repl参数被直接使用，不会使用Expand进行扩展。

Example

#### func (*Regexp) [ReplaceAll](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#538) 

```
func (re *Regexp) ReplaceAll(src, repl []byte) []byte
```

ReplaceAllLiteral返回src的一个拷贝，将src中所有re的匹配结果都替换为repl。在替换时，repl中的'$'符号会按照Expand方法的规则进行解释和替换，例如$1会被替换为第一个分组匹配结果。

#### func (*Regexp) [ReplaceAllString](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#444) 

```
func (re *Regexp) ReplaceAllString(src, repl string) string
```

ReplaceAllLiteral返回src的一个拷贝，将src中所有re的匹配结果都替换为repl。在替换时，repl中的'$'符号会按照Expand方法的规则进行解释和替换，例如$1会被替换为第一个分组匹配结果。

Example

#### func (*Regexp) [ReplaceAllFunc](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#566) 

```
func (re *Regexp) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte
```

ReplaceAllLiteral返回src的一个拷贝，将src中所有re的匹配结果（设为matched）都替换为repl(matched)。repl返回的切片被直接使用，不会使用Expand进行扩展。

#### func (*Regexp) [ReplaceAllStringFunc](https://github.com/golang/go/blob/master/src/regexp/regexp.go?name=release#468) 

```
func (re *Regexp) ReplaceAllStringFunc(src string, repl func(string) string) string
```

ReplaceAllLiteral返回src的一个拷贝，将src中所有re的匹配结果（设为matched）都替换为repl(matched)。repl返回的字符串被直接使用，不会使用Expand进行扩展。