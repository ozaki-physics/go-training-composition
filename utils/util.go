package utils

import (
	"log"
	"time"
)

// InitLog log の初期設定
func InitLog(s string) {
	// ログに接頭辞を付けられる
	log.SetPrefix(s + " ")
	// エラーの行数をつける 呼び出し元か
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// InitTimeZone タイムゾーンの初期設定
func InitTimeZone() {
	// タイムゾーンの変更
	const location = "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

// ErrCheck エラーチェックを雑にやってログに出力する
func ErrCheck(err error) {
	if err != nil {
		// 第1引数が 1だと 実際にエラーの位置
		// 第1引数が 2だと 呼び出し元の位置
		log.Output(2, "エラー発生元")
		log.Fatal(err)
	}
}

// TypeCheck 型を確認してログに出す
func TypeCheck(example interface{}) {
	log.Output(2, "型の確認元")
	log.Printf("%T\n", example)
}
