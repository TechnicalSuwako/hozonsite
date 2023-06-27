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
)

type (
  Page struct {
    Tit, Err, Lan, Ver, Ves, Url, Body string
    Ext []Exist // 既に存在する場合
  }
  Stat struct { // APIのみ
    Url, Ver string
  }
  Exist struct {
    Date, Url string
  }
)

// 日本語か英語 TODO：複数言語対応
func initloc (r *http.Request) string {
  cookie, err := r.Cookie("lang")
  if err == nil && cookie.Value == "en" {
    return "en"
  }
  return "ja"
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

// ホームページ
func siteHandler (cnf Config) func (http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {
    ftmpl := []string{cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"}
    data := &Page{Ver: version, Ves: strings.ReplaceAll(version, ".", "")}

    lang := initloc(r)

    data.Lan = lang
    ftmpl[0] = cnf.webpath + "/view/index.html"
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))

    data.Tit = getloc("top", lang)
    if r.Method == "POST" {
      err := r.ParseForm()
      if err != nil {
        fmt.Println(err)
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
      }

      // クッキー
      if r.PostForm.Get("langchange") != "" {
        cookie, err := r.Cookie("lang")
        if err != nil || cookie.Value == "ja" {
          http.SetCookie(w, &http.Cookie {Name: "lang", Value: "en", MaxAge: 31536000, Path: "/"})
        } else {
          http.SetCookie(w, &http.Cookie {Name: "lang", Value: "ja", MaxAge: 31536000, Path: "/"})
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
      }

      var exist []string

      if r.PostForm.Get("hozonsite") != "" {
        url := r.PostForm.Get("hozonsite")
        // HTTPかHTTPSじゃない場合
        if !checkprefix(url) {
          data.Err = getloc("errfuseiurl", lang)
          ftmpl[0] = cnf.webpath + "/view/404.html"
        } else {
          eurl := stripurl(url)
          exist = checkexist(eurl, cnf.datapath)
          if len(exist) == 0 || r.PostForm.Get("agree") == "1" {
            path := mkdirs(eurl, cnf.datapath)
            getpage(url, path)
            scanpage(path, eurl, cnf.datapath)
            http.Redirect(w, r, cnf.domain + strings.Replace(path, cnf.datapath, "", 1), http.StatusSeeOther)
          } else if len(exist) > 0 {
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
          } else {
            data.Err = getloc("errfusei", lang)
            ftmpl[0] = cnf.webpath + "/view/404.html"
          }
        }
      }
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
    ftmpl := []string{cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"}
    data := &Page{Ver: version, Ves: strings.ReplaceAll(version, ".", "")}
    lang := initloc(r)

    data.Lan = lang
    ftmpl[0] = cnf.webpath + "/view/index.html"
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
    path := strings.TrimPrefix(r.URL.Path, "/archive/")

    if strings.Contains(path, "/static/") {
      if !strings.HasSuffix(path, ".css") && !strings.HasSuffix(path, ".png") && !strings.HasSuffix(path, ".jpeg") && !strings.HasSuffix(path, ".jpg") && !strings.HasSuffix(path, ".webm") && !strings.HasSuffix(path, ".gif") && !strings.HasSuffix(path, ".js") {
        http.NotFound(w, r)
        return
      }

      fpath := cnf.datapath + "/archive/" + path
      http.ServeFile(w, r, fpath)
    } else {
      pth := r.URL.Path
      if !strings.HasSuffix(pth, "/") && !strings.HasSuffix(pth, "index.html") {
        pth += "/index.html"
      } else if strings.HasSuffix(pth, "/") && !strings.HasSuffix(pth, "index.html") {
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
      tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/archive.html"))
      tmpl.Execute(w, data)
      data = nil
    }
  }
}

// サーバー
func serv (cnf Config, port int) {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  http.HandleFunc("/api/", apiHandler(cnf))
  http.HandleFunc("/archive/", archiveHandler(cnf))
  http.HandleFunc("/", siteHandler(cnf))

  fmt.Println(fmt.Sprint("http://127.0.0.1:", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(":", port), nil)
}
