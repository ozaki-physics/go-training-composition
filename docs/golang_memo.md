## Go 言語の仕様
### 構造体を比較するとき 自動で中身の値比較をしてくれる
Go の言語仕様として __構造体の比較は ポインタの比較じゃなくて構造体の中の値の比較__ になっているらしい!  
[Comparison operators](https://go.dev/ref/spec#Comparison_operators)  

>構造体の値は、そのすべてのフィールドが比較可能であれば、比較可能である。  
>2つの構造体値が等しいのは、対応する非空白フィールドが等しい場合である。  

だから わざわざ Equals() を定義する必要はない  

僕は `==` では インスタンスの番号的なもの(アドレス)が等しくないと true にならないと思っていた  
`type Point struct {x,y int}` として  
`Point{x: 10, y: 20} == Point{x: 10, y: 20}` は struct としては別物だから false になると思っていた  
しかし Go では 言語仕様として struct の値で比較してくれるので 上記のコードは true になる  
等しいという概念には 等値 と 等価 が存在する  
- 等値: アドレスまで同じ  
- 等価: struct の中の値が等しい  
[golangのequalityの評価について](https://pod.hatenablog.com/entry/2016/07/30/204357)  

同様に 等しいという概念は struct だけでなく インスタンスにも拡張できる
>インスタンス同士が比較可能であるためには以下の2つの条件がが必要である。  
>- インスタンスの型が同一（等価）であること  
>- インスタンスの型が比較可能であること  
>
>たとえばある型を別の型に再定義しただけの場合でも等価とは見なされず，コンパイルエラーになる。  

注意点は キャスト と type alias では等価の挙動に少し違いがあるみたい  
キャスト: Number01 型を `type Number01 int` と定義して Number01 型と int 型に1を代入して `==` すると compile error になる  
ただし Number01 型を int にキャストして `==` すると true になる  
type alias: Number02 型を `type Number02 = int` と定義して Number02 型と int 型に1を代入して `==` すると true になる  
[インスタンスの比較可能性](https://text.baldanders.info/golang/comparability/)  

```go
type Number01 int
var i1 int = 1
var n1 Number01 = 1
fmt.Println(i1 == n1) // compile error
fmt.Println(i1 == int(n1)) // true
type Number02 = int
var n2 Number02 = 1
fmt.Println(i1 == n2) // true
```

参考として 公式リファレンスに どういうときに 等しい扱いになるかっぽいことが書かれている  
[Package reflect](https://pkg.go.dev/reflect@go1.17.6#DeepEqual)  

```go
import (
  "reflect"
)

type A struct {
  Value string
}

func main() {
  a1 := A{"value"}
  a2 := A{"value"}

  println(reflect.DeepEqual(a1, a2))   // true
  println(reflect.DeepEqual(&a1, &a2)) // true
  println(a1 == a2)                    // true
  println(&a1 == &a2)                  // false
}
```


### Go 言語の継承の考え方
[公式ドキュメント Embedding(埋め込み)](https://go.dev/doc/effective_go#embedding)  
理解できた部分だけ 抽出してまとめた  

>Go does not provide the typical, type-driven notion of subclassing, but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface.  
>Go は典型的な型駆動型のサブクラス化の概念を提供しませんが、構造体やインターフェースに型を埋め込んで、実装の一部を「借りる」機能があります。  

そもそも 継承ではなく インタフェースを使うべきと言われてるっぽい  
埋め込みは struct をラップする感じで使う  
[練習したコード](.././trainingEmbedding)  

埋め込みには3パターンある  
1. インタフェース に インタフェース を埋め込む
2. インタフェース に struct を埋め込む(できない)
3. struct に インタフェース を埋め込む
4. struct に struct を埋め込む

インタフェース の type には インタフェースしか埋め込めない  
```go
// ReadWriter is the interface that combines the Reader and Writer interfaces.
// ReadWriter は Reader と Writer の2つのインターフェースを統合したインターフェースです。
type ReadWriter interface {
    Reader
    Writer
}
// 構造体を定義するときに 名前の後ろに型を書かないのはなんで?
// -> Reader, Writer もインタフェースであり 型そのものだから 名前と思っているものが むしろ型
```

また struct に struct を埋め込むときは ポインタにするらしい

>The embedded elements are pointers to structs and of course must be initialized to point to valid structs before they can be used.  
>The ReadWriter struct could be written as  
>埋め込み要素は構造体へのポインタであり、もちろん、使用する前に有効な構造体を指すように初期化する必要があります。
>ReadWriter 構造体は、次のように記述できます。

埋め込み元(Parent) と 埋め込み先(Child) で同じ メソッド名やフィールド名があったとき  
埋め込み先(Child)のメソッド名やフィールド名が優先される  

```go
type A struct {
  Name string
  Age  int
}

func (a A) Print() {
  println("struct A", "name:", a.Name, "age:", a.Age)
}

// 埋め込み
type B struct {
  A
  // ポインタの場合は *A にして &A{} を渡す
}

func (b B) Print() {
  // 埋め込みで A の Name と Age が使える
  println("name:", b.Name, ", age:", b.Age)
  // 以下でも同じ
  println("name:", b.A.Name, ", age:", b.A.Age)
}

// オーバーライド
type C struct {
  Name string // A の Name は上書きされる
  A
}

// A の Print 関数は上書きされる
func (c C) Print() {
  println("struct C", "name:", c.Name, "age:", c.Age)
  // 埋め込んだ A のフィールドや関数は c.A.Name とか c.A.Print() で使える
}

func main() {
  b := B{A{"name01", 31}}
  b.Print() // name: name01, age: 31

  c := C{Name: "name02", A: A{"name01", 31}}
  c.Print()   // struct C name:  name02 age: 31
  // A の Print を使う
  c.A.Print() // struct A name: name01 age: 31
}
```

[Go言語(golang) 構造体の定義と使い方](https://golang.hateblo.jp/entry/golang-how-to-use-struct)  

### Go の toString 的なメソッド
任意の表記にして出力する方法は その struct に String メソッドを定義すること  
String メソッドのレシーバをポインタにすると ポインタとして渡された構造体だけに フォーマットが適用される  

```go
// int 型に 異なるラベルを付ける
type anotherInt int

func (a *anotherInt) String() string {
	return fmt.Sprintf("anotherInt: %d\n", a)
}

func main() {
	var i anotherInt
	i = 10
	fmt.Println(i) // 10
	fmt.Println(&i) // 0xc0000be000
  // String を定義した後
	fmt.Println(i) // 10
	fmt.Println(&i) // anotherInt: 824634474496 変な値にはなるがフォーマットは適用されている
}
```

### レシーバをポインタにするなら 一度インスタンスを生成しないと そのメソッドが使えない感じがする
Factory Method パターン と相性が悪い?  
