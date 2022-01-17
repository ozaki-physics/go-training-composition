## Go 言語の仕様
### 構造体を比較するとき 自動で中身の値比較をしてくれる
Go の言語仕様として __構造体の比較は ポインタの比較じゃなくて構造体の中の値の比較__ になっているらしい!  
[Comparison operators](https://go.dev/ref/spec#Comparison_operators)  

>構造体の値は、そのすべてのフィールドが比較可能であれば、比較可能である。  
>2つの構造体値が等しいのは、対応する非空白フィールドが等しい場合である。  

だから わざわざ Equals() を定義する必要はない  

#### 他にも参考資料  
[Package reflect](https://pkg.go.dev/reflect@go1.17.6#DeepEqual)  

[golangのequalityの評価について](https://pod.hatenablog.com/entry/2016/07/30/204357)  

```go
type Point struct {
  x,y int
}

// 等値
pt := Point{x: 10, y: 20}
t.Errorf("%v", pt == pt) // => true

// 等価
// (先入観でこれはfalseだと思っていた)
t.Errorf("%v", Point{x: 10, y: 20} == Point{x: 10, y: 20}) // => true
```

[Go言語(golang) 構造体の定義と使い方](https://golang.hateblo.jp/entry/golang-how-to-use-struct)  

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

[インスタンスの比較可能性](https://text.baldanders.info/golang/comparability/)  
>インスタンス同士が比較可能であるためには以下の2つの条件がが必要である。  
>- インスタンスの型が同一（等価）であること  
>- インスタンスの型が比較可能であること  
>
>たとえばある型を別の型に再定義しただけの場合でも等価とは見なされず，コンパイルエラーになる。  

```go
type Number int
var c1 int = 1
var c2 Number = 1
fmt.Println(c1 == c2) // compile error

// 等価な型にキャスト可能であれば エラーにならない
type Number int
var c1 int = 1
var c2 Number = 1
fmt.Println(c1 == int(c2)) // true

// また type alias であれば等価
type Number = int
var c1 int = 1
var c2 Number = 1
fmt.Println(c1 == c2) // true
```

### Go 言語の継承の考え方
そもそも 継承ではなく インタフェースを使うべきと言われている  
[Go言語(golang) 構造体の定義と使い方](https://golang.hateblo.jp/entry/golang-how-to-use-struct)  

以下のように 他言語のクラスの継承はGo言語では埋め込みで対応できるっぽい  

```go
type A struct {
  Name string
  Age  int
}

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

func main() {
  b := B{A{"Tanaka", 31}}
  b.Print() // name: Tanaka, age: 31
}
```

以下のように オーバーライドもできるっぽい  
ただ インタフェース の実装で対応すべきだと思う  

```go
type A struct {
  Name string
  Age  int
}

func (a A) Print() {
  println("struct A", "name:", a.Name, "age:", a.Age)
}

type B struct {
  Name string // A の Name は上書きされる
  A
}

// A の Print 関数は上書きされる
func (b B) Print() {
  println("struct B", "name:", b.Name, "age:", b.Age)
  // 埋め込んだ A のフィールドや関数は b.A.Name とか b.A.Print() で使える
}

func main() {
  b := B{Name: "Suzuki", A: A{"Tanaka", 31}}

  b.Print()   // struct B name:  Suzuki age: 31

  // A の Print を使う
  b.A.Print() // struct A name: Tanaka age: 31
}
```

[公式ドキュメント Embedding(埋め込み)](https://go.dev/doc/effective_go#embedding)  
インタフェース の struct には インタフェースしか埋め込めない  
構造体を定義するときに型を書かないのはなんで?  
