# Go の勉強
Go を勉強するのに使えそうなサイト  
1. https://gihyo.jp/dev/feature/01/go_4beginners
2. https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/index.html
3. https://www.slideshare.net/takuyaueda967/2016-go?ref=https://kirohi.com/to_study_golang
4. https://go-tour-jp.appspot.com/list

4番目の公式の Go Tour で勉強することにする  

## Go の基本
- 実行:  
`# go run ファイル名.go`  

- フォーマットを整える:  
`# go fmt hello.go`  
Go のインデントはハードタブ 気にせず fmt すればいい  

- コンパイルしてバイナリファイルを作る  
`# go build ファイル名.go`  
対応するOS用にコンパイルしたら `./ファイル名` など コマンドが無くても実行できるようになる  

- バージョン確認  
`go version`  

Go では package と ディレクトリ は同じ概念っぽい  

>Goでは，1つのパッケージは1つのディレクトリに格納します

# Go Tour
## プログラムの構成
Goのプログラムは パッケージ(package)で構成される  
インポートパスが "math/rand" のパッケージは package rand ステートメント(宣言)で始まるファイル群で構成されている  

```go:ex01.go
func main() {
  fmt.Println("HelloWorld")
}
```

## import の書き方
import が複数になるときは 前者の書き方が推奨されている  
前者の書き方を factored インポート ステートメント という  
import 文が書いてあるのに ソースコード内で使われてなかったらエラーになる  

```go
import (
        "fmt"
        "time"
)

import "fmt"
import "time"
```

## 外部公開な名前
最初の文字が大文字で始まる名前は 外部のパッケージから参照できるエクスポート(公開)された名前(exported name)  
Println も大文字スタートになっている = 外部から参照できる  
小文字なら同じパッケージなら参照できる  

```go
// 外部から参照できる
math.Pi
// 外部から参照できない
math.pi
// つまり

// 動く
fmt.Println(training.Message)
// 動かない
fmt.Println(training.message)
```

## 関数 の定義
引数を渡すときは 変数名の後ろに型名を書く  

```go:ex02.go
func add(x int, y int) int {
  return x + y
}

func add02(x, y int) int {
  return x + y
}
```

`x int, y int`を`x, y int`に省略できる  

複数の戻り値を返すことができる  
```go:ex03.go
func swap(x, y string) (string, string) {
  return y, x
}

func main() {
  a, b := swap("hello", "world")
  fmt.Println(a, b)
}
```

関数を定義する段階で戻り値を設定することができる  
具体的には 変数に名前をつけて戻り値にする(named return value)ことができる  

```go:ex04.go
// out01, out02 が戻り値になる
func split(in int) (out01, out02 int) {
  out01 = in * 4
  out02 = in + 4
  return
}
```

戻り値の意味を示す名前とすることで 関数のドキュメントとして表現できる  
return ステートメントに何も書かずに戻すことができる("naked" return という)  
ただし naked return ステートメントは 短い関数でのみ利用すべき  

func の中などで 定義と初期化を同時に行う(わざわざ 別で func を作らないで済む)  
```go
func aaa() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	err := なんかのメソッド()
	check(err)
}
```

## 変数の宣言
var ステートメントは変数(variable)を宣言する  
`var 変数名 型`が基本 例:`var i int`  

var 宣言では 変数毎に初期化子(initializer)を与えることができる  
`var i, j int = 1, 2`  

初期化子が与えられている場合 型を省略できる(自動で型を決めてくれる)  
関数の中では var 宣言の代わりに := の代入文を使い 暗黙的な型宣言ができる  
暗黙的な型宣言は 関数の外では使用できない  
変数に初期値を与えずに宣言すると ゼロ値(zero value)が与えられる  

```go:ex05.go
func main() {
  // ゼロ値
  var a int
  fmt.Println(a) // 0
  a = 1
  fmt.Println(a) // 1

  var i, j = 1, 2
  fmt.Println(i, j) // 1, 2

  b := 3
  fmt.Println(b) // 3

  var c = 4
  fmt.Println(c) // 4

  var z complex128 = cmplx.Sqrt(-5 + 12i)
  fmt.Println(z) // 2+3i (-5+12i の平方根)
}
```

特別な理由がない限り 整数の変数は int を使うべき  
型は 以下など がある  
- bool
- string
- int
- uint(正の整数型)
- uintptr(ポインタ)
- byte
- rune(Javaのchar)
- float32
- float64
- complex64(複素数)
- complex128

int や uint の中にも `int8, int16, int32, uint8, uint16, uint32` などあるが  
`int, uint, uintptr 型` は OS の bit 数に合わせて調整してくれるっぽい  

まとめて変数を宣言することができる  

```go:ex05.go
  var (
    ToBe bool = false
    MaxInt uint64 = 1<<64 - 1
  )
```

## ヒアドキュメント
バッククオートで囲むことで 複数行に渡る string(ヒアドキュメント)を記述できる  

```go:ex06.go
  var a = `
  hello
  world
  `
```

フォーマットを使った書き方  
`Printf()` は改行をしないで フォーマットをはめ込む用の関数  
型は `%T` で 値は `%v` ではめ込む  

```go:ex06.go
func main(){
  b := false
  fmt.Printf("Type: %T Value: %v\n", b, b)
  // Type: bool Value: false
}
```

