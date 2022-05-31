package part_2

import (
    "io"
    "net/http"
)

var (
    serv *http.ServeMux
)

func showHeaders(w http.ResponseWriter, r *http.Request) {
    for k := range r.Header {
        io.WriteString(w, k)

    }
}

func NewServer() {
    serv = http.NewServeMux()

    serv.HandleFunc("/headers", showHeaders)

    http.ListenAndServe(":9002", serv)
}
