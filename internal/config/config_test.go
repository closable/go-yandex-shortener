package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstValue(t *testing.T) {
	tests := []struct {
		name string
		var1 string
		var2 string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "must v1",
			var1: "v1",
			var2: "",
			want: "v1",
		},
		{
			name: "must v2",
			var1: "",
			var2: "v2",
			want: "v2",
		},
		{
			name: "must v1",
			var1: "v1",
			var2: "v2",
			want: "v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, firstValue(&tt.var1, &tt.var2), tt.want)

		})
	}
}

func TestUpdateFromConfig(t *testing.T) {
	var config = &config{}
	assert.Empty(t, config.BaseURL)
	assert.Empty(t, config.ServerAddress)
	updateFromConfig(config)
	assert.NotEmpty(t, config.BaseURL)
	assert.NotEmpty(t, config.ServerAddress)

}
