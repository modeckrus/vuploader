package changelog

import (
	"testing"
)

func TestVersion_Int64(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want int64
	}{
		{
			"Default",
			"1.2.3+456",
			int64(100200300456),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseVersion(tt.v); got != tt.want {
				t.Errorf("Version.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}
