package main

import (
  "net/http"
  "fmt"
  "bytes"
  "os"
  "io"
  "io/ioutil"
  "time"
  "math/rand"
  "strings"
  "github.com/gorilla/mux"
  "github.com/patrickmn/go-cache"
)

var c *cache.Cache

func init() {
  c = cache.New(5*time.Minute, 30*time.Second)
}

func AddImage(w http.ResponseWriter, r *http.Request) {
  if !Validate(r.Header.Get("Authorization")) {
    w.WriteHeader(401)
    return
  }

  file, header, err := r.FormFile("image")
  if err != nil {
    w.WriteHeader(500)
    fmt.Println(err)
    return
  }
  defer file.Close()

  filename := header.Filename
  extension := filename[strings.LastIndex(filename, "."):]

  out, err := os.Create("./images/" + RandomFilename() + extension)
  if err != nil {
    w.WriteHeader(500)
    fmt.Println(err)
    return
  }

  // verify file data
  _, err = io.Copy(out, file)
  if err != nil {
    w.WriteHeader(500)
    fmt.Println(err)
  }

  w.WriteHeader(201)
}

func RandomFilename() string {
  const CHAR_LENGTH = 10
  const chars = "abcdefghijklmnopqrstuvwxyz" +
                "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
                "0123456789"

  rand.Seed(time.Now().UTC().UnixNano())
  result := make([]byte, CHAR_LENGTH)

  for i := 0; i<CHAR_LENGTH; i++ {
    result[i] = chars[rand.Intn(len(chars))]
  }

  return string(result)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
  filename := mux.Vars(r)["filename"]
  filepath := "./images/" + filename
  fmt.Printf("filename = %s\n", filename)

  data, found := c.Get(filename)
  if found {
    http.ServeContent(w, r, filename, time.Time{}, bytes.NewReader(data.([]byte)))
    return
  }

  buf, err := ioutil.ReadFile(filepath)
  if err != nil {
    fmt.Println("file not found")
    fmt.Printf("err %v\n", err)
    w.WriteHeader(404)
    return
  }
  c.Set(filename, buf, cache.DefaultExpiration)
  http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(buf))
}