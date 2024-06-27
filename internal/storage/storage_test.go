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

func TestPing(t *testing.T) {
	store, _ := NewMemory()
	isPing := store.Ping()
	assert.Equal(t, true, isPing)
}

func TestFindExistingKey(t *testing.T) {
	store, _ := NewMemory()
	store.Urls["abcde"] = "http://mail.ru"
	store.Urls["abcd1"] = "http://mail1.ru"
	store.Urls["abcd2"] = "http://mail2.ru"

	tests := []struct {
		name string
		key  string
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Find key",
			key:  "abcde",
			want: true,
		},
		{
			name: "Find key",
			key:  "abcd1",
			want: true,
		},
		{
			name: "Not found key",
			key:  "***",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := store.FindExistingKey(tt.key)
			assert.Equal(t, found, tt.want)
		})
	}
}

func TestGetShortener(t *testing.T) {
	store, _ := NewMemory()
	store.Urls["abcde"] = "http://mail.ru"
	store.Urls["abcd1"] = "http://mail1.ru"
	store.Urls["abcd2"] = "http://mail2.ru"

	tests := []struct {
		name string
		url  string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Find key",
			url:  "http://mail.ru",
			want: "abcde",
		},
		{
			name: "Find key",
			url:  "http://mail1.ru",
			want: "abcd1",
		},
		{
			name: "Not found key",
			url:  "http://yandex.ru",
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != "test" {
				key, _ := store.GetShortener(1, tt.url)
				assert.Equal(t, key, tt.want)
			} else {
				key, _ := store.GetShortener(1, tt.url)
				assert.NotEqual(t, key, tt.want)
				assert.NotEmpty(t, key)
			}

		})
	}
}

func TestLength(t *testing.T) {
	store, _ := NewMemory()
	store.Urls["abcde"] = "http://mail.ru"
	store.Urls["abcd1"] = "http://mail1.ru"
	store.Urls["abcd2"] = "http://mail2.ru"
	assert.Equal(t, store.Length(), 3)

	store.Urls["abcd3"] = "http://mail3.ru"
	assert.Equal(t, store.Length(), 4)

}

func TestAddItem(t *testing.T) {
	store, _ := NewMemory()
	store.Urls["abcde"] = "http://mail.ru"
	assert.Equal(t, store.Length(), 1)
	store.AddItem("1", "test1")
	assert.Equal(t, store.Length(), 2)

	_, found := store.FindExistingKey("1")
	assert.Equal(t, found, true)

	assert.Equal(t, store.Length(), 2)
}
