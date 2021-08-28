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
//     整形済みテキスト
// 
// 次のやつはリンク
// https://golang.org/
package package02

import (
	"net/http"
)

var Const_big string = "OK_const pkg02_file01"
var const_small string = "NO_const pkg02_file01"

func Sample_server() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.ListenAndServe(":8080", nil)
}
