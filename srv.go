package main

import (
  "text/template"
  "fmt"
  "net/http"
)

type Page struct {
  Tit string
  Err string
  Lan string
  Ver string
  Ext []string // 既に存在する場合
  Url string // 確認ページ用
}

func serv (cnf Config, port int) {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  ftmpl := []string{cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"}

  /*http.HandleFunc("/exist", func(w http.ResponseWriter, r *http.Request) {
    data := &Page{Tit: "トップ", Ver: version}
    cookie, err := r.Cookie("lang")
    if err != nil {
      data.Lan = "ja"
    } else {
      data.Lan = cookie.Value
    }

    if data.Lan == "en" {
      data.Tit = "Top"
    }
  })*/

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data := &Page{Tit: "トップ", Ver: version}
    cookie, err := r.Cookie("lang")
    if err != nil {
      data.Lan = "ja"
    } else {
      data.Lan = cookie.Value
    }

    if data.Lan == "en" {
      data.Tit = "Top"
    }
    //tmpl := template.Must(template.ParseFiles(cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))

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
        url := r.PostForm.Get("hozonsite")
        // HTTPかHTTPSじゃない場合
        if !checkprefix(url) {
          if data.Lan == "ja" {
            data.Err = "URLは「http://」又は「https://」で始めます。"
          } else {
            data.Err = "The URL should start with \"http://\" or \"https://\"."
          }
          tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
        } else {
          //if r.PostForm.Get("sosin") != "" {}
        }
      }
    }

    tmpl.Execute(w, data)
    data = nil
  })

  fmt.Println(fmt.Sprint("http://127.0.0.1:", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(":", port), nil)
}
