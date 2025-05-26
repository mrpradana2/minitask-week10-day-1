# Stage 1: Build
FROM golang:1.24.3-alpine3.21 AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .
COPY ./public ./public

# Build the application
RUN go build -o server ./cmd

# Stage 2: Run
FROM alpine:3.21.3

# Set working directory
WORKDIR /app

# Copy the built binary from builder stage
COPY --from=builder /app/server .

# Expose port (ganti sesuai port aplikasi)
EXPOSE 8080

# Command to run the binary
CMD ["./server"]
