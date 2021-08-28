# Go でのドキュメントの書き方
そもそもコードを読んだだけで理解できるような書き方をすべきだけど ドキュメントを書きたくなる
ただ コードとドキュメントを分けて書くと参照しながら開発したり メンテするのが面倒だから
Javadoc のように ソースファイルに ドキュメントを書きつつ ドキュメントの自動生成 もできる方法を調べる

使えそうなツールは
- Go 標準の サブコマンド `go doc` 感覚的には `go fmt` と同じ
- 標準じゃない GoDoc ツール(`go get golang.org/x/tools/cmd/godoc`)

`go doc` サブコマンドは CLI しか対応してない
`godoc` ツールは Web サーバ立てるこもできる

__僕は godoc のコマンドがうまく実行できないけど 公式が godoc ツール使ってるから godoc の書き方をする__

go の公式っぽいサイト
- [トップページ01](https://golang.org/)
- [トップページ02](https://go.dev/)
どちららも[Packages ドキュメント](https://pkg.go.dev/std)に繋がる

## `go doc` サブコマンド の使い方
### 前提
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

ドキュメントのルール
- 連続した行は一つの段落になる
- 段落を区切りたい場合は空白行を間に入れる必要がある
- 英大文字で始まり、直前が句読点ではない単一行の段落は見出し
- 字下げすると整形済みテキスト
- URL はリンク

```go
// ここはパッケージコメントの最初になるから見出しではない
// 
// Aなど 英大文字で始まり単一行かつ句読点なしかつ前が見出しではないのでこれは見出し
// 
// 段落の開始
// 内容
// 段落の終了
// 
// 次の段落の開始
// 内容
// 次の段落の終了
// 
//     整形済みテキスt
// 
// 次のやつはリンク
// https://golang.org/
// 
// BUG(who): 保存機能は未実装です
// Deprecated: 非推奨を表す
// See 説明文: https://xxx
// TODO(who): 実装する
package sample
```

他にも Example を書くための それ用のメソッドを作って `// Output` って書くと `go test` でテストが実行できるらしい
また今後元気があるときに調べよう

## 公式も `godoc` ツールを使っているみたい
公式ドキュメントのブログ[Godoc: documenting Go code](https://go.dev/blog/godoc)で
>we have developed the godoc documentation tool.
`godoc` ツールを開発したって言ってるから 公式ツールっぽい

>Godoc is conceptually related to Python’s Docstring and Java’s Javadoc but its design is simpler.
GodocはPythonのDocstringやJavaのJavadocと概念的には似ていますが、デザインはよりシンプルです。
>The comments read by godoc are not language constructs (as with Docstring) nor must they have their own machine-readable syntax (as with Javadoc).
Godocで読まれるコメントは、Docstringのような言語構造ではなく、Javadocのような機械で読める独自の構文を持つ必要もありません。
>Godoc comments are just good comments, the sort you would want to read even if godoc didn’t exist.
Godocのコメントは、たとえgodocが存在しなくても読みたくなるような、優れたコメントなのです。

つまり Javadoc でいう @Param とか @return みたいなドキュメントの書き方は無いかも
機能としてあるのは 
`// BUG(who): 保存機能は未実装です`
`// Deprecated: 非推奨を表す`
>The “who” part should be the user name of someone who could provide more information.
who の部分には、より多くの情報を提供できる人のユーザー名を入力してください。
面倒だから毎回 who って書いてもいいかな?w 一応 ozaki って書いとくべき?

慣例として
`// See 説明文: https://xxx`

また[GoDocドキュメントの書き方](https://blog.lufia.org/entry/2018/05/14/150400)によると
拡張することで
`// TODO(who): 実装する`
など作れるらしい

`go doc` サブコマンドの[go/doc の Documentation](https://pkg.go.dev/go/doc)で
Javadocのような機械で読める独自の構文がないか調べたけど よく分からなかった
逆に `func ToHTML(w io.Writer, text string, words map[string]string)` があって html 作れるのか? って思った
どうやら `godoc` ツールで `ToHTML()` を使っているらしい

## 他の参考文献
[go doc の使い方・コメントを書いて、ちゃんと読む](https://ayasuda.github.io/pages/introduction_go_doc.html)
