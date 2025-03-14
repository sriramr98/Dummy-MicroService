# Start from the official Go image
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files first (if they exist)
COPY go.mod go.sum* ./

# Download dependencies (if go.mod exists)
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server

# Use a minimal alpine image for the final container
FROM alpine:3.19

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /go/bin/server .

# Expose the port the service runs on
EXPOSE 8080

# Run the binary
CMD ["./server"]
