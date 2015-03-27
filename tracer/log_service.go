package tracer

import (
  "io"
  "fmt"
  "log"
  "time"
  "strconv"
  "net/http"
  "encoding/json"
  "encoding/base64"

  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

type DbLogger struct {
  connection *gorm.DB

  IdsCache *TerminalIdsCache
}

func (logger DbLogger) storeEntry(termianlId uint64, msg string) {
  encodedMsg := base64.StdEncoding.EncodeToString([]byte(msg))

  entry := LogEntry{TerminalId : termianlId,
                    Timestamp : time.Now().String(),
                    Message : encodedMsg}

  go logger.connection.Create(&entry)
  go logger.IdsCache.AppendId(termianlId)
}

func NewDbLogger(dbName string) (DbLogger, error) {
  conn, err := gorm.Open("sqlite3", dbName)
  return DbLogger{&conn, nil}, err
}

func (handler DbLogger) PrepareDB() {
  if !handler.connection.HasTable(&LogEntry{}) {
    handler.connection = handler.connection.CreateTable(&LogEntry{})
  }
}

func (handler DbLogger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if req.Method == "GET" {
    handler.GetLogs(res, req)
    return
  } else if req.Method == "POST" {
    handler.AddNewEntry(res, req)
    return
  } else {
    fmt.Fprintf(res, "404")
    res.WriteHeader(404)
    return
  }
}

func (handler DbLogger) AddNewEntry(res http.ResponseWriter,
                                    req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  type message struct {
    Msg string `json:"message"`
  }

  var msg message
  {
    err := decoder.Decode(&msg)
    if err != nil && err != io.EOF {
      log.Println("Unable to decode terminal message. ", err)
      res.WriteHeader(503)
      return
    }
  }

  vars := mux.Vars(req)
  id_str := vars["id"]
  
  if len(msg.Msg) == 0 {
    log.Println("Empty message! For terminal ", id_str)
    res.WriteHeader(503)
    return
  }

  id, err := strconv.ParseUint(id_str, 10, 64)
  if err != nil {
    log.Println(err)
    res.WriteHeader(503)
    return
  }

  go handler.storeEntry(id, msg.Msg)
}

func (handler DbLogger) GetLogs(res http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  id_str := vars["id"]
  id, err := strconv.ParseUint(id_str, 10, 64)
  if err != nil {
    log.Println(err)
    res.WriteHeader(503)
    return
  }

  type entryType struct {
    Timestamp string `json:"timestamp"`
    Message string `json:"message"`
  }
  type responseType struct {
    Entries []entryType `json:"entries"`
  }

  entries := responseType{}
  sync := make(chan bool)
  go func() {
    logEntries := make([]LogEntry, 0)
    handler.connection.Where("terminal_id = ?", id).Find(&logEntries)

    entries.Entries = make([]entryType, len(logEntries))

    for i, v := range logEntries {
      entries.Entries[i] = entryType{ v.Timestamp, v.Message }
    }

    sync <- true
  }()


  encoder := json.NewEncoder(res)
  <- sync
  encoder.Encode(entries)
}
