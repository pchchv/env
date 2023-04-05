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
