package validate

import (
	"fmt"

	"github.com/wavly/shawty/utils"
)

type TooLong struct {
	lenght uint
}

type TooShort struct {
	lenght uint
}

type NotAplphaOrNum struct {
	text string
}

func (code *TooLong) Error() string {
	return fmt.Sprintf("Max lenght of the code is 8, but got %v", code.lenght)
}

func (_ *TooShort) Error() string {
	return "Min lenght of the code is 2"
}

func (code *NotAplphaOrNum) Error() string {
	return fmt.Sprintf("Only aplphabet characters and numbers are allowed in the code, but got %s", code.text)
}

func CustomCodeValidate(code string) error {
	if len(code) > 8 {
		return &TooLong{lenght: uint(len(code))}
	} else if len(code) < 2 {
		return &TooShort{}
	}

	if !utils.IsAplphabetOrNum(code) {
		return &NotAplphaOrNum{text: code}
	}

	return nil
}
