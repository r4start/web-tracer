package tracer

import (
  "net/http"

  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

type IdLister struct {
  connection *gorm.DB
}

func NewIdLister(dbName string) (IdLister, error) {
  conn, err := gorm.Open("sqlite3", dbName)
  return IdLister{&conn}, err
}

func (handler IdLister) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    res.WriteHeader(404)
    return
  }
}
