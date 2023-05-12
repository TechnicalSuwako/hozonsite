package main

import (
  "os"
  "time"
  "fmt"
  "strings"
  "path/filepath"
  "net/http"
  "io"
  "regexp"
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

func scanpage (path string, domain string, thisdomain string) {
  fn, err := os.ReadFile(path + "/index.html")
  if err != nil {
    fmt.Println("ファイルを開けられなかった：", err)
    return
  }

  /* 削除 */
  var script = regexp.MustCompile(`(<script.*</script>)`).ReplaceAllString(string(fn), "")
  var noscript = regexp.MustCompile(`(<noscript.*</noscript>)`).ReplaceAllString(string(script), "")
  var audio = regexp.MustCompile(`(<audio.*</audio>)`).ReplaceAllString(string(noscript), "")
  var video = regexp.MustCompile(`(<video.*</video>)`).ReplaceAllString(string(audio), "")
  var iframe = regexp.MustCompile(`(<iframe.*</iframe>)`).ReplaceAllString(string(video), "")
  /* 追加ダウンロード＋ローカル化 */
  var ass = regexp.MustCompile(`(<img.*src="|<meta.*content="|<link.*href=")(.*\.)(png|webm|jpg|jpeg|gif|css)`)

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
      f, err := os.Create(path + "/" + filename)
      if err != nil {
        fmt.Println("2. 作成失敗：", err)
        return
      }
      defer f.Close()

      af := domain + ss[0]
      if !strings.HasPrefix(ss[0], "/") {
        af = domain + "/" + ss[0]
      }
      i, err := http.Get(af)
      if err != nil {
        fmt.Println("2. ダウンロードに失敗：", err)
        return
      }
      defer i.Body.Close()
      _, err = io.Copy(f, i.Body)
      if err != nil {
        fmt.Println("2. コピーに失敗：", err)
        return
      }

      iframe = strings.Replace(iframe, ss[0], "/" + filename, -1)
    }

    err := os.WriteFile(path + "/index.html", []byte(iframe), 0644)
    if err != nil {
      fmt.Println("書込に失敗")
      return
    }
  }
}
