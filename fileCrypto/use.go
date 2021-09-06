package fileCrypto

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ozaki-physics/go-training-composition/trainingIo"
	"github.com/ozaki-physics/go-training-composition/utils"
	"io"
	"os"
)

type crypto struct {
	readpath  string
	writepath string
}

// RunFileEnCrypt ファイルの中身を暗号化する
func RunFileEnCrypt() {
	utils.InitLog("[ファイルの中身を暗号化]")
	utils.StartLog()

	cp := getReadWriteEnCrypt()
	key := getKey()

	fr, err := os.Open(cp.readpath)
	utils.ErrCheck(err)
	defer fr.Close()

	fw, err := os.Create(cp.writepath)
	utils.ErrCheck(err)
	defer fw.Close()

	s := bufio.NewScanner(fr)
	for i := 1; s.Scan(); i++ {
		oneLine := s.Text()
		// 暗号化
		cipherByte := enCrypt(oneLine, key)
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

// RunFileDeCrypt ファイルから復号する
func RunFileDeCrypt() {
	utils.InitLog("[ファイルの中身を復号]")
	utils.StartLog()

	cp := getReadWriteDnCrypt()
	key := getKey()

	fr, err := os.Open(cp.readpath)
	utils.ErrCheck(err)
	defer fr.Close()

	fw, err := os.Create(cp.writepath)
	utils.ErrCheck(err)
	defer fw.Close()

	s := bufio.NewScanner(fr)
	for i := 1; s.Scan(); i++ {
		oneLine := s.Text()
		// 16進数を byte へ
		plainByte, err01 := hex.DecodeString(oneLine)
		utils.ErrCheck(err01)
		// 復号
		plainText := deCrypt(plainByte, key)
		fmt.Printf("復号文(string) %03d行目: %s\n", i, plainText)
		// 書き出す
		_, err02 := fmt.Fprintln(fw, plainText)
		utils.ErrCheck(err02)
	}
	utils.ErrCheck(s.Err())

	utils.EndLog()
}

// getReadWriteEnCrypt 暗号化用の読み込みファイルと書き込みファイルの struct を作る
func getReadWriteEnCrypt() crypto {
	readfile := "fileCrypto/plain/plainText.md"
	// ファイルが存在するか確認
	trainingIo.SearchFile(readfile)
	writefile := "fileCrypto/cipher/cipherText.md"
	cp := crypto{readpath: readfile, writepath: writefile}

	return cp
}

// getReadWriteDnCrypt 復号用の読み込みファイルと書き込みファイルの struct を作る
func getReadWriteDnCrypt() crypto {
	readfile := "fileCrypto/cipher/cipherText.md"
	// ファイルが存在するか確認
	trainingIo.SearchFile(readfile)
	writefile := "fileCrypto/decode/decodeText.md"
	cp := crypto{readpath: readfile, writepath: writefile}

	return cp
}

// getKey 外部ファイルから key を取得する
func getKey() []byte {
	keyfile := "fileCrypto/key/AES-CTR.key"
	// ファイルが存在するか確認
	trainingIo.SearchFile(keyfile)

	file, err01 := os.Open(keyfile)
	utils.ErrCheck(err01)
	defer file.Close()

	content := make([]byte, 16)
	_, err02 := file.Read(content)

	if err02 == io.EOF {
		return content
	} else {
		utils.ErrCheck(err02)
	}

	return content
}

// enCrypt 暗号化
func enCrypt(plainText string, key []byte) []byte {
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

// deCrypt 復号
func deCrypt(cipherByte []byte, key []byte) string {
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
