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
                                 termId <- chan uint64,
                                 startTime <- chan string,
                                 quitSignal chan bool) {
  type entryType struct {
    Timestamp string `json:"timestamp"`
    Message string `json:"message"`
  }
  type responseType struct {
    Entries []entryType `json:"entries"`
  }

  entries := responseType{}

  var id uint64
  var timePoint string

  var status int = 0
  const allGot = 2

Loop:
  for {
    select {
      case id = <-termId:
        status += 1
        if status == allGot {
          break Loop
        }
        continue

      case <- quitSignal:
        return
      
      case timePoint = <-startTime:
        status += 1
        if status == allGot {
          break Loop
        }
        continue
    }
  }

  sync := make(chan bool)
  go func() {
    logEntries := make([]LogEntry, 0)
    if len(timePoint) == 0 {
      logger.connection.Where("terminal_id = ?", id).Find(&logEntries)
    } else {
      log.Println(timePoint)
      logger.connection.
        Where("terminal_id = ? and timestamp >= ?", id, timePoint).
        Find(&logEntries)
    }

    entries.Entries = make([]entryType, len(logEntries))

    for i, v := range logEntries {
      entries.Entries[i] = entryType{ v.Timestamp, v.Message }
    }

    sync <- true
  }()


  encoder := json.NewEncoder(res)
  <- sync
  encoder.Encode(entries)
  quitSignal <- true
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
  idSender := make(chan uint64)
  timeSender := make(chan string)
  quitSignal := make(chan bool)

  go handler.buildLogs(res, idSender, timeSender, quitSignal)

  vars := mux.Vars(req)
  id_str := vars["id"]
  id, err := strconv.ParseUint(id_str, 10, 64)
  if err != nil {
    log.Println(err)
    res.WriteHeader(http.StatusInternalServerError)
    quitSignal <- true
    return
  }

  requestValues := req.URL.Query()
  _, ok := requestValues["since"]

  idSender <- id

  if ok {
    timeSender <- requestValues["since"][0]
  } else {
    timeSender <- ""
  }

  <- quitSignal
}
