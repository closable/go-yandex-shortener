package utils

import (
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
