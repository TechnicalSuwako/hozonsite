package main

import (
  "encoding/json"
  "fmt"
)

func getlist (lang string) []byte {
  var jloc = []byte(`{
    "top": "トップ",
    "errfuseiurl": "URLは「http://」又は「https://」で始めます。",
    "errfusei": "不正なエラー。"
  }`)
  var eloc = []byte(`{
    "top": "Top",
    "errfuseiurl": "The URL should start with \"http://\" or \"https://\".",
    "errfusei": "Unknown error."
  }`)

  if lang == "en" { return eloc }
  return jloc
}

func getloc (str string, lang string) string {
  var payload map[string]interface{}
  err := json.Unmarshal(getlist(lang), &payload)
  if err != nil {
    fmt.Println("loc:", err)
    return ""
  }

  for k, v := range payload {
    if str == k {
      return v.(string)
    }
  }

  return ""
}
