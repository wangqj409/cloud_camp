package main

import (
    "fmt"
    "github.com/prometheus/client_golang/prometheus"
    "log"

    "github.com/prometheus/client_golang/prometheus/promhttp"
    "io"
    "math/rand"
    "net/http"
    "os"
    "time"
)

const (
    MetricsNamespace = "httpserver"
)

var (
    port = ":80" // listen port

    functionLatency = CreateExecutionTimeMetric(MetricsNamespace,
        "Time spent.")
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

func main() {
    serv := http.NewServeMux()

    serv.HandleFunc("/responseHeaders", responseHeaders)
    serv.HandleFunc("/version", version)
    serv.HandleFunc("/healthz", healthz)
    serv.Handle("/metrics", promhttp.Handler())

    http.ListenAndServe(port, serv)
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
    return prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Namespace: namespace,
            Name:      "execution_latency_seconds",
            Help:      help,
            Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
        }, []string{"step"},
    )
}

func images(w http.ResponseWriter, r *http.Request) {
    timer := NewTimer()
    defer timer.ObserveTotal()
    randInt := rand.Intn(2000)
    time.Sleep(time.Millisecond * time.Duration(randInt))
    w.Write([]byte(fmt.Sprintf("<h1>%d<h1>", randInt)))
}

type ExecutionTimer struct {
    histo *prometheus.HistogramVec
    start time.Time
    last  time.Time
}

func Register() {
    log.Print("register function latency prometheus metric, %v", functionLatency)
    err := prometheus.Register(functionLatency)
    if err != nil {
        log.Printf("error register: %v", err)
    } else {
        log.Print(err)
    }

}

func NewTimer() *ExecutionTimer {
    return NewExecutionTimer(functionLatency)
}

func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
    now := time.Now()
    return &ExecutionTimer{
        histo: histo,
        start: now,
        last:  now,
    }
}

func (t *ExecutionTimer) ObserveTotal() {
    log.Print("report execution time after function done")
    (*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
