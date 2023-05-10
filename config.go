package main

import (
  "io/ioutil"
  "fmt"
  "runtime"
  "encoding/json"
)

type Config struct {
  configpath string
  webpath string
  datapath string
  domain string
}

func getconf () Config {
  var payload map[string]interface{}
  var cnf Config

  prefix := "/usr"
  if runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" {
    prefix += "/local"
  }

  cnf.configpath = "/etc/hozonsite.json"
  //_, err = os.Stat(cnf.configpath)
  cnf.datapath = prefix + "/share/hozonsite"

  if runtime.GOOS == "freebsd" {
    cnf.configpath = prefix + cnf.configpath
  }

  data, err := ioutil.ReadFile(cnf.configpath)
  if err != nil {
    fmt.Println("エラー：", err)
  }
  json.Unmarshal(data, &payload)
  cnf.webpath = payload["webpath"].(string)
  cnf.domain = payload["domain"].(string)
  payload = nil

  return cnf
}
