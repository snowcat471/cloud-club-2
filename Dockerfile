# Build Stage
FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /app
COPY ./src .

ENV GO111MODULE=on

RUN go get 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o app .

# Running Stage
FROM alpine

WORKDIR /app
COPY --from=builder /app .

ENTRYPOINT [ "./app" ]