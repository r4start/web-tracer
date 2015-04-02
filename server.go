package main

import (
  "flag"

  "github.com/r4start/web-tracer/tracer"
)

func getServeAddress() tracer.ServerParameters {
  var params tracer.ServerParameters

  flag.StringVar(&params.Host, "host", "localhost", "IP address for listening")
  flag.StringVar(&params.Port, "port", "80", "Port number")
  flag.StringVar(&params.DbName, "dbname", ":memory:", "Database name or path")
  flag.StringVar(&params.SiteRoot, "site-root", "", "Path to site root folder")

  flag.Parse()

  return params
}

func main() {
  params := getServeAddress()

  tracer.StartServer(params)
}
