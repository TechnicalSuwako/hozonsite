package main

import (
  "fmt"
  "os"
  "strings"
  "strconv"
)

var version = "1.0.0"

func help () {
  fmt.Println("使い方：");
  fmt.Println("hozonsite -v               ：バージョンを表示");
  fmt.Println("hozonsite -s [ポート番号]  ：ポート番号でウェブサーバーを実行（デフォルト＝9920）");
  fmt.Println("hozonsite -h               ：ヘルプを表示");
  fmt.Println("hozonsite <URL>            ：コマンドラインでウェブサイトを保存");
}

func main () {
  cnf := getconf()
  args := os.Args
  if len(args) == 2 {
    if args[1] == "-v" {
      fmt.Println("hozonsite-" + version)
      return
    } else if args[1] == "-s" {
      serv(cnf, 9920)
    } else if args[1] == "-h" {
      help()
      return
    } else {
      if checkprefix(args[1]) {
        eurl := stripurl(args[1])
        exist := checkexist(eurl, cnf.datapath)
        var confirm string
        if len(exist) > 0 {
          fmt.Println("このページが既に保存されているみたいです。")
          fmt.Println("本当に手続きましょうか？ [y/N]")
          for _, ex := range exist {
            fmt.Println(strings.Replace(ex, cnf.datapath, cnf.domain, 1))
          }
          fmt.Scanf("%s", &confirm)
        }
        if len(exist) == 0 || confirm == "y" || confirm == "Y" {
          path := mkdirs(eurl, cnf.datapath)
          getpage(args[1], path)
          scanpage(path, eurl, cnf.datapath)
          fmt.Println(cnf.domain + strings.Replace(path, cnf.datapath, "", 1))
        }
        return
      } else {
        fmt.Println("URLは不正です。終了…")
        return
      }
    }
  } else if len(args) == 3 && args[1] == "-s" {
    if port, err := strconv.Atoi(args[2]); err != nil {
      fmt.Printf("%qは数字ではありません。\n", args[2])
      return
    } else {
      serv(cnf, port)
    }
  } else {
    help()
    return
  }
}
