package main

import (
  "os"
  "fmt"
  "log"
  "flag"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/r4start/web-tracer/tracer"
  "github.com/r4start/web-tracer/sitecache"
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

func notFoundPage(res http.ResponseWriter, req *http.Request) {
  res.Header().Add("Location", "http://" + req.Host)
  res.WriteHeader(302)
}

// if app can`t find a site root, then it will show this stub html.
func siteUnderReconstructionPage(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(res, "<!DOCTYPE html><html><head><title>Under" +
                   " construction</title></head><body>" +
                   "The site is under construction.<br>Please visit" +
                   " us later.</body></html>")
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

    {
      _, err := sitecache.NewSiteCache(params.SiteRoot)
      if err != nil {
        log.Fatal(err)
      }
    }

    {
      writeHandler, err := tracer.NewDbLogger(params.DbName)
      if err != nil {
        log.Fatal(err)
      } else {
        router.Handle("/terminal/{id:[0-9]+}", writeHandler)
      }
    }

    {
      app, err := tracer.NewApp(params.DbName)
      if err != nil {
        log.Fatal(err)
      } else {
        router.Handle("/", app)
      }
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
