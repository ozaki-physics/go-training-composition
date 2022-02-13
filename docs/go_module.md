# go 言語での外部ライブラリの扱い
## まとめ
開発を始めるときは ディレクトリを github.com/アカウント名/リポジトリ名 で作る  
リポジトリ名のディレクトリに移動して `$ go mod init github.com/アカウント名/リポジトリ名` を実行する  
モジュールを追加するときは `go install <package>@<version>`  
go.mod を編集したいときは `go mod tidy`  

## ライブラリの管理
go で パッケージの依存関係を管理する方法は go.mod ファイル  

go 1.11 以前の GOPATH モードしか無いときは `/go/src` 配下で書いた go のコードしか実行できなかった  
でも Go module モードは `/go/src` 以外のディレクトリに コードを書いても go.mod さえ生成しておけば実行できるようになったらしい  

開発を始めるときは 最初に `$ go mod init パス` を実行する  
パスは / つまり root からのパスと思われる  
このパスは慣例的に `github.com/アカウント名/リポジトリ名` になっていることが多い  
なぜなら go.mod で `module github.com/ozaki-physics/go-training-composition` にしたいから  
なぜ これにしたいかというと 他人がパッケージをインストールするときは github.com からインストールする感覚だからだと思われる  

Go module モードのおかけで `/go/src` の配下に go のコードを書かなくても良くなったが  
習慣のように `/go/src/github.com/アカウント名/リポジトリ名` に go のコードを書く人が多いらしい  

