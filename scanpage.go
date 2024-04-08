package main

import (
  "os"
  "fmt"
  "strings"
  "net/http"
  "net/url"
  "io"
  "regexp"
  "errors"
  "path/filepath"
)

func scanpage (path string, domain string, thisdomain string) error {
  // 先に保存したページを読み込む
  fn, err := os.ReadFile(path + "/index.html")
  if err != nil { return err }

  // 要らないタグを削除
  var script = regexp.MustCompile(
    `(<script.*</script>)`).ReplaceAllString(string(fn), "",
  )
  var noscript = regexp.MustCompile(
    `(<noscript.*</noscript>)`).ReplaceAllString(string(script), "",
  )
  var audio = regexp.MustCompile(
    `(<audio.*</audio>)`).ReplaceAllString(string(noscript), "",
  )
  var video = regexp.MustCompile(
    `(<video.*</video>)`).ReplaceAllString(string(audio), "",
  )
  var iframe = regexp.MustCompile(
    `(<iframe.*</iframe>)`).ReplaceAllString(string(video), "",
  )
  // 追加ダウンロード＋ローカル化
  var ass = regexp.MustCompile(
    // ルールに違反けど、長いからしょうがない・・・
    `(<img.*src=['"]|<meta.*content=['"]|<link.*href=['"])(.*\.)(png|webp|jpg|jpeg|gif|css|js|ico|svg|ttf|woff2)(\?[^'"]*)?`,
  )

  // 必要であれば、ページ内のURLを修正
  spath := "static/"
  if !strings.HasSuffix(path, "/") { spath = "/" + spath }
  spath = path + spath

  // また、追加ダウンロードのファイルに上記のフォルダを創作
  err = os.Mkdir(spath, 0755)
  if err != nil { return err }

  repmap := make(map[string]string)

  for _, cssx := range ass.FindAllString(iframe, -1) {
    // ページ内のURLを受け取る
    s := regexp.MustCompile(
      `(.*src=['"]|.*content=['"]|.*href=['"])`).Split(cssx, -1,
    )
    ss := regexp.MustCompile(`(['"].*)`).Split(s[1], -1)

    ogurl := ss[0] // 変わる前に元のURLを保存して
    // URLは//で始まるは愛
    if strings.HasPrefix(ss[0], "//") {
      ss[0] = "https:" + ss[0]
    }

    // ファイル名を見つけて
    fss := strings.Split(ss[0], "/")
    assdom := ""
    filename := fss[len(fss)-1]

    // httpかhttpsで始まる場合
    if strings.HasPrefix(ss[0], "http://") || 
       strings.HasPrefix(ss[0], "https://") {
      assdom = fss[2]
    }

    // フォルダの創作
    asspath := path + "/static/" + assdom
    err = os.MkdirAll(asspath, 0755)
    // 出来なければ、死ね
    if err != nil { return err }

    // ファイル名がなければ、次に値にスキップしてね
    if filename == "" { continue }

    // httpかhttpsで始まったら、ダウンロードだけしよう
    if strings.HasPrefix(ss[0], "http://") ||
       strings.HasPrefix(ss[0], "https://") {
      err = dlres(ss[0], filepath.Join(asspath, filename))
      if err != nil { return err }
    } else {
      // ローカルファイルなら、ちょっと変更は必要となるかしら
      u, err := url.Parse(domain)
      if err != nil { return err }

      rel, err := url.Parse(ss[0])
      if err != nil { return err }

      af := u.ResolveReference(rel).String()
      err = dlres(af, filepath.Join(asspath, filename))
      if err != nil { return err }
    }

    repmap[ogurl] = filepath.Join("/static", assdom, filename)
    if assdom == "" {
      repmap[ogurl] = filepath.Join("/static", filename)
    }

    if err != nil {
      fmt.Println(err)
      return errors.New("ダウンロードに失敗：")
    }
  }

  // URLをローカル化
  for ourl, lurl := range repmap {
    aurl := strings.ReplaceAll(path, thisdomain, "") + stripver(lurl)
    iframe = strings.ReplaceAll(iframe, ourl, aurl)
  }

  // index.htmlファイルを更新する
  err = os.WriteFile(path + "/index.html", []byte(iframe), 0644)
  if err != nil {
    fmt.Println(err)
    return errors.New("書込に失敗")
  }

  // エラーが出なかったから、返すのは不要
  return nil
}

// 画像、JS、CSS等ファイルのURLでパラメートルがある場合
func stripver (durl string) string {
  u, err := url.Parse(durl)
  if err != nil {
    fmt.Println("エラー：", err)
    return ""
  }

  u.RawQuery = ""
  return u.Path
}

func dlres (durl string, dest string) error {
  // ダウンロード
  res, err := http.Get(durl)
  if err != nil { return err }
  defer res.Body.Close()

  // URLでパラメートルがあれば、消す
  dest = stripver(dest)

  // MIMEタイプを確認
  ct := res.Header.Get("Content-Type")
  for mime, ext := range getmime() {
    if strings.Contains(ct, mime) && !strings.HasSuffix(dest, ext) {
      dest += ext
      break
    }
  }

  // ファイルを作成
  f, err := os.Create(dest)
  if err != nil { return err }
  defer f.Close()

  // ファイルを書き込む
  _, err = io.Copy(f, res.Body)
  if err != nil { return err }

  return nil
}
