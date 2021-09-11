package fileCrypto

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ozaki-physics/go-training-composition/trainingIo"
	"github.com/ozaki-physics/go-training-composition/utils"
	"io"
	"os"
)

// crypto 暗号化や復号で使う json の入れ物
type crypto struct {
	Encode filepath
	Decode filepath
	Key    string
}

// filepath encode または decode のときに使う ファイルの path
type filepath struct {
	Read  string
	Write string
}

// getPathCrypto JSON から encode に使う情報を取り出す
func getPathCrypto() crypto {
	path := "fileCrypto/json/filepath.json"
	// path := "fileCrypto/json/filepath_product.json"
	// ファイルが存在するか確認
	trainingIo.SearchFile(path)
	bytes, err01 := os.ReadFile(path)
	utils.ErrCheck(err01)

	var cp crypto
	err02 := json.Unmarshal(bytes, &cp)
	utils.ErrCheck(err02)

	return cp
}

// RunFileEnCrypto ファイルの中身を暗号化する
func RunFileEnCrypto() {
	utils.InitLog("[ファイルの中身を暗号化]")
	utils.StartLog()

	cp := getPathCrypto()
	// ファイルが存在するか確認
	trainingIo.SearchFile(cp.Encode.Read)
	fp := cp.Encode
	key := getKey(cp.Key)

	fr, err := os.Open(fp.Read)
	utils.ErrCheck(err)
	defer fr.Close()

	fw, err := os.Create(fp.Write)
	utils.ErrCheck(err)
	defer fw.Close()

	s := bufio.NewScanner(fr)
	for i := 1; s.Scan(); i++ {
		oneLine := s.Text()
		// 暗号化
		cipherByte := enCrypto(oneLine, key)
		// byte を 16進数へ
		cipherString := hex.EncodeToString(cipherByte)
		fmt.Printf("暗号文(16進数) %03d行目: %s\n", i, cipherString)
		// 書き出す
		_, err := fmt.Fprintln(fw, cipherString)
		utils.ErrCheck(err)
	}
	utils.ErrCheck(s.Err())
	utils.EndLog()
}

// RunFileDeCrypto ファイルから復号する
func RunFileDeCrypto() {
	utils.InitLog("[ファイルの中身を復号]")
	utils.StartLog()

	cp := getPathCrypto()
	// ファイルが存在するか確認
	trainingIo.SearchFile(cp.Decode.Read)
	fp := cp.Decode
	key := getKey(cp.Key)

	fr, err := os.Open(fp.Read)
	utils.ErrCheck(err)
	defer fr.Close()

	fw, err := os.Create(fp.Write)
	utils.ErrCheck(err)
	defer fw.Close()

	s := bufio.NewScanner(fr)
	for i := 1; s.Scan(); i++ {
		oneLine := s.Text()
		// 16進数を byte へ
		plainByte, err01 := hex.DecodeString(oneLine)
		utils.ErrCheck(err01)
		// 復号
		plainText := deCrypto(plainByte, key)
		fmt.Printf("復号文(string) %03d行目: %s\n", i, plainText)
		// 書き出す
		_, err02 := fmt.Fprintln(fw, plainText)
		utils.ErrCheck(err02)
	}
	utils.ErrCheck(s.Err())

	utils.EndLog()
}

// getKey 外部ファイルから key を取得する
func getKey(keyfile string) []byte {
	// ファイルが存在するか確認
	trainingIo.SearchFile(keyfile)

	file, err01 := os.Open(keyfile)
	utils.ErrCheck(err01)
	defer file.Close()

	content := make([]byte, 32)
	_, err02 := file.Read(content)

	if err02 == io.EOF {
		return content
	} else {
		utils.ErrCheck(err02)
	}

	return content
}

// enCrypto 暗号化
func enCrypto(plainText string, key []byte) []byte {
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	block, err01 := aes.NewCipher(key)
	utils.ErrCheck(err01)

	// 文字列を byte に変換
	plainByte := []byte(plainText)
	// 暗号文を入れる変数の用意
	cipherByte := make([]byte, aes.BlockSize+len(plainByte))
	// 暗号文の先頭ブロック(IV)の参照を取り出す
	iv := cipherByte[:aes.BlockSize]
	// iv がランダムなビット列する
	_, err02 := io.ReadFull(rand.Reader, iv)
	utils.ErrCheck(err02)

	// 暗号化用オブジェクトを用意
	encryptStream := cipher.NewCTR(block, iv)
	// 暗号化 CTR モード
	encryptStream.XORKeyStream(cipherByte[aes.BlockSize:], plainByte)

	return cipherByte
}

// deCrypto 復号
func deCrypto(cipherByte []byte, key []byte) string {
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	block, err01 := aes.NewCipher(key)
	utils.ErrCheck(err01)

	// 復号文を入れる変数の用意
	// ユニークな必要はあるが 安全な必要はないので 暗号文の先頭に差し込んである
	decryptedText := make([]byte, len(cipherByte[aes.BlockSize:]))

	// 復号化用オブジェクトを用意
	decryptStream := cipher.NewCTR(block, cipherByte[:aes.BlockSize])
	// 復号する 先頭の IV を除いた部分だけ
	decryptStream.XORKeyStream(decryptedText, cipherByte[aes.BlockSize:])

	return string(decryptedText)
}
