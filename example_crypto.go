package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"time"
)

// init 関数は main 関数の前に実行される初期化関数
func init() {
	initTimeZone()
	initLog()
}

// log の初期設定
func initLog() {
	// ログに接頭辞を付けられる
	log.SetPrefix("[暗号化の実験]")
	// エラーの行数をつける 呼び出し元か
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func errCheck(err error) {
	if err != nil {
		// 第1引数が 1だと 実際にエラーの位置
		// 第1引数が 2だと 呼び出し元の位置
		log.Output(2, "エラー発生元")
		log.Fatal(err)
	}
}

// タイムゾーンの初期設定
func initTimeZone() {
	// タイムゾーンの変更
	const location = "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

// 型確認
func typeCheck(example interface{}) {
	log.Output(2, "型の確認元")
	log.Printf("%T\n", example)
}

func main() {
	log.Println("main 開始")
	// example01()
	// example02()
	example03()
	log.Println("main 終了")
}

// 共通鍵暗号方式
// ブロック暗号化方式の AES
// 16byteの固定視長の平文しか暗号化できず使えない
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example01() {
	// 平文の用意
	plainText := []byte("This is 16 bytes")

	// 鍵の用意
	// 鍵の長さは 16, 24, 32 バイトのどれかにしないとエラー
	// これは32バイト(アルファベット1個が1バイトだから)
	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	block, err := aes.NewCipher(key)
	errCheck(err)

	// 暗号文を入れる変数の用意
	cipherText := make([]byte, len(plainText))
	// AES で暗号化
	block.Encrypt(cipherText, plainText)
	// 16進数で出力 結果は暗号化されている
	fmt.Printf("暗号文(16進数): %x\n", cipherText)

	// 復号文を入れる変数の用意
	decryptedText := make([]byte, len(cipherText))
	// 復号する
	block.Decrypt(decryptedText, cipherText)
	// 結果は元の平文が得られる
	fmt.Printf("復号文(string): %s\n", string(decryptedText))
}

// 共通鍵暗号方式
// ブロック暗号化方式の AES
// CBC モード
// See Go 言語で学ぶ『暗号技術入門』Part 3 -CBC Mode-: https://skatsuta.github.io/2016/03/06/hyuki-crypt-book-go-3/
func example02() {
	// 平文の用意 長さが 16 バイトの整数倍でない場合はパディングする必要がある
	plainText := []byte("secret text 9999")

	// 鍵の用意
	key := []byte("secret-key-12345")
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	block, err01 := aes.NewCipher(key)
	errCheck(err01)

	// 暗号文を入れる変数の用意
	// 先頭に初期化ベクトル (IV) を入れるため、1ブロック分余計に確保する
	// CBC モードは1つ前の暗号ブロックを使って暗号化するから 一番最初はランダムで用意する
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	// 暗号文の戦闘ブロック(IV)の参照を取り出す
	iv := cipherText[:aes.BlockSize]
	// iv がランダムなビット列する
	_, err02 := io.ReadFull(rand.Reader, iv)
	errCheck(err02)

	// 暗号化用オブジェクトを用意
	block01 := cipher.NewCBCEncrypter(block, iv)
	// AES で暗号化 CBC モード
	block01.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	// 16進数で出力 結果は暗号化されている
	fmt.Printf("暗号文(16進数): %x\n", cipherText)

	// 復号化用オブジェクトを用意
	block02 := cipher.NewCBCDecrypter(block, iv)
	// 復号文を入れる変数の用意
	decryptedText := make([]byte, len(plainText))
	// 復号する 先頭の IV を除いた部分だけ
	block02.CryptBlocks(decryptedText, cipherText[aes.BlockSize:])
	// 結果は元の平文が得られる
	fmt.Printf("復号文(string): %s\n", decryptedText)
}

// 共通鍵暗号方式
// ブロック暗号化方式の AES
// CTR モード
// だから ストリーム暗号とみなせる
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example03() {
	// 平文の用意
	plainText := []byte("Bob loves Alice. But Alice hate Bob...")

	// 鍵の用意
	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	block, err01 := aes.NewCipher(key)
	errCheck(err01)

	// 暗号文を入れる変数の用意
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	// 暗号文の戦闘ブロック(IV)の参照を取り出す
	iv := cipherText[:aes.BlockSize]
	// iv がランダムなビット列する
	_, err02 := io.ReadFull(rand.Reader, iv)
	errCheck(err02)

	// 暗号化用オブジェクトを用意
	encryptStream := cipher.NewCTR(block, iv)
	// AES で暗号化 CTR モード
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	// 16進数で出力 結果は暗号化されている
	fmt.Printf("暗号文(16進数): %x \n", cipherText)

	// 復号化用オブジェクトを用意
	// ユニークな必要はあるが 安全な必要はないので 暗号文の先頭に差し込んである
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	// 復号文を入れる変数の用意
	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	// 復号する 先頭の IV を除いた部分だけ
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
	// 結果は元の平文が得られる
	fmt.Printf("復号文(string): %s\n", string(decryptedText))
}
