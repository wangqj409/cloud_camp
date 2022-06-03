package part_2

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

func responseHeaders(w http.ResponseWriter, r *http.Request) {
    for k := range r.Header {
        val := r.Header.Get(k)

        // w.Header().Set(k, val)
        // 🤔️？？？？  Set 字段不一定能被设置上，Add 很稳定
        w.Header().Add(k, val)
        io.WriteString(w, fmt.Sprintf("Got Header:%s value:%v\n", k, val))
    }

    fmt.Fprintf(os.Stdout, "%s %s %s %d\n", time.Now().Format(time.RFC3339), r.URL.String(), r.RemoteAddr, http.StatusOK)
}

func healthz(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "200")
}

func version(w http.ResponseWriter, r *http.Request) {

    // k := r.PostFormValue("key")
    // w.Write([]byte(os.Getenv(k)))

    version := os.Getenv("VERSION")
    w.Write([]byte(version))

    w.Header().Add("VERSION", version)

}