## 型変換
変数 v を T 型へ変換する `a := T(v)`  
また 右側の変数が型を持っている場合 左側の新しい変数は同じ型になる  
`b := a` で b は a の型になる  
型変換と似たような操作に 型アサーションがある

```go:ex07.go
func main() {
  var a int
  b := a
  fmt.Printf("%T\n", b)
  // int
  i := 42
  j := 3.142
  k := 0.867 + 0.5i
  fmt.Printf("%T, %T, %T\n", i, j, k)
  // int, float64, complex128
}
```

## 定数(constant)
const キーワードを使って変数と同じように宣言する  
文字(character)型, 文字列(string), boolean, 数値(numeric)のみで使える  
定数の頭文字は大文字を使う(定数なら外部に公開しても問題ないため)  
型のない定数は その状況によって必要な型を取る  
`:=`で宣言することができない  

```go:ex07.go
  const Pi = 3.14
  fmt.Println(Pi)
```

## for ステートメント

```go:ex08.go
func main() {
  var sum int
  for i := 0; i < 10; i++ {
    sum += i
  }
  fmt.Println(sum)
}
```

for 文の宣言部分の初期化と 後処理ステートメントの記述は任意  

```go:ex08.go
func main() {
  var sum02 = 1
  for ; sum02 < 10; {
    sum02++
  }
  fmt.Println(sum02)
}
```

while 文と等価にするには セミコロン(;)を省略する  

```go:ex08.go
func main() {
  var sum03 = 1
  for sum03 < 10 {
    sum03++
  }
  fmt.Println(sum03)
}
```

ループ条件を省略すれば 無限ループ(infinite loop)になる  

```go
for {
}
```

## if ステートメント

```go:ex09.go
func main() {
  a := -2
  if a < 0 {
    fmt.Println("負")
  }else{
    fmt.Println("正")
  }
}
```

if 文は for 文のように 条件の前に評価するための簡単なステートメントを書くことができる
if ステートメントで宣言された変数は else ブロック内でも使える

```go:ex09.go
func aaa(x, n, lim float64) float64 {
  if v := math.Pow(x, n); lim < v {
    return v
  }
  return 10
}

func bbb(x, n, lim float64) float64 {
  if v := math.Pow(x, n); lim < v {
    return v
  } else if v == lim {
    return v*2
  }
  return 10
}
```

戻り値の定義した関数を書いたら必ず通る `return` を書かないとエラーになる  

実行順の注意  
`fmt.Println()` は 2つの pow が先に実行された後に実行される  

```go:ex09.go
  fmt.Println(
    math.Pow(3, 2),
    math.Pow(3, 3),
  )
```

ざっくりと平方根を求める関数  
ニュートン法が使われている  
平方根で特に有効らしい

```go
func Sqrt(x float64) float64 {
  z := 1.0
  for i := 0; i < 10; i++ {
    z -= (z*z - x) / (2*z)
  }
  return z
}
```

## switch ステートメント
switch ステートメントは if - else ステートメントのシーケンスを短く書く方法  
選択された case だけを実行してそれに続く全ての case は実行されない  
つまり break を書かなくて良い  
上から下へ case を評価する  
cate は変数でも int 型以外でも使える  

```go:ex10.go
func main() {
  switch os := runtime.GOOS; os {
  case "darwin":
          fmt.Println("OS X.")
  case "linux":
          fmt.Println("Linux.")
  default:
          fmt.Printf("%s.\n", os)
  }
}
```

if-elseif-else の長くなりやすいつながりを Go では条件のない switch(switch true)で書く  

```go:ex10.go
func main() {
  t := time.Now()
  switch {
  case t.Hour() < 12:
    fmt.Println("Good morning!")
  case t.Hour() < 17:
    fmt.Println("Good afternoon.")
  default:
    fmt.Println("Good evening.")
  }
  // 以下と同じこと
  if t := time.Now(); t.Hour() < 12 {
    fmt.Println("Good morning!")
  } else if t.Hour() < 17 {
    fmt.Println("Good afternoon.")
  } else {
    fmt.Println("Good evening.")
  }
}
```

## defer ステートメント
defer へ渡した関数の実行を 呼び出し元の関数の終わり(returnする)まで遅延させる  

```go:ex11.go
func main() {
  defer fmt.Println("world")
  fmt.Println("hello")
  // hello worldと出力される
}
```

defer へ渡した関数が複数ある場合 その呼び出しはスタック(stack)される  
呼び出し元の関数が return するとき defer へ渡した関数は LIFO(last-in-first-out) の順番  

```go:ex11.go
for i := 0; i < 10; i++ {
        defer fmt.Println(i)
}
// 結果は9,8,7,6,5,4,3,2,1,0
```

## ポインタ
Go はポインタを扱う  
ポインタとは 値のメモリアドレスを表す  
変数 `v` のポインタは `*v` 型で ゼロ値は`nil`  
`&` オペレータは そのオペランド(operand)へのポインタを引き出す  
これは "dereferencing" または "indirecting" としてよく知られている  
`var p *int` と宣言することもできる  
`*` でポインタでなく値を取り出せる  

```go:ex12.go
  p := &i
  // p は i のポインタ(アドレスが格納されている)
  // *p は i と同義
  *p = 16
  // i に 16 を代入するのと同義
  // Println(*p) はできても Println(*i) はエラー
```

## struct (構造体)
ストラクトはフィールド(field)の集まり  

