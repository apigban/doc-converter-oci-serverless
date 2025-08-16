# ---- Build Stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy modules and download dependencies first for build caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the function binary
RUN CGO_ENABLED=0 go build -o /app/func ./functions/${FUNCTION_NAME}/main.go

# ---- Final Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/func .
ENTRYPOINT ["/app/func"]