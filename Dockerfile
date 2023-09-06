FROM golang:1.19.8-alpine3.17 as builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/scheduler-extender main.go


FROM alpine:3.17 as runner

COPY --from=builder /app/bin/scheduler-extender ./scheduler-extender

ENTRYPOINT ["./scheduler-extender"]
CMD ["-port=8081"]