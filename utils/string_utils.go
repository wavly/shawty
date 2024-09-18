package utils

import "unicode"

func IsAplphabet(text string) bool {
	for _, s := range text {
		if !unicode.IsLetter(s) {
			return false
		}
	}
	return true
}

func IsASCII(text string) bool {
	for i := 0; i < len(text); i++ {
		if text[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
