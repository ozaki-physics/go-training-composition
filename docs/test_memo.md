### テストコードはどこに書くのが良いか
テストで1つのパッケージにするのか?
[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md) によると
/test ディレクトリ は
>追加の外部テストアプリとテストデータ
らしいからテストだけをまとめたパッケージを作るのは違う?

[標準パッケージの fmt](https://github.com/golang/go/tree/master/src/fmt) では
わざわざ test をまとめたパッケージは見当たらない
同じパッケージ内に *_test.go を作っている
ただ `package fmt_test` という名前がついている


[他言語プログラマが最低限、気にすべきGoのネーミングルール](https://zenn.dev/keitakn/articles/go-naming-rules)
標準パッケージでテストコードには `パッケージ名_test` と書かれていた
