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

func newBadAddLogRecordRequest(terminalId uint64, t *testing.T) *http.Request {
  url := "/terminal/" + strconv.FormatUint(terminalId, 10)
  jsonRequest := []byte(`{ "mssage" : "aslkdj" }`)

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

func sendRequestToServer(srv *mux.Router,
                         req *http.Request,
                         expectedCode int,
                         t *testing.T) *httptest.ResponseRecorder {
  recorder := httptest.NewRecorder()
  srv.ServeHTTP(recorder, req)

  if recorder.Code != expectedCode {
    t.Fatal("Server returned code ", recorder.Code)
  }

  return recorder
}

func decodeMessages(logResponse *responseType, t *testing.T) {
  for i, v := range logResponse.Entries {
    decodeBytes, err := base64.StdEncoding.DecodeString(v.Message)
    if err != nil {
      t.Fatal(err)
    }

    logResponse.Entries[i].Message = string(decodeBytes)
  }
}

func TestSendRecvMessage(t *testing.T) {
  var termId uint64 = 177
  msg := "helo youg"

  testServer := newTestEnv(t)
  req := newAddLogRecordRequest(msg, termId, t)

  go sendRequestToServer(testServer, req, http.StatusOK, t)

  // Timeout is necessary due the server must have some
  // extra time for adding a new record.
  time.Sleep(5000)

  req = newGetLogsRequest(termId, t)

  recorder := sendRequestToServer(testServer, req, http.StatusOK, t)
  t.Log("Got response ", recorder.Body.String())

  logEntries := decodeResponse(recorder.Body, t)

  if len(logEntries.Entries) != 1 {
    t.Fatal("Log entries length doesn`t equal")
  }

  decodeMessages(logEntries, t)

  t.Log("Decoded json ", *logEntries)

  if logEntries.Entries[0].Message != msg {
    t.Fatal("Message text is wrong!\n Got ", logEntries.Entries[0].Message,
            "\nShould be ", msg)
  }
}

func TestLogBadMessage(t *testing.T) {
  var id uint64 = 17
  
  testServer := newTestEnv(t)
  req := newBadAddLogRecordRequest(id, t)

  recorder := httptest.NewRecorder()
  testServer.ServeHTTP(recorder, req)

  if recorder.Code != http.StatusBadRequest {
    t.Fatal("Wrong response status. Expected 400. Got ", recorder.Code)
  }

  t.Log("Response code is ", recorder.Code)
}

func TestSeveralMessages(t *testing.T) {
  msgCount := 10
  msg := "Dude roll it over"
  var termId uint64 = 604

  testServer := newTestEnv(t)

  for i := 0; i < msgCount; i++ {
    req := newAddLogRecordRequest(msg, termId, t)
    sendRequestToServer(testServer, req, http.StatusOK, t)
  }

  // Timeout is necessary due the server must have some
  // extra time for adding a new record.
  time.Sleep(5000)

  req := newGetLogsRequest(termId, t)
  recorder := sendRequestToServer(testServer, req, http.StatusOK, t)

  t.Log("Got response ", recorder.Body.String())

  logEntries := decodeResponse(recorder.Body, t)
  if len(logEntries.Entries) != msgCount {
    t.Fatal("Log entries length doesn`t equal")
  }

  decodeMessages(logEntries, t)

  t.Log("Decoded json ", *logEntries)
  for _, v := range logEntries.Entries {
    if v.Message != msg {
      t.Fatal("Message text is wrong!\n Got ", v.Message,
              "\nShould be ", msg)
    }
  }
}
