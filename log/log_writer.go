package log

import (
  "fmt"
  "net/http"
)

type DbLogHandler struct {}

func (handler DbLogHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Log record was added with new handler %p!", &handler)
}
