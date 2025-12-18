# -------- Build stage --------
FROM golang:1.23-alpine AS builder

# Log base image
RUN echo "🐹 Using golang:1.23-alpine as builder image"

WORKDIR /app
RUN echo "📁 Working directory set to /app"

# Copy go.mod and go.sum first (better caching)
RUN echo "📦 Copying go.mod and go.sum"
COPY go.mod go.sum ./

RUN echo "⬇️  Downloading Go module dependencies"
RUN go mod download

# Copy the rest of the source
RUN echo "📂 Copying application source code"
COPY . .

# Build the binary
RUN echo "🔨 Building Go binary (linux/amd64, static)"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -o gobank ./command/server

RUN echo "✅ Build complete: /app/gobank"

# -------- Runtime stage --------
FROM alpine:latest

RUN echo "🪶 Using minimal alpine runtime image"

WORKDIR /app
RUN echo "📁 Runtime working directory set to /app"

# Copy binary from builder
RUN echo "📤 Copying binary from builder stage"
COPY --from=builder /app/gobank .

# Expose API port
RUN echo "🌐 Exposing port 8081"
EXPOSE 8081

# Runtime log (visible when container starts)
CMD echo "🚀 Starting GObank API on port 8081" && ./gobank
