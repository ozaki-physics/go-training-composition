# Go言語の勉強
golangを勉強するのに使えそうなサイト
1. https://gihyo.jp/dev/feature/01/go_4beginners
2. https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/index.html
3. https://www.slideshare.net/takuyaueda967/2016-go?ref=https://kirohi.com/to_study_golang
4. https://go-tour-jp.appspot.com/list
主に1番目、2番目を使う。
でも1番目は、2014年6月にリリースされた最新バージョン1.3をベースにGoについて徹底解説→古い。
2番目は、分からん。
3番目は2017年でバージョン1.8
4番目は公式だから、一番いいかもしれないが、サイト上で完結するように設定されてしまっている。


本当にgolang勉強する時、どういう環境構成にしようかな。
Dockerかローカルか
Cloud9か、vimか
Cloud9で良かろう。下手に会社にログが残る方が嫌かも。しかも、ソースコードの管理が煩雑だし、私的保存もできないやん。
アニメファンドと被るの嫌だから、新たにCloud9(EC2)を立てよう。
そんで、ローカルでも遊べるように、Docker-composeにしよう。
Dockerでやるとコンパイルした結果のバイナリファイルは、Docker用になるんだろうか。

## Cloud9 で環境を作った
`$ sudo curl -L "https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-`compose
`$ sudo chmod +x /usr/local/bin/docker-compose`
`$ docker image pull golang:1.15`
`$ docker-compose build`
`$ docker-compose up -d`
`$ docker-compose exec go_training bash`

## golang の練習
`# go run hello.go`
`# go fmt hello.go`
ボリュームしたままで、フォーマット整えることができる
`# go build hello.go`でコンパイルして、バイナリファイルが作れる。
`# ./hello`とバイナリファイルを実行することもできる。
`# go version`でどのOSか確認しておくこと
他のOS用のバイナリファイルを作ることもできる。それをクロスコンパイルという。

packageはディレクトリと同義っぽい。
>Goでは，1つのパッケージは1つのディレクトリに格納します。
そしてディレクトリを作らず、環境変数GOPATHも指定しなかったら、エラーになった。
myprojectディレクトリにgosample.go main.goを作った。
```
main.go:5:2: cannot find package "gosample" in any of:
    /usr/local/go/src/gosample (from $GOROOT)
    /go/src/gosample (from $GOPATH)
```
Golang開発のタスクランナーとしてMakefileを利用することが多いらしい。
でも、GOPATHを設定すると、のディレクトリの命名規則を用いて，Makefileなどの構成ファイルは一切なしで，依存関係を解決してビルドできるらしい。
1.1から、GOPATHを必ず設定するようになった。
golangのバージョンが1.11以上だと、GoModules（バージョン管理）を使用した場合の方法がディファクトスタンダードらしい?
むしろGOPATHを使わない?
時代の遷移があったっぽい。
相対パスimportはGoコマンドとGOPATHがまだ無い時代に利用されていた。
↓
GOPATHが採用された後、GOPATH内では絶対パス指定(fully-qualified path)推奨、GOPATH外では相対パス指定するしかない(しかしビルドキャッシュが無かったのでビルドが遅かった)。
↓
Go Modulesでは相対パスか絶対パスのどちらかしか選択できない状況となった。マジョリティをサポートするために絶対パスコードの互換性保証を選んだ。
_概念は獲得できるだろうが、最新版の記事を探した方がいいのか?_
あまり良さそうな記事がなかった。つまり書籍しかないのか?
みんなのGoを見たが、たしかに良書そうな雰囲気を感じた。


### 4番目の公式のチュートリアルをWebでやる
```
func main(){
        fmt.Println("HelloWorld")
}
```
Goのプログラムは、パッケージ( package )で構成される。
インポートパスが "math/rand" のパッケージは、 package rand ステートメントで始まるファイル群で構成する。
ステートメントとは、宣言のこと。
```
import (
        "fmt"
        "time"
)
import "fmt"
import "time"
```
importが複数になるときは、前者の書き方が推奨。
前者の書き方を、factoredインポートステートメントという。
import文が書いてあるのに、ソースコード内で使われてなかったらエラーになる。

