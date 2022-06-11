FROM golang:1.17 as builder

WORKDIR /www/src/

COPY http_server.go /www/src/
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN GOOS=linux go build -installsuffix cgo -o http_server http_server.go


FROM scratch
COPY --from=builder /www/src/http_server /

ENTRYPOINT ["/http_server"]


