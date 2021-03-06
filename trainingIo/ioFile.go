// ファイル での入出力の勉強
package trainingIo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/ozaki-physics/go-training-composition/utils"
	"io"
	"io/fs"     // file system に特化
	"io/ioutil" // io の util
	"log"
	"os"
)

// ThisDirFile 読み込みファイルと書き込みファイルを保持する
type ThisDirFile struct {
	ReadName  string
	WriteName string
}

// SearchFile ファイルやディレクトリを探す
// See ファイルが存在するかの判定: https://pkg.go.dev/os@go1.17#IsNotExist
// 昔は os.IsNotExist() だったけど 今は errors.Is(err, fs.ErrPermission) を使う
// See 今のファイルが存在するか判定: https://pkg.go.dev/errors@go1.17#Is
func SearchFile(fileName string) {
	_, err := os.Stat(fileName)
	// ファイルが存在するか
	if errors.Is(err, fs.ErrNotExist) {
		log.Println("ファイルが存在しません")
		log.Fatal(err)
	} else {
		// ファイルが存在しない以外のエラー
		utils.ErrCheck(err)
	}
}

// OpenFile ファイルを開く
// 読み込み専用　ファイルがなければエラー
// See: https://pkg.go.dev/os@go1.17#Open
// 内部的には os.File.OpenFile(name, O_RDONLY, 0) だけ
// 必ず閉じること
// See ファイルが存在するか判定: https://pkg.go.dev/errors@go1.17#Is
func OpenFile(fileName string) {
	f, err := os.Open(fileName)
	// ファイルが存在するか
	if errors.Is(err, fs.ErrNotExist) {
		log.Println("ファイルが存在しません")
	} else {
		// ファイルが存在しない以外のエラー
		utils.ErrCheck(err)
	}
	fmt.Println("ファイルを開きました")
	// file path を出力
	filePath := f.Name()
	fmt.Printf("ファイルの path: %s\n", filePath)

	// defer は return するまで defer に渡した関数の実行を遅延させる
	// つまり エラーで止まっても 正常にメソッドの最後に到達しても必ず実行される
	defer f.Close()
}

// OpenFile02 ファイルを開く
// 読み書きはフラグを渡す
// ファイルがなければ作られる
// ほとんどのユーザーは Open() か Create() を使う
// See 基本使わないらしい: https://pkg.go.dev/os@go1.17#OpenFile
// See 読み書きのパラメータ: https://pkg.go.dev/os@go1.17#pkg-constants
// See 内部的には os.Open() してる: https://github.com/golang/go/blob/master/src/os/file.go
func OpenFile02(fileName string) {
	var permission fs.FileMode = 0755
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, permission)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ファイルを開きました")

	// 閉じる
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ファイルを閉じました")
}

// CreateFile ファイルを作成する
// ファイルがなければ作られる
// ファイルがあれば開く
// 内部的には os.File.OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666) だけ
// See 何も処理してない: https://github.com/golang/go/blob/master/src/os/file.go
func CreateFile(fileName string) {
	f, err := os.Create(fileName)
	utils.ErrCheck(err)
	defer f.Close()
}

// AllDataReadFileName ファイルから一度にすべてのデータを読み込む
// 内部では os.ReadFile(filename) を呼んでるだけ
// See: https://github.com/golang/go/blob/master/src/io/ioutil/ioutil.go
// See: https://pkg.go.dev/io/ioutil@go1.17#ReadFile
func AllDataReadFileName(fileName string) {
	data, err := ioutil.ReadFile(fileName)
	utils.ErrCheck(err)
	// 読み込んだデータを一気に出力
	// 戻り値が []byte だから string にする
	fmt.Println(string(data))
}

// AllDataReadFile は allDataReadFileObject を呼び出すためだけ
func AllDataReadFile(fileName string) {
	f, err := os.Open(fileName)
	utils.ErrCheck(err)
	defer f.Close()
	allDataReadFileObject(f)
}

