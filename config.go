package main

import (
  "os"
  "fmt"
  "runtime"
  "encoding/json"
)

type Config struct {
  configpath, webpath, datapath, domain string
}

func getconf () Config {
  var payload map[string]interface{}
  var cnf Config

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
  data, err := os.ReadFile(cnf.configpath)
  if err != nil {
    fmt.Println("エラー：", err)
  }
  json.Unmarshal(data, &payload)
  cnf.webpath = payload["webpath"].(string) // データパス
  cnf.domain = payload["domain"].(string) // ドメイン名
  payload = nil // もういらなくなった

  return cnf
}
