package filesystem

import "fmt"

func FormatStringForCsv(s string) string {
	s = fmt.Sprintf("%q", s)
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
