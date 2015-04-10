package tracer

import (
  "io"
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

func (logger DbLogger) buildLogs(res http.ResponseWriter,
                                 term_id chan uint64,
                                 start_time chan string,
                                 quit_signal chan bool) {
  type entryType struct {
    Timestamp string `json:"timestamp"`
    Message string `json:"message"`
  }
  type responseType struct {
    Entries []entryType `json:"entries"`
  }

  entries := responseType{}

  var id uint64

  select {
    case id = <-term_id:
      break
    case <- quit_signal:
      return
    case _ = (<-start_time):
      break
  }

  sync := make(chan bool)
  go func() {
    logEntries := make([]LogEntry, 0)
    logger.connection.Where("terminal_id = ?", id).Find(&logEntries)

    entries.Entries = make([]entryType, len(logEntries))

    for i, v := range logEntries {
      entries.Entries[i] = entryType{ v.Timestamp, v.Message }
    }

    sync <- true
  }()


  encoder := json.NewEncoder(res)
  <- sync
  encoder.Encode(entries)
  quit_signal <- true
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
    res.WriteHeader(http.StatusBadRequest)
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
      res.WriteHeader(http.StatusBadRequest)
      return
    }
  }

  vars := mux.Vars(req)
  id_str := vars["id"]
  
  if len(msg.Msg) == 0 {
    log.Println("Empty message! For terminal ", id_str)
    res.WriteHeader(http.StatusBadRequest)
    return
  }

  id, err := strconv.ParseUint(id_str, 10, 64)
  if err != nil {
    log.Println(err)
    res.WriteHeader(http.StatusInternalServerError)
    return
  }

  go handler.storeEntry(id, msg.Msg)
}

func (handler DbLogger) GetLogs(res http.ResponseWriter, req *http.Request) {
  id_sender := make(chan uint64)
  time_sender := make(chan string)
  quit_signal := make(chan bool)

  go handler.buildLogs(res, id_sender, time_sender, quit_signal)

  vars := mux.Vars(req)
  id_str := vars["id"]
  id, err := strconv.ParseUint(id_str, 10, 64)
  if err != nil {
    log.Println(err)
    res.WriteHeader(http.StatusInternalServerError)
    quit_signal <- true
    return
  }

  id_sender <- id
  <- quit_signal
}
