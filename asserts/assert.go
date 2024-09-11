package asserts

import "log"

// NoErr logs a fatal error message along with any provided error details.
//
// Checks if an error is non-nil and logs both a custom message and the error details,
// then terminates the program.
func NoErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}