```go:ex13.go
type Vertex struct {
	X int
	Y int
}

func main() {
  fmt.Println(Vertex{1, 2})
  // 出力 {1 2}
}
```

struct の field は ドット(.)を用いてアクセス  

```go:ex13.go
func main() {
  v := Vertex{1, 2}
  v.X = 4
  fmt.Println(v.X)
  // 出力 4
}
```

struct の field は struct のポインタを通してアクセスすることもできる  
`(*p).X`と書けるが 面倒だから`p.X`と書くことが多い  

```go:ex13.go
  w := Vertex{3, 5}
  p := &w
  fmt.Println((*p).X)
```

struct リテラルは field の値を列挙することで新しい struct の初期値を割り当てることができる  
一部だけ列挙して初期化することもできる  
列挙の順番は関係ない  

```go:ex13.go
  v2 := Vertex{X: 6, Y: 7}
  fmt.Println(v2)
  // 出力 {6 7}
  v3 := Vertex{X:8}
  fmt.Println(v3)
  // 出力 {8 0}
  v4 := Vertex2{Z: 9, Y: 10}
  fmt.Println(v4)
  // 出力 {0 9 10}
```

`&`オペレータを頭に付けると 新しく割り当てられた struct へのポインタを渡す  

```go:ex13.go
  v5 := &Vertex{2, 1}
  fmt.Println(v5)
  // 出力 &{2 1}
```

func の中などで 定義と初期化を同時に行う  
```go
func aaa() {
	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}
}
```

## 配列
`[n]T` は `T` 型の `n` 個の変数の配列(array)を表す  
固定長だから 個数(配列のサイズ)の情報も含めて1つの型になる  
つまり配列のサイズを途中で変えることはできない  
関数に配列を渡す場合は値渡しとなり 配列のコピーが渡される  

```go:ex14.go
  var a [3]int
  fmt.Println(a[0], a[1])
  // 出力 0 0
  primes := [6]int{1, 2, 3, 4, 5, 6}
  fmt.Println(primes[1], primes[3])
  // 出力 2 4
```

## スライス(Slices)
`[]T` は `T` 型のスライスを表す  
配列は固定長だが スライスは可変長  
スライスは配列よりもより一般的に使う  
Java の ArrayList な感じ  
コロンで区切られた2つのインデックスで境界を指定する  
最初の要素は含むが 最後の要素は除いた半開区間を選択  

```go:ex14.go
  var s []int = primes[1:4]
  fmt.Println(s)
  // 出力 [2 3 4]
```

スライスは配列への参照のようなもの  
スライスの要素を変更すると その元となる配列の対応する要素が変更される  

```go:ex14.go
  names := [3]string{
    "Alice",
    "Bob",
    "Carol",
  }
  // 改行して配列を作った場合は最後にもコンマがいる
  fmt.Println(names)
  // 出力 [Alice Bob Carol]
  b := names[1:3]
  fmt.Println(b)
  // 出力 [Bob Carol]
  b[0] = "XXX"
  fmt.Println(names)
  // 出力 [Alice XXX Carol]
```

Slice リテラルは長さのない Array リテラルのようなもの  
`[]bool{true, true, false}` は `[3]bool{true, true, false}` の配列リテラルを作成し 同時に配列リテラルを参照するスライスを作成する  
また struct でも配列でもスライスも作れる  
struct を改行しながら格納するときは 最後の要素に","が必要  

```go:ex14.go
  c := []struct {
    i int
    b bool
  }{
    {2, true},
    {3, true},
    {5, true},
    {7, false},
    {11, false},
    {13, true},
  }
  fmt.Println(c)
  // 出力 [{2 true} {3 true} {5 true} {7 false} {11 false} {13 true}]
```

スライスするときは `a[0:]`など 上限または下限を省略できる  
省略した場合の下限は0 上限はスライスの長さになる  

スライスは長さ(length)と容量(capacity)を持つ  
- スライスの長さ: スライスに含まれる要素の数  
- スライスの容量: スライスの最初の要素から数えて 元となる配列の最後まで要素数  

容量より大きい長さは指定できない  
スライス s の長さは `len(s)` 容量は `cap(s)` で取得する  
必ず len <= cap になる  

`s[左:右]`にて  
- 右だけ書いてあるパターンは cap が保存される, スライスの最初が0番目になるため
- 左だけ書いてあるパターンは cap は減る, スライスの最初が左の番目になるため

スライスのゼロ値は nil  
nil スライスは 0の長さと 0の容量を持っているが 元となる配列は持っていない。  

スライスは 組み込み(パッケージ不要)の `make()` を使用して作成できる  
これは 動的サイズの配列を作成する方法でもある  
`make([]型, 長さ, 容量)` で作成する  
make 関数は内部的にはゼロ化された配列を生成し その配列を指すスライスを返す  

```go:ex15.go
  a := make([]int, 5)
  fmt.Println(a, len(a), cap(a))
  // 出力 [0 0 0 0 0] 5 5
  b := make([]int, 0, 5)
  fmt.Println(b, len(b), cap(b))
  // 出力 [] 0 5
```

スライスは 他のスライスを含む任意の型を含むことができる  
つまりn次元配列を作れる  

