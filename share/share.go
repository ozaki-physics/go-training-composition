package share

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// terminalArgs 実行時にコマンド引数から ファイルのパスを受け取る
func TerminalArgs() string {
	flag.Parse()
	args := flag.Args()

	var filepath string
	if len(args) > 0 {
		// コマンド引数の1個目を取得
		filepath = args[0]
		// ファイルが存在するか確認
		searchFile(filepath)
	} else {
		// 標準入力で JSON の filapath を受け取る
		log.Fatal("実行時にパスが渡されませんでした")
	}

	log.Printf("%s を加工します\n", filepath)
	return filepath
}

// ファイルの有無を調べる
func searchFile(fileName string) {
	_, err := os.Stat(fileName)
	// ファイルが存在するか
	if errors.Is(err, fs.ErrNotExist) {
		log.Println("ファイルが存在しません")
		log.Fatal(err)
	}
}

// ターミナルを開いたままにする
func StopTerminal() {
	fmt.Println("何か入力してください。Enterキーを押すと終了します。")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // ユーザーの入力を待ちます
	fmt.Println("入力された内容:", scanner.Text())
}

// RenameFile は、与えられた oldPath から newPath へファイル名を変更します。
func RenameFile(oldPath, newPath string) error {
	// os.Rename を使用してファイル名を変更
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return nil
}

// SplitPath は与えられたパスをディレクトリとファイル名に分割します。
func SplitPath(path string) (dir, file string) {
	// ディレクトリとファイル名を分割
	dir = filepath.Dir(path)
	file = filepath.Base(path)
	return
}

// utf8ToShiftJIS は UTF-8 文字列 を Shift-JIS バイト列に変換します
func UTF8ToShiftJIS(utf8Str string) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(utf8Str)), japanese.ShiftJIS.NewEncoder())
	return io.ReadAll(reader)
}

// ShiftJISToUTF8 は Shift-JIS バイト列を UTF-8 文字列 に変換します
func ShiftJISToUTF8(sjisBytes []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(sjisBytes), japanese.ShiftJIS.NewDecoder())
	return io.ReadAll(reader)
}

// 渡されたファイルの最後に修飾子をつける
// path := "aaa/bbb/ccc.txt", decoration := _old
// new_path := "aaa/bbb/ccc_old.txt"
func SuffixString(path string, decoration string) string {
	dir, file := SplitPath(path)

	ext := filepath.Ext(file)
	name := file[:len(file)-len(ext)]
	newFile := name + decoration + ext

	new_path := filepath.Join(dir, newFile)

	return new_path
}
