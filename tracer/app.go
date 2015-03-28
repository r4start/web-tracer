package tracer

import (
  "os"
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
)

type ServerParameters struct {
  Host string
  Port string
  DbName string
  SiteRoot string
}

func notFoundPage(res http.ResponseWriter, req *http.Request) {
  res.Header().Add("Location", "http://" + req.Host + "/404.html")
  res.WriteHeader(404)
}

func isSiteRootExists(path string) bool {
  _, err := os.Stat(path)
  if err == nil { return true }
  return false
}

func StartServer(params ServerParameters) {
  router := mux.NewRouter()
  router.NotFoundHandler = http.HandlerFunc(notFoundPage)

  var idsCache = NewTerminalIdsCache()

  {
    writeHandler, err := NewDbLogger(params.DbName)
    if err != nil {
      log.Fatal(err)
    } else {
      writeHandler.IdsCache = &idsCache
      router.Handle("/terminal/{id:[0-9]+}", writeHandler)
    }

    writeHandler.PrepareDB()
  }

  {
    idsLister := NewIdLister(&idsCache)
    router.Handle("/ids", idsLister)  
  }

  idsCache.AppendIds(LoadIdsFromDb(params.DbName))

  if isSiteRootExists(params.SiteRoot) {
    router.PathPrefix("/").Handler(http.FileServer(http.Dir(params.SiteRoot)))
  }
 
  http.Handle("/", router)

  bind := fmt.Sprintf("%s:%s", params.Host, params.Port)
  
  fmt.Printf("Listening on %s. Use database %s\n", bind, params.DbName)
  
  err := http.ListenAndServe(bind, nil)
  
  if err != nil {
    log.Fatal(err)
  }
}
