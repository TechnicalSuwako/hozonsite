package main

import (
  "os/exec"
  "time"
  "fmt"
  "strings"
  "path/filepath"
)

func checkexist (url string, prefix string) []string {
  res, err := filepath.Glob(prefix + "/archive/*" + url2path(url))
  if err != nil {
    fmt.Println("Err:", err)
  }
  return res
  //fmt.Println("Exist?", res)
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
  cmd := exec.Command("mkdir", "-p", path)

  cmd.Run()

  return path
}

//func getpage (url string, path string) {}

//func scanpage (path string) {}
