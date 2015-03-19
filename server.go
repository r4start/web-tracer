package main

import (
  "os"
  "fmt"
  "log"
  "flag"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/r4start/web-tracer/tracer"
)

type ServerParameters struct {
  Host string
  Port string
  DbName string
  SiteRoot string
}

func getServeAddress() ServerParameters {
  var params ServerParameters

  flag.StringVar(&params.Host, "host", "localhost", "IP address for listening")
  flag.StringVar(&params.Port, "port", "80", "Port number")
  flag.StringVar(&params.DbName, "dbname", "tracer.db", "Database name or path")
  flag.StringVar(&params.SiteRoot, "site-root", "www/", "Path to site root folder")

  flag.Parse()

  return params
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

func main() {
  params := getServeAddress()

  if isSiteRootExists(params.SiteRoot) {
    router := mux.NewRouter()
    router.NotFoundHandler = http.HandlerFunc(notFoundPage)

    var idsCache = tracer.NewTerminalIdsCache()

    {
      writeHandler, err := tracer.NewDbLogger(params.DbName)
      if err != nil {
        log.Fatal(err)
      } else {
        writeHandler.IdsCache = &idsCache
        router.Handle("/terminal/{id:[0-9]+}", writeHandler)
      }

      writeHandler.PrepareDB()
    }

    {
      idsLister, err := tracer.NewIdLister(params.DbName)
      if err != nil {
        log.Fatal(err)
      } else {
        idsLister.IdsCache = &idsCache
        router.Handle("/ids", idsLister)
      }
    }

    idsCache.AppendIds(tracer.LoadIdsFromDb(params.DbName))

    router.PathPrefix("/").Handler(http.FileServer(http.Dir(params.SiteRoot)))

    http.Handle("/", router)
  } else {
    log.Fatal("Please specify site root.")
  }

  bind := fmt.Sprintf("%s:%s", params.Host, params.Port)
  
  fmt.Printf("Listening on %s. Use database %s\n", bind, params.DbName)
  
  err := http.ListenAndServe(bind, nil)
  
  if err != nil {
    panic(err)
  }
}
