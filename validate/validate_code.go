package validate

import (
	"fmt"

	"github.com/wavly/shawty/utils"
)

type TooLong struct {
	lenght uint
}

type NotASCII struct {}

func (code *TooLong) Error() string {
	return fmt.Sprintf("Max lenght of the code is 8, but got %v", code.lenght)
}

func (_ *NotASCII) Error() string {
	return "Only ASCII characters are allowed in the code"
}

func CustomCodeValidate(code string) error {
	if len(code) > 8 {
		return &TooLong{lenght: uint(len(code))}
	} else if !utils.IsASCII(code) {
		return &NotASCII{}
	}
	return nil
}
