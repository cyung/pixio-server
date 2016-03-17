package main

import (
  "io/ioutil"
  "log"
  "encoding/json"
)

type Configuration struct {
  ChrisKey string `json:"CHRIS_KEY"`
}

var key string

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

  key = config.ChrisKey
}

func GetKey() string {
  return key
}

func Validate(client_key string) bool {
  return client_key == key
}