```go:ex15.go
  import("strings")
  board := [][]string{
    []string{"-", "-", "-"},
    []string{"-", "-", "-"},
    []string{"-", "-", "-"},
  }
  for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
  }
  // 出力
  // - - -
  // - - -
  // - - -
```

スライスへ新しい要素を追加するには Goの組み込みの `append()` を使う  

```go:ex15.go
  var c []int
  fmt.Println(c, len(c), cap(c))
  // 出力 [] 0 0
  c = append(c, 0, 1, 2)
  fmt.Println(c, len(c), cap(c))
  // 出力 [0 1 2] 3 4
  // 追加するリストの容量が 元のリストの容量より大きい場合 メモリ上はリストを割り当て直すみたい
  // その時に 容量が変化するのかも
```

for ループに利用する range は スライスや マップ(map)をひとつずつ反復処理するために使う  
スライスを range で繰り返す場合 range は反復毎に2つの変数を返す  
1つ目の変数はインデックス(index) 2つ目の変数はインデックスの場所の要素のコピー  

```go:ex15.go
  d := make([]int, 3)
  d[1] = 10
  for i, value := range d {
    fmt.Printf("%d: %d\n", i, value)
  }
  // 出力
  // 0: 0
  // 1: 10
  // 2: 0
```

インデックスや値は "_"(アンダースコア)へ代入することで捨てることができる  

```go:ex15.go
for i, _ := range d
for _, value := range d
```

もし インデックスだけが必要なのであれば 2つ目の値を省略できる `for i := range d`  

## マップ(map)
map はキーと値とを関連付ける  
マップのゼロ値は nil  
nil のマップはキーを持っておらず キーの追加もできない  
`make()` は初期化され使用可能な指定された型のマップを返す  
もし map に渡すトップレベルの型が単純な型名である場合は リテラルの要素から推定できるため 型名を省略できる  

```go:ex16.go
type Vertex struct {
  x, y float64
}

func main() {
  // key は string, value は Vertex で宣言
  // まだ nil な map で使えない
  var m map[string]Vertex
  // make() で使える map にする
  m = make(map[string]Vertex)
  m["key01"] = Vertex{
    3.14, -2.71,
  }
  m["key02"] = Vertex{1, 3}
  fmt.Println(m["key01"])
  // 出力 {3.14 -2.71}
  fmt.Println(m)
  // 出力 map[key01:{3.14 -2.71} key02:{1 3}]

  // make() を使わない方法
  var p = map[string]Vertex{
    "key01":Vertex{1, 2},
    "key02":Vertex{3, 4},
  }
  fmt.Println(p)
  // 出力 map[key01:{1 2} key02:{3 4}]

  // 型を省略した書き方
  var r = map[string]Vertex{
    "key01": {1.41, 1.73},
    "key02": {5, 6},
  }
  fmt.Println(r)
  // 出力 map[key01:{1.41 1.73} key02:{5 6}]
}
```

map 型 m への操作  
要素(element)の挿入や更新: `m[key] = element`  
要素の取得: `element = m[key]`  
要素の削除: `delete(m, key)`  
キーに対する要素が存在するかどうか: `element, ok := m[key]`  
もし m に key があれば 変数 ok は true となり 存在しなければ ok は false になる  
また map に key が存在しない場合 element は map の要素の型のゼロ値になる  

```go:ex16.go
  a := make(map[string]int)
  a["Answer"] = 42
  fmt.Println(a)
  // 出力 map[Answer:42]
  delete(a, "Answer")
  v, ok := a["Answer"]
  fmt.Println(v, ok)
  // 出力 0 false
```

計算量を無視した単語カウンターの自作  

```go:ex17.go
import (
  "fmt"
  "strings"
)

func WordCount(s string) map[string]int {
  a := strings.Fields(s)
  b := make(map[string]int)
  for i := 0; i < len(a); i++ {
    count := 0
    for j := 0; j < len(a); j++ {
      if a[j] == a[i] {
        count += 1
      }
    }
    b[string(a[i])] = count
  }
  return b
}

func main() {
  str := "Hello World Hello"
  fmt.Println(WordCount(str))
  // 出力 map[Hello:2 World:1]
}
```
## 関数値(function value)
関数も変数であり 他の変数のように関数を渡すことができる  
使う関数を引数で変える的な。pythonでもあったな  

```go:18.go
import (
  "fmt"
  "math"
)

func compute(fn func(float64, float64) float64) float64 {
  return fn(3, 4)
}

func main() {
  hypot := func(x, y float64) float64 {
    return math.Sqrt(x*x + y*y)
  }

  fmt.Println(hypot(3, 4))
  // 出力 5 なぜなら3^2+4^2=5^2 の平方根
  fmt.Println(compute(hypot))
  // 出力 5 なぜなら渡す値はcomputeで既に決まっているから
  fmt.Println(compute(math.Pow))
  // 出力 81 なぜなら3^4=81
}
```

## Goの関数は クロージャ(closure)
クロージャ(関数閉包)とはプログラミング言語における関数オブジェクトの一種  
ラムダ式や無名関数の概念  

`adder()` はクロージャを返す  
そして 各クロージャ(ここでは for 文の中の pos)は それ自身(adder)の sum 変数へバインドされます  

```go:ex19.go
func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func main() {
	// pos には sum += x の関数が入っている
	pos := adder()
	for i := 0; i < 3; i++ {
		// pos(i) は sum
		// i は x
		fmt.Println(pos(i))
	}
	// 出力
	// 0
	// 1
	// 3

	fmt.Println(pos(10))
	// 出力 13 なぜなら今保持してる3が足されるから
}
```

