package storage

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKeyByValue(t *testing.T) {
	store, _ := NewMemory()
	store.Urls["abcde"] = "http://mail.ru"
	tests := []struct {
		name   string
		txtURL string
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "Find key",
			txtURL: "http://mail.ru",
			want:   "abcde",
		},
		{
			name:   "key not found",
			txtURL: "http://yandex.ru/1233",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := store.FindKeyByValue(tt.txtURL)
			assert.Equal(t, key, tt.want)

		})
	}
}

func BenchmarkFindKeyByValue(b *testing.B) {
	store, _ := NewMemory()
	store.Urls["abcd0"] = "http://0mail.ru"
	store.Urls["abcd1"] = "http://1mail.ru"
	store.Urls["abcd2"] = "http://2mail.ru"
	store.Urls["abcd3"] = "http://3mail.ru"
	store.Urls["abcd4"] = "http://4mail.ru"
	store.Urls["abcd5"] = "http://5mail.ru"
	store.Urls["abcd6"] = "http://6mail.ru"
	store.Urls["abcd7"] = "http://7mail.ru"
	store.Urls["abcd8"] = "http://8mail.ru"
	store.Urls["abcd9"] = "http://98mail.ru"

	for i := 0; i < b.N; i++ {
		store.FindKeyByValue(getKeyforBench())
	}
}

func getKeyforBench() string {
	r := rand.Intn(9)
	return fmt.Sprintf("abcd%d", r)
}
