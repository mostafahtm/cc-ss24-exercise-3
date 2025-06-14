FROM golang:1.22 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY post_service/ .

RUN GOOS=linux GOARCH=amd64 go build -o post_service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/post_service .

EXPOSE 3030

CMD ["./post_service"]
