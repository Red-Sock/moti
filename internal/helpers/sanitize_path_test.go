package helpers

import (
	"testing"
)

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "simple path",
			source: "foo/bar",
			want:   "foo-bar",
		},
		{
			name:   "nested path",
			source: "github.com/redsock/moti",
			want:   "github.com-redsock-moti",
		},
		{
			name:   "no slash",
			source: "moti",
			want:   "moti",
		},
		{
			name:   "empty",
			source: "",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizePath(tt.source); got != tt.want {
				t.Errorf("SanitizePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
