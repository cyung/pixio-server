package main

import (
  "io/ioutil"
  "log"
  "encoding/json"
)

type Configuration struct {
  Key string `json:"key"`
  Url string `json:"url"`
}

var _key string
var _url string

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
  _url = config.Url
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