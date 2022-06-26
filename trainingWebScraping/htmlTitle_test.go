package trainingWebScraping

import (
	"fmt"
	"testing"
)

func TestSearchTitle(t *testing.T) {

	patterns := []struct {
		html     string
		expected string
	}{
		{"<title>hello</title>", "hello"},
		{"<title img=\"aaa\">hello</title data-aaa=\"aaa\">", "hello"},
	}

	for idx, pattern := range patterns {
		t.Run(fmt.Sprintf("いけた?"), func(t *testing.T) {
			actual := searchTitle(pattern.html)

			if actual != pattern.expected {
				t.Errorf("インデックス %d で期待は %s なのに実際は %s だった", idx, pattern.expected, actual)
			}
		})
	}
}
