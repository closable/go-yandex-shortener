package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileFindKeyByValue(t *testing.T) {
	store, _ := NewFile("/tmp/short-url-db.json")
	sh1, _ := store.GetShortener(1, "http://main.ru")
	sh2, _ := store.GetShortener(1, "http://yandex.ru")

	tests := []struct {
		name   string
		txtURL string
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "Find key",
			txtURL: "http://main.ru",
			want:   sh1,
		},
		{
			name:   "Find key",
			txtURL: "http://yandex.ru",
			want:   sh2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := store.FindKeyByValue(tt.txtURL)
			assert.Equal(t, key, tt.want)

		})
	}
}

func TestFileFindExistingKey(t *testing.T) {
	store, _ := NewFile("/tmp/short-url-db.json")
	sh1, _ := store.GetShortener(1, "http://main.ru")
	sh2, _ := store.GetShortener(1, "http://yandex.ru")

	tests := []struct {
		name string
		key  string
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Find key",
			key:  sh1,
			want: true,
		},
		{
			name: "Find key",
			key:  sh2,
			want: true,
		},
		{
			name: "Find key false",
			key:  "12345",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, find := store.FindExistingKey(tt.key)
			assert.Equal(t, find, tt.want)

		})
	}
}

func TeasFileGetShortener(t *testing.T) {
	store, _ := NewFile("/tmp/short-url-db.json")
	tests := []struct {
		name  string
		url   string
		short bool
	}{
		// TODO: Add test cases.
		{
			name:  "Get shorterner",
			url:   "http://mail.ru",
			short: true,
		},
		{
			name:  "Get shorterner",
			url:   "http://yandex.ru",
			short: true,
		},
		// {
		// 	name: "Find key false",
		// 	key:  "12345",
		// 	want: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh, _ := store.GetShortener(1, tt.url)
			assert.Equal(t, len(sh) > 0, tt.short)

		})
	}
}
