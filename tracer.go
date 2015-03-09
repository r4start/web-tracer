package main

import (
  _ "os"
  "fmt"
  "flag"
  "net/http"
)

func getServeAddress() (string, string) {
  var host string
  var port string

  flag.StringVar(&host, "host", "localhost", "IP address for listening")
  flag.StringVar(&port, "port", "4000", "Port number")

  flag.Parse()

  return host, port
}

func main() {
  host, port := getServeAddress()

  http.HandleFunc("/", hello)
  bind := fmt.Sprintf("%s:%s", host, port)
  
  fmt.Printf("listening on %s...", bind)
  
  err := http.ListenAndServe(bind, nil)
  
  if err != nil {
    panic(err)
  }
}

func hello(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "Database status is")
}
