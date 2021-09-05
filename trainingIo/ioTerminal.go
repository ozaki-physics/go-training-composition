// ターミナル での入出力の勉強
package trainingIo

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ozaki-physics/go-training-composition/utils"
	"os"
	"strings"
)

// OutTerminal 型を気にしない出力方法
// func Fprintln(w io.Writer, a ...interface{}) (n int, err error)
// See: https://pkg.go.dev/fmt@go1.17#Fprintln
func OutTerminal() {
	const name, age = "ozaki", 25
	// os.Stdout はターミナルへの標準出力 *os.File と同じ扱いになる
	// See: https://pkg.go.dev/os@go1.17#pkg-variables
	fmt.Fprintln(os.Stdout, name, "is", age, "years old.")
}

// UseFmtScan 一番簡単な方法
// func Scan(a ...interface{}) (n int, err error)
// 半角で区切って取得する
// エンターも半角スペース扱い
// See: https://pkg.go.dev/fmt@go1.17#Scan
func UseFmtScan() {
	fmt.Print("fmt.Scan の 入力待ち > ")
	// fmt.Scan() にポインタを渡すのが重要
	// ポインタを渡すから先に ゼロ値で定義しておく
	var input01, input02 string
	// 入力数が定義より多くてもエラーにならない
	// 多く入力された分は 実行された後のターミナルに勝手に書き込まれるのが厄介
	n, err := fmt.Scan(&input01, &input02)
	utils.ErrCheck(err)
	fmt.Printf("入力個数: %d\n", n)
	fmt.Printf("入力01: %s, 入力02: %s\n", input01, input02)
}

// UseBufioScanner01 2番目に簡単
// エンターで区切って入力される
// 終わるためには ファイルの最後 や エラーを発生させるしかない
// だから ターミナルでパイプを使ってファイルを入力するときは相性いいかも
// os.Stdin ターミナルへの標準入力 *os.File と同じ扱いになる
// See: https://pkg.go.dev/os@go1.17#pkg-variables
func UseBufioScanner01() {
	fmt.Print("bufio.Scanner の 入力待ち > ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("入力されたもの: %s\n", input)
	}

	// scanner.Scan() が false になってから Scanner に err が格納されるらしい
	err := scanner.Err()
	utils.ErrCheck(err)
}

// UseBufioScanner02 は UseBufioScanner01 を終われるようにしたもの
// ただ goto を使うために scanner.Scan() を for から if にしたし
// goto は良くない設計だから改善の必要が大いにある
func UseBufioScanner02() {
forStart:
	fmt.Print("bufio.Scanner の 入力待ち > ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "":
			fmt.Println("何も入力されていません")
			goto forStart
		case "n":
			fmt.Println("終了します")
			goto forEnd
		default:
			fmt.Printf("入力されたもの: %s\n", input)
			goto forStart
		}
	}
forEnd:
}

// UseBufioReader いい感じに終われるし goto が1個で済む
// os.Stdin が *os.File なことをいい感じ使っている
func UseBufioReader() {
	r := bufio.NewReader(os.Stdin)

loop:
	fmt.Print("bufio.Reader の 入力待ち > ")
	input, err := r.ReadString('\n')
	utils.ErrCheck(err)

	// \n を削除して 全部小文字にする
	// See: https://pkg.go.dev/strings@go1.17#Trim
	// See: https://pkg.go.dev/strings@go1.17#ToLower
	input = strings.ToLower(strings.Trim(input, "\n"))
	switch input {
	case "":
		fmt.Println("何も入力されていません")
		goto loop
	case "n":
		fmt.Println("終了します")
	default:
		fmt.Printf("入力されたもの: %s\n", input)
		goto loop
	}
}

// TerminalArgsFlag ターミナル引数を受け取る
// flag を使うバージョン
// `go run example_ioFile.go -opt01=aaa -opt02=bbb hello` と 先にオプションを書かないといけない
func TerminalArgsFlag() {
	// ターミナル引数にオプションを設定方法は2つ
	// func String(name string, value string, usage string) *string
	// f := flag.String("オプション名", "初期値", "説明")
	// 戻り値がポインタに注意
	// See: https://pkg.go.dev/flag@go1.17#String
	opt01 := flag.String("opt01", "初期値01", "説明01")
	// func StringVar(p *string, name string, value string, usage string)
	// flag.StringVar(&str, "オプション名", "初期値", "説明")
	// See: https://pkg.go.dev/flag@go1.17#StringVar
	var opt02 string
	flag.StringVar(&opt02, "opt02", "初期値02", "説明02")

	// プログラムが flag にアクセスし始める前に flag.Parse() を実行する必要がある
	// See: https://pkg.go.dev/flag@go1.17#Parse
	// 内部的には CommandLine.Parse(os.Args[1:]) で os.Args[1:] と引数を取得してる
	flag.Parse()
	// ターミナル引数を出力する 先にオプションを書いた場合 オプションは引数に含まれない
	args := flag.Args()
	fmt.Printf("引数: %s\n", args)
	// オプション引数
	fmt.Printf("オプション01: %s\n", *opt01)
	fmt.Printf("オプション02: %s\n", opt02)
}

// TerminalArgsOs ターミナル引数を受け取る
// os.Args を使うバージョン
// See: https://pkg.go.dev/os@go1.17#pkg-variables
func TerminalArgsOs() {
	args := os.Args
	fmt.Printf("引数の数: %d\n", len(args))
	fmt.Printf("引数の中身: %#v\n", args)

	// 引数を順に出力
	for i, v := range args {
		fmt.Printf("引数[%d] -> %s\n", i, v)
	}
}

// TerminalArgsFile 引数にファイル path を渡して ファイルの内容を出力する
func TerminalArgsFile() {
	flag.Parse()

	var filename string

	args := flag.Args()
	if len(args) > 0 {
		filename = args[0]
	}
	fmt.Printf("%s を読み込みます\n", filename)
	// ファイルが存在するか確認
	SearchFile(filename)

	file, err := os.Open(filename)
	utils.ErrCheck(err)
	defer file.Close()

	// 内容を出力
	allDataReadFileObject(file)
}
