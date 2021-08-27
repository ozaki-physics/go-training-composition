# Go でのドキュメントの書き方
そもそもコードを読んだだけで理解できるような書き方をすべきだけど ドキュメントを書きたくなる
ただ コードとドキュメントを分けて書くと参照しながら開発したり メンテするのが面倒だから
Javadoc のように ソースファイルに ドキュメントを書きつつ ドキュメントの自動生成 もできる方法を調べる

使えそうなツールは
- Go 標準の サブコマンド `go doc` 感覚的には `go fmt` と同じ
- 標準じゃない GoDoc ツール(`go get golang.org/x/tools/cmd/godoc`)

`go doc` サブコマンドは CLI しか対応してない
`godoc` ツールは Web サーバ立てるこもできる

go の公式っぽいサイト
- [トップページ01](https://golang.org/)
- [トップページ02](https://go.dev/)

どちららも[Packages ドキュメント](https://pkg.go.dev/std)に繋がる
そして [Packages ドキュメント](https://pkg.go.dev) は `godoc` ツールを使ってそう

## `go doc` サブコマンド の使い方
### 前提
`$ go version`
go version go1.15.6 linux/amd64
`$ pwd`
/go/src/github.com/ozaki-physics/go-training-composition
`$ go help doc`
doc の help が見れる
### ドキュメントの書き方
`/* ドキュメント */` を使って書くことが多いらしい
```go
/*
Sample_server はメソッドです

これで改行して書ける
*/
func Sample_server() {
  }

// または

// go doc サブコマンド でのパッケージのドキュメント
// 
// 改行される
package package01
```
### 例
`$ go doc package02`
パッケージのドキュメントだけが出てくる
`$ go doc package02.Sample_server`
メソッドのドキュメントだけが出てくる

Flags は 6種ある
1. `$ go doc -all package02`
→ ドキュメント ありで パッケージもメソッドもすべて出てくる
2. `$ go doc -c package02`
3. `$ go doc -cmd package02`
4. `$ go doc -short package02`
→ ドキュメント なしで 1行で出力される
5. `$ go doc -src package02`
6. `$ go doc -u package02`
→ ドキュメント なしで private なメソッドや定数も表示される

組み合わせることもできる
`$ go doc -short -u package02`
→ private も含めて1行で出力

`$ go help doc` で出力される `go doc encoding/json` を調べる
package のドキュメントとして
```bash
/go/src/github.com/ozaki-physics/go-training-composition# go doc encoding/json
Package json implements encoding and decoding of JSON as defined in RFC
7159. The mapping between JSON and Go values is described in the
documentation for the Marshal and Unmarshal functions.

See "JSON and Go" for an introduction to this package:
https://golang.org/doc/articles/json_and_go.html
```
と出力されて この文字列が書かれているのは
```bash
/usr/local/go/src/encoding/json# head --line 15 encode.go
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package json implements encoding and decoding of JSON as defined in
// RFC 7159. The mapping between JSON and Go values is described
// in the documentation for the Marshal and Unmarshal functions.
//
// See "JSON and Go" for an introduction to this package:
// https://golang.org/doc/articles/json_and_go.html
package json

import (
        "bytes"
        "encoding"
```
encoding/json パッケージには 色々なファイルがあるが encode.go だけ package json の直前かつ改行も無くドキュメントが書かれてた
だから `go doc encoding/json`をしたときに このドキュメントが出ると思われる

## `godoc` ツール の使い方
[自作したパッケージの詳細をGoDocで確認できるようにする](https://hodalog.com/show-details-of-go-package-using-godoc/)を見ながらやってみようとした
`go get golang.org/x/tools/cmd/godoc` をやって go.mod に足してみたけど
`godoc` コマンドが見つからないって言われてて面倒になった笑

ドキュメントの書き方
`// ドキュメント` を使って書くことが多いらしい
また ドキュメントその要素名で
[Packages ドキュメント](https://pkg.go.dev) の UI が `godoc` ツール と似ている
かつ 標準 package を見ると `// ドキュメント` を多用してるし `doc.go` ファイルもあるから `godoc` ツールを使ってそう?

他の参考文献
[GoDocドキュメントの書き方](https://blog.lufia.org/entry/2018/05/14/150400)
[go doc の使い方・コメントを書いて、ちゃんと読む](https://qiita.com/ayasuda/items/53933c83d0fb7152c7e9)
