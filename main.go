package main

import (
  "net/http"
  "fmt"

  "github.com/a-h/templ"
)

func main() {
  component := hello("John")
  http.Handle("/", templ.Handler(component))

  fmt.Println("Listening on :3000")
  http.ListenAndServe(":3000", nil)
}
