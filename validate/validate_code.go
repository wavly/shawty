package validate

import (
	"fmt"

	"github.com/wavly/surf/utils"
)

type TooLong struct {
	length uint
}

type TooShort struct {
	length uint
}

type NotAlphaOrNum struct {
	text string
}

func (code *TooLong) Error() string {
	return fmt.Sprintf("Max length of the code is 8, but got %v", code.length)
}

func (_ *TooShort) Error() string {
	return "Min length of the code is 2"
}

func (code *NotAlphaOrNum) Error() string {
	return fmt.Sprintf("Only alphabetical characters and numbers are allowed in the code, but got %s", code.text)
}

func CustomCodeValidate(code string) error {
	if len(code) > 8 {
		return &TooLong{length: uint(len(code))}
	} else if len(code) < 2 {
		return &TooShort{}
	}

	if !utils.IsAlphabetOrNum(code) {
		return &NotAlphaOrNum{text: code}
	}

	return nil
}
