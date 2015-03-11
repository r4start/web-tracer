package tracer

import (
  "fmt"
  "net/http"

  "code.google.com/p/go-sqlite/go1/sqlite3"

  "github.com/r4start/web-tracer/sitecache"
)

type DbLogger struct {
  connection *sqlite3.Conn
  cache *sitecache.SiteCache
}

func NewDbLogger(dbName string) (DbLogger, error) {
  handler := DbLogger{nil, nil}

  conn, err := sqlite3.Open(dbName)
  if err == nil {
    handler.connection = conn
  }
  
  return handler, err
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
