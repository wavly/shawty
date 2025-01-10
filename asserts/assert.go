package asserts

import (
	"os"

	prettylogger "github.com/wavly/surf/pretty-logger"
)

var logger = prettylogger.GetLogger(nil)

// Exists if error is not nil
func NoErr(err error, msg string) {
	if err != nil {
		logger.Error(msg, "error", err)
		os.Exit(1)
	}
}

// Exits if the check is true
func AssertEq(check bool, msg string) {
	if check {
		logger.Error(msg)
		os.Exit(1)
	}
}
