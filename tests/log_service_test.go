package tracercheck

import (
  _ "fmt"
  "io"
  "time"
  "bytes"
  "testing"
  "net/http"
  "encoding/json"
  "encoding/base64"
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
  messageText := "helo youg"
  jsonRequest := []byte(`{ "message" : "` + messageText + `" }`)

  req, e := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
  if  e != nil {
    t.Fatal(e)
  }

  go router.ServeHTTP(recorder, req)

  // Timeout is necessary due the server must have some
  // extra time for adding a new record.
  time.Sleep(5000)

  req, e = http.NewRequest("GET", url, nil)
  if e != nil {
    t.Fatal(e)
  }

  recorder = httptest.NewRecorder()

  router.ServeHTTP(recorder, req)

  t.Log("Got responce ", recorder.Body.String())
  
  type entryType struct {
    Timestamp string `json:"timestamp"`
    Message string `json:"message"`
  }
  type responseType struct {
    Entries []entryType `json:"entries"`
  }

  jsonDecoder := json.NewDecoder(recorder.Body)

  var logEntries responseType
  err = jsonDecoder.Decode(&logEntries)
  if err != nil && err != io.EOF {
    t.Fatal(err)
  }

  if len(logEntries.Entries) != 1 {
    t.Fatal("Log entries length doesn`t equal")
  }

  for i, v := range logEntries.Entries {
    decodeBytes, err := base64.StdEncoding.DecodeString(v.Message)
    if err != nil {
      t.Fatal(err)
    }

    logEntries.Entries[i].Message = string(decodeBytes)
  }

  t.Log("Decoded json ", logEntries)

  if logEntries.Entries[0].Message != messageText {
    t.Fatal("Message text is wrong!\n Got ", logEntries.Entries[0].Message,
            "\nShould be ", messageText)
  }
}
