package tracerdb

import (
  "fmt"
  "net/http"

  "code.google.com/p/go-sqlite/go1/sqlite3"
)

type DbLogHandler struct {
  connection *sqlite3.Conn
}

func CreateDbLogHandler(db_name string) (DbLogHandler, error) {
  handler := DbLogHandler{nil}

  conn, err := sqlite3.Open(db_name)
  if err == nil {
    handler.connection = conn
  }
  
  return handler, err
}

func (handler DbLogHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Log record was added with new handler %p!", &handler)
}

func (handler DbLogHandler) AddNewEntry(res http.ResponseWriter, req *http.Request) {
  
}
