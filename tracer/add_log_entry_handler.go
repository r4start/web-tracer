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
}

func (logger DbLogger) storeEntry(termianlId uint64, msg string) {
  encodedMsg := base64.StdEncoding.EncodeToString([]byte(msg))

  entry := LogEntry{TerminalId : termianlId,
                    Timestamp : time.Now().String(),
                    Message : encodedMsg}

  go logger.connection.Create(&entry)
  log.Printf("Adding entry for %d, content: %s\n", termianlId, encodedMsg)
}

func NewDbLogger(dbName string) (DbLogger, error) {
  conn, err := gorm.Open("sqlite3", dbName)
  return DbLogger{&conn}, err
}

func (handler DbLogger) PrepareDB() {
  if !handler.connection.HasTable(&LogEntry{}) {
    handler.connection = handler.connection.CreateTable(&LogEntry{})
  }
}

func (handler DbLogger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    fmt.Fprintf(res, "404")
  } else {
    handler.AddNewEntry(res, req)
  }
}

func (handler DbLogger) AddNewEntry(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  type message struct {
    Msg string `json:"message"`
  }

  var msg message
  {
    err := decoder.Decode(&msg)
    if err != nil && err != io.EOF {
      log.Println("Unable to decode terminal message. ", err)
      return
    }
  }

  vars := mux.Vars(req)
  id_str := vars["id"]
  
  if len(msg.Msg) == 0 {
    log.Println("Empty message! For terminal ", id_str)
    return
  }

  id, err := strconv.ParseUint(id_str, 10, 64)
  if err != nil {
    log.Println(err)
    return
  }

  go handler.storeEntry(id, msg.Msg)

  res.WriteHeader(200)
}
