package main

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
)

func main() {
  router := NewRouter()
  port := 3001
  addr := ":" + strconv.Itoa(port)

  fmt.Printf("Listening on port %d\n", port)
  log.Fatal(http.ListenAndServe(addr, router))
}