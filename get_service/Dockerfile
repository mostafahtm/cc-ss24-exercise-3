# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum from root
COPY go.mod go.sum ./
RUN go mod download

# Copy only the service source
COPY get_service/ .                  

RUN GOOS=linux GOARCH=amd64 go build -o get_service .

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/get_service .

EXPOSE 3030

CMD ["./get_service"]
