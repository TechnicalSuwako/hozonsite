package main

import (
  "text/template"
  "fmt"
  "net/http"
  "strings"
  "strconv"
  "time"
  "os"
  "encoding/json"
  "gitler.moe/suwako/goliblocale"
)

type (
  Page struct {
    Err, Lan, Ver, Ves, Url, Body string
    i18n map[string]string
    Ext []Exist // 既に存在する場合
  }
  Stat struct { // APIのみ
    Url, Ver string
  }
  Exist struct {
    Date, Url string
  }
)

var ftmpl []string
var data *Page

func (p Page) T (key string) string {
  return p.i18n[key]
}

// 言語設定、デフォルト＝ja
func initloc (r *http.Request) string {
  supportedLanguages := map[string]bool{
    "ja": true,
    "en": true,
  }

  cookie, err := r.Cookie("lang")
  if err != nil {
    return "ja"
  }

  if _, ok := supportedLanguages[cookie.Value]; ok {
    return cookie.Value
  } else {
    return "ja"
  }
}

func tspath (p string) string {
  pc := strings.Split(p, "/")

  for i := len(pc) - 1; i >= 0; i-- {
    if _, err := strconv.Atoi(pc[i]); err == nil {
      return pc[i]
    }
  }

  return ""
}

func handleStatic(path string, cnf Config, w http.ResponseWriter, r *http.Request) {
    if !strings.HasSuffix(path, ".css") &&
       !strings.HasSuffix(path, ".png") &&
       !strings.HasSuffix(path, ".jpeg") &&
       !strings.HasSuffix(path, ".jpg") &&
       !strings.HasSuffix(path, ".webm") &&
       !strings.HasSuffix(path, ".gif") &&
       !strings.HasSuffix(path, ".js") {
    http.NotFound(w, r)
    return
  }

  fpath := cnf.datapath + "/archive/" + path
  http.ServeFile(w, r, fpath)
}

func handlePost(w http.ResponseWriter, r *http.Request, cnf Config) {
  err := r.ParseForm()
  if err != nil {
    fmt.Println(err)
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }

  // 言語変更
  if lang := r.PostFormValue("lang"); lang != "" {
    http.SetCookie(
      w,
      &http.Cookie{Name: "lang", Value: lang, MaxAge: 31536000, Path: "/"},
    )
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
  }

  var exist []string
  langu := initloc(r)
  i18n, err := goliblocale.GetLocale(cnf.webpath + "/locale/" + langu)
  if err != nil {
    fmt.Printf("liblocaleエラー：%v", err)
    return
  }

  if r.PostForm.Get("hozonsite") == "" {
    data.Err = i18n["errfusei"]
    ftmpl[0] = cnf.webpath + "/view/404.html"
    return
  }

  url := r.PostForm.Get("hozonsite")
  // HTTPかHTTPSじゃない場合
  if !checkprefix(url) {
    data.Err = i18n["errfuseiurl"]
    ftmpl[0] = cnf.webpath + "/view/404.html"
    return
  }

  eurl := stripurl(url)
  exist = checkexist(eurl, cnf.datapath)
  if len(exist) == 0 || r.PostForm.Get("agree") == "1" {
    path := mkdirs(eurl, cnf.datapath)
    getpage(url, path)
    scanpage(path, eurl, cnf.datapath)
    http.Redirect(
      w,
      r,
      cnf.domain + strings.Replace(path, cnf.datapath, "", 1),
      http.StatusSeeOther,
    )
    return
  }

  ftmpl[0] = cnf.webpath + "/view/check.html"
  data.Url = url
  var existing []Exist
  e := Exist{}
  for _, ex := range exist {
    ti, err := strconv.ParseInt(tspath(ex), 10, 64)
    if err != nil {
      fmt.Println(err)
      http.Redirect(w, r, "/", http.StatusSeeOther)
      return
    }

    t := time.Unix(ti, 0)
    e.Date = t.Format("2006年01月02日 15:04:05")
    e.Url = strings.Replace(ex, cnf.datapath, cnf.domain, 1)
    existing = append(existing, e)
  }

  data.Ext = existing
}

// ホームページ
func siteHandler (cnf Config) func (http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {
    ftmpl = []string{
      cnf.webpath + "/view/index.html",
      cnf.webpath + "/view/header.html",
      cnf.webpath + "/view/footer.html",
    }
    data = &Page{Ver: version, Ves: strings.ReplaceAll(version, ".", "")}

    lang := initloc(r)
    data.Lan = lang

    i18n, err := goliblocale.GetLocale(cnf.webpath + "/locale/" + lang)
    if err != nil {
      fmt.Printf("liblocaleエラー：%v", err)
      return
    }
    data.i18n = i18n
    ftmpl[0] = cnf.webpath + "/view/index.html"
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))

    if r.Method == "POST" {
      handlePost(w, r, cnf)
    }

    tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
    tmpl.Execute(w, data)
  }
}

// /api TODO
func apiHandler (cnf Config) func (http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(200)
    buf, _ := json.MarshalIndent(&Stat{Url: cnf.domain, Ver: version}, "", "  ")
    _, _ = w.Write(buf)
  }
}

// /archive
func archiveHandler (cnf Config) func (http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {
    ftmpl := []string{
      cnf.webpath + "/view/index.html",
      cnf.webpath + "/view/header.html",
      cnf.webpath + "/view/footer.html",
    }
    data := &Page{Ver: version, Ves: strings.ReplaceAll(version, ".", "")}
    lang := initloc(r)
    data.Lan = lang

    i18n, err := goliblocale.GetLocale(cnf.webpath + "/locale/" + lang)
    if err != nil {
      fmt.Printf("liblocaleエラー：%v", err)
      return
    }
    data.i18n = i18n
    ftmpl[0] = cnf.webpath + "/view/index.html"
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
    path := strings.TrimPrefix(r.URL.Path, "/archive/")

    if strings.Contains(path, "/static/") {
      handleStatic(path, cnf, w, r)
      return
    }

    pth := r.URL.Path
    if !strings.HasSuffix(pth, "/") &&
       !strings.HasSuffix(pth, "index.html") {
      pth += "/index.html"
    } else if strings.HasSuffix(pth, "/") &&
              !strings.HasSuffix(pth, "index.html") {
      pth += "index.html"
    }

    file := cnf.datapath + pth
    if _, err := os.Stat(file); os.IsNotExist(err) {
      http.Redirect(w, r, "/404", http.StatusSeeOther)
      return
    }

    bdy, err := os.ReadFile(file)
    if err != nil {
      http.Redirect(w, r, "/404", http.StatusSeeOther)
      return
    }

    data.Body = string(bdy)
    tmpl = template.Must(
      template.ParseFiles(cnf.webpath + "/view/archive.html"),
    )
    tmpl.Execute(w, data)
    data = nil
  }
}

// サーバー
func serv (cnf Config, port int) {
  http.Handle(
    "/static/",
    http.StripPrefix("/static/",
    http.FileServer(http.Dir(cnf.webpath + "/static"))),
  )

  http.HandleFunc("/api/", apiHandler(cnf))
  http.HandleFunc("/archive/", archiveHandler(cnf))
  http.HandleFunc("/", siteHandler(cnf))

  fmt.Println(fmt.Sprint(
    "http://" + cnf.ip + ":",
    port,
    " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"),
  )
  http.ListenAndServe(fmt.Sprint(cnf.ip + ":", port), nil)
}
