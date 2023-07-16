package main

import (
  "fmt"
  "os"
  "strings"
  "strconv"
)

var version = "1.1.0"

func help () {
  fmt.Println("使い方：");
  fmt.Println("hozonsite -v               ：バージョンを表示");
  fmt.Println("hozonsite -s [ポート番号]  ：ポート番号でウェブサーバーを実行（デフォルト＝9920）");
  fmt.Println("hozonsite -h               ：ヘルプを表示");
  fmt.Println("hozonsite <URL>            ：コマンドラインでウェブサイトを保存");
}

func main () {
  cnf := getconf() // コンフィグファイル
  args := os.Args // コマンドラインのパラメートル
  if len(args) == 2 {
    if args[1] == "-v" { // バージョンを表示
      fmt.Println("hozonsite-" + version)
      return
    } else if args[1] == "-s" { // :9920でウェブサーバーを実行
      serv(cnf, 9920)
    } else if args[1] == "-h" { // ヘルプを表示
      help()
      return
    } else { // コマンドラインでウェブサイトを保存
      if checkprefix(args[1]) { // プロトコールはあってるかどうか
        eurl := stripurl(args[1]) // パラメートルの文字（?、=等）を削除
        exist := checkexist(eurl, cnf.datapath) // 既に/usr/share/hozonsite/archiveに存在するかどうか
        var confirm string // 既に存在したら、使う
        if len(exist) > 0 { // あ、既に存在する
          fmt.Println("このページが既に保存されているみたいです。")
          fmt.Println("本当に手続きましょうか？ [y/N]")
          for _, ex := range exist { // 既に存在するページのURLを表示
            fmt.Println(strings.Replace(ex, cnf.datapath, cnf.domain, 1))
          }
          fmt.Scanf("%s", &confirm)
        }
        if len(exist) == 0 || confirm == "y" || confirm == "Y" { // 存在しない OR 「本当に手続きましょうか？」でYを入力した場合
          path := mkdirs(eurl, cnf.datapath)
          getpage(args[1], path) // ページをダウンロード
          scanpage(path, eurl, cnf.datapath) // 色々の必須な編集
          fmt.Println(cnf.domain + strings.Replace(path, cnf.datapath, "", 1)) // 新しいURLを表示
        }
        return
      } else { // 結局HTTPかHTTPSじゃないわね…
        fmt.Println("URLは不正です。終了…")
        return
      }
    }
  } else if len(args) == 3 && args[1] == "-s" { // 好みなポート番号でウェブサーバーを実行
    if port, err := strconv.Atoi(args[2]); err != nil { // でも、数字じゃないかもしん
      fmt.Printf("%qは数字ではありません。\n", args[2])
      return
    } else { // OK、実行しよ〜
      serv(cnf, port)
    }
  } else { // パラメートルは不明の場合、ヘルプを表示
    help()
    return
  }
}
