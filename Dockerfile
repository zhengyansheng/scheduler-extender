FROM golang:1.19.8-alpine3.17 as builder

ARG APP
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${APP} main.go


FROM alpine:3.17 as runner

ARG APP

COPY --from=builder /app/${APP} /app

ENTRYPOINT ["/app"]
CMD ["-port=8000"]