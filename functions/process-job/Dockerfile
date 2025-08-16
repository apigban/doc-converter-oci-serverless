# ---- Build Stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy modules and download dependencies first for build caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the function binary
# This assumes the Docker build context is the repository root
RUN cd functions/${FUNCTION_NAME} && CGO_ENABLED=0 go build -o /app/func .

# ---- Final Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/func .
ENTRYPOINT ["/app/func"]