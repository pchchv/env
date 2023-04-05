package env

import (
	"bytes"
	"unicode"
)

func indexOfNonSpaceChar(src []byte) int {
	return bytes.IndexFunc(src, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
}

func isChar(char rune) func(rune) bool {
	return func(v rune) bool {
		return v == char
	}
}