### もし GOPATH モード と Go module モード を切り替える
go の環境変数 GO111MODULE で変更できる  
- off だと GOPATH mode  
- on だと module-aware mode  
- auto だと `$GOPATH/src` の外に対象のリポジトリがある かつ go.mod が存在する場合は module-aware mode, そうでない場合は GOPATH mode  
[Golang - Go Modulesで開発環境の用意する](https://qiita.com/so-heee/items/56f5317b42cec3d94383)  

### モジュールをインストールするとき
[Go1.17における go get の変更点](https://future-architect.github.io/articles/20210818a/)  
>コマンドのインストールで go get を使うのは非推奨  
>コマンドのインストールは以下のように go install を使いましょう  
ビルドしてインストールする機能が `go install` と重複するから `go get` は非推奨になったらしい  

`go get` だと go.mod に変更が入り go.sum も生成される  
だけど `go install` だと go.mod には変更が入らない  
`go install` が推奨にも関わらず go.mod に変更が反映されないのはよいのか?  
[go mod完全に理解した](https://zenn.dev/optimisuke/articles/105feac3f8e726830f8c) より  
ソースコードの import に書いていたら `go mod tidy` で go.mod にも反映されるらしい  
つまり 毎回 `go get` で意識的に go.mod を変更しなくても `go install` でモジュールをインストールして使わなければ go.mod には反映されないという便利さになった?  

[Go1.16からの go get と go install について](https://qiita.com/eihigh/items/9fe52804610a8c4b7e41) より  
>「バイナリのビルドとインストールのための go install」、「go.mod 編集のための go get」と役割が整理されていく  
>go build や go test で自動的に go.mod が更新されることがなくなりました。go.mod の編集は go get , go mod tidy , あるいは手作業で行います。  
今までは ツールのインストールをしたいなら `go get` が使われてきたらしい  
`go install <package>@<version>` をすると Module 配下に関係なく `$GOPATH/bin` にインストールされるらしい  
go.mod を編集したいなら 4パターンあるが 4パターンとも __`go mod tidy` で包括できる__  
1. コードに import を書いてから `go get` (新規モジュールの追加のみ)
2. コードに import を書いてから `go mod tidy` (新規モジュールの追加と不要モジュールの削除)
3. コードに import を書いてないから `go get <package>[@<version>]`
4. コードに import を書いてないから go.mod を手作業で編集

[cmd/go: 'go install' should install executables in module mode outside a module #40276](https://github.com/golang/go/issues/40276)
>Ideally, we would have one command (go install) for installing executables and one command (go get) for changing dependencies.  
>理想的には、実行ファイルをインストールするための1つのコマンド(go install)と、依存関係を変更するための1つのコマンド(go get)を私たちは使う  

[Goプロジェクトのはじめかたとおすすめライブラリ8.5選。ひな形にも使えるサンプルもあるよ。](https://qiita.com/yagi_eng/items/65cd812107362d36ae86) によると  
>各ソースコードの import 文にライブラリのパスを記載しておくと go run の時に自動でインストールしてくれます。  
これはされなくなった可能性があるけど まぁ深く考えるのはやめる  

### 公式による説明 [Using Go Modules](https://go.dev/blog/using-go-modules)  
>go mod init creates a new module, initializing the go.mod file that describes it.  
>`go mod init` は新しいモジュールを作成し、そのモジュールを記述した go.mod ファイルを初期化します。  
>go build, go test, and other package-building commands add new dependencies to go.mod as needed.  
>`go build` , `go test` , その他のパッケージ構築コマンドは、必要に応じて新しい依存関係を go.mod に追加します。  
>go list -m all prints the current module’s dependencies.  
>`go list -m all` は、現在のモジュールの依存関係を表示します。  
>go get changes the required version of a dependency (or adds a new dependency).  
>`go get` は依存関係の必要なバージョンを変更します (または新しい依存関係を追加します)。  
>go mod tidy removes unused dependencies.  
>`go mod tidy` は使われていない依存関係を削除します。  

## モジュールとパッケージは概念が別
モジュールは 1個の go.mod で構成される  
パッケージは 1個のモジュールの中に複数作ることができる  

## その他
### `go mod tidy` のオブション
[Goメモ-161 (go.mod の 内容を Go 1.17 に調整する)](https://devlights.hatenablog.com/entry/2021/09/02/112354)  
>1.17 から go mod tidy に -go フラグが追加された  
`$ go mod tidy -go=1.17`  

### 参考文献
[Tutorial: Get started with Go](https://go.dev/doc/tutorial/getting-started)  
[Go Modules Reference](https://go.dev/ref/mod)  
公式の[Go Modules](https://github.com/golang/go/wiki/Modules)  

## 外部モジュールをインストールしたときの記録
[Gin の公式サイトの Quickstart](https://gin-gonic.com/docs/quickstart/)  
`$ go get -u github.com/gin-gonic/gin`  
今後は `$ go get -d -v -u パッケージ名` が良さそう  
__go get は使われない方法になるから 今なら `go install` で `go get tidy` かな__  

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
バイナリをインストールするのがダメなだけで go.mod を編集するコマンドは `go get -d` のままみたい  

開発環境など 環境的なパッケージは `go install` で使えるようにして  
プロダクトで使うパッケージは go.mod に書くために `go get -d` を使えば良さそう  
__今では わざわざ `go get -d` しなくても `go get tidy` で良い__  

### go get コマンド
golang では -なんとか を フラグという  
-v フラグ は 詳細な進行状況とデバッグ出力を有効にする, 削除されたモジュールを出力する  
-u フラグ は マイナーリリースなど依存パッケージまでネットワークから更新する 使用しているパッケージの更新にも使える  
-d フラグ は ソースをダウンロードし ビルド(インストール)はされないらしい  
-go フラグ は go のバージョン変更を go.mod に反映してくれるらしい(`go mod tidy -v -go=1.17`)  
そもそも `go mod tidy` で 使われていない依存モジュールを削除から 1.15 -> 1.17 に移行するときにも使った  
詳しくは `go help mod tidy`  

[Go1.17における go get の変更点](https://future-architect.github.io/articles/20210818a/)  
[Goメモ-161 (go.mod の 内容を Go 1.17 に調整する)](https://devlights.hatenablog.com/entry/2021/09/02/112354)  

### go install の挙動
path: /go/src/github.com/ozaki-physics/go-training-composition で 
`go install github.com/ramya-rao-a/go-outline` をやろうとしたら  
```bash
no required module provides package github.com/ramya-rao-a/go-outline; to add it:
        go get github.com/ramya-rao-a/go-outline
```
path: / で 
`go install github.com/ramya-rao-a/go-outline` をやろうとしたら  
```bash
go install: version is required when current directory is not in a module
        Try 'go install github.com/ramya-rao-a/go-outline@latest' to install the latest version
```
path: / で 
`go install github.com/ramya-rao-a/go-outline@latest` を実行したら  
/go/bin に go-outline が追加されて go.mod は編集されなかった  
__go install では go.mod は変更されないから合っている__  
