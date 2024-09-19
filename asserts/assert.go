package asserts

import (
    "log"
)

// NoErr logs a fatal error message along with any provided error details.
//
// This function checks if an error is non-nil. If the error is non-nil, it logs both
// a custom message and the error details, then terminates the program.
//
// Parameters:
// - err: The error to check.
// - msg: A custom error message.
// - log.Fatalf: Logs the fatal error message details and in this case your custom message and the error
func NoErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %v", msg, err)
    }
}