`pos:=adder()` とやると `sum:=0` までは実行されて `return func(int) int` が `pos` に入ってる感じ  
だから `pos(10)` とかは動く  
そして for 文2回目の pos とかは for 文1回目の pos の sum が入っている  

```go:ex20.go
func adder() func(int) int {
	sum := 5
	fmt.Println("-1-", sum)
	return func(x int) int {
		fmt.Println("-2-", sum)
		sum += x
		fmt.Println("-3-", x)
		fmt.Println("-4-", sum)
		return sum
	}
}

func main() {
	fmt.Println("when 01")
	pos := adder()
	fmt.Println("when 02")
	fmt.Println(pos(10))
	fmt.Println("when 03")
	for i := 0; i < 3; i++ {
		fmt.Println("when 04")
		fmt.Println(pos(i))
	}
	// 出力
	// when 01
	// -1- 5
	// when 02
	// -2- 5
	// -3- 10
	// -4- 15
	// 15
	// when 03
	// when 04
	// -2- 15
	// -3- 0
	// -4- 15
	// 15
	// when 04
	// -2- 15
	// -3- 1
	// -4- 16
	// 16
	// when 04
	// -2- 16
	// -3- 2
	// -4- 18
	// 18
}
```

クロージャを使ったフィボナッチ数列の自作  

```go:ex21.go
func fibonacci() func() int {
	a := 0
	b := 1
	return func() int {
		out := a
		tmp := a + b
		a = b
		b = tmp
		return out
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

## メソッド
メソッドは クラスや型などに 属すもの  
関数は クラスや型などに 属さないもの  

Go にはクラス(class)の仕組みが無い  
しかし 型にメソッド(method)を定義できる  
親となる型?対象?をレシーバという  

書き方は メソッドを定義するに レシーバを func とメソッド名の間に 自身の引数リストで記述する  
`Abs1()` も `Abs2()` も どっちで書いても機能的には同じに見えるが 内部が違う  

```go:ex22.go
type Vertex struct {
	X, Y float64
}

// メソッドの作成
// Abs01 メソッドは v という名前の Vertex 型のレシーバを持つ。型は任意。
func (v Vertex) Abs01() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// 関数の作成
func Abs02(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// レシーバをポインタした書き方
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func main() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs01())
	// 出力 5
	fmt.Println(Abs02(v))
	// 出力 5
	v.Scale(10)
	fmt.Println(v)
	// 出力 {30 40}
}
```

レシーバには struct 型だけでなく 任意の型が使える  
`Vertex` から `type MyFloat float64` で宣言した `MyFloat` 型などでも良い  

レシーバを伴うメソッドの宣言は レシーバ型が同じパッケージにある必要がある  
つまり 他のパッケージに定義している型をレシーバとしたメソッドは作れない  

レシーバは ポインタにすることもできる  
レシーバ自身を更新することが多いため 変数レシーバよりもポインタレシーバの方が一般的  

レシーバは 第0引数って感覚らしい  

```go
// 関数宣言で 引数がポインタ
func ScaleFunc(v *Vertex, f float64) {}
ScaleFunc(v, 5)  // Compile error!
ScaleFunc(&v, 5) // OK
// → ポインタを渡す。一致してないとダメ(引数ポインタなんだから)

// 関数宣言で 引数が値
func AbsFunc(v Vertex) float64 {}
fmt.Println(AbsFunc(v))  // OK
fmt.Println(AbsFunc(&v)) // Compile error!
// → 値を渡す。一致してないとダメ(フツーの関数)

// レシーバ宣言で レシーバ(第0引数)がポインタ
func (v *Vertex) Scale(f float64) {}
v.Scale(5)  // OK Go が気を効かせてるだけで内部的には (&v).Scale(5) として解釈
&v.Scale(10) // OK
// → 第0引数には 値でもポインタでもok

// レシーバ宣言で レシーバ(第0引数)が値
func (v Vertex) Abs() float64 {}
fmt.Println(v.Abs()) // OK
fmt.Println(&v.Abs()) // OK Go が気を効かせてるだけで内部的には v.Abs() として解釈
// → 第0引数には 値でもポインタでもok
```

関数宣言で 引数がポインタ → ポインタを渡す。一致してないとダメ(引数ポインタなんだから)  
関数宣言で 引数が値 → 値を渡す。一致してないとダメ(フツーの関数)  
レシーバ宣言で レシーバ(第0引数)がポインタ → 第0引数には 値でもポインタでも ok  
レシーバ宣言で レシーバ(第0引数)が値 → 第0引数には 値でもポインタでも ok  

```go
func (p Person) Greet(msg string) {} // メソッド宣言
func Person.Greet(p Person, msg string) // 関数宣言
// 上の2個は同じこと
func (pp *Person) Shout(msg string) // メソッド宣言
func (*Person).Shout(pp *Person, msg string) // 関数宣言
// 上の2個は同じこと
```

ポインタレシーバを使う理由  
- メソッドがレシーバが指す先の変数を変更するため
- メソッドの呼び出し毎に変数のコピーを避けるため

構造体が大きいと 呼び出した時に全部コピーするのは処理が重くなるから  

変数レシーバというより 値レシーバという名前の方が一般的?  
普通に引数で渡すと コピーされたことになるから メソッド内で値を変更しても 元のデータには影響しない  

非同期処理で安全とも言い変えることができる  
- メソッド呼び出しごとにレシーバの値はコピーされる
- 値レシーバの値はメソッド内で書き換えても元のレシーバの値にはまったく影響がない
- ポインタレシーバにするとメソッド内でレシーバの値を書き換えられる
- ポインタレシーバのメソッド内でフィールドを呼び出す場合には その前に nil チェックをすべき
値レシーバにしたほうがよい場合  
- int, string などのプリミティブ型をベースとする型
- map, chan といった参照型をベースとする型(最初からポインタだから)
- 不変型
- 小さい構造体

## interface(インタフェース)
interface(インタフェース)型は メソッドのシグニチャの集まりで定義する  
メソッドの集まりを実装した値を interface 型の変数へ持たせることができる  
型にメソッドを実装していくことによって インタフェースを実装する  

```go
// インタフェース宣言
type I interface {
	M()
}
// struct 宣言 T 型は string を持つ
type T struct {
	S string
}
// レシーバとして T 型を指定して メソッド M() を実装
func (t T) M() {
	fmt.Println(t.S)
}

