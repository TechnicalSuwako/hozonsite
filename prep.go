package main

import (
  "os"
  "time"
  "fmt"
  "strings"
  "path/filepath"
  "net/http"
  "io"
)

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

func getpage (url string, path string) {
  curl, err := http.Get(url)
  if err != nil {
    fmt.Println("CURLエラー：", err)
    return
  }

  defer curl.Body.Close()
  body, err2 := io.ReadAll(curl.Body)
  if err2 != nil {
    fmt.Println("読込エラ：", err2)
    return
  }

  fn, err3 := os.Create(path + "/index.html")
  if err3 != nil {
    fmt.Println("ファイルの創作エラー：", err3)
    return
  }

  defer fn.Close()
  _, err4 := fn.WriteString(string(body))
  if err4 != nil {
    fmt.Println("ファイル書込エラー：", err4)
  }
}

//func scanpage (path string) {}
