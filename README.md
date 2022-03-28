# go-training-composition
Go 言語で様々なことをやってみる練習リポジトリ

## 目的 Overview
### Go でアプリを作るための以下の練習

- それなりに練習した  
  - [暗号理論(crypto パッケージ)の勉強](./docs/cryptography_memo.md): [サンプル](./trainingCrypto)  
  - [ファイル I/O の勉強](./docs/ioFile_memo.md): [サンプル](./trainingIo/ioFile.go)  
    (おかげで 公式ドキュメントを読む力がついたと思う)  
  - ターミナル I/O の勉強: [サンプル](./trainingIo/ioTerminal.go)  
  - レイヤードアーキテクチャの勉強: [サンプル01](./ddd01), [サンプル02](./ddd02)  
  - [Go 言語でのテスト方法](./docs/test_memo.md): [サンプル](./pkg03), [テスト駆動開発 書籍の要約](./docs/tdd_summary.md), [サンプル](./trainingTest)  


- ふつうに練習した  
  - [Go Tour](./docs/go_tour.md): [サンプル](./goTour)  
  - [シンプルなディレクトリ構成](./docs/directory_memo.md)  
  - [外部パッケージのインストール](./docs/go_module.md)  
  - [JSON パッケージの勉強](./docs/json_memo.md): [サンプル](./trainingJson)  
  - [VS Code in Container の設定](./.devcontainer/devcontainer.json)  
  - [GitHub の issue, pull Request の テンプレート作成](./.github)  
  - GitHub でリポジトリにタグをつける方法と リリースノートの自動生成  
  - [http リクエスト の GET と POST](./docs/http_memo.md): [サンプル](./trainingWebScraping)  

- ちょっとだけ練習した  
  - [Log パッケージの勉強](./docs/err_memo.md): [サンプル](./utils/util.go)  
  - [TimeZone パッケージの勉強](./trainingTimeZone)  
  - [godoc の書き方](./docs/godoc_memo.md)  
  - [Go 言語の仕様](./docs/effective_go.md)  
  - [データ圧縮](./docs/compress_memo.md): [サンプル](./trainingCompress)  

### 作ったもの

- [テキストファイルの中身を暗号化するツール](./fileCrypto/use.go)  
  - パスワード, 入力ファイルの path, 出力ファイルの path は JSON ファイルに書いて 読み込ませる  
  - ファイル, ターミナル I/O と crypto パッケージ  
