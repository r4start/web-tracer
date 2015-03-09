package log

import (
  "fmt"
  "net/http"
)

func AddToLog(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Log record was added!")
}
