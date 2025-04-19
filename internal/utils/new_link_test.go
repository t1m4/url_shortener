package utils

import (
	"testing"
	"url_shortener/configs"
)

func FuzzGetNewLink(f *testing.F) {
	testcases := []int{0, 1, 2}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, lenght int) {
		if lenght < 0 {
			t.Skip()
		}
		actual := GetNewLink(lenght)
		if len(actual) != lenght {
			t.Errorf("ERORR %d, %d", len(actual), configs.UrlLenght)
		}
	})
}
