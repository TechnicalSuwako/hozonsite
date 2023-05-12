package main

import (
  "os"
  "fmt"
  "strings"
  "net/http"
  "io"
  "regexp"
  "errors"
)

func scanpage (path string, domain string, thisdomain string) error {
  fn, err := os.ReadFile(path + "/index.html")
  if err != nil {
    fmt.Println(err)
    return errors.New("ファイルを開けられなかった：")
  }

  /* 削除 */
  var script = regexp.MustCompile(`(<script.*</script>)`).ReplaceAllString(string(fn), "")
  var noscript = regexp.MustCompile(`(<noscript.*</noscript>)`).ReplaceAllString(string(script), "")
  var audio = regexp.MustCompile(`(<audio.*</audio>)`).ReplaceAllString(string(noscript), "")
  var video = regexp.MustCompile(`(<video.*</video>)`).ReplaceAllString(string(audio), "")
  var iframe = regexp.MustCompile(`(<iframe.*</iframe>)`).ReplaceAllString(string(video), "")
  /* 追加ダウンロード＋ローカル化 */
  var ass = regexp.MustCompile(`(<img.*src="|<meta.*content="|<link.*href=")(.*\.)(png|webm|jpg|jpeg|gif|css|js)`)
  spath := "static/"
  if !strings.HasSuffix(path, "/") {
    spath = "/" + spath
  }
  spath = path + spath
  err1 := os.Mkdir(spath, 0755)
  if err1 != nil {
    fmt.Println(err1)
    return errors.New("失敗：")
  }

  for _, cssx := range ass.FindAllString(iframe, -1) {
    s := regexp.MustCompile(`(.*src="|.*content="|.*href=")`).Split(cssx, -1)
    ss := regexp.MustCompile(`(".*)`).Split(s[1], -1)
    if strings.HasPrefix(ss[0], "http://") || strings.HasPrefix(ss[0], "https://") {
      // TODO
    } else {
      fss := strings.Split(ss[0], "/")
      filename := fss[len(fss)-1]

      if filename == "" {
        continue
      }
      f, err := os.Create(spath + filename)
      if err != nil {
        fmt.Println(err)
        return errors.New("2. 作成失敗：")
      }
      defer f.Close()

      af := domain + ss[0]
      if !strings.HasPrefix(ss[0], "/") {
        af = domain + "/" + ss[0]
      }
      i, err := http.Get(af)
      if err != nil {
        fmt.Println(err)
        return errors.New("2. ダウンロードに失敗：")
      }
      defer i.Body.Close()
      if strings.HasSuffix(filename, "css") || strings.HasSuffix(filename, "js") {
        body, err := io.ReadAll(i.Body)
        if err != nil {
          fmt.Println(err)
          return errors.New("2. 読込エラー：")
        }

        _, err2 := f.WriteString(string(body))
        if err2 != nil {
          fmt.Println(err)
          return errors.New("2. ファイル書込エラー：")
        }
      } else {
        _, err = io.Copy(f, i.Body)
        if err != nil {
          fmt.Println(err)
          return errors.New("2. コピーに失敗：")
        }
      }

      iframe = strings.Replace(iframe, ss[0], "/static/" + filename, -1)
    }

    err := os.WriteFile(path + "/index.html", []byte(iframe), 0644)
    if err != nil {
      fmt.Println(err)
      return errors.New("書込に失敗")
    }
  }

  return nil
}
