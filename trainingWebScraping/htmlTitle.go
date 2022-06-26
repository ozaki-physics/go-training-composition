package trainingWebScraping

import (
	"regexp"
)
func searchTitle(html string) string {
	r01 := regexp.MustCompile("<title.*>.*</title.*")
	str := r01.FindString(html)
	r02 := regexp.MustCompile("</title.*>")
	str = r02.ReplaceAllString(str, "")
	r03 := regexp.MustCompile("<title.*>")
	str = r03.ReplaceAllString(str, "")
	return str
}
