package main

import (
  "os/exec"
  "time"
  "fmt"
  "strings"
)

func mkdirs (url string, prefix string) string {
  rep := ""
  t := time.Now().Unix()

  if strings.HasPrefix(url, "https:/") {
    rep = strings.Replace(url, "https:/", "", 1)
  } else {
    rep = strings.Replace(url, "http:/", "", 1)
  }

  if strings.HasSuffix(rep, "/") {
    rep = strings.TrimSuffix(rep, "/")
  }

  path := fmt.Sprint(prefix, "/archive/", t, rep)
  cmd := exec.Command("mkdir", "-p", path)

  cmd.Run()

  return path
}

//func getpage (url string, path string) {}

//func scanpage (path string) {}
