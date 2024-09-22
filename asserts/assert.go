package asserts

import "log"

// logs the provided message and error, then exits if the error is not nil.
func NoErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
