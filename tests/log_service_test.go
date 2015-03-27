package tracercheck

import (
  "io"
  "time"
  "bytes"
  "testing"
  "strconv"
  "net/http"
  "encoding/json"
  "encoding/base64"
  "net/http/httptest"

  "github.com/gorilla/mux"

  "github.com/r4start/web-tracer/tracer"
)

type entryType struct {
  Timestamp string `json:"timestamp"`
  Message string `json:"message"`
}

type responseType struct {
  Entries []entryType `json:"entries"`
}

func newTestEnv(t *testing.T) *mux.Router {
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

  return router
}

func newAddLogRecordRequest(message string,
                            terminalId uint64,
                            t *testing.T) *http.Request {
  url := "/terminal/" + strconv.FormatUint(terminalId, 10)
  jsonRequest := []byte(`{ "message" : "` + message + `" }`)

  req, e := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
  if  e != nil {
    t.Fatal(e)
  }

  return req
}

func newGetLogsRequest(terminalId uint64,
                       t *testing.T) *http.Request {
  url := "/terminal/" + strconv.FormatUint(terminalId, 10)

  req, e := http.NewRequest("GET", url, nil)
  if e != nil {
    t.Fatal(e)
  }

  return req
}

func decodeResponse(data io.Reader, t *testing.T) *responseType {
  jsonDecoder := json.NewDecoder(data)

  var logEntries responseType
  err := jsonDecoder.Decode(&logEntries)
  if err != nil && err != io.EOF {
    t.Fatal(err)
  }

  return &logEntries
}

func TestSendRecvMessage(t *testing.T) {
  var termId uint64 = 177
  msg := "helo youg"
  testServer := newTestEnv(t)
  req := newAddLogRecordRequest(msg, termId, t)

  recorder := httptest.NewRecorder()
  go testServer.ServeHTTP(recorder, req)

  // Timeout is necessary due the server must have some
  // extra time for adding a new record.
  time.Sleep(5000)

  req = newGetLogsRequest(termId, t)

  recorder = httptest.NewRecorder()

  testServer.ServeHTTP(recorder, req)

  t.Log("Got response ", recorder.Body.String())

  logEntries := decodeResponse(recorder.Body, t)

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

  t.Log("Decoded json ", *logEntries)

  if logEntries.Entries[0].Message != msg {
    t.Fatal("Message text is wrong!\n Got ", logEntries.Entries[0].Message,
            "\nShould be ", msg)
  }
}

func TestLogBadMessage(t *testing.T) {
  // testServer := newTestEnv(t)
}
