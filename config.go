package main

import (
  "os"
  "fmt"
  "runtime"
  "encoding/json"
  "io/ioutil"
  "errors"
)

type Config struct {
  configpath, webpath, datapath, domain, ip string
}

var cnf Config

func getconf () (Config, error) {
  // バイナリ、データ、及びFreeBSDとNetBSDの場合、コンフィグ
  prefix := "/usr"
  // BSDだけはただの/usrではない
  if runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" {
    prefix += "/local"
  } else if runtime.GOOS == "netbsd" {
    prefix += "/pkg"
  }

  // コンフィグファイル
  cnf.configpath = "/etc/hozonsite/config.json"
  cnf.datapath = prefix + "/share/hozonsite"

  // また、FreeBSDとNetBSDだけは違う場所だ。OpenBSDは正しい場所
  // FreeBSD = /usr/local/etc/hozonsite/config.json
  // NetBSD  = /usr/pkg/etc/hozonsite/config.json
  if runtime.GOOS == "freebsd" || runtime.GOOS == "netbsd" {
    cnf.configpath = prefix + cnf.configpath
  }

  // コンフィグファイルがなければ、死ね
  data, err := ioutil.ReadFile(cnf.configpath)
  if err != nil {
    fmt.Println("confif.jsonを開けられません：", err)
    return cnf, errors.New(
      "コンフィグファイルは " + cnf.configpath + " に創作して下さい。",
    )
  }

  var payload map[string]interface{}
  json.Unmarshal(data, &payload)
  if payload["webpath"] == nil {
    return cnf, errors.New("「webpath」の値が設置していません。")
  }
  if payload["domain"] == nil {
    return cnf, errors.New("「domain」の値が設置していません。")
  }
  if payload["ip"] == nil {
    return cnf, errors.New("「ip」の値が設置していません。")
  }
  if _, err := os.Stat(payload["webpath"].(string)); err != nil {
    fmt.Printf("%v\n", err)
    return cnf, errors.New(
      "mkdiorコマンドをつかって、 " + payload["webpath"].(string),
    )
  }
  cnf.webpath = payload["webpath"].(string) // データパス
  cnf.domain = payload["domain"].(string) // ドメイン名
  cnf.ip = payload["ip"].(string) // IP
  payload = nil // もういらなくなった

  return cnf, nil
}
