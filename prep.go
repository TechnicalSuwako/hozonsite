package main

import (
  "os"
  "time"
  "fmt"
  "strings"
  "path/filepath"
)

func checkprefix (url string) bool {
  return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func checkexist (url string, prefix string) []string {
  res, err := filepath.Glob(prefix + "/archive/*" + url2path(url))
  if err != nil {
    fmt.Println("Err:", err)
  }
  return res
}

func url2path (url string) string {
  res := ""
  if strings.HasPrefix(url, "https:/") {
    res = strings.Replace(url, "https:/", "", 1)
  } else {
    res = strings.Replace(url, "http:/", "", 1)
  }

  if strings.HasSuffix(res, "/") {
    res = strings.TrimSuffix(res, "/")
  }

  return res
}

func mkdirs (url string, prefix string) string {
  rep := url2path(url)
  t := time.Now().Unix()

  path := fmt.Sprint(prefix, "/archive/", t, rep)
  err := os.MkdirAll(path, 0755)
  if err != nil {
    fmt.Println("失敗：", err)
  }

  return path
}
