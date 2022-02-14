## Go のプロジェクト構成を調べる
go:1.15 には /go ディレクトリの下に /bin /src の2個がある  
/go の直下に /pkg は見当たらない  
`go env`で go で使うような path が全部確認できる  

[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)を 参考にする人が多いが 細かすぎるのでは? という意見もある  

[標準パッケージの fmt](https://github.com/golang/go/tree/master/src/fmt)では 1つのパッケージの直下に .go をめっちゃ置く  

[Go 言語のプロジェクト構成](https://blog.tokoyax.com/entry/go/project)  
testdata や _ で始まるディレクトリは Go のパッケージとはみなされないから Go 以外のコードやファイルを置くといいらしい  

## 定義
この git リポジトリを
docker 内の`/go/src/github.com/ozaki-physics/go-training-composition`にマウントした
モジュール名 github.com/ozaki-physics/go-training-composition
パッケージ名 package01
パッケージ名 package02

1つのリポジトリには 1個のモジュールしか入れられない
他のモジュールを使いたい場合は go mod を使って /pkg とかにインストールするらしい
パッケージは複数作って良い

小文字で1単語(単数形)が推奨
キャメルケースやスネークケースは使わない

ファイル名 は 有名OSS や 標準パッケージ でスネークケースを使うことが多いらしい
個人的には パッケージ名 もスネークケースになる気がするのだが...笑

ディレクトリ名はケバブケースを使っているらしい

## ビルドについて
`go build` は go.mod があるディレクトリじゃないとできない
-o フラグで build したバイナリをどこに保存するか決められる(相対パスも絶対パスも使える)
ビルドするファイル名を指定することもできる 指定したファイル名がバイナリファイル名になる
指定しないとプロジェクトのディレクトリ名が バイナリファイル名になる

以下は具体例
/go で `go build` したら
```bash
go: go.mod file not found in current directory or any parent directory; see 'go help modules'
```
go.mod が見つからないって言われたから
/go/src/github.com/ozaki-physics/go-training-composition で `go build` したら
/go/src/github.com/ozaki-physics/go-training-composition に `go-training-composition` というバイナリができた
/go/src/github.com/ozaki-physics/go-training-composition で `go build main.go` したときと同じバイナリ

/go/src/github.com/ozaki-physics/go-training-composition で
`go install github.com/ozaki-physics/go-training-composition` ってやったら
/go/bin に go-training-composition のバイナリが生成された
そして使えた
だけど build しないで install するのは良くない気がしている

/go/bin は `go install` のための場所で 自分のプロジェクトの build したバイナリを置く場所じゃないかも
でも /go/src/github.com/ozaki-physics/go-training-composition の中に /bin を作って
ビルドしたら /go-training-composition/bin の中に入れるのも違う気がする
だって マウントするときに /go-training-composition ごと上書きされてしまうから
じゃあ ソースコードを全部 /go-training-composition/src に入れて /go-training-composition/src だけマウントする?
たしかにできるが プロジェクトの中に /bin と /src が存在していいのか?
-> どーせ build した時に /go-training-composition/bin にバイナリファイルを置く指定をするなら /go/bin を流用する

/go/src/github.com/ozaki-physics/go-training-composition で
`go build -o /go/bin` ってやったら
/go/bin に go-training-composition のバイナリが生成された
そして使えた
`go build -o /go/bin main.go` ってやったら
/go/bin に main のバイナリが生成された
`go build -o ../../../../bin` ってやったら
/go/bin に go-training-composition のバイナリが生成された
相対パスでも書けるが 絶対の方が見やすい

## Go での package の扱いを理解する
よく見る`package main`は何を意味しているのか?
他の go ファイルを import する方法

Go のプログラムはパッケージ(package)で構成されている
プログラムは main パッケージから開始される
プログラム実行時の処理開始位置の main パッケージあるいは main 関数を エントリポイント という
インポートパスが "math/rand" のパッケージは package rand ステートメントで始まるファイル群で構成される
main パッケージが大きくなってきた時、処理を複数の関数あるいはファイルに分割したい。
その時 import 文では 新しく作った独自パッケージへのパスとして $GOPATH/src 直下のディレクトリから パッケージファイル(*.go)の直上のディレクトリまでを指定する

Go Modules のとき importするパスは 自分のプロジェクトルートディレクトリ直下のディレクトリ名から importしたいパッケージの直上のディレクトリまで

Go Modules(vgo)を使うならばimportにて相対パスは使用不可なので絶対パスを指定しよう
import 自体は相対 path でも書けるが Go の初期が絶対 path しか使えなかったため 慣例的に絶対 path を使うらしい

わかったこと
- 対象ディレクトリ配下に複数のmainパッケージのモジュールを置く
→ エラー
`package main`のファイルは 1プロジェクトに1個っぽい
main にかけなくなった分のコードが外に出される
- 異なる名前のパッケージ宣言を持つモジュールを同じディレクトリの中に置く
→ エラー
1つのディレクトリの中にあるファイルの package ステートメントが異なっていたどき(片方 package hello01, 片方 package hello02)のとき コンパイルエラーになる

ディレクトリ名は goodnight でも`package hello`にすることができる
→ ディレクトリ名とパッケージ名は同じにしなくてよい
そのとき main ファイルで import するとき 記述は ディレクトリ名だが 実際に 〇〇.△△() で使うときは hello.△△() になる
混乱のもとになるので ディレクトリ名と同じ package 名にする
または
```go
import goodnight "github.com/〇〇/goodnight"
```
と import する package に別名, 別のラベルを付けると goodnight.△△() で使えるようになる

Go Modules の魅力は /go/src 以外のディレクトリでも go get して 異なるディレクトリにダウンロードできることかもしれない
始めは go mod init したら記録される path が基本の場所で そこから import の path を書いたらいいと思っていたけどそうじゃないみたい
package のバージョンやなんの package を使っているか管理できるだけで 基本は src から path 書かないといけないっぽい
あと これ以上は公式リファレンスを読まないといけないくて 時間対進捗が良くないと思われる
とりあえず Go Modules についてはここらへんにして 書籍が終わるまでは Go Modules 下(実質 GOPATH)で開発する

## 疑問
go.mod にどういう意味で記述されたのか?
go.sum の意味は?
どこに インストールされたのか?

イメージを作り直さないと up -d するたびに go get される

test コードを書くときはどういうディレクトリ構成にすればいいのか

### import の書き方
import のときに ドットを使って path を書くと 使うときにクラス名を省略して書ける
```go
import (
  . "fmt"
)
func main() {
  Println("hello")
}
```

使ってない パッケージ を import に書くと build で怒られるが アンダースコア を使うと怒られなくなる


[GitHubリポジトリ作成時の定形作業をTemplate Repositoryで省力化する](https://devblog.thebase.in/entry/2020/06/23/131444)
を見てると 自分のディレクトリ構成が合っているか不安になってくる
特に docker の中で /go/src/github.com/ozaki-physics/go-training-composition がカレントディレクトリになるのが気になる

`go fmt ./...` パッケージ全部に対して fmt をする
