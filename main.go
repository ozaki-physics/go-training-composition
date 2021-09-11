package main

import (
	"fmt"
	"github.com/ozaki-physics/go-training-composition/fileCrypto"
	"github.com/ozaki-physics/go-training-composition/package01"
	"github.com/ozaki-physics/go-training-composition/package02"
	"github.com/ozaki-physics/go-training-composition/trainingCrypto"
	"github.com/ozaki-physics/go-training-composition/trainingIo"
	"github.com/ozaki-physics/go-training-composition/trainingJson"
	"github.com/ozaki-physics/go-training-composition/trainingTimeZone"
	"github.com/ozaki-physics/go-training-composition/utils"
	"log"
)

// init システム的に main 関数の前に実行される初期化関数
func init() {
	utils.InitTimeZone()
}

func main() {
	// mainPkg()
	// mainTimeZone()
	// mainCrypto()
	// ioFileVersion()
	// ioTerminalVersion()
	mainFileCrypto()
	// mainJson()
}

// mainPkg ディレクトリ構成を試す
func mainPkg() {
	fmt.Println("Hello World!")
	fmt.Println(package01.Const_big)
	// これはエラー
	// fmt.Println(package01.const_small)
	fmt.Println(package02.Const_big)
	// package02.Sample_server()
}

// fileVersion ファイル入出力の勉強
func ioFileVersion() {
	utils.InitLog("[入出力の実験]")
	log.Println("ファイル版 開始")
	f := trainingIo.ThisDirFile{"trainingIo/read.md", "trainingIo/write.md"}
	fmt.Printf("f := %v\n", f)
	trainingIo.SearchFile(f.ReadName)
	trainingIo.OpenFile(f.ReadName)
	trainingIo.OpenFile02(f.WriteName)
	trainingIo.CreateFile("trainingIo/aaa")
	trainingIo.AllDataReadFileName(f.ReadName)
	trainingIo.AllDataReadFile(f.ReadName)
	trainingIo.DataReadFile(f.ReadName)
	trainingIo.ScanDataReadFile(f.ReadName)
	trainingIo.CopyFile(f.ReadName, f.WriteName)
	trainingIo.RenameFile("trainingIo/aaa", "trainingIo/bbb")
	trainingIo.WriteAllData(f.WriteName)
	trainingIo.WriteDataFile(f.WriteName)
	trainingIo.WriteDataWriter(f.WriteName)
	trainingIo.Reader01DataReadFile(f.ReadName)
	trainingIo.Reader02DataReadFile(f.ReadName)
	log.Println("ファイル版 終了")
}

// terminalVersion ターミナル入出力の勉強
func ioTerminalVersion() {
	utils.InitLog("[入出力の実験]")
	log.Println("ターミナル版 開始")
	trainingIo.OutTerminal()
	trainingIo.UseFmtScan()
	// trainingIo.UseBufioScanner01() // go run example_ioFile.go < trainingIo/read.md が正常に終了する
	trainingIo.UseBufioScanner02()
	trainingIo.UseBufioReader()
	trainingIo.TerminalArgsFlag()
	trainingIo.TerminalArgsOs()
	// trainingIo.TerminalArgsFile()
	log.Println("ターミナル版 終了")
}

// mainTimeZone タイムゾーンの勉強
func mainTimeZone() {
	trainingTimeZone.MainTimeZone()
}

// mainCrypto 暗号理論の勉強
func mainCrypto() {
	utils.InitLog("[暗号化の実験]")
	utils.StartLog()
	// 共通鍵暗号方式 ブロック暗号化方式の AES
	trainingCrypto.Example01()
	// 共通鍵暗号方式 ブロック暗号化方式の AES CBC モード
	trainingCrypto.Example02()
	// 共通鍵暗号方式 ブロック暗号化方式の AES CTR モード だから ストリーム暗号とみなせる
	trainingCrypto.Example03()
	// 公開鍵暗号方式 RSA-PKCS1v15 で暗号化
	trainingCrypto.Example04()
	// ハッシュ SHA-2 の SHA-512
	trainingCrypto.Example05()
	// メッセージ認証コード(MAC) 否認ができず 送信者の証明ができない
	trainingCrypto.Example06()
	// デジタル署名 公開鍵暗号の応用なため 今回は 楕円曲線暗号 を使う
	trainingCrypto.Example07()
	// 証明書(x509) 自己署名証明書を作ってみる
	trainingCrypto.Example08()
	// TLS
	trainingCrypto.Example09()
	// パスワード
	trainingCrypto.Example10()
	utils.EndLog()
}

// mainFileCrypto ファイルの中身を暗号化するツール
func mainFileCrypto() {
	fileCrypto.RunFileEnCrypto()
	fileCrypto.RunFileDeCrypto()
}

// mainJson JSON 読み込みの勉強
func mainJson() {
	utils.InitLog("[JSONの実験]")
	utils.StartLog()
	// trainingJson.Example()
	trainingJson.ReadJson01()
	utils.EndLog()
}
