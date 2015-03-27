package tracercheck

import (
  _ "fmt"
  "time"
  "bytes"
  "testing"
  "net/http"
  _ "encoding/base64"
  "net/http/httptest"

  "github.com/gorilla/mux"

  "github.com/r4start/web-tracer/tracer"
)

func TestLogSingleMessage(t *testing.T) {
  dbName := ":memory:"
  logger, err := tracer.NewDbLogger(dbName)
  if err != nil {
    t.Fatal(err);
  }

  var idsCache = tracer.NewTerminalIdsCache()
  logger.IdsCache = &idsCache

  logger.PrepareDB()

  router := mux.NewRouter()
  router.Handle("/terminal/{id:[0-9]+}", logger)

  recorder := httptest.NewRecorder()
  url := "/terminal/177"
  jsonRequest := []byte(`{ "message" : "helo youg" }`)

  req, e := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
  if  e != nil {
    t.Fatal(e)
  }

  go router.ServeHTTP(recorder, req)

  // Timeout necessary, because the server must have some
  // extra time for addind a new record.
  time.Sleep(5000)

  req, e = http.NewRequest("GET", url, nil)
  if e != nil {
    t.Fatal(e)
  }

  recorder = httptest.NewRecorder()

  router.ServeHTTP(recorder, req)

  t.Log("Request and responce are not equal! ",
        string(jsonRequest),
        recorder.Body.String())
  
}
