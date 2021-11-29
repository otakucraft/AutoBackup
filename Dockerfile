FROM golang:1.17 AS builder

WORKDIR /go/src/app
COPY . .

RUN apt update && apt install -y upx
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o backup
RUN upx backup

FROM debian:10 AS runner

WORKDIR /go/bin
EXPOSE 8088
COPY --from=builder /go/src/app/backup /go/bin/backup
HEALTHCHECK --timeout=10s CMD curl --silent --fail http://127.0.0.1:8088/healthcheck

CMD ["./backup", "-p", "8088"]
