package main

import (
  "fmt"
  "log"
  "net/http"
)

func main() {
  router := NewRouter()

  fmt.Printf("Listening on port %d\n", 3000)
  log.Fatal(http.ListenAndServe(":3000", router))
}