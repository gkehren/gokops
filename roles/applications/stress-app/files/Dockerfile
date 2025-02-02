# Stage 1: Build the Go binary.
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies first for caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Stage 2: Run the binary in a minimal container.
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/main .

# Expose the port your application listens on.
EXPOSE 4242

# Run the binary.
CMD ["./main"]
