package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {

	tests := []struct {
		name   string
		txtURL string
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "Check URL is valid",
			txtURL: "http://yandex.ru",
			want:   true,
		},
		{
			name:   "Check protocol",
			txtURL: "yandex.ru",
			want:   false,
		},
		{
			name:   "Check URL",
			txtURL: "dgdgdgdgdfgdfg",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ValidateURL(tt.txtURL), tt.want)

		})
	}
}

func TestGetRandomKey(t *testing.T) {

	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
		{
			name: "Add shortener",
			want: int(6), // default key length
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := GetRandomKey(tt.want)
			assert.Equal(t, len(shortener), tt.want)
		})
	}
}

func ExampleMakeServerAddres() {

	out1, _ := MakeServerAddres("localhost:8080", "true")
	fmt.Println(out1)

	out2, _ := MakeServerAddres("localhost:8080", "")
	fmt.Println(out2)

	out3, _ := MakeServerAddres("10.10.10.10", "true")
	fmt.Println(out3)

	//Output
	//:443
	//localhost:8080
	//10.10.10.10:443
}

func TestMakeServerAddres(t *testing.T) {

	tests := []struct {
		name      string
		txtURL    string
		flagHTTPS string
		want      string
	}{
		// TODO: Add test cases.
		{
			name:      "Check URL is valid",
			txtURL:    "localhost:8080",
			flagHTTPS: "true",
			want:      ":443",
		},
		{
			name:      "Check protocol",
			txtURL:    "localhost:8080",
			flagHTTPS: "",
			want:      "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := MakeServerAddres(tt.txtURL, tt.flagHTTPS)
			assert.Equal(t, res, tt.want)

		})
	}
}