- [PayPal の Sandbox 環境へリクエスト](./trainingWebScraping/paypal.go): [API ドキュメント トップ(外部リンク)](https://developer.paypal.com/home/), [API ドキュメント(外部リンク)](https://developer.paypal.com/docs/checkout/advanced/integrate)  
  - ClientID と Secret を JSON ファイル から取得する  
  - ClientID と Secret で BASIC 認証を通して Access Token を取得する  
  - Access Token をリクエストヘッダーに含めて Client Token を取得する  
- [GMO コインへリクエストして 暗号資産のレートを取得する](./trainingWebScraping/gmoCoin.go): [API ドキュメント トップ(外部リンク)](https://api.coin.z.com/docs/#outline), [API ドキュメント(外部リンク)](https://api.coin.z.com/docs/#ticker)  
  - API 叩いて シンボル(BTC, ETH など12種)の情報(価格など)を JSON で取得する  
- [CoinMarketCap へリクエストして 暗号資産のレートを取得する](./docs/CoinMarketCap_api_memo.md): [サンプル](./requestCoinMarketCap), [API ドキュメント トップ(外部リンク)](https://coinmarketcap.com/api/), [API ドキュメント(外部リンク)](https://coinmarketcap.com/api/documentation/v1#section/Standards-and-Conventions)  
  - API 叩いて シンボル(基本なんでも)の情報(価格など)を JSON で取得して構造体に格納する  


## インストール方法 Install

## 環境 Requirement
- Docker

## 使い方 Usage
VS Code の拡張機能 Remote - Containers(識別子: ms-vscode-remote.remote-containers) を使って開発する  
コンテナ内で VS Code を起動し go 言語のための VS Code の拡張機能 Go(識別子: golang.go) を使う  
[golang.go](https://marketplace.visualstudio.com/items?itemName=golang.Go)  

golang.go は 様々なモジュールをインストールする必要がある  
そのモジュールを Docker image に含めないようにするため  
コンテナを起動した後の コンテナ内 VS Code で変更を加える  
VS Code の通知より install All をする  

コンテナに変更を加えたため 基本は コンテナ削除をしない  
(削除した場合は 再度 install All をすればいいだけ)  

コンテナ内の git では日本語が使えないため コミットするときは ローカルの git bash 等を使う  
```bash
$ docker-compose build
$ docker-compose up -d
# VS Code より Remote - Containers で接続する
# たまに .devcontainer\devcontainer.json の差分を検知して rebuild するような通知が来る
# その時は docker image も作り直されて 古い方の image が <none> になるため削除する

# VS Code より Remote - Containers で接続する
# VS Code の通知(golang.go)より install All をする

# 基本は VS Code 内のターミナルで良いが ローカルの PowerShell からアクセスしたくなった場合
$ docker-compose exec go_training bash

# 終えるとき
# VS Code より Remote - Containers で接続をやめる
$ docker-compose stop
# 再開するとき
$ docker-compose start
# VS Code より Remote - Containers で接続する
```

ちなみに image の時点で go build は済んでおり  
image から直接 run または docker-compose.yml の command をコメントアウトで確認できる  
```bash
$ docker container run --rm -d -p 8080:8080 --name check_go_training go1.17:training_composition_vscode_in_container
$ docker container stop check_go_training
```

```bash
/go/src/github.com/ozaki-physics/go-training-composition# go mod init $REPOSITORY
```

### 外部モジュールのバージョンアップ
例として github.com/gin-gonic/gin をバージョンアップする  
1. コンテナにアタッチする
2. モジュールのバージョンアップして go.mod を更新する
3. コンテナを削除してもバージョンアップが反映されるように docker image を作り直す

```bash
$ docker-compose up -d
$ docker-compose exec go_training bash

# モジュールのバージョンアップ
/go/src/github.com/ozaki-physics/go-training-composition# go get -d -v -u github.com/gin-gonic/gin
# 不要モジュールの削除
/go/src/github.com/ozaki-physics/go-training-composition# go mod tidy -v

$ docker-compose down
$ docker image rm go1.17:training_composition_vscode_in_container
# docker image の作り直し
$ docker-compose build
```

### golang.go でインストールされるモジュール
- Installing github.com/uudashr/gopkgs/v2/cmd/gopkgs (/go/bin/gopkgs) SUCCEEDED  
[gopkgs](https://github.com/uudashr/gopkgs)  
インポートできるパッケージのリストを表示するツール  
- Installing github.com/ramya-rao-a/go-outline (/go/bin/go-outline) SUCCEEDED  
[Go Outline](https://github.com/ramya-rao-a/go-outline)  
JSON 表現を抽出するためのシンプルなユーティリティ  
- Installing github.com/cweill/gotests/gotests (/go/bin/gotests) SUCCEEDED  
[gotests](https://github.com/cweill/gotests)  
テスト生成ツール  
- Installing github.com/fatih/gomodifytags (/go/bin/gomodifytags) SUCCEEDED  
[gomodifytags](https://github.com/fatih/gomodifytags)  
golang の struct に タグを追加したり更新したりする  
- Installing github.com/josharian/impl (/go/bin/impl) SUCCEEDED  
[impl](https://github.com/josharian/impl)  
インターフェースを実装するためのメソッドスタブを生成  
- Installing github.com/haya14busa/goplay/cmd/goplay (/go/bin/goplay) SUCCEEDED  
[goplay - The Go Playground Client](https://github.com/haya14busa/goplay)  
The Go Playground にコードを貼り付けつつ Web ページへ遷移する  
- Installing github.com/go-delve/delve/cmd/dlv (/go/bin/dlv) SUCCEEDED  
[delve](https://github.com/go-delve/delve)  
デバッガ  
- Installing github.com/go-delve/delve/cmd/dlv@master (/go/bin/dlv-dap) SUCCEEDED  
[delve](https://github.com/go-delve/delve)  
デバッガ  
- Installing honnef.co/go/tools/cmd/staticcheck (/go/bin/staticcheck) SUCCEEDED  
[staticcheck](https://pkg.go.dev/honnef.co/go/tools/staticcheck)  
リンター linter  
- Installing golang.org/x/tools/gopls (/go/bin/gopls) SUCCEEDED  
[gopls, the Go language server](https://pkg.go.dev/golang.org/x/tools/gopls)  
Go チームによって開発された公式の Go 言語サーバーです  
LSP 互換のエディターに IDE 機能を提供します  
コードの自動補完補完ツールらしい  

## 参考文献 References
[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)  
[Go の公式 github](https://github.com/golang/go)
