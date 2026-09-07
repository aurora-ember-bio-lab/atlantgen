package main

import (
  "encoding/json"
  "log"
  "net/http"
)

type MigrationRequest struct {
  Source string `json:"source"`
}

func main() {
  http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"status":"ok"}`))
  })

  http.HandleFunc("/migrations", func(w http.ResponseWriter, r *http.Request) {
    var req MigrationRequest
    json.NewDecoder(r.Body).Decode(&req)

    log.Println("Queue migration:", req.Source)

    w.Write([]byte(`{"queued":true}`))
  })

  log.Println("API listening on :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
