package src

import (
  "os"
  "time"
  "fmt"
  "strings"
  "path/filepath"
)

// HTTPかHTTPSの確認
func Checkprefix (url string) bool {
  return strings.HasPrefix(
    url, "http://") || strings.HasPrefix(url, "https://",
  )
}

// ページは既に存在するの？
func Checkexist (url string, prefix string) []string {
  res, err := filepath.Glob(prefix + "/archive/*" + url2path(url))
  if err != nil {
    fmt.Println("Err:", err)
  }
  return res
}

// http:/かhttps:/はいらない。最後の「/」は必要
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

// 必要なフォルダの創作
func Mkdirs (url string, prefix string) string {
  rep := url2path(url)
  t := time.Now().Unix()

  path := fmt.Sprint(prefix, "/archive/", t, rep)
  err := os.MkdirAll(path, 0755)
  if err != nil {
    fmt.Println("失敗：", err)
  }

  return path
}
