package extractor

import "strings"

func clean(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