func main() {
  // T 型を インタフェース型に入れて インタフェース型のメソッド M() を実行
	var i I = T{"hello"}
	i.M()
}
```

見かけの型と実際に保持してる型が異なるのは インタフェースのメリット  

```go:ex23.go
func output(i I) {
	fmt.Printf("値: %v, 型: %T\n", i, i)
}

// インタフェース宣言
type I interface {
	M()
}
// struct 宣言 T 型は string を持つ
type T struct {
	S string
}
// レシーバとして *T 型を指定して メソッド M() を実装
func (t *T) M() {
	fmt.Println(t.S)
}
// struct 宣言 F 型は float64 を持つ
type F float64
// レシーバとして F 型を指定して メソッド M() を実装
func (f F) M() {
	fmt.Println(f)
}

func main() {
	// インタフェース型の宣言
	var i I
	// T 型のポインタを i に格納
	i = &T{"Hello"}
	// (value, type) を確認
	output(i)
	// 出力 値: &{Hello}, 型: *main.T
	i.M()
	// 出力 Hello

	// F 型を i に格納
	i = F(3.14)
	output(i)
	// 出力 値: 3.14, 型: main.F
	i.M()
	// 出力 3.14
}
```

nil レシーバ  
インタフェースの実体?(インタフェースを見かけ上の型だから 実装元の型のこと笑)が nil ならメソッドは nil をレシーバーとして呼び出される  
値が`<nil>`的な  
nil インタフェースの値は 値も具体的な型も持たない  

フツーの言語だと nil を参照しただけでエラーになるけど Go はならない  
nil インタフェースのメソッドを呼び出すと Go ではランタイムエラーになる  
nil を適切に処理するメソッドを書くのが一般的  

```go:ex24.go
func (t *T) M() {
	// nil だったときの処理
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}
```

ゼロ個のメソッドを指定されたインタフェース型は  空のインタフェース という  
`interface{}`  
空のインタフェースは 任意の型の値を保持できる  

```go
func output(i interface{}) {
	fmt.Printf("値: %v, 型: %T\n", i, i)
}

func main() {
	var i interface{}
	output(i)
	// 出力 値: <nil>, 型: <nil>

	i = 42
	output(i)
	// 出力 値: 42, 型: int

	i = "hello"
	output(i)
	// 出力 値: hello, 型: string
}
```

## 型アサーション
アサーション とは 表明 のこと プログラムの前提条件を示すのに使う  
型アサーション を使うと インタフェースの実体の型がなにか明らかにすることができる  
空のインタフェースは 任意の型を保持できるが実際に使うときに動的にチェックしないといけないから そのときに型アサーションを使う  
単純に型チェックなど デバッグでも使いそう  
型アサーション は インタフェースの値の素になる具体的な値を利用する手段を提供する  
アサーションは コードのその箇所で必ず真であるべき式の形式をとる  
`t := i.(T)`  
インタフェースの値 i が具体的な型 T を保持し 基になる T 型の値を変数 t に代入する  
ただ i が T を持っていないとき`panic`というエラーになる  

型を保持しているか確認するために 戻り値を2個受け取ることができる  
`t, ok := i.(T)`  
i が T を保持していれば t は基になる値になり ok は true になる  
i が T を保持していなければ t は T のゼロ値 ok は false になる  
ok を受け取る書き方だと panic は起こらない  

```go:ex26.go
func main() {
	// 実体は string 型の 任意の型を保持できる空インタフェースを宣言
	var i interface{} = "hello"
	
	// 型アサーション
	s, ok := i.(string)
	fmt.Println(s, ok)
	// 出力 hello true

	// 型アサーション
	f, ok := i.(float64)
	fmt.Println(f, ok)
	// 出力 0 false 実体は string だから
}
```

型 switch はいくつかの型アサーションを直列に使用できる構造  
型 switch は通常の switch 文と似ている  
型 switch の case は型を指定し その値は指定されたインターフェースの値が保持する値の型と比較する  
型 switch を使うとき 型アサーションの部分を type にする

```go:ex26.go
func do(i interface{}) {
	// i.(type) でとりあえず実体の型を受け取って case で判定する
	switch v := i.(type) {
	case int:
		fmt.Printf("型: %T, 値: %v, 2乗: %v\n", v, v, v*2)
	case string:
		fmt.Printf("型: %T, 値: %v, %q is %v bytes long\n", v, v, v, len(v))
	default:
		fmt.Printf("%T型なんて知らん!\n", v)
	}
}
func main() {
	do(3)
	// 出力 型: int, 値: 3, 2乗: 6
	do("hello")
	// 出力 型: string, 値: hello, "hello" is 5 bytes long
	do(true)
	// 出力 bool型なんて知らん!
}
```

インタフェースの例  
もっともよく使われている interface の1つに fmt パッケージ に定義されている Stringer がある  
Java でいう toString メソッドな感じ? その型を `Println()` した時のデフォルト表記的な  
単純に言うと `String()` というメソッドを任意のレシーバに定義すると `Println()` したときのフォーマットを定義できる  

```go
type Stringer interface {
    String() string
}
```

`String()`の実装練習  

```go:ex27.go
type IPAddr [4]byte

