package trainingCompress

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
)

// compressZlib zlib で圧縮する
// 日本では 1980年代に 圧縮を 凍結という流れもあったが普及しなかった(でも解凍だけは残ってしまった)
// See: https://pkg.go.dev/compress/zlib@go1.17.8#example-NewWriter
func compressZlib(c string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return &buf, err
	}
	zw.Write([]byte(c))
	defer zw.Close()

	return &buf, nil
}

// extractZlib zlib を逆圧縮する
// 日本では 解凍, 展開, 抽出 など定まっていない(英語でも 抽出の extract を使うことはあるらしい)
func decompressZlib(b io.Reader) (io.Reader, error) {
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return r, nil
}

// MainZlib 圧縮して逆圧縮する
func MainZlib() {
	fmt.Println("zlib による圧縮")
	// 圧縮
	zr, err := compressZlib(content)
	if err != nil {
		log.Println(err)
	}
	b := zr.Bytes()
	fmt.Printf("%d bytes: %v\n", len(b), b)

	// 逆圧縮
	r, err := decompressZlib(zr)
	if err != nil {
		log.Println(err)
	}
	io.Copy(os.Stdout, r)
}
