// See: Go 言語でデータ圧縮と解凍 https://text.baldanders.info/golang/compress-data/
package trainingCompress

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"log"
	"os"
)

// const content = "Hello world\n"
const content = "AAAAAAAAAA,AAAAAAAAAA,AAAAAAAAAA,AAAAAAAAAA,AAAAAAAAAA\n"

// NormalByte 圧縮してないときのバイト
func NormalByte() {
	fmt.Println("圧縮してない")
	b := []byte(content)
	fmt.Printf("%d bytes: %v\n", len(b), b)
	fmt.Printf(content)
}

// compressDeflate deflate で圧縮
func compressDeflate(c string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	zw, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	zw.Write([]byte(c))
	defer zw.Close()

	return &buf, nil
}

// decompressDeflate deflate で逆圧縮
func decompressDeflate(b io.Reader) (io.Reader, error) {
	r := flate.NewReader(b)
	defer r.Close()
	return r, nil
}

func MainDeflate() {
	fmt.Println("deflate による圧縮")
	// 圧縮
	zr, err := compressDeflate(content)
	if err != nil {
		log.Println(err)
	}
	b := zr.Bytes()
	fmt.Printf("%d bytes: %v\n", len(b), b)
	// 逆圧縮
	r, err := decompressDeflate(zr)
	if err != nil {
		log.Println(err)
	}
	io.Copy(os.Stdout, r)
}
