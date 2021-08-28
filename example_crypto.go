package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
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
	// 共通鍵暗号方式 ブロック暗号化方式の AES
	// example01()
	// 共通鍵暗号方式 ブロック暗号化方式の AES CBC モード
	// example02()
	// 共通鍵暗号方式 ブロック暗号化方式の AES CTR モード だから ストリーム暗号とみなせる
	// example03()
	// 公開鍵暗号方式 RSA-PKCS1v15 で暗号化
	// example04()
	// ハッシュ SHA-2 の SHA-512
	// example05()
	// メッセージ認証コード(MAC) 否認ができず 送信者の証明ができない
	// example06()
	// デジタル署名 公開鍵暗号の応用なため 今回は 楕円曲線暗号 を使う
	// example07()
	// 証明書(x509) 自己署名証明書を作ってみる
	// example08()
	// TLS
	example09()
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
	// 暗号化
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
	// 暗号化 CBC モード
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
	// 暗号化 CTR モード
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

// 公開鍵暗号方式
// RSA-PKCS1v15 で暗号化
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example04() {
	// 平文の用意
	plainText := []byte("Bob loves Alice.")

	// size of key (bits)
	// 2048は2030年まで使うことができる
	keySize := 2048

	// 秘密鍵と公開鍵を生成
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	errCheck(err)
	// 公開鍵を取得 privateKey 構造体の中に 秘密鍵と対応した公開鍵がある
	publicKey := &privateKey.PublicKey

	// 暗号化
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	errCheck(err)
	// 16進数で出力 結果は暗号化されている
	fmt.Printf("暗号文(16進数): %x\n", cipherText)

	// 復号する
	decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	errCheck(err)
	// 結果は元の平文が得られる
	fmt.Printf("復号文(string): %s\n", decryptedText)
}

// ハッシュ
// SHA-2 の SHA-512
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example05() {
	msg := []byte("foo")
	checksum512 := sha512.Sum512(msg)
	fmt.Println(checksum512)
}

// メッセージ認証コード(MAC)
// 否認ができず 送信者の証明ができない
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example06() {
	msg := []byte("Bob loves Alice.")
	key := []byte("passw0rd")

	// HMAC は任意の hash.Hash 関数を使うことができる
	// 今回は SHA-512 を使う
	h1 := hmac.New(sha512.New, key)
	h1.Write(msg)
	mac1 := h1.Sum(nil)
	fmt.Printf("MAC1(16進数): %x\n", mac1)

	h2 := hmac.New(sha512.New, key)
	h2.Write(msg)
	mac2 := h2.Sum(nil)
	fmt.Printf("MAC2(16進数): %x\n", mac2)

	fmt.Printf("2個の MAC 値は同じか? %v\n", hmac.Equal(mac1, mac2))
}

// デジタル署名
// 公開鍵暗号の応用なため 今回は 楕円曲線暗号 を使う
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example07() {
	// 秘密鍵と公開鍵を生成
	// 利用できる曲線は P-224, P-256, P-384, P-521
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	errCheck(err)
	// 公開鍵を取得 privateKey 構造体の中に 秘密鍵と対応した公開鍵がある
	publicKey := &privateKey.PublicKey

	// 任意の長さのハッシュ値を生成
	// 本当は メッセージをハッシュ化したものを使うが 簡略化のため 直接生成する
	hash := []byte("This is message.")
	// 署名する
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	errCheck(err)

	// r, s が何か分からない
	fmt.Printf("r の値: %d\n", r)
	fmt.Printf("s の値: %d\n", s)

	// 署名が正しいか 公開鍵で確認する
	if ecdsa.Verify(publicKey, hash, r, s) {
		fmt.Println("確認できた")
	} else {
		fmt.Println("確認できなかった")
	}
}

// 証明書(x509)
// 自己署名証明書を作ってみる
// 公開鍵暗号としては楕円曲線暗号を使い PEM形式でファイルに保存する(ca.pem)
// 証明書の検証はしてない
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example08() {
	// 秘密鍵と公開鍵を生成
	privateKey, err01 := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	errCheck(err01)
	// 公開鍵を取得 privateKey 構造体の中に 秘密鍵と対応した公開鍵がある
	publicKey := &privateKey.PublicKey

	// 証明書のテンプレート
	ca := x509.Certificate{
		IsCA:         true,
		SerialNumber: big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Japan"},
			Organization: []string{"TCNKSM ECDSA CA Inc."},
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 証明書を作成
	derBytes, err02 := x509.CreateCertificate(rand.Reader, &ca, &ca, publicKey, privateKey)
	errCheck(err02)

	certOut, err03 := os.Create("ca.pem")
	errCheck(err03)

	defer certOut.Close()

	err04 := pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	errCheck(err04)
}

// TLS
// See Go言語と暗号技術(AESからTLS): https://deeeet.com/writing/2015/11/10/go-crypto/
func example09() {
	fmt.Println("また今後 元気があるときに調べたい")
}
