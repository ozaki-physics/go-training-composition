## Go のプロジェクト構成を調べる
go:1.15 には /go ディレクトリの下に /bin /src の2個がある
/go の直下に /pkg は見当たらない

go.mod はどの path に作るのがいいのか?

`go env`で go で使うような path が全部確認できる

### 記事を読んで
パッケージを公開することを考えて src/github.com/ozaki-physics/go-training-chat にするらしい
github の リポジトリ作ると /github.com/ozaki-physics/リポジトリ名 になる
go のパッケージ思想的に 外部パッケージをインストールして使うとき github から取ってくるから
go.mod は パッケージに1個? src に1個?
リポジトリに1個っぽい? つまり1パッケージに1個?
すると クリーンアーキテクチャではレイヤーごとにパッケージにするから go.mod 多くなるけどいいの?
ってなるから たぶんリポジトリというかアプリあたり1個だろう

[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)を
参考にする人が多いが 細かすぎるのでは? という意見もある
[標準パッケージの fmt](https://github.com/golang/go/tree/master/src/fmt)では
1つのパッケージの直下に .go をめっちゃ置く
[Go 言語のプロジェクト構成](https://blog.tokoyax.com/entry/go/project)
testdata や _ で始まるディレクトリは Go のパッケージとはみなされないから Go 以外のコードやファイルを置くといいらしい
[Golang - Go Modulesで開発環境の用意する](https://qiita.com/so-heee/items/56f5317b42cec3d94383)
公式の[Go Modules](https://github.com/golang/go/wiki/Modules)より
`$ go mod init github.com/my/repo`
[Goプロジェクトのはじめかたとおすすめライブラリ8.5選。ひな形にも使えるサンプルもあるよ。](https://qiita.com/yagi_eng/items/65cd812107362d36ae86)
>各ソースコードの import 文にライブラリのパスを記載しておくと go run の時に自動でインストールしてくれます。
なるほど だから import は github.com/ って書いていくのか
この記事の モジュール名と github のリポジトリ名は同義かな

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


## 試しに外部モジュールを インストールしてみる
[Gin の公式サイトの Quickstart](https://gin-gonic.com/docs/quickstart/)
`$ go get -u github.com/gin-gonic/gin`
今後は `$ go get -d -v -u パッケージ名` が良さそう

[Go1.16からの go get と go install について](https://qiita.com/eihigh/items/9fe52804610a8c4b7e41)
>Go Module追加編集のためのgo get、ツールなどのバイナリインストールのgo installと住み分けることができそう
go 1.16 からは go install と go get の役割が整理されていくらしい
バイナリのビルドとインストールのための go install
go.mod 編集のための go get

2021/10/09 go:1.17 で go get 使ってインストールしようとしたら非推奨だと言われた
```bash
go get: installing executables with 'go get' in module mode is deprecated.
        To adjust and download dependencies of the current module, use 'go get -d'.
        To install using requirements of the current module, use 'go install'.
        To install ignoring the current module, use 'go install' with a version,
        like 'go install example.com/cmd@latest'.
        For more information, see https://golang.org/doc/go-get-install-deprecation
        or run 'go help get' or 'go help install'.
go get: モジュールモードで 'go get' を使って実行ファイルをインストールすることは非推奨です。
        現在のモジュールの依存関係を調整してダウンロードするには、'go get -d'を使ってください。
        現在のモジュールの要件を使ってインストールするには、'go install'を使ってください。
        現在のモジュールを無視してインストールするには、'go install'にバージョンを指定して
        'go install example.com/cmd@latest' のようにしてください
        詳細については、https://golang.org/doc/go-get-install-deprecation
        または 'go help get' や 'go help install' を実行してください。
```
go では インストールとダウンロードは別の概念
バイナリをインストールするのがダメなだけで go.mod を編集するコマンドは go get -d のままみたい

開発環境など 環境的なパッケージは `go install` で使えるようにし
プロダクトで使うパッケージは go.mod に書くために `go get -d` を使えば良さそう

### go get コマンド
golang では -なんとか を フラグという
-v フラグ は 詳細な進行状況とデバッグ出力を有効にする
-u フラグ は マイナーリリースなど依存パッケージまでネットワークから更新する 使用しているパッケージの更新にも使える
[Go1.17における go get の変更点](https://future-architect.github.io/articles/20210818a/)
-d フラグ は ソースをダウンロードし ビルド(インストール)はされないっぽい

### go mod tidy
`go mod tidy` で、使われていない依存モジュールを削除できるから
1.15 -> 1.17 に移行した

[Goメモ-161 (go.mod の 内容を Go 1.17 に調整する)](https://devlights.hatenablog.com/entry/2021/09/02/112354)
`go mod tidy -go=1.17` と -go フラグを付けると go のバージョン変更を go.mod に反映してくれるらしい

-v フラグ は 削除されたモジュールを出力する
詳しくは `go help mod tidy`

### go install の挙動
/go/src/github.com/ozaki-physics/go-training-composition で
`go install github.com/ramya-rao-a/go-outline` をやろうとしたら
```bash
no required module provides package github.com/ramya-rao-a/go-outline; to add it:
        go get github.com/ramya-rao-a/go-outline
```
と言われて
/ で
`go install github.com/ramya-rao-a/go-outline` をやろうとしたら
```bash
go install: version is required when current directory is not in a module
        Try 'go install github.com/ramya-rao-a/go-outline@latest' to install the latest version
```
って言われて
`go install github.com/ramya-rao-a/go-outline@latest` を実行したら
/go/bin に go-outline が追加されて go.mod は編集されなかった

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
