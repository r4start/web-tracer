package main

import (
  "fmt"
  "net/http"
  "os"
  "database/sql"
  "github.com/lib/pq"
)

var is_success bool = false

func main() {
  db_host := os.Getenv("OPENSHIFT_POSTGRESQL_DB_HOST")
  db_port := os.Getenv("OPENSHIFT_POSTGRESQL_DB_PORT")
  db_user := "admin4egsj8v"
  db_pwd := "yAdiWMWAicBu"
  db_name := "tracer"

  db_conn_line := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
                              db_host, db_port, db_user, db_pwd, db_name)

  _, err := sql.Open("postgres", db_conn_line)
  if err == nil {
    is_success = true
  }

  http.HandleFunc("/", hello)
  bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
  fmt.Printf("listening on %s...", bind)
  err = http.ListenAndServe(bind, nil)
  if err != nil {
    panic(err)
  }
}

func hello(res http.ResponseWriter, req *http.Request) {
  var msg string
  if is_success {
    msg = "opened"
  } else {
    msg = "not found"
  }

  fmt.Fprintf(res, "Database status is %s", msg)
}
