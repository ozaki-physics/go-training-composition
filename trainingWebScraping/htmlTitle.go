package trainingWebScraping

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func GetTitle() {
	s := []string{
		"https://github.com/ozaki-physics",
	}

	for _, target := range s {
		a := getHtml(target)
		aa := fmt.Sprintf("[%s](%s)", a, target)
		fmt.Println(aa)
	}
}

func getHtml(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyByte, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		fmt.Printf("Response の status code: %d, 該当の url: %s\n", res.StatusCode, url)
	}
	if err != nil {
		log.Fatal(err)
	}

	html := string(bodyByte)
	return searchTitle(html)
}

func searchTitle(html string) string {
	r01 := regexp.MustCompile("<title.*>.*</title.*")
	str := r01.FindString(html)
	r02 := regexp.MustCompile("</title.*>")
	str = r02.ReplaceAllString(str, "")
	r03 := regexp.MustCompile("<title.*>")
	str = r03.ReplaceAllString(str, "")
	return str
}
