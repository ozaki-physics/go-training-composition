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
