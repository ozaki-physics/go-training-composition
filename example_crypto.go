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
	// log.Println("main 開始")
	// example01()
	// example02()
	// example03()
	// log.Println("main 終了")
}

func example01() {
	plainText := []byte("This is 16 bytes")

	// 鍵の長さは 16, 24, 32 バイトのどれかにしないとエラー
	// これは32バイト(アルファベット1個が1バイトだから)
	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	block, err := aes.NewCipher(key)
	errCheck(err)

	// AES で暗号化
	cipherText := make([]byte, len(plainText))
	block.Encrypt(cipherText, plainText)
	// 16進数で出力 結果は暗号化されている
	fmt.Printf("暗号文(16進数): %x\n", cipherText)

	// 復号する
	decryptedText := make([]byte, len(cipherText))
	block.Decrypt(decryptedText, cipherText)
	// 結果は元の平文が得られる
	fmt.Printf("復号文(string): %s\n", string(decryptedText))
}

func example02() {
	// 平文。長さが 16 バイトの整数倍でない場合はパディングする必要がある
	plainText := []byte("secret text 9999")
	// 暗号化データ。先頭に初期化ベクトル (IV) を入れるため、1ブロック分余計に確保する
	encrypted := make([]byte, aes.BlockSize+len(plainText))

	// IV は暗号文の先頭に入れておくことが多い
	iv := encrypted[:aes.BlockSize]
	// iv がランダムなビット列になる
	_, err := io.ReadFull(rand.Reader, iv)
	errCheck(err)

	// ブロック暗号として AES を使う場合
	key := []byte("secret-key-12345")
	block, err := aes.NewCipher(key)
	errCheck(err)

	// CBC モードで暗号化する
	mode01 := cipher.NewCBCEncrypter(block, iv)
	mode01.CryptBlocks(encrypted[aes.BlockSize:], plainText)
	fmt.Printf("encrypted: %x\n", encrypted)

	// 復号するには復号化用オブジェクトが別に必要
	mode02 := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(plainText))
	// 先頭の IV を除いた部分を復号する
	mode02.CryptBlocks(decrypted, encrypted[aes.BlockSize:])
	fmt.Printf("decrypted: %s\n", decrypted)
	// Output: decrypted: secret text 9999
}

func example03() {
	plainText := []byte("Bob loves Alice. But Alice hate Bob...")

	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

	// Create new AES cipher block
	block, err01 := aes.NewCipher(key)
	errCheck(err01)

	// Create IV
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	_, err02 := io.ReadFull(rand.Reader, iv)
	errCheck(err02)

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	fmt.Printf("Cipher text: %x \n", cipherText)

	// Decrpt
	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
	fmt.Printf("Decrypted text: %s\n", string(decryptedText))
}
