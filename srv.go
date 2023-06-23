package main

import (
  "text/template"
  "fmt"
  "net/http"
  "strings"
  "os"
  "encoding/json"
)

type (
  Page struct {
    Tit string
    Err string
    Lan string
    Ver string
    Ves string
    Ext []string // 既に存在する場合
    Url string // 確認ページ用
    Body string // 保存したページ用
  }
  Stat struct {
    Url string
    Ver string
  }
)

func initloc (r *http.Request) string {
  cookie, err := r.Cookie("lang")
  if err != nil {
    return "ja"
  } else {
    return cookie.Value
  }
}

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
      if err != nil { fmt.Println(err) }

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

      if r.PostForm.Get("hozonsite") != "" {
        fmt.Println("sasa")
        url := r.PostForm.Get("hozonsite")
        fmt.Println(url)
        // HTTPかHTTPSじゃない場合
        if !checkprefix(url) {
          data.Err = getloc("errfuseiurl", lang)
          ftmpl[0] = cnf.webpath + "/view/404.html"
        } else {
          //if r.PostForm.Get("sosin") != "" {}
        }
      }
    }

    tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
    tmpl.Execute(w, data)
  }
}

func apiHandler (cnf Config) func (http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(200)
    buf, _ := json.MarshalIndent(&Stat{Url: cnf.domain, Ver: version}, "", "  ")
    _, _ = w.Write(buf)
  }
}

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

      bdy, err := os.ReadFile(cnf.datapath + pth)
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

func serv (cnf Config, port int) {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  http.HandleFunc("/api/", apiHandler(cnf))
  http.HandleFunc("/archive/", archiveHandler(cnf))
  http.HandleFunc("/", siteHandler(cnf))

  fmt.Println(fmt.Sprint("http://127.0.0.1:", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(":", port), nil)
}
