# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Set necessary environment variables for CGO and Go modules
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# The output binary will be named 'godoctor' and placed in /app/bin
RUN go build -o bin/godoctor ./cmd/godoctor

# Stage 2: Runner
FROM golang:1.24-alpine

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/godoctor .

# Expose any necessary ports (if applicable, e.g., for a web server)
# EXPOSE 8080

# Define the entrypoint for the application
ENTRYPOINT ["./godoctor"]

