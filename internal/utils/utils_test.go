package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {

	tests := []struct {
		name   string
		txtUrl string
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "Check URL is valid",
			txtUrl: "http://yandex.ru",
			want:   true,
		},
		{
			name:   "Check protocol",
			txtUrl: "yandex.ru",
			want:   false,
		},
		{
			name:   "Check URL",
			txtUrl: "dgdgdgdgdfgdfg",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ValidateURL(tt.txtUrl), tt.want)

		})
	}
}
