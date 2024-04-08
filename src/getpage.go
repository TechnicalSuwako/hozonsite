package src

import (
  "os"
  "fmt"
  "net/http"
  "io"
  "strings"
)

// URLでパラメートル（?、=等）がある場合
func Stripurl (url string) string {
  res := strings.ReplaceAll(url, "?", "")
  res = strings.ReplaceAll(res, "=", "")
  return res
}

func Getpage (url string, path string) {
  // ページを読み込む
  curl, err := http.Get(url)
  if err != nil {
    fmt.Println("CURLエラー：", err)
    return
  }
  defer curl.Body.Close() // ソフトの終了する時に実行する

  // ページの内容を読み込む
  body, err2 := io.ReadAll(curl.Body)
  if err2 != nil {
    fmt.Println("読込エラー：", err2)
    return
  }

  // 空index.htmlファイルを創作する
  fn, err3 := os.Create(path + "/index.html")
  if err3 != nil {
    fmt.Println("ファイルの創作エラー：", err3)
    return
  }
  defer fn.Close() // ソフトの終了する時に実行する

  // あのindex.htmlファイルに内容をそのまま書き込む
  _, err4 := fn.WriteString(string(body))
  if err4 != nil {
    fmt.Println("ファイル書込エラー：", err4)
  }
}
