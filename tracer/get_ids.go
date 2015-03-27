package tracer

import (
  "net/http"
  "encoding/json"
)

type IdLister struct {
  idsCache *TerminalIdsCache
}

func NewIdLister(cache *TerminalIdsCache) IdLister {
  return IdLister{cache}
}

func (handler IdLister) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    res.WriteHeader(http.StatusBadRequest)
    return
  } else if handler.idsCache == nil {
    res.WriteHeader(http.StatusInternalServerError)
    return
  }

  type idsType struct {
    Ids []uint64 `json:"ids"`
  }

  var ids idsType
  signal := make(chan bool)

  go func() {
    ids.Ids = handler.idsCache.GetIds()
    signal <- true
  }()

  encoder := json.NewEncoder(res)

  <- signal
  encoder.Encode(ids)
}
