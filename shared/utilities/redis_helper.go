package utilities

import "strings"

func GetFullKey(keys ...string) string {
	return strings.Join(keys, ":")
}
