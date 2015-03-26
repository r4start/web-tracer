package tracercheck

import (
  _ "fmt"
  "bytes"
  "testing"
  "net/http"
  "net/http/httptest"

  "github.com/r4start/web-tracer/tracer"
)

func TestLogSingleMessage(t *testing.T) {
  dbName := ""
  logger, err := tracer.NewDbLogger(dbName)
  if err != nil {
    t.Fatal(err);
  }

  logger.PrepareDB()

  recorder := httptest.NewRecorder()
  url := "http://localhost/terminal/177"
  jsonRequest := []byte(`{ "message" : "helo youg" }`)

  req, e := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
  if  e != nil {
    t.Fatal(e)
  }

  logger.ServeHTTP(recorder, req)

  req, e = http.NewRequest("GET", url, nil)
  if e != nil {
    t.Fatal(e)
  }

  recorder = httptest.NewRecorder()
  logger.ServeHTTP(recorder, req)
}
