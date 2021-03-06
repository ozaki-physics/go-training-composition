## ファイルの入出力を勉強する
[Go言語(golang) ファイルの読み書き、作成、存在確認、一行ずつ処理、コピー など](https://golang.hateblo.jp/entry/2018/11/09/163000)

- ファイルやディレクトリの存在を確認
`_, err := os.Stat(fileName)`
エラーもファイルが存在しないエラーと それ以外のエラーを使い分ける
- ファイルを開く
`f, err := os.Open(fileName)`
読み込み専用 ファイルがなければ作成されずエラー
開いたら必ず `defer r.Close()` する
defer は return するまで defer に渡した関数の実行を遅延させる
つまり エラーで止まっても 正常にメソッドの最後に到達しても必ず実行される
- ファイルを開く
フラグを渡しつつ 名前付きでファイルを開く
ほとんどのユーザーは Open() か Create() を使う
`func OpenFile(name string, flag int, perm FileMode) (*File, error)`
- ファイルを作成
`func Create(name string) (*File, error)`
読み書きの両方に対応
内部的には os.OpenFile が呼ばれるだけ
OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
O_RDWR: ファイルの読み込みと書き込み両方
O_CREATE: ファイルがなければ作成
O_TRUNC: ファイルを開くときに内容をすべて切り詰める
- ファイルから一度にすべてのデータを読み込む
`data, err := ioutil.ReadFile(fileName)`
内部的には go 1.16 より os.ReadFile() を呼び出すだけ
呼び出しが成功すると err == EOF じゃなく err == nil が返る
[公式ドキュメント](https://pkg.go.dev/io/ioutil@go1.17#ReadFile) の example を見ると Close は必要ないらしい
- 開いてるファイルから一度にすべてのデータを読み込む
`func ReadAll(r io.Reader) ([]byte, error)`
`data, err := ioutil.ReadAll(file)`
内部的には go 1.16 より io.ReadAll() を呼び出すだけ
呼び出しが成功すると err == EOF じゃなく err == nil が返る
[公式ドキュメント](https://pkg.go.dev/io/ioutil@go1.17#ReadFile) の example を見ると Close は必要ないらしい
引数は io.Reader と書いてあるが file *os.File を渡しても大丈夫だった
[io.Reader](https://pkg.go.dev/io@go1.17#Reader)を見ても 
```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```
で Read も見つけられないし
[os.File](https://pkg.go.dev/os@go1.17#File)を見ても
```go
type File struct {
	// contains filtered or unexported fields
}
```
でよく分からん
なぜなら *os.File は struct で io.Reader は interface だから
io.Reader メソッドを持っている *os.File は引数で渡せる

http.Response を読み出すときなどにもよく使うらしい
- ファイルのデータをコピーする
`_, err = io.Copy(w, r)`
`func Copy(dst Writer, src Reader) (written int64, err error)`
io.Copy の引数は io.Writer と io.Reader
つまり その要件を満たすもの(*os.File など)を渡す
要件とは io.Writer と io.Reader のように
それぞれの interface に定義してある関数を持っているか
- ファイル名の変更
`os.Rename(oldPath, newPath)`
ファイル名の変更 つまり 実質 path の移動
`func Rename(oldpath, newpath string) error`
newpath が存在して ディレクトリじゃない場合 すでに存在したファイルは削除され置き換えられる
エラーの時は os パッケージの中の *LinkError 型が返ってくる
- 何バイトずつ読むか指定してファイルを読み込む
`func (f *File) Read(b []byte) (n int, err error)`
引数に渡した len(b) のバイトまで読み込む
戻り値は 読んだバイト数とエラー
ファイルの終わりでは 0, io.EOF を返す
`n, err := file.Read(fileContent)`
- ファイルにデータを一気に書き込む
`func WriteFile(filename string, data []byte, perm fs.FileMode) error`
go 1.16 より 内部的には `os.WriteFile(filename, data, perm)` を呼んでるだけ
- ファイルを1行ずつ読み込む
`func NewScanner(r io.Reader) *Scanner`
Scan() は 行がある限り true を返す
Text() で一行ずつ取得できる
```go
	s := bufio.NewScanner(file)
	for s.Scan() {
	oneLine := s.Text()
```

[Goでテキストファイルを読み書きする時に使う標準パッケージ](https://qiita.com/qt-luigi/items/2c13ad68e7d9f8f8c0f2)
>バイト配列や文字列の単位で読み書きするなら os パッケージ
>バッファリングしながら読み書きするなら bufio パッケージ
>一括で読み書きするなら ioutil パッケージ


## 感想
*os.File をよく使うイメージ
1行ずつ扱うなら bufio かな

## 整理する
主に使ったパッケージ: bufio, io, os
特に os.File を圧倒的に使った
他にも
bufio.NewScanner
bufio.NewReader
bufio.NewWriter

os.File は struct
io.Reader は interface
os.File は io.Reader を持っているらしい

Fprint(w io.Writer, a ...interface{}) で 第1引数に os.File だったり bufio.Writer を渡してる
つまり os.File も bufio.Writer も io.Writer interface を持っている
io.Writer に Write() や WriteString() が定義されてるのかな

- 読み込み
  -　一気に
  `ioutil.ReadFile(string)` 実体は os.ReadFile(string)
  `ioutil.ReadAll(os.File)` 実体は io.ReadAll(io.Reader)
  - バイトずつ
  `os.File.Read([]byte)`
  - 1行ずつ
  `bufio.NewScanner(os.File).Scan()`
  `bufio.NewScanner(os.File).ReadString(byte)`
  `bufio.NewScanner(os.File).ReadLine()`
- 書き込み
  - os.File 系
  `os.File.WriteString(string)`
  `os.File.Write([]byte)`
  `fmt.Fprintln(os.File, interface{})`
  - bufio.Writer 系
  `bufio.Writer.WriteString(string)`
  `bufio.Writer.Write([]byte)`
  `fmt.Fprintln(bufio.Writer, interface{})`
  - 一気に
  `ioutil.WriteFile(string, []byte, fs.FileMode)` 実体は os.WriteFile(string, []byte, fs.FileMode)
- ファイルやディレクトリを探す
`os.Stat(string)`
- ファイルを開く
`os.Open(string)` 実体は os.File.OpenFile(name, O_RDONLY, 0)
`os.OpenFile(string, os.O_RDWR|os.O_CREATE, fs.FileMode)` 基本使わないで Open, Create を使う
- ファイルを作成する
`os.Create(strubg)` 実体は os.File.OpenFile(string, O_RDWR|O_CREATE|O_TRUNC, fs.FileMode)
- ファイルのコピー
`io.Copy(os.Create(string), os.Open(string))`
- ファイル名の変更
`os.Rename(oldPath, newPath string)`


[[1分学習]Goで標準入力を受け付ける](https://zenn.dev/tomatomato/articles/onemin_golang_stdin)
[Goで標準入力とファイル読み込みを可能にするインタフェース](https://medium.com/eureka-engineering/go-read-from-stdin-or-file-bb7a9197b904)
[[Go] ファイルや標準入力から一行ずつ読み込む](https://qiita.com/hnakamur/items/a53b701c8827fe4bfec7)
[Go言語(golang)のループ for, for..range, foreach, while](https://golang.hateblo.jp/entry/2019/10/07/171630)
[Go 言語で標準入力から読み込む競技プログラミングのアレ --- 改訂第二版](https://qiita.com/tnoda_/items/b503a72eac82862d30c6)
>簡単に書きたいとき
>`fmt.Scan` を使う
>たくさん (> 10^5) 読み込みたいとき
>`bufio の Scanner` を使う
>長い行を (> 64x10^3) 読み込みたいとき
>`bufio の ReadLine` を使う

[golang でコマンドライン引数を使う](https://qiita.com/nakaryooo/items/2d0befa2c1cf347800c3)
>- os パッケージの Args を使う
>- flag パッケージを使う
>flag を使うのが良さそう
>型指定やオプショナル引数、デフォルト値の設定等、色々できるから
[Go言語 - os.Argsでコマンドラインパラメータを受け取る](https://blog.y-yuki.net/entry/2017/04/30/000000)


よく使う暗号とか log とか ちゃんとパッケージ化しておいた方がいいかもな
てか main に色々書くの良くなかった笑
ちゃんと main から パッケージを呼び出して使うって形にしないとなぁ
テキトーにコード書いたから後でちゃんと整理したい

## 疑問
`OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)` とかメソッドの引数に パイプがあるが意味が分からない
だいたい or みたいな意味で使われる気はしてるが

fileName って僕は書いてるが go のライブラリ見てると filename って書いてる
僕が好んで file_name や fileName を使うのは 英語ネイティブじゃないから filename では 読めないだけ?
oldpath, newpath string も 全部小文字だった

golang って オーバーロードをどう思ってるのか
-> golang にはオーバーロードは無い コンパイルしたら怒られるっぽい
[名前違いの関数をいくつも提供する](https://future-architect.github.io/articles/20190713/)

file の EOF ってなんだ?
-> End Of File の略
-> ファイルの終端を表す特殊な記号やデータなどのこと

*os.File, *bufio.Writer の両方で
WriteString, Write メソッド使えるから
インタフェースとして持っている?
-> 調べたけど分からんかった
[os.File の WriteString](https://github.com/golang/go/blob/master/src/os/file.go)
```go
// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (f *File) WriteString(s string) (n int, err error) {
	var b []byte
	hdr := (*unsafeheader.Slice)(unsafe.Pointer(&b))
	hdr.Data = (*unsafeheader.String)(unsafe.Pointer(&s)).Data
	hdr.Cap = len(s)
	hdr.Len = len(s)
	return f.Write(b)
}

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (f *File) Write(b []byte) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	n, e := f.write(b)
	if n < 0 {
		n = 0
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	epipecheck(f, e)

	if e != nil {
		err = f.wrapErr("write", e)
	}

	return n, err
}
```
[bufio の WriteString](https://github.com/golang/go/blob/master/src/bufio/bufio.go)
```go
// WriteString writes a string.
// It returns the number of bytes written.
// If the count is less than len(s), it also returns an error explaining
// why the write is short.
func (b *Writer) WriteString(s string) (int, error) {
	nn := 0
	for len(s) > b.Available() && b.err == nil {
		n := copy(b.buf[b.n:], s)
		b.n += n
		nn += n
		s = s[n:]
		b.Flush()
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], s)
	b.n += n
	nn += n
	return nn, nil
}
```
