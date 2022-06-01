package part_2

import (
    "fmt"
    "io"
    "net/http"
)

var (
    serv *http.ServeMux
)

func showHeaders(w http.ResponseWriter, r *http.Request) {
    for k := range r.Header {
        io.WriteString(w, fmt.Sprintf("%s %v\n", k, r.Header.Get(k)))
    }
}
func healthz(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "200")
}

func NewServer() {
    serv = http.NewServeMux()

    serv.HandleFunc("/headers", showHeaders)
    serv.HandleFunc("/healthz", healthz)

    http.ListenAndServe(":9002", serv)
}
