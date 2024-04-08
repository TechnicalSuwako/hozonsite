package main

import (
  "fmt"
  "os"
  "strings"
  "strconv"
)

var sofname = "hozonsite"
var version = "1.1.1"

func help() {
  fmt.Println("使い方：")
  fmt.Println(sofname + " -v               ：バージョンを表示")
  fmt.Println(
    sofname +
    " -s [ポート番号]  ：ポート番号でウェブサーバーを実行（デフォルト＝9920）",
  )
  fmt.Println(sofname + " -h               ：ヘルプを表示")
  fmt.Println(sofname + " <URL>            ：コマンドラインでウェブサイトを保存")
}

func saveurlcmd(url string, cnf Config) {
  // 結局HTTPかHTTPSじゃないわね…
  if !checkprefix(url) {
    fmt.Println("URLは不正です。終了…")
    return
  }

  // パラメートルの文字（?、=等）を削除
  eurl := stripurl(url)

  // 既に/usr/local/share/hozonsite/archiveに存在するかどうか
  exist := checkexist(eurl, cnf.datapath)

  // 既に存在したら、使う
  var confirm string

  // あ、既に存在する
  if len(exist) > 0 {
    fmt.Println("このページが既に保存されているみたいです。")
    fmt.Println("本当に手続きましょうか？ [y/N]")

    // 既に存在するページのURLを表示
    for _, ex := range exist {
      fmt.Println(strings.Replace(ex, cnf.datapath, cnf.domain, 1))
    }
    fmt.Scanf("%s", &confirm)
  }

  // 存在しない OR 「本当に手続きましょうか？」でYを入力した場合
  if len(exist) == 0 || confirm == "y" || confirm == "Y" {
    path := mkdirs(eurl, cnf.datapath)
    // ページをダウンロード
    getpage(url, path)
    // 色々の必須な編集
    scanpage(path, eurl, cnf.datapath)
    // 新しいURLを表示
    fmt.Println(cnf.domain + strings.Replace(path, cnf.datapath, "", 1))
  }
}

func main() {
  // コンフィグファイル
  cnf, err := getconf()
  if err != nil {
    fmt.Println(err)
    return
  }

  // コマンドラインのパラメートル
  args := os.Args

  if len(args) == 2 {
    // バージョンを表示
    if args[1] == "-v" {
      fmt.Println(sofname + "-" + version)
      return
    } else if args[1] == "-s" { // :9920でウェブサーバーを実行
      serv(cnf, 9920)
    } else if args[1] == "-h" { // ヘルプを表示
      help()
      return
    } else {
      // コマンドラインでウェブサイトを保存
      saveurlcmd(args[1], cnf)
      return
    }
  } else if len(args) == 3 && args[1] == "-s" {
    // 好みなポート番号でウェブサーバーを実行
    // でも、数字じゃないかもしん
    if port, err := strconv.Atoi(args[2]); err != nil {
      fmt.Printf("%qは数字ではありません。\n", args[2])
      return
    } else {
      // OK、実行しよ〜
      serv(cnf, port)
    }
  } else {
    // パラメートルは不明の場合、ヘルプを表示
    help()
    return
  }
}
