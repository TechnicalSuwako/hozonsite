package main

import (
  "fmt"
  "os"
  "strings"
  "strconv"
)

var version = "1.0.0"

func checkprefix (url string) bool {
  return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

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
        path := mkdirs(args[1], cnf.datapath)
        // TODO: ページの保存
        //getpage(args[1], path)
        // TODO: ページの確認
        //scanpage(path)
        fmt.Println(cnf.domain + strings.Replace(path, cnf.datapath, "", 1))
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
