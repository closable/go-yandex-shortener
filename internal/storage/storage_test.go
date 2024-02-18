package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKeyByValue(t *testing.T) {
	store := New()
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
