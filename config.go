package main

import (
  "io/ioutil"
  "log"
  "encoding/json"
)

type Configuration struct {
  Key string `json:"key"`
}

var _key string
const _url string = "http://localhost:3001"

func init() {
  file, err := ioutil.ReadFile("./config.json")
  if err != nil {
    log.Fatal(err)
  }

  var config Configuration
  err = json.Unmarshal(file, &config)
  if err != nil {
    log.Fatal(err)
  }

  _key = config.Key
}

func GetKey() string {
  return _key
}

func BaseUrl() string {
  return _url
}

func Validate(client_key string) bool {
  return client_key == _key
}