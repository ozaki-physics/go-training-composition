package trainingCompress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

// compressGzip gzip で圧縮する
func compressGzip(c string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write([]byte(c))
	if err != nil {
		return nil, err
	}
	defer zw.Close()

	return &buf, nil
}

// extractGzip gzip で逆圧縮する
func decompressGzip(b io.Reader) (io.Reader, error) {
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return r, nil
}

// MainGzip 圧縮して逆圧縮する
// See: https://pkg.go.dev/compress/gzip@go1.17.8#NewWriterLevel
func MainGzip() {
	fmt.Println("gzip による圧縮")
	// 圧縮
	zr, err := compressGzip(content)
	if err != nil {
		log.Println(err)
	}
	b := zr.Bytes()
	fmt.Printf("%d bytes: %v\n", len(b), b)
	// 逆圧縮
	r, err := decompressGzip(zr)
	if err != nil {
		log.Println(err)
	}
	io.Copy(os.Stdout, r)
}
