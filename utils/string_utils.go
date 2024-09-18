package utils

import (
	"unicode"
	"unicode/utf8"
)

func IsAplphabet(text string) bool {
	for _, s := range text {
		if !unicode.IsLetter(s) {
			return false
		}
	}
	return true
}

func IsASCII(text string) bool {
	return utf8.RuneCount([]byte(text)) == len(text)
}
