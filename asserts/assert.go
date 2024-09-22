package asserts

import "log"

// logs the provided message and error, then exits if the error is not nil.
func NoErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// logs the message and exits if the boolean check is true
func AssertEq(check bool, msg ...any) {
	if check {
		log.Fatalln(msg...)
	}
}
