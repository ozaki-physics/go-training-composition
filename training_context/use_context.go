package training_context

import (
	"context"
	"fmt"
	"time"
)

func child(ctx context.Context) {
	// 処理を始める最初に context の状態を確認
	if err := ctx.Err(); err != nil {
		return
	}
	// 処理を進める
	fmt.Println("キャンセルされてない")
}

// CancelAlert 2回 出力されるが 3回目は出力されない
func CancelAlert() {
	ctx, cancel := context.WithCancel(context.Background())
	child(ctx)
	child(ctx)
	cancel()
	child(ctx)
}

// ExampleCancel
// see: [公式](https://pkg.go.dev/context?utm_source=gopls#example-WithCancel)
func ExampleCancel() {
	// 整数を生成し 整数の入った goroutine を返す
	gen := func(ctx context.Context) <-chan int {
		// 永遠と整数が格納されていく
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					// goroutine を漏らさない
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	// 整数を使い終わったら スコープの終了と同時に キャンセル する
	defer cancel()

	for n := range gen(ctx) {
		fmt.Println(n)
		// 一致したら ExampleCancel を終了させる
		// つまり defer が実行されて gen(ctx) の中の Done() が動くと思われる
		if n == 3 {
			break
		}
	}
}

// ExampleWithTimeout
// see: [公式](https://pkg.go.dev/context#example-WithTimeout)
func ExampleWithTimeout() {
	// 1sec より短いから 先に ctx.Done() が動く
	const shortDuration = 1 * time.Millisecond
	// 1sec より長いから time.After() が動く
	// const shortDuration = 10000 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("停止しすぎ")
	case <-ctx.Done():
		// "context deadline exceeded" と出力される
		fmt.Println(ctx.Err())
	}
}

// WaitCancel Done() で Cancel 通知待ちする
func WaitCancel() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() { fmt.Println("別の goroutine") }()
	fmt.Println("停止")
	// 通知待ちして timeout して動き出す
	<-ctx.Done()
	fmt.Println("そして時は動き出す")
}

// WaitCancelAndDo キャンセル待ちしながら別処理もする
// はじめは キャンセルされていない が出力(別処理)されるが 1sec 後からは get 整数 が出力される
func WaitCancelAndDo() {
	ctx, cancel := context.WithCancel(context.Background())
	task := make(chan int)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case i := <-task:
				fmt.Println("get", i)
			default:
				fmt.Println("キャンセルされてない")
			}
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(time.Second)
	for i := 0; i < 5; i++ {
		task <- i
	}
	cancel()
}

// ContextValue context にデータを入れて取り出す
func ContextValue() {
	const zeroTraceID = ""
	// 型アサーションするときに プリミティブ型でもいいが データの意味を分かりやすくするため
	type traceID string
	// 他のパッケージと衝突を避けるため
	type traceIDKey struct{}

	// 外にメソッド作るのが面倒だったから 内部で定義した
	getTraceID := func(ctx context.Context) traceID {
		if v, ok := ctx.Value(traceIDKey{}).(traceID); ok {
			return v
		}
		return zeroTraceID
	}

	// 外にメソッド作るのが面倒だったから 内部で定義した
	setTraceID := func(ctx context.Context, tid traceID) context.Context {
		return context.WithValue(ctx, traceIDKey{}, tid)
	}

	ctx := context.Background()
	fmt.Printf("トレースID = %q\n", getTraceID(ctx))
	ctx = setTraceID(ctx, "test-id")
	fmt.Printf("トレースID = %q\n", getTraceID(ctx))
}