// allDataReadFileObject ファイルから一度にすべてのデータを読み込む
// ただ ファイルのオブジェクトを引数に渡す
// 内部では io.ReadAll(r) を呼んでるだけ
// See: https://github.com/golang/go/blob/master/src/io/ioutil/ioutil.go
// 引数 r が io.Reader なのに 自作メソッドは *os.File を渡しても動く理由が分からない
// なぜなら *os.File は struct で io.Reader は interface だから
// io.Reader メソッドを持っている *os.File は引数で渡せる
// See: https://pkg.go.dev/io/ioutil@go1.17#ReadAll
// See io.Reader: https://pkg.go.dev/io@go1.17#Reader
// See os.File: https://pkg.go.dev/os@go1.17#File
func allDataReadFileObject(file *os.File) {
	data, err := ioutil.ReadAll(file)
	utils.ErrCheck(err)

	fmt.Println(string(data))
}

// DataReadFile は dataReadFileName を呼び出すためだけ
func DataReadFile(fileName string) {
	f, err := os.Open(fileName)
	utils.ErrCheck(err)
	defer f.Close()
	dataReadFileName(f, 128)
}

// dataReadFileName 何バイトずつ読むか指定してファイルを読み込む
// See: https://pkg.go.dev/os@go1.17#File.Read
func dataReadFileName(file *os.File, volume int) {
	// 読み込みたい量を指定
	content := make([]byte, volume)
	for {
		// 引数に渡した len(content) のバイトまで読み込む
		// content に読み込んだ内容が書かれる
		// 1回で読み終わらなかったら 前回の続きからもう一度読み込む
		// そのとき 新たに読み込んだ内容が content の先頭に足され 前回の内容が後ろにスライドしてく
		readByteCount, err := file.Read(content)

		if err == io.EOF {
			fmt.Println("読み込み終わり")
			break
		} else {
			utils.ErrCheck(err)
		}

		// 後ろに前回読み込んだ内容がスライドしてるから 新たに読み込んだ分だけ表示する
		fmt.Println(string(content[:readByteCount]))
		fmt.Printf("%dバイト読んだ\n", readByteCount)
	}
}

// ScanDataReadFile は dataReadFileScan を呼び出すだけ
func ScanDataReadFile(fileName string) {
	f, err := os.Open(fileName)
	utils.ErrCheck(err)
	defer f.Close()
	dataReadFileScan(f)
}

// scanReadFile ファイルを1行ずつ読み込む
// See: https://pkg.go.dev/bufio@go1.17#NewScanner
// 分割関数のデフォルトは bufio.ScanLines()
// bufio.ScanLines() は改行コードを消しつつ テキストの各行を返す
// See: https://pkg.go.dev/bufio#ScanLines
// func (s *Scanner) Scan() bool は *Scanner を次のトークン(行)へ進める
// ファイルの最後 or エラーで *Scanner が停止すると false が返る
// false を返した後に *Scanner に発生したエラーが入るので やっと Err() で取得できる
// io.EOF のときは Err() は nil を返す
// See: https://pkg.go.dev/bufio@go1.17#Scanner.Scan
// func (s *Scanner) Scan() bool は *Scanner が最初に検知した EOF 以外のエラーを返す
// See: https://pkg.go.dev/bufio@go1.17#Scanner.Err
// func (s *Scanner) Text() string は 今のトークンから文字列を返す
// See: https://pkg.go.dev/bufio@go1.17#Scanner.Text
// func (s *Scanner) Bytes() []byte で byte で返すこともできる
// 区切り文字を変更するには Scanner.Split() に 自作した SplitFunc を渡す
func dataReadFileScan(file *os.File) {
	s := bufio.NewScanner(file)
	// 区切り文字の変更
	// s.Split(scanColon)
	// 出力
	i := 1
	for s.Scan() {
		oneLine := s.Text()
		fmt.Printf("%d行目: %s\n", i, oneLine)
		i++
	}
	// s.Scan() が false になってから Scanner に err が格納されるらしい
	err := s.Err()
	utils.ErrCheck(err)
}

// ScanLines デフォルトの SplitFunc
// dropCR は \r を削除するためにありそう
// プライベートメソッドだし自作するときは要らないのかも
// See: https://github.com/golang/go/blob/master/src/bufio/scan.go
/*
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
*/

