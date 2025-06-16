FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY delete_service/ .

RUN GOOS=linux GOARCH=amd64 go build -o delete_service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/delete_service .

EXPOSE 3030

CMD ["./delete_service"]
