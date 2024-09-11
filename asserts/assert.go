package asserts

import "log"

func NoErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}
