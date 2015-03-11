package main

import (
  "fmt"
  "log"
  "flag"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/r4start/web-tracer/tracerdb"
)

func getServeAddress() (string, string, string) {
  var host string
  var port string
  var db_name string

  flag.StringVar(&host, "host", "localhost", "IP address for listening")
  flag.StringVar(&port, "port", "4000", "Port number")
  flag.StringVar(&db_name, "dbname", "tracer.db", "Database name or path")

  flag.Parse()

  return host, port, db_name
}

func main() {
  host, port, db_name := getServeAddress()

  router := mux.NewRouter()
  router.HandleFunc("/", mainPage)
  router.NotFoundHandler = http.HandlerFunc(notFoundPage)

  writeHandler, err := tracerdb.CreateDbLogHandler(db_name)
  if err != nil {
    log.Fatal(err)
  } else {
    router.Handle("/terminal/{id:[0-9]+}", writeHandler)
  }

  http.Handle("/", router)

  bind := fmt.Sprintf("%s:%s", host, port)
  
  fmt.Printf("Listening on %s. Use database %s", bind, db_name)
  
  err = http.ListenAndServe(bind, nil)
  
  if err != nil {
    panic(err)
  }
}

func mainPage(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Web tracer main page.")
}

func notFoundPage(res http.ResponseWriter, req *http.Request) {
  res.Header().Add("Location", "http://" + req.Host)
  res.WriteHeader(302)
}