Goでは、最初の文字が大文字で始まる名前は、外部のパッケージから参照できるエクスポート(公開)された名前( exported name )
math.piは動かないが、math.Piは動く
引数を渡すことができる。しかし、変数名の 後ろ に型名を書く
```
func add(x int, y int) int {
	return x + y
}
```
`x int, y int`を`x, y int`に省略できる
関数は複数の戻り値を返すことができる
```
func swap(x, y string) (string, string) {
	return y, x
}

func main() {
	a, b := swap("hello", "world")
	fmt.Println(a, b)
}
```
戻り値となる変数に名前をつける( named return value )ことができる。
戻り値に名前をつけると、関数の最初で定義した変数名として扱われる
```
func split(sum int) (x, y int) {
	x = sum * 4
	y = sum - x
	return
}
```
x, yが名前付き戻り値になっている。
戻り値の意味を示す名前とすることで、関数のドキュメントとして表現するように使用する
return ステートメントに何も書かずに戻すことができ、"naked" return という。
naked returnステートメントは、短い関数でのみ利用すべき
var ステートメントは変数( variable )を宣言する
`var i int`
var 宣言では、変数毎に初期化子( initializer )を与えることができる。
`var i, j int = 1, 2`
初期化子が与えられている場合、型を省略でき、自動で型を決めてくれる。
関数の中では var 宣言の代わりに := の代入文を使い、暗黙的な型宣言ができる。
関数の外では使用できない。

int, uint, uintptr 型は、OSのbit数に合わせて調整してくれるっぽい。
特別な理由がない限り、整数の変数はintを使うべき。
Goでは文字そのものを表すためにruneという言葉を使う
```
var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)
```
変数もまとめて宣言が可能
型は`bool, string, int, uint, byte, rune(Javaのchar), float32, float64, complex64, complex128`
変数に初期値を与えずに宣言すると、ゼロ値( zero value )が与えられる
バッククオートで囲むことで，複数行に渡るstring（ヒアドキュメント）を記述できる

`fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)`
%Tは型で、%vが値っぽい
型変換では、変数 v 、型 T があった場合、 T(v) は、変数 v を T 型へ変換
右側の変数が型を持っている場合、左側の新しい変数は同じ型になる
```
var i int
j := i // j is an int
```
右側が型を指定しない数値である場合、左側の新しい変数は右側の定数の精度に基いて int, float64, complex128 の型になる
```
i := 42           // int
f := 3.142        // float64
g := 0.867 + 0.5i // complex128
```
定数( constant )は、 const キーワードを使って変数と同じように宣言する。
定数は、文字(character)、文字列(string)、boolean、数値(numeric)のみで使える
:= は使えない。
`const Pi = 3.14`
定数の頭文字は大文字っぽい
型のない定数は、その状況によって必要な型を取る
for文
```
for i := 0; i < 10; i++ {
	sum += i
}
```
初期化と後処理ステートメントの記述は任意
```
sum := 1
for ; sum < 1000; {
        sum += sum
}
```
セミコロン(;)を省略することもできる。while文と等価になる。
```
for sum < 1000 {
        sum += sum
}
```
ループ条件を省略すれば、無限ループ( infinite loop )になる。
```
for {
}
```
if ステートメント
```
func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}
```
if ステートメントは、 for のように、条件の前に、評価するための簡単なステートメントを書くことができる
```
if v := math.Pow(x, n); v < lim {
        return v
}
```
if ステートメントで宣言された変数は、 else ブロック内でも使える
```
if v := math.Pow(x, n); v < lim {
    return v
} else if v == lim {
	return v*2
}
} else {
    fmt.Printf("%g >= %v\n", v, lim)
}
```
実行順の注意。fmt.Println は、2つの pow が先に実行されてから
```
func main() {
        fmt.Println(
                pow(3, 2, 10),
                pow(3, 3, 20),
        )
}
```
ざっくりと平方根を求める関数。ニュートン法が使われている。平方根で特に有効らしい。
```
func Sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2*z)
	}
	return z
}
```
switch ステートメントは if - else ステートメントのシーケンスを短く書く方法
選択された case だけを実行してそれに続く全ての case は実行されない。つまりbreakを書かなくて良い
上から下へcaseを評価する。
cateは変数もint型以外も使える
```
switch os := runtime.GOOS; os {
case "darwin":
        fmt.Println("OS X.")
case "linux":
        fmt.Println("Linux.")
default:
        fmt.Printf("%s.\n", os)
}
```
"if-then-else"のつながりをgolangでは、 switch true (条件のないswitch)と書く
```
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
}
```
defer ステートメントは、 defer へ渡した関数の実行を、呼び出し元の関数の終わり(returnする)まで遅延させるもの
```
func main() {
	defer fmt.Println("world")

	fmt.Println("hello")
}
hello worldと出力される
```
defer へ渡した関数が複数ある場合、その呼び出しはスタック( stack )される。
呼び出し元の関数がreturnするとき、 defer へ渡した関数は LIFO(last-in-first-out) の順番
```
for i := 0; i < 10; i++ {
        defer fmt.Println(i)
}
結果は9,8,7,6,5,4,3,2,1,0
```
Goはポインタを扱います。 ポインタは値のメモリアドレスを表す
変数 T のポインタは、 *T 型で、ゼロ値は nil
`var p *int`
& オペレータは、そのオペランド( operand )へのポインタを引き出す
`i := 42`からの`p = &i`
これは "dereferencing" または "indirecting" としてよく知られている
struct (構造体)は、フィールド( field )の集まり
```
type Vertex struct {
	X int
	Y int
}

func main() {
	fmt.Println(Vertex{1, 2})
}
出力は、{1 2}
```
structのフィールドは、ドット( . )を用いてアクセス
```
func main() {
	v := Vertex{1, 2}
	v.X = 4
        fmt.Println(v.X)
}
出力は4
```
structのフィールドは、structのポインタを通してアクセスすることもできる
`(*p).X`と書けるが、面倒だから`p.X`と書ける
```
func main() {
	v := Vertex{1, 2}
	p := &v
	p.X = 1e9
	fmt.Println(v)
}
```
structリテラルは、フィールドの値を列挙することで新しいstructの初期値の割り当て
インスタンスが作られる感じ?でもgolangとオブジェクト指向は違うって聞いたことある。
X: 1 として X だけを初期化することもできる
`v2 = Vertex{X: 1}`
& を頭に付けると、新しく割り当てられたstructへのポインタを戻す
`p  = &Vertex{1, 2}`出力は`&{1 2}`
[n]T 型は、型 T の n 個の変数の配列( array )を表す
固定長だから、個数の情報も含めて1つの型になる
関数に配列を渡す場合は値渡しとなり，配列のコピーが渡される
`var a [10]int`
配列の長さは、型の一部分。よって、配列のサイズを変えることはできない。
`fmt.Println(a[0], a[1])`とアクセスできる。
`primes := [6]int{2, 3, 5, 7, 11, 13}`と同時に代入することもできる
配列は固定長。スライスは可変長。スライスは配列よりもより一般的
型 []T は 型 T のスライスを表す
`var s []int = primes[1:4]`
最初の要素は含むが、最後の要素は除いた半開区間を選択。0,1,2,3,4で1,2,3となる
スライスは配列への参照のようなもの
スライスの要素を変更すると、その元となる配列の対応する要素が変更される
スライスのリテラルは長さのない配列リテラルのようなもの
`[]bool{true, true, false}`は、`[3]bool{true, true, false}`の配列リテラルを作成し、配列リテラルを参照するスライスを作成する
```
s := []struct {
        i int
        b bool
}{
        {2, true},
        {3, false},
        {5, true},
        {7, true},
        {11, false},
        {13, true},
}
```
構造体でもいけるっぽい
スライスするときは、上限または下限を省略できる。`a[0:]`とか
スライスは長さ( length )と容量( capacity )を持つ
スライスの長さは、それに含まれる要素の数。
スライスの容量は、スライスの最初の要素から数えて、元となる配列の最後まで要素数。
スライス s の長さと容量は len(s) と cap(s) という式を使用して得る
必ずlen <= capっぽい。
`s[左:右]`コロンの右だけ書いてあるパターンは、capが保存される。左だけ書いてあるパターンは、capは減る。
スライスのゼロ値は nil 。nil スライスは 0 の長さと容量を持っており、元となる配列を持っていない。
スライスは、組み込みの make 関数を使用して作成できる
これは、動的サイズの配列を作成する方法
make 関数はゼロ化された配列を割り当て、その配列を指すスライスを返す
`a := make([]int, 5)`は、len=5 cap=5 [0 0 0 0 0]
`b := make([]int, 0, 5)`は、len=0 cap=5 []
スライスは、他のスライスを含む任意の型を含むことができる。2次元配列的な。
```
board := [][]string{
        []string{"_", "_", "_"},
        []string{"_", "_", "_"},
        []string{"_", "_", "_"},
}
```
スライスへ新しい要素を追加するには、Goの組み込みの append を使う
`func append(s []T, vs ...T) []T`vs は、追加する T 型の変数群
`s = append(s, 2, 3, 4)`
lenの増加の法則は単純にvsの分が増えて分かるが、capの増加の法則が分からない。
for ループに利用する range は、スライスや、マップ( map )をひとつずつ反復処理するために使う。
スライスをrangeで繰り返す場合、rangeは反復毎に2つの変数を返す。
1つ目の変数はインデックス( index )で、2つ目はインデックスの場所の要素のコピー
```
pow := make([]int, 10)
for i, v := range pow {
        fmt.Printf("2**%d = %d\n", i, v)
}
```
インデックスや値は、 " _ "(アンダーバー) へ代入することで捨てることができる。
```
for i, _ := range pow
for _, value := range pow
```
もしインデックスだけが必要なのであれば、2つ目の値を省略できる。`for i := range pow`
map はキーと値とを関連付ける
マップのゼロ値は nil。 nil マップはキーを持っておらず、キーの追加もできない。
make 関数は初期化され使用できるようにした指定された型のマップを返す
```
type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

func main() {
	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{
		40.68433, -74.39967,
	}
        m["aaa"] = Vertex{
		1, 3,
	}
	fmt.Println(m["Bell Labs"])
}
```
```
var m = map[string]Vertex{
	"Bell Labs": Vertex{
		40.68433, -74.39967,
	},
	"Google": Vertex{
		37.42202, -122.08408,
	},
}
```
もし、mapに渡すトップレベルの型が単純な型名である場合は、リテラルの要素から推定できますので、その型名を省略することができる
```
var m = map[string]Vertex{
	"Bell Labs": {40.68433, -74.39967},
	"Google":    {37.42202, -122.08408},
}
```
map m の操作
m へ要素(elem)の挿入や更新:`m[key] = elem`
要素の取得:`elem = m[key]`
要素の削除:`delete(m, key)`
キーに対する要素が存在するかどうか:`elem, ok := m[key]`
もし、 m に key があれば、変数 ok は true となり、存在しなければ、 ok は false となります。
なお、mapに key が存在しない場合、 elem はmapの要素の型のゼロ値になる
```
m := make(map[string]int)
m["Answer"] = 42
delete(m, "Answer")
v, ok := m["Answer"]
```
計算量を無視した単語カウンターの自作
```
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
```
関数も変数。他の変数のように関数を渡すことができる。使う関数を引数で変える的な。pythonでもあったな。
```
hypot := func(x, y float64) float64 {
        return math.Sqrt(x*x + y*y)
}
func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}
fmt.Println(hypot(3, 4)) // 出力5 なぜなら3^2+4^2=5^2
fmt.Println(compute(hypot)) // 出力5 なぜなら渡す値はcomputeで既に決まっているから
fmt.Println(compute(math.Pow)) // 出力81 なぜなら3^4=81
```
Goの関数は クロージャ( closure )
adder 関数はクロージャを返す。クロージャ(関数閉包)とはプログラミング言語における関数オブジェクトの一種。ラムダ式や無名関数の概念
そして、各クロージャ(ここではfor文の中のpow)は、それ自身(adder)の sum 変数へバインドされます。
```
func adder() func(int) int {
	sum := 0
	return func(int) int {
		return sum
	}
}
func main() {
	pos:=adder()
	for i:=0; i<3; i++ {
		pos(i)
	}
}
```
`pos:=adder()`とやると、`sum:=0`までは実行されて、`return func(int) int`がposに入ってる感じ
だからpos(10)とかは動く。そしてfor文2回目のposとかは、for文1回目のposのsumが入っている。
```
func adder() func(int) int {
	sum := 5
	fmt.Println("iti", sum)
	return func(x int) int {
		fmt.Println("ni", sum)
		sum += x
		fmt.Println("san", x)
		fmt.Println("si", sum)
		return sum
	}
}

func main() {
	fmt.Println("when1")
	pos := adder()
	fmt.Println("when2")
	fmt.Println(pos(10))
	fmt.Println("when3")
	for i := 0; i < 3; i++ {
		fmt.Println("when0")
		fmt.Println(pos(i))
	}
}
```
クロージャを使ったフィボナッチ数列の自作
```
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
Goには、クラス( class )のしくみはありませんが、型にメソッド( method )を定義できる。
型に属するメソッド的な。親となる型?をレシーバという。
書き方は、メソッドに、レシーバを付ける。
付け方は、func キーワードとメソッド名の間に自身の引数リストで表す。
Abs メソッドは v という名前の Vertex 型のレシーバを持つ。型は任意。
```
type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
```
Abs1()もAbs2()も、どっちで書いても機能的には同じに見えるが、内部が違う。
レシーバを伴うメソッドの宣言は、レシーバ型が同じパッケージにある必要がある
```

