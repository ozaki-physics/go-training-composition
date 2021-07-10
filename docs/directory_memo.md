## Go のプロジェクト構成を調べる
go:1.15 には /go ディレクトリの下に /bin /src の2個がある
/go の直下に /pkg は見当たらない

go.mod はどの path に作るのがいいのか?
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
