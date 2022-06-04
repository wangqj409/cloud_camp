package part_2

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

var (
    port = ":9002" // listen port
)

/**
- 接收客户端 request，并将 request 中带的 header 写入 response header
- 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
- Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
- 当访问 localhost/healthz 时，应返回 200
*/

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


func NewServer() {
    serv := http.NewServeMux()

    serv.HandleFunc("/responseHeaders", responseHeaders)
    serv.HandleFunc("/version", version)
    serv.HandleFunc("/healthz", healthz)

    http.ListenAndServe(port, serv)
}