type MyFloat float64

func (f MyFloat) Abs1() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}
func Abs2(f MyFloat) float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func main() {
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs1())
	fmt.Println(Abs2(f))
}

```
レシーバは、ポインタにすることもできる。
レシーバ自身を更新することが多いため、変数レシーバよりもポインタレシーバの方が一般的
```
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}
```
レシーバは、第0引数って感覚らしい
関数でポインタ→一致してないとダメ
レシーバでポインタ→第0引数には、値でもポインタでもok
関数で値→一致してないとダメ
レシーバで値→第0引数には、値でもポインタでもok
```
func (p Person) Greet(msg string) {}
func Person.Greet(p Person, msg string)
と同じこと
func (pp *Person) Shout(msg string)
func (*Person).Shout(pp *Person, msg string)
と同じこと
```
ポインタレシーバを使う理由
- メソッドがレシーバが指す先の変数を変更するため
- メソッドの呼び出し毎に変数のコピーを避けるため
構造体が大きいと、呼び出した時に全部コピーするのは処理が重くなるから。

変数レシーバというより、値レシーバという名前の方が一般的?
普通に引数で渡すと、コピーされたことになるから、メソッド内で値を変更しても、元のデータには影響しない。
非同期処理で安全とも言い変えることができる。
- メソッド呼び出しごとにレシーバの値はコピーされる
- 値レシーバの値はメソッド内で書き換えても元のレシーバの値にはまったく影響がない
- ポインタレシーバにするとメソッド内でレシーバの値を書き換えられる
- ポインタレシーバのメソッド内でフィールドを呼び出す場合には、その前に nil チェックをすべき
値レシーバにしたほうがよい場合
- int, string などのプリミティブ型をベースとする型
- map, chan といった参照型をベースとする型(最初からポインタだから)
- 不変型
- 小さい構造体

interface(インタフェース)型は、メソッドのシグニチャの集まりで定義
メソッドの集まりを実装した値を、interface型の変数へ持たせることができる
型にメソッドを実装していくことによって、インタフェースを実装する。
```
type I interface {
	M()
}
type T struct {
	S string
}
func (t T) M() {
	fmt.Println(t.S)
}
func main() {
	var i I = T{"hello"}
	i.M()
}
```
インターフェース自体の中にある具体的な値が nil の場合、メソッドは nil をレシーバーとして呼び出される
値が`<nil>`的な。フツーの言語だと、nilを参照しただけでエラーになるが。
nil インターフェースの値は、値も具体的な型も持たない。
nil インターフェースのメソッドを呼び出すと、ランタイムエラー
ゼロ個のメソッドを指定されたインターフェース型は、 空のインターフェース という。`interface{}`
空のインターフェースは、任意の型の値を保持できる
```
func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
func main() {
	var i interface{}
	describe(i)

	i = 42
	describe(i)

	i = "hello"
	describe(i)
}
```
型アサーション は、インターフェースの値の基になる具体的な値を利用する手段を提供する
アサーションとは、表明のこと。プログラムの前提条件を示すのに使う。
表明は、プログラムのその箇所で必ず真であるべき式の形式をとる
デバッグで使いそう。テストコードとはまた違うテストだな。てか単純に型チェックだわ。
`t := i.(T)`
インターフェースの値 i が具体的な型 T を保持し、基になる T の値を変数 t に代入する
`t, ok := i.(T)`
i が T を保持していれば、 t は基になる値になり、 ok は真(true)
1個しか値を返さないかと思えば、真偽値もちゃんと返すのね。
型switch はいくつかの型アサーションを直列に使用できる構造
型switchは通常のswitch文と似ている。
型switchのcaseは型を指定し、その値は指定されたインターフェースの値が保持する値の型と比較する
型switchを使うとき、型アサーションの部分をtypeにする
```
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}
```

もっともよく使われているinterfaceの一つに fmt パッケージ に定義されている Stringer がある。
JavaでいうtoStringメソッドな感じ?その型をprintlnした時のデフォルト表記的な。
```
type Stringer interface {
    String() string
}
```
```
func (i IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])
}
```

Goのプログラムは、エラーの状態を error 値で表す。
error 型は fmt.Stringer に似た組み込みのインタフェース
Goでは，関数が多値を返せることを利用して，内部で発生したエラーを戻り値で表現します。
```
type error interface {
    Error() string
}
```
```
func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
```
error がnilなら成功したことを示し、 error が nilでない(何かしらエラーメッセージを呼び出し元に返した場合)は失敗したことを示す。
関数の処理に成功した場合はエラーはnilにし，異常があった場合はエラーだけに値が入り，他方はゼロ値


io パッケージは、データストリームを読むことを表現する io.Reader インタフェースを規定している。
io.Reader インタフェースは Read メソッドを持つ`func (T) Read(b []byte) (n int, err error)`
よくあるパターンは、別の io.Reader をラップし、ストリームの内容を何らかの方法で変換するio.Reader
例えば、 gzip.NewReader は、 io.Reader (gzipされたデータストリーム)を引数で受け取り、 *gzip.Reader を返します。
その *gzip.Reader は、 io.Reader (展開したデータストリーム)を実装している

image パッケージは、以下の Image インタフェースを定義している
```
package image

