package src

import (
  "os"
  "fmt"
  "net/http"
  "io"
  "regexp"
  "strings"

  "golang.org/x/text/encoding/japanese"
  "golang.org/x/text/transform"
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
  // ソフトの終了する時に実行する
  defer curl.Body.Close()

  // ページの内容を読み込む
  body, err2 := io.ReadAll(curl.Body)
  if err2 != nil {
    fmt.Println("読込エラー：", err2)
    return
  }

  // Content-TypeヘッダーはUTF-8又は駄目のエンコーディングかの確認
  checkJis := `(?i)<meta.*?charset=(["']?)shift[_-]?jis`
  jisRegex, errr := regexp.Compile(checkJis)
  if errr != nil {
    fmt.Println(errr)
    return
  }

  checkEuc := `(?i)<meta.*?charset=(["']?)euc[_-]?jp`
  eucRegex, erre := regexp.Compile(checkEuc)
  if erre != nil {
    fmt.Println(erre)
    return
  }

  // 文字エンコーディングを変換する
  if jisRegex.Match(body) {
    shiftJISDecoder := japanese.ShiftJIS.NewDecoder()
    utf8Reader := transform.NewReader(
      strings.NewReader(string(body)),
      shiftJISDecoder,
    )
    utf8Body, err3 := io.ReadAll(utf8Reader)
    if err3 != nil {
      fmt.Println("文字エンコーディング変換エラー：", err3)
      return
    }

    body = utf8Body
  } else if eucRegex.Match(body) {
    eucJPDecoder := japanese.EUCJP.NewDecoder()
    utf8Reader := transform.NewReader(
      strings.NewReader(string(body)),
      eucJPDecoder,
    )
    utf8Body, err3 := io.ReadAll(utf8Reader)
    if err3 != nil {
      fmt.Println("文字エンコーディング変換エラー：", err3)
      return
    }

    body = utf8Body
  }

  // 空index.htmlファイルを創作する
  fn, err4 := os.Create(path + "/index.html")
  if err4 != nil {
    fmt.Println("ファイルの創作エラー：", err4)
    return
  }
  // ソフトの終了する時に実行する
  defer fn.Close()

  // あのindex.htmlファイルに内容をそのまま書き込む
  _, err5 := fn.WriteString(string(body))
  if err5 != nil {
    fmt.Println("ファイル書込エラー：", err5)
  }
}
