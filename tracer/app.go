package tracer

import (
  "fmt"
  "net/http"

  "code.google.com/p/go-sqlite/go1/sqlite3"

  "github.com/r4start/web-tracer/sitecache"
)

const (
  mainPageStub = "<!DOCTYPE html><html><head><title>Under" +
                 " construction</title></head><body>" +
                 "Dear guest, the site is under construction.<br>Please visit" +
                 " us later.</body></html>"
)

type App struct {
  connection *sqlite3.Conn
  cache *sitecache.SiteCache
}

func NewApp(dbName string) (App, error) {
  newApp := App{nil, nil}

  conn, err := sqlite3.Open(dbName)
  if err == nil {
    newApp.connection = conn
  }
  
  return newApp, err
}

func (handler App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if req.Method == "POST" {
    fmt.Fprintf(res, "503")
    res.WriteHeader(503)
  } else if req.Method == "GET" {
    fmt.Fprintf(res, mainPageStub)
  } else {
    fmt.Fprintf(res, "503")
    res.WriteHeader(503)
  }
}