func (i IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}

	for name, ip := range hosts {
		// 本来なら loopback: [127 0 0 1] とただの配列として出力される
		fmt.Printf("%v: %v\n", name, ip)
		// String() を IPAddr 型に実装したことで
		// loopback: 127.0.0.1 と出力されるようになった
	}
}
```

## 戻り値としてのエラー
Go のプログラムは エラーの状態を error 値で表す  
error 型は fmt.Stringer に似た組み込みのインタフェース  
Go では 関数が複数の値値を返せることを利用して 内部で発生したエラーを戻り値で表現する  
この値を使ってエラーハンドリングする  

```go
type error interface {
    Error() string
}
```

error が nil なら成功  
error が nilでない(何かしらエラーメッセージを呼び出し元に返した場合)は失敗  

関数の処理に成功した場合 エラーは nil 異常があった場合 エラーだけに値が入り他方はゼロ値になる  

```go:ex28.go
type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("いつ: %v, 何が: %s", e.When, e.What)
}

func run() error {
	return &MyError{
		time.Now(),
		"実行できません",
	}
}

func main() {
	// err が nil じゃない つまり エラーが起こっている
	if err := run(); err != nil {
		fmt.Println(err)
		// 出力 いつ: 2021-01-16 00:01:35.7540284 +0000 UTC m=+0.000046701, 何が: 実行できません
	}
}
```

## io パッケージ
io パッケージは データストリームを読むことを表現する `io.Reader` インタフェースを規定している  
Go の標準ライブラリには ファイル, ネットワーク接続, 圧縮, 暗号化 などで このインタフェースが実装されていることが多い  
`io.Reader` インタフェースは Read メソッドを持つ  
`func (T) Read(b []byte) (n int, err error)`  
Read は データを与えられたバイトのスライスに入れ 入れたバイトのサイズとエラーの値を返す  
ストリームの終端は `io.EOF` というエラーを返す  

```go:ex29.go
import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("Hello, Reader!")

	// ゼロ値が8個のバイト型の配列
	b := make([]byte, 8)
	fmt.Println(b, len(b), cap(b))
	// 出力 [0 0 0 0 0 0 0 0] 8 8
	for {
		// 8個のバイトごとに取り出す感じ
		n, err := r.Read(b)
		fmt.Printf("n: %v, err: %v\n", n, err,)
		fmt.Printf("b[:n]: %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
	// 出力
	// n: 8, err: <nil>
	// b[:n]: "Hello, R"
	// n: 6, err: <nil>
	// b[:n]: "eader!"
	// n: 0, err: EOF
	// b[:n]: ""
}
```

よくあるパターンは `io.Reader` をラップし ストリームの内容を何らかの方法で変換する `io.Reader` を作る  
例えば gzip.NewReader は `io.Reader` (gzipされたデータストリーム)を引数で受け取り `*gzip.Reader` を返す  
その `*gzip.Reader` は `io.Reader` (展開したデータストリーム)を実装している  

## image パッケージ
image パッケージは 以下の Image インタフェースを定義している  

```go
package image

type Image interface {
	ColorModel() color.Model
	Bounds() Rectangle
	At(x, y int) color.Color
}
```

`image.Rect(0, 0, width, height)` のようにして `image.Rectangle` を返す  
ColorModel は `color.RGBAModel` を返す  
At は ひとつの色を返す  

`color.Model` は定義済みの `color.RGBAModel` で代用できる
`color.Color` は定義済みの `color.RGBA` で代用できる
代用することで `color.Color` と `color.Model` の2個のインタフェースを無視できる

```go
package main

import (
	"fmt"
	"image"
)

func main() {
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(m.Bounds())
	// 出力 (0,0)-(100,100)
	fmt.Println(m.At(0, 0).RGBA())
	// 出力 0 0 0 0
}
```

## goroutine(ゴルーティン, Goルーティン)
goroutine (ゴルーティン)は Goのランタイムに管理される軽量なスレッド  

`go f(x, y, z)`と書けば 新しい goroutine が実行する  
f , x , y , z の評価自体は 実行元(カレント:current)の goroutine で実行される  
fメソッド の実行は 新しい goroutine で実行する  
goroutine は 同じアドレス空間で実行されるため 共有メモリへのアクセスは必ず同期する必要がある  

チャネル(Channel)型は チャネルオペレータの `<-` を用いて値の送受信ができる通り道  
(データは 矢印の方向に流れる)  

```go
ch <- v    // v をチャネル ch へ送信する
v := <-ch  // ch から受信した変数を v へ割り当てる
```

基本は片方が準備できるまで送受信はブロックされる  
つまり 明確なロックや条件変数がなくても goroutine の同期が簡単にできる  

チャネルを使うときは マップとスライスと同様に 使う前に生成する必要がある  
`ch := make(chan int)`  

```go:ex30.go
// int配列 を受け取るのと 返す先のチャネルを指定されて 指定されたチャネルへ合計値を返す
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	// c というチャネルを作成する
	c := make(chan int)
	// int配列を2等分して c チャネルへ送る
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	// c チャネルから受け取る
	x, y := <-c, <-c // receive from c
	// 受け取った値を出力する
	fmt.Println(x, y, x+y)
	// 出力 -5 17 12
}
```

チャネルは バッファ(buffer)としても使える  
バッファを持つチャネルを初期化するには make の2つ目の引数にバッファの長さを与える  
`ch := make(chan int, 100)`  
バッファが詰まった時は チャネルへの送信をブロックする  
バッファが空の時には チャネルの受信(チャネルから変数へ送りきる)をブロック  

```go:ex31.go
func main() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(<-ch)
	// 先に ch から受け取ってからじゃないとエラーになる
	ch <- 4
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

