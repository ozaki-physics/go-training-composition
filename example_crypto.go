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
	// タイムゾーンの変更
	const location = "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	log.SetPrefix("[暗号化の実験]")
}

func main() {
	// log.Println("main 開始")
	// example01()
	// example02()
	// example03()
	// log.Println("main 終了")
	// 型確認
	// log.Printf("%T\n", "string")
}

func errCheck(err error) {
	if err != nil {
		log.Fatal("エラー発生", err)
	}
}

func example01() {
	plainText := []byte("This is 16 bytes")

	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	// Encrypt
	cipherText := make([]byte, len(plainText))
	block.Encrypt(cipherText, plainText)
	fmt.Printf("Cipher text: %x\n", cipherText)

	// Decrypt
	decryptedText := make([]byte, len(cipherText))
	block.Decrypt(decryptedText, cipherText)
	fmt.Printf("Decrypted text: %s\n", string(decryptedText))
}

func example02() {
	// 平文。長さが 16 バイトの整数倍でない場合はパディングする必要がある
	plainText := []byte("secret text 9999")
	// 暗号化データ。先頭に初期化ベクトル (IV) を入れるため、1ブロック分余計に確保する
	encrypted := make([]byte, aes.BlockSize+len(plainText))

	// IV は暗号文の先頭に入れておくことが多い
	iv := encrypted[:aes.BlockSize]
	// IV としてランダムなビット列を生成する
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	// ブロック暗号として AES を使う場合
	key := []byte("secret-key-12345")
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// CBC モードで暗号化する
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[aes.BlockSize:], plainText)
	fmt.Printf("encrypted: %x\n", encrypted)

	// 復号するには復号化用オブジェクトが別に必要
	mode = cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(plainText))
	// 先頭の IV を除いた部分を復号する
	mode.CryptBlocks(decrypted, encrypted[aes.BlockSize:])
	fmt.Printf("decrypted: %s\n", decrypted)
	// Output:
	// decrypted: secret text 9999
}

func example03() {
	// 鍵の長さは 16, 24, 32 バイトのどれかにしないとエラー
	key := []byte("aes-secret-key-1")
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 暗号化される平文の長さは 16 バイト (128 ビット)
	plainText := []byte("secret plain txt")
	// 暗号化されたバイト列を格納するスライスを用意する
	encrypted := make([]byte, aes.BlockSize)
	// AES で暗号化をおこなう
	c.Encrypt(encrypted, plainText)
	// 結果は暗号化されている
	fmt.Println(string(encrypted))
	// Output:
	// #^ϗ~:f9��˱�1�

	// 復号する
	decrypted := make([]byte, aes.BlockSize)
	c.Decrypt(decrypted, encrypted)
	// 結果は元の平文が得られる
	fmt.Println(string(decrypted))
	// Output:
	// secret plain txt
}
