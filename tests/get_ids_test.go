package tracercheck

import (
  "io"
  "testing"
  "net/http"
  "encoding/json"
  "net/http/httptest"

  "github.com/r4start/web-tracer/tracer"
)

type idsTest struct{}

type idsType struct {
  Ids []uint64 `json:"ids"`
}

func (r idsTest) NewGetIdsRequest(t *testing.T) *http.Request {
  url := "/ids/"

  req, e := http.NewRequest("GET", url, nil)
  if  e != nil {
    t.Fatal(e)
  }

  return req
}

func (r idsTest) SendRequestToServer(srv *tracer.IdLister,
                                     req *http.Request,
                                     expectedCode int,
                                     t *testing.T) *httptest.ResponseRecorder {
  recorder := httptest.NewRecorder()
  srv.ServeHTTP(recorder, req)

  if recorder.Code != expectedCode {
    t.Fatal("Server returned code ", recorder.Code)
  }

  return recorder
}

func (r idsTest) DecodeResponse(data io.Reader, t *testing.T) *idsType {
  jsonDecoder := json.NewDecoder(data)

  var ids idsType
  err := jsonDecoder.Decode(&ids)
  if err != nil && err != io.EOF {
    t.Fatal(err)
  }

  return &ids
}

func TestEmptyIds(t *testing.T) {
  p := idsTest{}

  cache := tracer.NewTerminalIdsCache()
  idLister := tracer.NewIdLister(&cache)
  req := p.NewGetIdsRequest(t)

  recorder := p.SendRequestToServer(&idLister, req, http.StatusOK, t)
  response := p.DecodeResponse(recorder.Body, t)

  if len(response.Ids) != 0 {
    t.Fatal("Ids must be empty! Got response ", response)
  }
}

func TestNotEmptyIds(t *testing.T) {
  p := idsTest{}
  testIds := []uint64{10, 43, 24234, 49856, 11, 17}

  cache := tracer.NewTerminalIdsCache()
  cache.AppendIds(testIds)

  idLister := tracer.NewIdLister(&cache)
  req := p.NewGetIdsRequest(t)

  recorder := p.SendRequestToServer(&idLister, req, http.StatusOK, t)
  response := p.DecodeResponse(recorder.Body, t)

  if len(response.Ids) != len(testIds) {
    t.Fatal("Ids count mismatched! Original ", len(testIds),
            "Got", len(response.Ids))
  }

  for _, v := range response.Ids {
    matched := false
    for _, k := range testIds {
      if v == k {
        matched = true
        break
      }
    }

    if !matched {
      t.Fatal("Can`t match value ", v)
    }
  }
}
