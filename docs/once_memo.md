# 相互排他ロック
## 概要

>sync パッケージは、相互排他ロックなどの基本的な同期プリミティブを提供します。  
>Once 型と WaitGroup 型以外のほとんどは、低レベルのライブラリルーチンで使用することを意図しています。  
>より高度な同期処理はチャネルや通信で行うのがよいでしょう。  
>このパッケージで定義された型を含む値は、コピーしてはいけません。  

## サンプル
[sync パッケージのサンプル](https://pkg.go.dev/sync#example-Once) より  
```go
func main() {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
// 出力
// Only once
```
`once.Do()` を使わないと10回出力される  

### 参考文献
[sync.Onceの挙動がなんとなくしか分からなかったので実際に試してみた](https://qiita.com/usk81/items/26041e9939f7fd89fcc9)
```go
func main() {
	foo := "no"
	bar := "no"

	var once sync.Once

	once.Do(func() {
		foo = "A"
	})
	once.Do(func() {
		bar = "B"
	})

	fmt.Printf("foo: %s\n", foo)
	fmt.Printf("bar: %s\n", bar)
}
// 出力
// foo: A
// bar: NoN
```
```go
func main() {
	foo := "no"
	bar := "no"

	var once01 sync.Once
	var once02 sync.Once

	once01.Do(func() {
		foo = "A"
	})
	once02.Do(func() {
		bar = "B"
	})

	fmt.Printf("foo: %s\n", foo)
	fmt.Printf("bar: %s\n", bar)
}
// 出力
// foo: A
// bar: B
```

つまり `sync.Once` の変数1個あたり1回の実行を保証してる感じ  
1回しか実行してほしくないからと 1個の `Do()` に詰め込んだりせず `sync.Once` を作ればよいと思われる  