type Image interface {
    ColorModel() color.Model
    Bounds() Rectangle
    At(x, y int) color.Color
}
```
goroutine (ゴルーチン)は、Goのランタイムに管理される軽量なスレッド
`go f(x, y, z)`と書けば、新しいgoroutineが実行する
f , x , y , z の評価自体は、実行元(current)のgoroutineで実行され、 fメソッド の実行は、新しいgoroutineで実行する
チャネル( Channel )型は、チャネルオペレータの `<-` を用いて値の送受信ができる通り道
(データは、矢印の方向に流れる。
```
ch <- v    // v をチャネル ch へ送信する
v := <-ch  // ch から受信した変数を v へ割り当てる
```
マップとスライスと同様に、チャネルは使う前に生成する必要がある。`ch := make(chan int)`
チャネルは、 バッファ ( buffer )として使える
バッファを持つチャネルを初期化するには、 make の２つ目の引数にバッファの長さを与える
`ch := make(chan int, 100)`
バッファが詰まった時は、チャネルへの送信をブロックする
 バッファが空の時には、チャネルの受信(チャネルから変数へ送りきる)をブロック
```
func main() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```
送り手は、これ以上の送信する値がない場合の状態を、チャネルを close するという。
`v, ok := <-ch`
受信する値がない、かつ、チャネルが閉じているなら、 ok の変数は、 false
実は、ループの for i := range c は、チャネルが閉じられるまで、チャネルから値を繰り返し受信し続けていた
チャネルは、通常closeする必要はない。closeするのは、これ以上値が来ないことを受け手が知る必要があるときにだけ
```
フィボナッチ数列
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
select ステートメントは、複数の通信操作のことで、goroutineを待たせる。つまり、複数のチャネルで受信を待てる。
```
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 4; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

```
一番最後の`fibonacci(c, quit)`が起動しつつも、別スレッドでfor文が回っていて、チャネルcに値が届いてたら、`fmt.Println(<-c)`をするのかな
どの case も準備ができていないのであれば、 select の中の default が実行される
```
func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
```
通信が必要ない。つまりコンフリクトを避けたい。つまり。一度に1つのgoroutineだけが変数にアクセスできるようにしたい場合は、
排他制御( mutual exclusion )と呼ばれ、このデータ構造を指す一般的な名前は mutex (ミューテックス)です。
Goの標準ライブラリは、排他制御をsync.Mutexと次の二つのメソッド(Lock, Unlock)で提供する。
Lock と Unlock で囲むことで排他制御で実行するコードを定義する
```
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}
```
mutexがUnlockされることを保証するために defer を使うこともできる
```
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return c.v[key]
}
```
Go言語のツアーを一通り終わらしたぞ!
これから、Go言語をバリバリ書いていこうかな!
トレンドだから楽しみだ!

Goでは多様な書き方を認めないことで，言語の仕様を小さく保つ方針
breakやcontinueは今でも使えるのかな
if/else文が繰り返す場合は，switch文を用いたほうがスッキリ書ける場合がある
Goのswitch文は非常に柔軟であり，値の比較だけでなく条件分岐にも使用できる
