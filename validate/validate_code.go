package validate

import (
	"fmt"

	"github.com/wavly/shawty/utils"
)

type TooLong struct {
	lenght uint
}

type NotAlphaOrNum struct {
	text string
}

func (code *TooLong) Error() string {
	return fmt.Sprintf("Max lenght of the code is 8, but got %v", code.lenght)
}

func (code *NotAlphaOrNum) Error() string {
	return fmt.Sprintf("Only aplphabet characters and numbers are allowed in the code, but got %s", code.text)
}

func CustomCodeValidate(code string) error {
	if len(code) > 8 {
		return &TooLong{lenght: uint(len(code))}
	}

	if !utils.IsAplphabetOrNum(code) {
		return &NotAlphaOrNum{text: code}
	}

	return nil
}
