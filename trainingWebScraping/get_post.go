package trainingWebScraping

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TryURLParameter URL パラメータを付与する方法
// 一番簡単に http リクエストをする
func TryGETURLParameter() {
	urlString := "http://localhost"
	// URL パラメータを作成
	values := url.Values{}
	values.Add("hello", "world")
	values.Add("how", "are")

	// URL の生成
	url := urlString + "?" + values.Encode()
	fmt.Println(url)
	// リクエストする
	// res, err := http.Get(url)
}

// TryGETRequestHeader http.Client を作って GET メソッドを Do で実行
func TryGET() {
	urlString := "http://localhost"

	// Request を生成
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		panic(err)
	}
	// 組み立てたクエリを生クエリ文字列に変換して設定
	values := url.Values{}
	values.Add("hello", "world")
	values.Add("how", "are")
	req.URL.RawQuery = values.Encode()
	fmt.Println(req.URL.RawQuery)

	// リクエストヘッダの追加
	req.Header.Set("foo", "bar")
	// BASIC 認証用のメソッドも存在する
	// req.SetBasicAuth("user", "pass")
	fmt.Println(req.Header)

	// 10秒でタイムアウトするクライアントを作成
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	// 実際にリクエストが飛ばないように コメントアウトしている
	/*
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(resp.Body)
	*/
	fmt.Println("Syntax Error 回避")
	fmt.Println(client)
}

// TryPOSTSimplePostForm
func TryPOSTSimplePostForm() {
	urlString := "http://localhost"

	values := url.Values{}
	values.Add("hello", "world")
	values.Add("how", "are")

	// 実際にリクエストが飛ばないように コメントアウトしている
	// resp, err := http.PostForm(urlString, values)
	fmt.Println("Syntax Error 回避" + urlString)
}

// TryGETRequestHeader http.Client を作って POST メソッドを Do で実行
func TryPOST() {
	urlString := "http://localhost"

	values := url.Values{}
	values.Add("hello", "world")
	values.Add("how", "are")

	// 10秒でタイムアウトするクライアントを作成
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	// Request を生成
	req, err := http.NewRequest("POST", urlString, strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 実際にリクエストが飛ばないように コメントアウトしている
	/*
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	*/
	fmt.Println("Syntax Error 回避")
	fmt.Println(client)
}
