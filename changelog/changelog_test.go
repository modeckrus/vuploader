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
			int64(1002003456),
		},
		{
			"Decimal",
			"1.2.30+456",
			int64(1002030456),
		},
		{
			"Decimal",
			"1.20.30+456",
			int64(1020030456),
		},
		{
			"Sent",
			"1.20.300+456",
			int64(1020300456),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.v); got != tt.want {
				t.Errorf("Version.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}
