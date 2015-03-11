package main

import (
  "os"
  "fmt"
  "log"
  "flag"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/r4start/web-tracer/tracerdb"
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
  flag.StringVar(&params.Port, "port", "4000", "Port number")
  flag.StringVar(&params.DbName, "dbname", "tracer.db", "Database name or path")
  flag.StringVar(&params.SiteRoot, "site-root", "www/", "Path to site root folder")

  flag.Parse()

  return params
}

func mainPage(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Web tracer main page.")
}

func notFoundPage(res http.ResponseWriter, req *http.Request) {
  res.Header().Add("Location", "http://" + req.Host)
  res.WriteHeader(302)
}

func siteUnderReconstructionPage(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Stub")
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
    router.HandleFunc("/", mainPage)
    router.NotFoundHandler = http.HandlerFunc(notFoundPage)

    writeHandler, err := tracerdb.CreateDbLogHandler(params.DbName)
    if err != nil {
      log.Fatal(err)
    } else {
      router.Handle("/terminal/{id:[0-9]+}", writeHandler)
    }

    http.Handle("/", router)
  } else {
    http.HandleFunc("/", siteUnderReconstructionPage)
  }

  bind := fmt.Sprintf("%s:%s", params.Host, params.Port)
  
  fmt.Printf("Listening on %s. Use database %s", bind, params.DbName)
  
  err := http.ListenAndServe(bind, nil)
  
  if err != nil {
    panic(err)
  }
}
