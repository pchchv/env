package env

import (
	"bytes"
	"unicode"
)

const charComment = '#'

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

// getStatementPosition returns the start position of the statement.
// It skips any comment string or character that is not a space.
func getStatementStart(src []byte) []byte {
	pos := indexOfNonSpaceChar(src)
	if pos == -1 {
		return nil
	}

	src = src[pos:]
	if src[0] != charComment {
		return src
	}

	// skip comment section
	pos = bytes.IndexFunc(src, isChar('\n'))
	if pos == -1 {
		return nil
	}

	return getStatementStart(src[pos:])
}

// isSpace tells whether the rune is a space character but not a line feed character,
// which is different from unicode.IsSpace, which also applies line feed as a space.
func isSpace(r rune) bool {
	switch r {
	case '\t', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

func isLineEnd(r rune) bool {
	if r == '\n' || r == '\r' {
		return true
	}
	return false
}