実は ループの `for i := range c` は チャネルが閉じられるまでチャネルから値を繰り返し受信し続けるという処理をしていた  

送り手は これ以上の送信する値がないことを伝えるために チャネルを close できる  
close されているかは チャネルから2個目の戻り値を設定して ok が false かで確かめる  
`v, ok := <-ch`  
受信する値がない かつ チャネルが閉じているなら ok は false  
チャネルは 通常 close する必要はない  
close するのは これ以上値が来ないことを受け手が知る必要があるときにだけ  
例えば range ループの終了など  

close するなら送り手のチャネルだけを close する  
もし close したチャネルへ送信すると panic を起こす  

```go:ex32.go
// チャネルを使った自作フィボナッチ数列
// close の扱いが微妙な気がする
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
func main() {
	c := make(chan int, 3)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
```

## select ステートメント
select ステートメントは 複数の通信操作のこと  
goroutine を待機させることができる  
つまり 複数のチャネルで受信を待機させることができる  
select は複数ある case のどれかが準備できるようになるまでブロックする  
準備ができた case を実行する  
もし複数の case の準備ができているなら case はランダムに選択される  
どの case も準備ができていないのであれば select の中の default が実行される  

```go:ex33.go
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		// どれかの case が実行できるようになるまで goroutine スレッドを待機させる
		select {
		case c <- x: // x を c チャネルに送信したら実行する
			x, y = y, x+y
		case <-quit: // quit チャネルを値を受け取ったら実行する
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	// 2個のチャネルを生成
	c := make(chan int)
	quit := make(chan int)
	// for文中は c チャネルに値を渡して 終わったら quit チャネルに値を渡す という goroutine スレッドを作る
	go func() {
		for i := 0; i < 4; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()

	fibonacci(c, quit)
	// 出力
	// 0
	// 1
	// 1
	// 2
	// quit
}
```

一番最後の `fibonacci(c, quit)` が起動しつつも 別スレッドでfor文が回っていて チャネル c で値を受け取っていたら `fmt.Println(<-c)` をする感じ  

```go:ex34.go
// 自作で select の default を実験した
func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("チック.")
		case <-boom:
			fmt.Println("ボーン!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
```

## 排他制御の mutex(ミューテックス)
チャネルは goroutine 間の通信の素晴らしい方法だった  
だが 通信が必要ない場合はどうするか  
一度に1つの goroutine だけが変数にアクセスできるようにしたい場合はどうするか  
つまり コンフリクトを避けたい  
そのときは 排他制御(mutual exclusion)を使う  
排他制御のデータ構造を示す一般的な名前は mutex (ミューテックス)  

Go の標準ライブラリは 排他制御を `sync.Mutex` と2つのメソッド(`Lock()`, `Unlock()`)で提供する  
`Lock()` と `Unlock()` で囲むことで排他制御で実行するコードを定義する  

```go:ex35.go
import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter 型は排他制御ができて key の数を保持する
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc メソッドは指定された key のカウンタを増やす
func (c *SafeCounter) Inc(key string) {
	// 一度に1つの goroutine しか map にアクセスできないようにロックする
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}

// Value メソッドは指定された key のカウンタ値を返す
func (c *SafeCounter) Value(key string) int {
	// 一度に1つの goroutine しか map にアクセスできないようにロックする
	c.mu.Lock()
	// mutex が Unlock されることを保証するために defer を使うこともできる
	defer c.mu.Unlock()
	return c.v[key]
}

func main() {
	c := SafeCounter{
		v: make(map[string]int),
	}
	for i := 0; i < 10; i++ {
		go c.Inc("key")
		fmt.Println("今は", i, "番目")
	}
	// main とは 別の goroutine の処理によって結果が変わるため 数秒待つ
	fmt.Println(time.Second)
	// 出力 1s
	time.Sleep(time.Second)
	fmt.Println("keyの数は", c.Value("key"))
	// 出力 keyの数は 10
}
```
