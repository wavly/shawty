package utils

import (
	"unicode"
	"unicode/utf8"
)

func IsAlphabet(text string) bool {
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

func IsAlphabetOrNum(text string) bool {
	for _, c := range text {
		if !IsValidChar(c) {
			return false
		}
	}

	return true
}

func IsValidChar(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9'
}
