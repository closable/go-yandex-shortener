package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShortener(t *testing.T) {

	tests := []struct {
		name   string
		txtURL string
		want   int
	}{
		// TODO: Add test cases.
		{
			name:   "Add shortener",
			txtURL: "http://yandex.ru",
			want:   int(6), // default key length
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := GetShortener(tt.txtURL)
			assert.Equal(t, len(shortener), tt.want)
		})
	}
}

func TestFindKeyByValue(t *testing.T) {
	Storage.Urls["abcde"] = "http://mail.ru"
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
			key := FindKeyByValue(tt.txtURL)
			assert.Equal(t, key, tt.want)

		})
	}
}