// scanColon 自作した SplitFunc
// コロン自体は削除されて取得される
func scanColon(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// ":" だと type byte にしてって怒られた
	if i := bytes.IndexByte(data, ':'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// CopyFile ファイルのデータをコピーする
// See: https://pkg.go.dev/io@go1.17#Copy
func CopyFile(readName, writeName string) {
	// 読み込みファイルの用意
	r, err := os.Open(readName)
	utils.ErrCheck(err)
	// 書き込みファイルの用意
	w, err := os.Create(writeName)
	utils.ErrCheck(err)

	// コピー
	_, err = io.Copy(w, r)
	utils.ErrCheck(err)
}

// RenameFile ファイル名の変更 つまり 実質 path の移動
// See: https://pkg.go.dev/os@go1.17#Rename
// エラーの時は os パッケージの中の *LinkError 型が返ってくる
// LinkError は struct であり リンク, シンボリックリンク, 名前変更システムコール中のエラー, エラー原因のパス
// See: https://pkg.go.dev/os@go1.17#LinkError
func RenameFile(oldPath, newPath string) {
	err := os.Rename(oldPath, newPath)
	utils.ErrCheck(err)
}

// WriteAllData ファイルにデータを一気に書き込む
// ファイルがなければ引数のパーミッションで作成する
// example 見ると Close は要らないようだ
// 内部では os.WriteFile(filename, data, perm) を呼んでるだけ
// See: https://github.com/golang/go/blob/master/src/io/ioutil/ioutil.go
// もっと深堀りすると os.OpenFile() して *os.file の Write() してるだけ
// See: https://github.com/golang/go/blob/master/src/os/file.go
func WriteAllData(fileName string) {
	content := "hello world !"
	var permission fs.FileMode = 0755

	err := ioutil.WriteFile(fileName, []byte(content), permission)
	utils.ErrCheck(err)
}

// WriteDataFile は 3つを呼び出すため
// writeFileString
// writeFileByte
// writeFileFprint
func WriteDataFile(fileName string) {
	f, err := os.Create(fileName)
	// 書き込み権限がないからちゃんとエラーになる
	// f, err := os.Open(fileName)
	utils.ErrCheck(err)
	defer f.Close()

	data01 := "hello world !\n"
	writeFileString(f, data01)
	writeFileByte(f, []byte(data01))
	writeFileFprint(f, []byte(data01))
	data02 := "how are you ?\n"
	writeFileString(f, data02)
	writeFileByte(f, []byte(data02))
	writeFileFprint(f, []byte(data02))

	// 直近で書き込んだ内容を Read するには Seek でファイルの先頭に戻る必要がある
	// 戻らないと何も読み込めなかった
	// O_APPEND でファイルを開いたら動かないらしい
	// 第1引数はファイル内のカーソルの場所 起点から何バイトに置くか
	// 第2引数は 0 はファイルの原点からの相対的な位置, 1 は現在のオフセットからの相対的な位置, 2 は終端からの相対的な位置
	// See: https://pkg.go.dev/os@go1.17#File.Seek
	f.Seek(0, 0)
	// 表示
	dataReadFileName(f, 256)
}

// writeFileString ファイルに書き込む
// サンプルも無いし正しい使い方か分からない
// func (f *File) Write(b []byte) (n int, err error) と似てるけど
// func (f *File) WriteString(s string) (n int, err error) は
// string で書ける点が違い
// See: https://pkg.go.dev/os@go1.17#File.WriteString
func writeFileString(file *os.File, data string) {
	// 第1引数 は ファイル内でカーソルのいる場所っぽい
	_, err := file.WriteString(data)
	utils.ErrCheck(err)
}

// writeFileByte ファイルに書き込む
// See: https://pkg.go.dev/os@go1.17#File.Write
func writeFileByte(file *os.File, data []byte) {
	_, err := file.Write(data)
	utils.ErrCheck(err)
}

// writeFileFprint ファイルに書き込む
// See: https://pkg.go.dev/fmt@go1.17#Fprint
func writeFileFprint(file *os.File, data ...interface{}) {
	// fmt.Fprintln の第1引数は io.Writer インタフェース持っていればいい
	_, err := fmt.Fprintln(file, data)
	utils.ErrCheck(err)
}

// WriteDataWriter 3つを呼び出すため
// writeWriterString
// writeWriterByte
// writeWriterFprint
// See: https://pkg.go.dev/bufio@go1.17#Writer
func WriteDataWriter(fileName string) {
	f, err := os.Create(fileName)
	utils.ErrCheck(err)
	defer f.Close()

	writer := bufio.NewWriter(f)

	data01 := "hello world !\n"
	writeWriterString(writer, data01)
	writeWriterByte(writer, []byte(data01))
	writeWriterFprint(writer, []byte(data01))
	data02 := "how are you ?\n"
	writeWriterString(writer, data02)
	writeWriterByte(writer, []byte(data02))
	writeWriterFprint(writer, []byte(data02))

	// Flush() を実行しないと書き込まない どころか 内容全部消しちゃう
	// バッファのデータを io.Writer で書き込む?
	// os.Open() にすると書き込まないのにエラーにもならない
	// See: https://pkg.go.dev/bufio@go1.17#Writer.Flush
	writer.Flush()

	// 書き込んだ後だから *os.File で読み込める
	f.Seek(0, 0)
	// 表示
	dataReadFileName(f, 256)
}

// writeWriterString ファイルに書き込む
func writeWriterString(w *bufio.Writer, data string) {
	_, err := w.WriteString(data)
	utils.ErrCheck(err)
}

// writeWriterByte ファイルに書き込む
func writeWriterByte(w *bufio.Writer, data []byte) {
	_, err := w.Write(data)
	utils.ErrCheck(err)
}

// writeWriterFprint ファイルに書き込む
func writeWriterFprint(w *bufio.Writer, data ...interface{}) {
	// fmt.Fprintln の第1引数は io.Writer インタフェース持っていればいい
	_, err := fmt.Fprintln(w, data)
	utils.ErrCheck(err)
}

// Reader01DataReadFile は dataReadFileReader01 を呼び出すだけ
func Reader01DataReadFile(fileName string) {
	f, err := os.Open(fileName)
	utils.ErrCheck(err)
	defer f.Close()
	dataReadFileReader01(f)
}

// dataReadFileReader01 ファイルを1行ずつ読み込む
// NewScanner と似てる
// See: https://pkg.go.dev/bufio@go1.17#Reader.ReadString
func dataReadFileReader01(file *os.File) {
	reader := bufio.NewReader(file)

	// byte("\n") も '\n' もダメだった
	delim := byte('\n')
	// 出力
	i := 1
	for {
		oneLine, err := reader.ReadString(delim)
		if err == io.EOF {
			fmt.Println("読み込み終わり")
			break
		} else {
			utils.ErrCheck(err)
		}

		fmt.Printf("%d行目: %s\n", i, oneLine)
		i++
	}
}

// Reader02DataReadFile は dataReadFileReader02 を呼び出すだけ
func Reader02DataReadFile(fileName string) {
	f, err := os.Open(fileName)
	utils.ErrCheck(err)
	defer f.Close()
	dataReadFileReader02(f)
}

// dataReadFileReader02 ファイルを1行ずつ読み込む
// NewScanner と似てる
// reader.ReadLine() は low-level の line-reading らしい
// バッファに対して行が長すぎる場合に isPrefix が true になるっぽい?
// See: https://pkg.go.dev/bufio@go1.17#Reader.ReadLine
func dataReadFileReader02(file *os.File) {
	reader := bufio.NewReader(file)

	// 出力
	i := 1
	for {
		oneLine, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("読み込み終わり")
			break
		} else {
			utils.ErrCheck(err)
		}

		fmt.Printf("%d行目: %s\n", i, oneLine)
		i++
		// isPrefix の使い方が分からない
		if isPrefix {
			fmt.Println("バッファに対して行が長すぎる?")
		}
	}
}
