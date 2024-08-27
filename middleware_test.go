package main

import (
	"regexp"
	"testing"
)

func TestSkipper(t *testing.T) {

	url := "http://127.0.0.1:8000/docs/index.html"

	skips := []string{"docs"}

	for _, pattern := range skips {
		re := regexp.MustCompile(pattern)

		isMatch := re.MatchString(url)

		if !isMatch {
			t.Fatal("not match")
		}
	}
}
