package main

import (
  "net/http"
  "fmt"

  "github.com/a-h/templ"
)

func main() {
  component := hello("world")

  http.HandleFunc("/", indexHandler)
  http.Handle("/contacts", templ.Handler(component))

  fmt.Println("Listening on :3000")
  http.ListenAndServe(":3000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/contacts", http.StatusFound)
  return
}
