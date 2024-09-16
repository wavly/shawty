package utils

import "unicode"

func IsAcsii(text string) bool {
	for _, s := range text {
		if !unicode.IsLetter(s) {
			return false
		}
	}
	return true
}
