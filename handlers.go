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
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/patrickmn/go-cache"
)

type Image struct {
  Url string `json:"url"`
}

var c *cache.Cache

func init() {
  c = cache.New(5*time.Minute, 30*time.Second)

  go func() {
    for {
      err := cleanOld()
      if err != nil {
        fmt.Println(err)
      }
      time.Sleep(5 * time.Second)
    }
  }()
}

// remove images older than a set time
func cleanOld() error {
  d, err := os.Open("./images")
  if err != nil {
    fmt.Println(err)
  }

  filenames, err := d.Readdirnames(-1)

  for _, filename := range filenames {
    old, err := isOld(filename)
    if err != nil {
      return err
    }

    if old {
      err := removeFile(filename)
      if err != nil {
        return err
      }
    }
  }

  return nil
}

// checks if a file is older than a set time
func isOld(filename string) (bool, error) {
  const MAX_TIME = 24.0 * 5 // max 5 days (in hours)
  f, err := os.Open("./images/" + filename)
  if err != nil {
    return false, err
  }
  defer f.Close()

  info, err := f.Stat()
  if err != nil {
    return false, err
  }

  age := time.Since(info.ModTime())
  return age.Hours() > MAX_TIME, nil
}

func removeFile(filename string) error {
  err := os.Remove("./images/" + filename)
  if err != nil {
    return err
  }

  fmt.Printf("removed %s\n", filename)
  return nil
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
  image_filename := RandomFilename() + extension

  out, err := os.Create("./images/" + image_filename)
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

  image := &Image{"http://ng-chris.com/img/" + image_filename}
  err = json.NewEncoder(w).Encode(image)
  if err != nil {
    panic(err)
  }
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