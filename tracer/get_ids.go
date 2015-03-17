package tracer

import (
  "net/http"
  "encoding/json"

  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

type IdLister struct {
  connection *gorm.DB

  IdsCache *TerminalIdsCache
}

func NewIdLister(dbName string) (IdLister, error) {
  conn, err := gorm.Open("sqlite3", dbName)
  return IdLister{&conn, nil}, err
}

func (handler IdLister) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    res.WriteHeader(404)
    return
  } else if handler.IdsCache == nil {
    res.WriteHeader(503)
    return
  }

  encoder := json.NewEncoder(res)

  type idsType struct {
    Ids []uint64 `json:"ids"`
  }

  var ids = idsType{handler.IdsCache.GetIds()}
  _ = encoder.Encode(ids)
}
