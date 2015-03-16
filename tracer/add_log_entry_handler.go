package tracer

import (
  "fmt"
  "net/http"

  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

type DbLogger struct {
  connection *gorm.DB
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
  fmt.Fprintf(res, "POST")
}
