package main

import (
  _ "os"
  "fmt"
  "log"
  "flag"
  "net/http"

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

  http.HandleFunc("/", hello)

  writeHandler, err := tracerdb.CreateDbLogHandler(db_name)
  if err != nil {
    log.Fatal(err)
  } else {
    http.Handle("/log", writeHandler)
  }

  bind := fmt.Sprintf("%s:%s", host, port)
  
  fmt.Printf("Listening on %s. Use database %s", bind, db_name)
  
  err = http.ListenAndServe(bind, nil)
  
  if err != nil {
    panic(err)
  }
}

func hello(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Database status is")
}
