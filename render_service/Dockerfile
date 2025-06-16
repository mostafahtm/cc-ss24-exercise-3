# Stage 1: Build
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum from the root
COPY go.mod go.sum ./
RUN go mod download

# Copy only the render_service files
COPY render_service/main.go .
COPY render_service/views ./views
COPY render_service/css ./css

# Build the service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o render_service .


# Stage 2: Final image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/render_service .
COPY --from=builder /app/views ./views
COPY --from=builder /app/css ./css

EXPOSE 3030

CMD ["./render_service"]
