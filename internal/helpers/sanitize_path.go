package helpers

import (
	"strings"
)

func SanitizePath(source string) string {
	return strings.ReplaceAll(source, "/", "-")
}
