ドキュメントとコードを分けて書くと参照するのが面倒だから
なるべく コードだけで理解できるようにすべき
java doc みたいな go doc はどうやって書けばいいんだろう
あと暗号化の example01, example02, example03 の違いをちゃんと書いておきたい


カレントディレクトリ
/go/src/github.com/ozaki-physics/go-training-composition
go doc コマンドは go のサブコマンドで
go fmt コマンドと同じ感覚

`$ go version`
go version go1.15.6 linux/amd64

`$ go help doc`
doc の help が見れる

`$ go doc`
コマンドで表示されるだけでブラウザでは使えなさそう

コメントの書き方
```go
/*
Sample_server はメソッドです

これで改行して書ける
*/
func Sample_server() {
}
```

`$ go doc package02`
パッケージのコメントだけが出てくる

`$ go doc package02.Sample_server`
メソッドのコメントだけが出てくる

Flags は 6種ある
1. `$ go doc -all package02`
コメント ありで パッケージもメソッドもすべて出てくる
2. `$ go doc -c package02`
3. `$ go doc -cmd package02`
4. `$ go doc -short package02`
コメント なしで 1行で出力される
5. `$ go doc -src package02`
6. `$ go doc -u package02`
コメント なしで private なメソッドや定数も表示される

組み合わせることもできる
`$ go doc -short -u package02`
private も含めて1行で出力

なんか惜しそう
/usr/local/go/src/encoding/json

/usr/local/go/src/encoding# go doc json
と
go help doc で出てくる sample と同じ

/usr/local/go/src/encoding/json# head encode.go
で
```go
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
```
たぶん 他のファイルでは package json の直前に書いてなくて
改行が入ってるから出力されないのかな

/usr/local/go/src/encoding/json# head --line 15 encode.go

公式ですら GoDoc をインストールして使っていそう
