FROM golang:1.19-alpine3.17 AS builder

WORKDIR /app

COPY . .

RUN go build -o /main ./cmd


FROM alpine:3.17

EXPOSE 9080

WORKDIR /app
COPY --from=builder /main .

ENTRYPOINT [ "/app/main" ]