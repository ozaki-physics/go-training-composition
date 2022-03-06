package trainingWebScraping

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func Main() {
	// url := "https://www.ozaki-physics.com"
	url := "http://www.google.com/robots.txt"
	html := getHtml(url)
	fmt.Println(html)
}

// getHtml
// 公式のサンプルを少し改良した
// see: 公式によるサンプル https://pkg.go.dev/net/http@go1.17.7#example-Get
func getHtml(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyByte, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, bodyByte)
	}
	if err != nil {
		log.Fatal(err)
	}

	// return string(bodyByte)
	return strconv.Itoa(res.StatusCode)
}
