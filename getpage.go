package main

import (
  "os"
  "fmt"
  "net/http"
  "io"
)

func getpage (url string, path string) {
  curl, err := http.Get(url)
  if err != nil {
    fmt.Println("CURLエラー：", err)
    return
  }

  defer curl.Body.Close()
  body, err2 := io.ReadAll(curl.Body)
  if err2 != nil {
    fmt.Println("読込エラー：", err2)
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
