FROM golang:1.22 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY put_service/ .

RUN GOOS=linux GOARCH=amd64 go build -o put_service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/put_service .

EXPOSE 3030

CMD ["./put_service"]
